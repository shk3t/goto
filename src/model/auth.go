package model

type User struct {
	Id       int
	Login    string
	Password string
	IsAdmin  bool
}