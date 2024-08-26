package models

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/fdemchenko/arcus/internal/validator"
)

const MaxUserNameLength = 50
const MinPasswordLength = 6

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
	v.Check(utf8.RuneCountInString(user.Name) <= MaxUserNameLength, "name", fmt.Sprintf("name should be less or equal than %d characters long", MaxUserNameLength))

	v.Check(user.Password.Plain != "", "password", "should not be empty")
	v.Check(CheckUserPassword(user.Password.Plain), "password", fmt.Sprintf("password should be at least %d characters long", MinPasswordLength))

}

func CheckUserPassword(password string) bool {
	return utf8.RuneCountInString(password) >= MinPasswordLength
}
