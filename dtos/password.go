package dtos

type PasswordReq struct {
	Title    string `json:"title"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
