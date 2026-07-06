package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"example/web-service-gin/database"
	"example/web-service-gin/dtos"
	"example/web-service-gin/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser godoc
// @Summary      Register a new user
// @Description  Register a new user in the system
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        payload  body      dtos.UserRegisterReq  true  "Registration payload"
// @Success      200      {object}  dtos.AuthRes
// @Failure      400      {object}  dtos.MessageRes "Could not create user"
// @Router       /auth/register [post]
func CreateUser(c *gin.Context) {
	var new_user dtos.UserRegisterReq
	// recieving the body and binding to the new_user struct
	if err := c.ShouldBind(&new_user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Couldnt create user"})
		return
	}
	// creating the uuid to store in db
	id := uuid.New().String()
	// pasword hashing and salting
	password_hash, err_p := bcrypt.GenerateFromPassword([]byte(new_user.Password), bcrypt.DefaultCost)
	if err_p != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Couldnt create user"})
		return
	}
	mpin_hash, err_m := bcrypt.GenerateFromPassword([]byte(new_user.Mpin), bcrypt.MinCost)
	if err_m != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Couldnt create user"})
		return
	}
	jwt_tok, err := helpers.GenerateJWT(id, new_user.Email)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Couldnt create user"})
		return
	}
	// Generate salt
	saltBytes := make([]byte, 16)
	if _, err := rand.Read(saltBytes); err != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Couldnt create user"})
		return
	}
	salt := hex.EncodeToString(saltBytes)

	// save the data in db and return the success message
	_, err = database.DB.Exec(
		"INSERT INTO users (id , email, password, mpin, salt) VALUES(?,?,?,?,?)",
		id,
		new_user.Email,
		password_hash,
		mpin_hash,
		salt,
	)
	// if for some reason the data is not saved in db send the error response
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Couldnt create user"})
		return
	}
	// send the confirm message
	c.IndentedJSON(http.StatusOK, dtos.AuthRes{Msg: "Success", Token: jwt_tok, Salt: salt})
}


// Login godoc
// @Summary      Authenticate user
// @Description  Verifies credentials and generates a session JWT Token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        payload  body      dtos.UserLoginReq  true  "Login credentials"
// @Success      200      {object}  dtos.AuthRes
// @Failure      400      {object}  dtos.MessageRes "Invalid input credentials"
// @Failure      401      {object}  dtos.MessageRes "Unauthorized"
// @Failure      500      {object}  dtos.MessageRes "Internal Server Error"
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var user_login dtos.UserLoginReq
	if err := c.ShouldBindJSON(&user_login); err != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Couldnt Login check valid credentials"})
		return
	}
	var dbID, dbEmail, dbPassword, dbMpin, dbSalt string
	err := database.DB.QueryRow("Select id, email, password, mpin, salt FROM users WHERE email = ? ", user_login.Email).Scan(&dbID, &dbEmail, &dbPassword, &dbMpin, &dbSalt)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Invalid credentials"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbMpin), []byte(user_login.Mpin))
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, dtos.MessageRes{
			Msg: "Invalid credentials",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(user_login.Password))
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, dtos.MessageRes{
			Msg: "Invalid credentials",
		})
		return
	}
	jwt_tok, err := helpers.GenerateJWT(dbID, dbEmail)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, dtos.MessageRes{
			Msg: "Could not generate token",
		})
		return
	}
	c.IndentedJSON(http.StatusOK, dtos.AuthRes{
		Msg:   "Login successful",
		Token: jwt_tok,
		Salt:  dbSalt,
	})
}

// GetUserById godoc
// @Summary      Get user profile
// @Description  Retrieve the logged-in user profile info using the user ID loaded from the JWT claims.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer JWT_TOKEN"
// @Success      200            {object}  dtos.UserProfileRes
// @Failure      400            {object}  dtos.MessageRes "Unauthorized"
// @Failure      401            {object}  dtos.MessageRes "Unauthorized"
// @Router       /protected/getuser [get]
func GetUserById(c *gin.Context) {
	var user dtos.UserProfileRes
	id_p := c.GetString("id")
	err := database.DB.QueryRow("SELECT id, email FROM users WHERE id = ?", id_p).Scan(&user.Id, &user.Email)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Couldnt get user You are not authorised!"})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}
