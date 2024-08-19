package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"time"

	"github.com/fdemchenko/arcus/internal/models"
	"github.com/fdemchenko/arcus/internal/services/mail"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrActivationTokenExpired = errors.New("activation token has expired")
)

type UsersRepository interface {
	Insert(user models.User) (int, error)
	Activate(userID int) error
}

type TokensRepository interface {
	Insert(models.Token) error
	GetByTokenHash(hash []byte, scope string) (*models.Token, error)
	DeleteAllForUser(userID int, scope string) error
}

type MailerProducer interface {
	Publish(command mail.SendEmailCommand[any]) error
}

type UsersService struct {
	usersRepository    UsersRepository
	tokensRepository   TokensRepository
	mailerProducer     MailerProducer
	logger             *slog.Logger
	activationTokenTTL time.Duration
}

func NewUserService(
	usersRepository UsersRepository,
	logger *slog.Logger,
	tokensRepository TokensRepository,
	mailerProducer MailerProducer,
	activationTokenTTL time.Duration,
) *UsersService {
	return &UsersService{
		usersRepository:    usersRepository,
		logger:             logger,
		tokensRepository:   tokensRepository,
		mailerProducer:     mailerProducer,
		activationTokenTTL: activationTokenTTL,
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

	activationToken, err := models.GenerateToken(models.ScopeActivation, us.activationTokenTTL, newUserID)
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

func (us *UsersService) Activate(tokenPlaintext string) error {
	const op = "services.UserService.Activate"
	logger := us.logger.With(slog.String("op", op))

	// decode to bytes
	tokenBytes := make([]byte, len(tokenPlaintext)/2)
	_, err := hex.Decode(tokenBytes, []byte(tokenPlaintext))
	if err != nil {
		logger.Error("failed to decode hex-encoded token", slog.String("error", err.Error()))
		return err
	}

	// get the token hash
	tokenSha256 := sha256.Sum256(tokenBytes)

	token, err := us.tokensRepository.GetByTokenHash(tokenSha256[:], models.ScopeActivation)
	if err != nil {
		logger.Error("failed to retrieve token from DB", slog.String("error", err.Error()))
		return err
	}

	defer func() {
		err := us.tokensRepository.DeleteAllForUser(token.UserID, models.ScopeActivation)
		if err != nil {
			logger.Error("failed to delete activation token", slog.String("error", err.Error()))
		}
	}()

	if token.ExpiresAt.Before(time.Now()) {
		logger.Info("activation token has expired")
		return ErrActivationTokenExpired
	}

	err = us.usersRepository.Activate(token.UserID)
	if err != nil {
		logger.Error("failed to update user record", slog.String("error", err.Error()))
		return err
	}
	return nil
}
