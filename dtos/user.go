package dtos

type UserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Enc_Key  string `json:"enc_key" binding:"required,min=5"`
	Mpin     string `json:"mpin" binding:"required,len=4"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type UserRes struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type UserSuccess struct {
	Msg   string `json:"msg"`
	Token string `json:"token"` // jwt token
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Mpin     string `json:"mpin" binding:"required,len=4"`
}
