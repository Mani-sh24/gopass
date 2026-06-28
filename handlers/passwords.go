package handlers

import (
	"example/web-service-gin/database"
	"example/web-service-gin/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreatePassword godoc
// @Summary      Create password record
// @Description  Create a password credential record for the authenticated user.
// @Tags         Passwords
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "Bearer JWT_TOKEN"
// @Param        payload        body      dtos.PasswordCreateReq  true  "Password payload"
// @Success      200            {object}  dtos.PasswordRes
// @Failure      400            {object}  dtos.MessageRes "Could not create password"
// @Failure      401            {object}  dtos.MessageRes "Unauthorized"
// @Router       /passwordcrud/create [post]
func CreatePassword(c *gin.Context) {
	// stores password hashed from client
	var password_req dtos.PasswordCreateReq
	p_id := uuid.New().String()
	if err := c.ShouldBind(&password_req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Couldnt create password"})
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
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Couldnt create password"})
		return
	}
	password_res := dtos.PasswordRes{
		Id:       p_id,
		UserId:   user_id,
		Email:    password_req.Email,
		Password: password_req.Password,
		Title:    password_req.Title,
	}

	c.IndentedJSON(http.StatusOK, password_res)
}

// GetAllPasswords godoc
// @Summary      List all password records
// @Description  Retrieve all password credential records belonging to the authenticated user.
// @Tags         Passwords
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer JWT_TOKEN"
// @Success      200            {array}   dtos.PasswordRes
// @Failure      400            {object}  dtos.MessageRes "Bad Request"
// @Failure      401            {object}  dtos.MessageRes "Unauthorized"
// @Failure      500            {object}  dtos.MessageRes "Internal Server Error"
// @Router       /passwordcrud/all [get]
func GetAllPasswords(c *gin.Context) {
	id := c.GetString("id")                                                                                     // userid from middleware
	rows, err := database.DB.Query("SELECT id , title , email , password FROM passwords WHERE user_id = ?", id) // Query will select every rows it will point at starting
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.MessageRes{
			Msg: "Could not retrieve passwords",
		})
		return
	}
	defer rows.Close() // telling db to free up resources
	var passwordList []dtos.PasswordRes

	for rows.Next() { // moves the pointer to next row
		var password dtos.PasswordRes
		err := rows.Scan(&password.Id, &password.Title, &password.Email, &password.Password)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Could not retrieve passwords"})
			return
		}
		password.UserId = id
		passwordList = append(passwordList, password)
	}
	c.IndentedJSON(http.StatusOK, passwordList)
}

// DeletePassword godoc
// @Summary      Delete a password record
// @Description  Deletes a specific password record by ID, verifying that it belongs to the authenticated user.
// @Tags         Passwords
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer JWT_TOKEN"
// @Param        id             path      string  true  "Password Record ID"
// @Success      200            {object}  dtos.MessageRes "Deleted"
// @Failure      400            {object}  dtos.MessageRes "Could not delete password"
// @Failure      401            {object}  dtos.MessageRes "Unauthorized"
// @Failure      404            {object}  dtos.MessageRes "Password not found or not authorized to delete"
// @Router       /passwordcrud/delete/{id} [delete]
func DeletePassword(c *gin.Context) {
	uid := c.GetString("id")
	id := c.Param("id") // param from path (post id)
	res, err := database.DB.Exec("DELETE FROM passwords WHERE id = ? AND user_id = ?", id, uid)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Couldnt delete password"})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, dtos.MessageRes{Msg: "Could not verify deletion"})
		return
	}
	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, dtos.MessageRes{Msg: "Password not found or not authorized to delete"})
		return
	}
	c.IndentedJSON(http.StatusOK, dtos.MessageRes{Msg: "Deleted"})
}

// UpdatePassword godoc
// @Summary      Update a password record
// @Description  Partially update a specific password record's title, email, or password.
// @Tags         Passwords
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "Bearer JWT_TOKEN"
// @Param        id             path      string                  true  "Password Record ID"
// @Param        payload        body      dtos.PasswordUpdateReq  true  "Update payload"
// @Success      200            {object}  dtos.MessageRes "Password updated successfully"
// @Failure      400            {object}  dtos.MessageRes "Bad Request"
// @Failure      401            {object}  dtos.MessageRes "Unauthorized"
// @Failure      404            {object}  dtos.MessageRes "Password record not found or unauthorized"
// @Router       /passwordcrud/update/{id} [patch]
func UpdatePassword(c *gin.Context) {
	uid := c.GetString("id") // User ID from JWT
	id := c.Param("id")       // Password record ID from path parameter

	var req dtos.PasswordUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Invalid request body"})
		return
	}

	// Build update statement dynamically
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
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Nothing to update"})
		return
	}

	// Remove the trailing comma and space, append WHERE clause
	query = query[:len(query)-2] + " WHERE id = ? AND user_id = ?"
	args = append(args, id, uid)

	res, err := database.DB.Exec(query, args...)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.MessageRes{Msg: "Could not update password record"})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, dtos.MessageRes{Msg: "Could not verify update"})
		return
	}
	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, dtos.MessageRes{Msg: "Password record not found or unauthorized"})
		return
	}

	c.IndentedJSON(http.StatusOK, dtos.MessageRes{Msg: "Password updated successfully"})
}
