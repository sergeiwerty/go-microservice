package models

type UserStorer interface {
}

type User struct {
	ID       int
	Email    string
	Password string
}
