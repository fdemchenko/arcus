package app

import (
	"log/slog"

	"github.com/fdemchenko/arcus/internal/models"
)

type Application struct {
	userService UserService
	logger      *slog.Logger
}

type UserService interface {
	Register(models.User) (int, error)
	Activate(token string) error
}

func New(userService UserService, logger *slog.Logger) *Application {
	return &Application{
		userService: userService,
		logger:      logger,
	}
}
