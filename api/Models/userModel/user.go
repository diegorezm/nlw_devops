package usermodel

import (
	"time"
)

type UserModel struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserModel(id int, username, email, password string, createdAt, updatedAt time.Time) *UserModel {
	return &UserModel{Id: id, Username: username, Password: password, CreatedAt: createdAt, UpdatedAt: updatedAt}
}
