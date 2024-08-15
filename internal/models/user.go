package models

import (
	"regexp"
	"time"

	"github.com/fdemchenko/arcus/internal/validator"
)

type User struct {
	ID        int
	Name      string
	Email     string
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
	v.Check(CheckUserPassword(user.Password.Plain), "password", "password should be min 8 characters, at least one letter, one number and one special character")

}

func CheckUserPassword(password string) bool {
	// Minimum eight characters, at least one letter, one number and one special character
	passwordRx := regexp.MustCompile(`^*[A-Za-z].*\\d.*[@$!%*#?&][A-Za-z\\d@$!%*#?&]{8,}$`)
	return passwordRx.MatchString(password)
}
