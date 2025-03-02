package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID (id int) (*User, error)
	CreateUser (User) error
}

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=100"`
}

type LoginUserPayload struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}