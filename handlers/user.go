package handlers

import (
	schema "example/web-service-gin/Models"
	"example/web-service-gin/database"
	"example/web-service-gin/dtos"
	"example/web-service-gin/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var new_user dtos.UserReq
	// recieving the body and binding to the new_user struct
	if err := c.ShouldBind(&new_user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt create user"})
		return
	}
	// creating the uuid to store in db
	id := uuid.New().String()
	// pasword hashing and salting
	password_hash, err_p := bcrypt.GenerateFromPassword([]byte(new_user.Password), bcrypt.DefaultCost)
	if err_p != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt create user"})
		return
	}
	mpin_hash, err_m := bcrypt.GenerateFromPassword([]byte(new_user.Mpin), bcrypt.MinCost)
	if err_m != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt create user"})
		return
	}
	jwt_tok, err := helpers.GenerateJWT(id, new_user.Email)
	// save the data in db and return the success message
	_, err = database.DB.Exec(
		"INSERT INTO users (id , email, password, mpin) VALUES(?,?,?,?)",
		id,
		new_user.Email,
		password_hash,
		mpin_hash,
	)
	// if for some reason the data is not saved in db send the error response
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt create user"})
		return
	}
	// send the confirm message

	c.IndentedJSON(http.StatusOK, dtos.UserSuccess{Msg: "Success", Token: jwt_tok}) // sendong id temporarily
}

func GetAllUsers(c *gin.Context) {
	rows, err := database.DB.Query("SELECT email , id FROM users") // Query will select every rows it will point at starting
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	defer rows.Close() // telling db to free up resources
	var users []dtos.UserRes

	for rows.Next() { // moves the pointer to next row
		var user dtos.UserRes
		err := rows.Scan(&user.Email, &user.Id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Bad Request"})
			return
		}
		users = append(users, user)
	}
	c.IndentedJSON(http.StatusOK, users)
}

func Login(c *gin.Context) {
	var user_login schema.UserModel
	if err := c.ShouldBindJSON(&user_login); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt Login check valid credentials"})
		return
	}
	var dbID, dbEmail, dbPassword, dbMpin string
	err := database.DB.QueryRow("Select id, email, password, mpin FROM users WHERE email = ? ", user_login.Email).Scan(&dbID, &dbEmail, &dbPassword, &dbMpin)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Invalid credentials"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbMpin), []byte(user_login.Mpin))
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"msg": "Invalid credentials",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(user_login.Password))
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"msg": "Invalid credentials",
		})
		return
	}
	jwt_tok, err := helpers.GenerateJWT(dbID, dbEmail)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"msg": "Could not generate token",
		})
		return
	}
	c.IndentedJSON(http.StatusOK, dtos.UserSuccess{
		Msg:   "Login successful",
		Token: jwt_tok,
	})
}

func GetUserById(c *gin.Context) {
	var user dtos.UserRes
	id_p := c.GetString("id")
	err := database.DB.QueryRow("SELECT email , password FROM users WHERE id = ?", id_p).Scan(&user.Email, &user.Id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt get user You are not authorised!"})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}
