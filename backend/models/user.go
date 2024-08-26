package models

type User struct {
	Name     string
	Email    string
	Password string // Hashed
	Id       int
}