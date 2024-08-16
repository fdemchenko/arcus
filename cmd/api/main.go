package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/fdemchenko/arcus/internal/app"
	"github.com/fdemchenko/arcus/internal/config"
	"github.com/fdemchenko/arcus/internal/repositories/postgres"
	"github.com/fdemchenko/arcus/internal/services"
	"github.com/fdemchenko/arcus/internal/services/mail"
	"github.com/fdemchenko/arcus/templates"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

func main() {
	cfg := config.MustLoad()
	logger := initLogger(cfg.Env)

	db, err := openDB(cfg.Storage)
	handleCriticalError(err, "failed to create DB connections pool", logger)
	logger.Info("connected to db successfully")

	conn, err := amqp.Dial(cfg.RabbitMQConnString)
	handleCriticalError(err, "failed to connect to RabbitMQ", logger)

	channel, err := conn.Channel()
	handleCriticalError(err, "failed to create RabbitMQ channel", logger)

	usersRepo := &postgres.UsersRepository{DB: db}
	tokensRepo := &postgres.TokensRepository{DB: db}

	mailerService := mail.NewMailSender(cfg.SMTPMailer, templates.TemplatesFS)
	consumer, err := mail.NewMailerConsumer(mailerService, channel)
	handleCriticalError(err, "failer to create mailer consumer", logger)

	err = consumer.StartConsuming()
	handleCriticalError(err, "failer to start mailer consumer", logger)

	producer, err := mail.NewMailerProducer(channel)
	handleCriticalError(err, "failer to create mailer producer", logger)

	userService := services.NewUserService(usersRepo, logger, tokensRepo, producer)
	application := app.New(userService, logger)

	// http server start
	address := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	server := http.Server{
		Addr:    address,
		Handler: application.Routes(),
	}

	logger.Info("starting application web server", slog.String("address", address))
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("http server error", slog.String("error", err.Error()))
	}
}

func openDB(dbConfig config.StorageConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConfig.DSN)
	if err != nil {
		return nil, err
	}

	conn, err := db.Conn(context.Background())
	if err != nil {
		return nil, err
	}

	err = conn.PingContext(context.Background())
	if err != nil {
		return nil, err
	}
	conn.Close()

	return db, nil
}

func initLogger(env string) *slog.Logger {
	switch env {
	case EnvDevelopment:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvProduction:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic("unknown environment: " + env)
	}
}

func handleCriticalError(err error, message string, logger *slog.Logger) {
	if err != nil {
		logger.Error(message, slog.String("error", err.Error()))
		os.Exit(1)
	}
}
