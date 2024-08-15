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
	_ "github.com/lib/pq"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

func main() {
	cfg := config.MustLoad()
	logger := initLogger(cfg.Env)

	db, err := openDB(cfg.Storage)
	if err != nil {
		logger.Error("failed to create DB connections pool", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("connected to db successfully")

	usersRepo := &postgres.UsersRepository{DB: db}
	userService := services.NewUserService(usersRepo)
	application := app.New(userService)

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
