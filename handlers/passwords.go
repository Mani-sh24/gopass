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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Couldnt create password"})
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
	id := c.GetString("id")                                                                                     // userid from middleware
	rows, err := database.DB.Query("SELECT id , title , email , password FROM passwords WHERE user_id = ?", id) // Query will select every rows it will point at starting
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Could not retrieve passwords",
		})
		return
	}
	defer rows.Close() // telling db to free up resources
	var passwordList []dtos.PasswordRes

	for rows.Next() { // moves the pointer to next row
		var password dtos.PasswordRes
		err := rows.Scan(&password.Id, &password.Title, &password.Email, &password.Password)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Could not retrieve passwords"})
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
	// sql doest return error if the deleted rows are 0 so we check how many rows were affected
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

func UpdatePassword(c *gin.Context) {
	uid := c.GetString("id") // User ID from JWT
	id := c.Param("id")      // Password record ID from path parameter

	var req dtos.UpdatePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body"})
		return
	}

	// 1. Build the update statement dynamically
	query := "UPDATE passwords SET "
	var args []interface{}

	if req.Title != nil {
		query += "title = ?, "
		args = append(args, *req.Title)
	}
	if req.Email != nil {
		query += "email = ?, "
		args = append(args, *req.Email)
	}
	if req.Password != nil {
		query += "password = ?, "
		args = append(args, *req.Password)
	}

	// If no fields were provided in the request
	if len(args) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Nothing to update"})
		return
	}

	// Remove the trailing comma and space from query, then append WHERE clause
	query = query[:len(query)-2] + " WHERE id = ? AND user_id = ?"
	args = append(args, id, uid)

	// 2. Execute the update
	res, err := database.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Could not update password record"})
		return
	}

	// 3. Verify that the row actually existed and belonged to the user
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Password record not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Password updated successfully"})
}
