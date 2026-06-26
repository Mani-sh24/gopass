package main

import (
	"example/web-service-gin/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserReq struct {
	Email    string `json:"email"`
	Mpin     string `json:"mpin"`
	Enc_Key  string `json:"enc_key"`
	Password string `json:"password"`
}
type UserRes struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}
type UserSuccess struct {
	Msg   string `json:"msg"`
	Token string `json:"token"` // jwt token
}

var users = []UserReq{
	{Email: "abc@gmail.com", Mpin: "2422", Enc_Key: "Manish1@"},
}

func createUser(c *gin.Context) {
	var new_user UserReq
	// recieving the body and binding to the new_user struct
	if err := c.BindJSON(&new_user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt create user"})
		return
	}
	// creating the uuid to store in db
	id := uuid.New()
	// save the data in db and return the success message
	_, err := database.DB.Exec(
		"INSERT INTO users (id , email, password, mpin , enc_key) VALUES(?,?,?,?,?)",
		id.String(),
		new_user.Email,
		new_user.Password,
		new_user.Mpin,
		new_user.Enc_Key,
	)
	// if for some reason the data is not saved in db send the error response
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt create user"})
		return
	}
	// send the confirm message
	c.IndentedJSON(http.StatusOK, UserSuccess{Msg: "Success", Token: id.String()}) // sendong id temporarily
}
func getAllUsers(c *gin.Context) {
	rows, err := database.DB.Query("SELECT email , id FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	defer rows.Close()
	var users []UserRes

	for rows.Next() {
		var user UserRes
		err := rows.Scan(&user.Id, &user.Email)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		users = append(users, user)
	}
	c.IndentedJSON(http.StatusOK, users)
}

func getUserById(c *gin.Context) {
	var user UserRes
	id_p := c.Param("id")
	err := database.DB.QueryRow("SELECT email , id FROM users WHERE id = ?", id_p).Scan(&user.Id, &user.Email)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt get user"})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}
func main() {
	database.Connect_to_db()
	database.Init()
	router := gin.Default()
	{
		auth := router.Group("/auth")
		auth.POST("/register", createUser)
		auth.GET("/getuser/:id", getUserById)
	}
	router.POST("/all", getAllUsers)

	router.Run("localhost:8080")
}
