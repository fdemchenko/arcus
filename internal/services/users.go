package services

import (
	"github.com/fdemchenko/arcus/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepository interface {
	Insert(user models.User) (int, error)
}

type UsersService struct {
	usersRepository UsersRepository
}

func NewUserService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

func (us *UsersService) Register(user models.User) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password.Plain), bcrypt.DefaultCost)
	if err != nil {
		return 0, nil
	}

	user.Password.Hash = hash
	newUserID, err := us.usersRepository.Insert(user)
	if err != nil {
		return 0, err
	}
	return newUserID, nil
}
