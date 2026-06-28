package dtos

type PasswordReq struct {
	Title    string `json:"title" binding:"required,min=1,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

type PasswordRes struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Email    string `json:"email"`
	Password string `json:"password"`
	User_id  string `json:"uid"`
}

type UpdatePasswordReq struct {
	Title    *string `json:"title" binding:"omitempty,min=1,max=100"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6,max=100"`
}
