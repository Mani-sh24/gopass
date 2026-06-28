package dtos

type PasswordReq struct {
	Title    string `json:"title"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PasswordRes struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Email    string `json:"email"`
	Password string `json:"password"`
	User_id  string `json:"uid"`
}
