package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/fdemchenko/arcus/internal/app"
	"github.com/fdemchenko/arcus/internal/config"
	"github.com/fdemchenko/arcus/internal/repositories/postgres"
	"github.com/fdemchenko/arcus/internal/services"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoad()
	db, err := openDB(cfg.Storage)
	if err != nil {
		log.Fatal(err)
	}

	usersRepo := &postgres.UsersRepository{DB: db}
	userService := services.NewUserService(usersRepo)
	application := app.New(userService)

	address := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	server := http.Server{
		Addr:    address,
		Handler: application.Routes(),
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
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
