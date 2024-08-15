package services

import (
	"log/slog"

	"github.com/fdemchenko/arcus/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepository interface {
	Insert(user models.User) (int, error)
}

type UsersService struct {
	usersRepository UsersRepository
	logger          *slog.Logger
}

func NewUserService(usersRepository UsersRepository, logger *slog.Logger) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
		logger:          logger,
	}
}

func (us *UsersService) Register(user models.User) (int, error) {
	const op = "services.UserService.Register"
	logger := us.logger.With("op", op)

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password.Plain), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to hash password with bcrypt", slog.String("error", err.Error()))
		return 0, nil
	}

	user.Password.Hash = hash
	newUserID, err := us.usersRepository.Insert(user)
	if err != nil {
		logger.Error("failed to create new user in DB", slog.String("error", err.Error()))
		return 0, err
	}
	return newUserID, nil
}
