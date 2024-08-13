package models

import "time"

type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
