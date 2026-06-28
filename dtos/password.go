package dtos

// PasswordCreateReq defines the payload for creating a credentials record
type PasswordCreateReq struct {
	Title    string `json:"title" binding:"required,min=1,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// PasswordRes defines the response payload for a password credential
type PasswordRes struct {
	Id       string `json:"id"`
	UserId   string `json:"userId"`
	Title    string `json:"title"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// PasswordUpdateReq defines the payload for partially updating a password record
type PasswordUpdateReq struct {
	Title    *string `json:"title" binding:"omitempty,min=1,max=100"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6,max=100"`
}
