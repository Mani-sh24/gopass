package dtos

// UserRegisterReq defines the input for user registration
type UserRegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Mpin     string `json:"mpin" binding:"required,len=4"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// UserLoginReq defines the input for user login
type UserLoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Mpin     string `json:"mpin" binding:"required,len=4"`
}

// AuthRes defines the response containing the authentication token
type AuthRes struct {
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

// UserProfileRes defines the user profile response
type UserProfileRes struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

// MessageRes defines a generic success or informational message response
type MessageRes struct {
	Msg string `json:"msg"`
}
