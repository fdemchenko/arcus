package services

import (
	"log/slog"
	"time"

	"github.com/fdemchenko/arcus/internal/models"
	"github.com/fdemchenko/arcus/internal/services/mail"
	"golang.org/x/crypto/bcrypt"
)

const ActivationTokenTTL = time.Hour * 2

type UsersRepository interface {
	Insert(user models.User) (int, error)
}

type TokensRepository interface {
	Insert(models.Token) error
}

type MailerProducer interface {
	Publish(command mail.SendEmailCommand[any]) error
}

type UsersService struct {
	usersRepository  UsersRepository
	tokensRepository TokensRepository
	mailerProducer   MailerProducer
	logger           *slog.Logger
}

func NewUserService(
	usersRepository UsersRepository,
	logger *slog.Logger,
	tokensRepository TokensRepository,
	mailerProducer MailerProducer,
) *UsersService {
	return &UsersService{
		usersRepository:  usersRepository,
		logger:           logger,
		tokensRepository: tokensRepository,
		mailerProducer:   mailerProducer,
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

	activationToken, err := models.GenerateToken(models.ScopeActivation, ActivationTokenTTL, newUserID)
	if err != nil {
		logger.Error("failed to create activation token", slog.String("error", err.Error()))
		return newUserID, nil
	}

	err = us.tokensRepository.Insert(*activationToken)
	if err != nil {
		logger.Error("failed to insert token to DB", slog.String("error", err.Error()))
		return newUserID, nil
	}

	command := mail.SendEmailCommand[any]{
		To:           user.Email,
		TemplateName: "user_welcome.tmpl",
		TemplateData: mail.UserWelcomeData{Token: activationToken.PlainText, Name: user.Name},
	}
	err = us.mailerProducer.Publish(command)
	if err != nil {
		logger.Error("failed to send email command", slog.String("error", err.Error()))
	}

	return newUserID, nil
}
