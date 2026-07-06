package models

type UserModel struct {
	Id       string
	Email    string
	Mpin     string
	Password string
	Salt     string
}
