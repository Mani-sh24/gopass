package handlers

import (
	"example/web-service-gin/database"
	"example/web-service-gin/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePassword(c *gin.Context) {
	// stores password hashed from client
	var password_req dtos.PasswordReq
	p_id := uuid.New().String()
	if err := c.ShouldBind(&password_req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt create paswword"})
		return
	}
	user_id := c.GetString("id")
	_, err := database.DB.Exec(
		"INSERT INTO passwords (user_id, id, title, email, password) VALUES(?,?,?,?,?)",
		user_id,
		p_id,
		password_req.Title,
		password_req.Email,
		password_req.Password,
	)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt create password"})
		return
	}
	password_res := dtos.PasswordRes{
		Id:       p_id,
		User_id:  user_id,
		Email:    password_req.Email,
		Password: password_req.Password,
		Title:    password_req.Title,
	}

	c.IndentedJSON(http.StatusOK, password_res)
}
func GetAllPasswords(c *gin.Context) {
	// var password_res dtos.PasswordRes
	id := c.GetString("id")                                                                                     // userid from middleware
	rows, err := database.DB.Query("SELECT id , title , email , password FROM passwords WHERE user_id = ?", id) // Query will select every rows it will point at starting
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	defer rows.Close() // telling db to free up resources
	var passwordList []dtos.PasswordRes

	for rows.Next() { // moves the pointer to next row
		var password dtos.PasswordRes
		err := rows.Scan(&password.Id, &password.Title, &password.Email, &password.Password)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		passwordList = append(passwordList, password)
	}
	c.IndentedJSON(http.StatusOK, passwordList)

}
func DeletePassword(c *gin.Context) {
	uid := c.GetString("id")
	id := c.Param("id") // param from path (post id)
	res, err := database.DB.Exec("DELETE FROM passwords WHERE id = ? AND user_id = ?", id, uid)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt delete password"})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "Could not verify deletion"})
		return
	}
	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Password not found or not authorized to delete"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"msg": "Deleted",
	})
}
