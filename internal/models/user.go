package models

import (
	"time"
	"unicode/utf8"

	"github.com/fdemchenko/arcus/internal/validator"
)

type User struct {
	ID        int
	Name      string
	Email     string
	Activated bool
	Password  Password
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Password struct {
	Plain string
	Hash  []byte
}

func (user *User) Validate(v validator.Validator) {
	v.Check(user.Email != "", "email", "should not be empty")
	v.Check(validator.IsValidEmail(user.Email), "email", "should be valid email")

	v.Check(user.Name != "", "name", "should not be empty")
	v.Check(len(user.Name) <= 50, "name", "name should be less or equal than 50 characters long")

	v.Check(user.Password.Plain != "", "password", "should not be empty")
	v.Check(CheckUserPassword(user.Password.Plain), "password", "password should be at least 6 characters long")

}

func CheckUserPassword(password string) bool {
	return utf8.RuneCountInString(password) >= 6
}
