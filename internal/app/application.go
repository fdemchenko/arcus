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
	SendActivationToken(user models.User) error
	GetByID(userID int) (*models.User, error)
}

func New(userService UserService, logger *slog.Logger) *Application {
	return &Application{
		userService: userService,
		logger:      logger,
	}
}
