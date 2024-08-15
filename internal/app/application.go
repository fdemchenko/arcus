package app

import "github.com/fdemchenko/arcus/internal/models"

type Application struct {
	userService UserService
}

type UserService interface {
	Register(models.User) (int, error)
}

func New(userService UserService) *Application {
	return &Application{
		userService: userService,
	}
}
