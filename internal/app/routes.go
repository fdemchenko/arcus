package app

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/fdemchenko/arcus/internal/api/request"
	"github.com/fdemchenko/arcus/internal/api/response"
	"github.com/fdemchenko/arcus/internal/models"
	"github.com/fdemchenko/arcus/internal/repositories"
	"github.com/fdemchenko/arcus/internal/validator"
	"github.com/justinas/alice"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", app.registerUser)
	mux.HandleFunc("PUT /auth/activate", app.activateUser)
	middlewares := alice.New(app.RecoveryMiddleware, app.LoggingMiddleware)

	return middlewares.Then(mux)
}

func (app *Application) registerUser(w http.ResponseWriter, r *http.Request) {
	const op = "app.routes.registerUser"
	logger := app.logger.With("op", op)
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := request.ReadJSON(r.Body, &input)
	if err != nil {
		logger.Error("failed to decode JSON user input", slog.String("error", err.Error()))
		response.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	newUser := models.User{
		Name:  strings.TrimSpace(input.Name),
		Email: strings.TrimSpace(input.Email)}
	newUser.Password.Plain = strings.TrimSpace(input.Password)

	v := validator.New()
	newUser.Validate(v)
	if !v.IsValid() {
		logger.Error("user model input is not valid")
		response.SendError(w, http.StatusBadRequest, v.Errors)
		return
	}

	id, err := app.userService.Register(newUser)
	if err != nil {
		logger.Error("failed to register user", slog.String("error", err.Error()))
		if errors.Is(err, repositories.ErrEmailAlreadyExists) {
			response.SendError(w, http.StatusConflict, err.Error())
		} else {
			response.SendServerError(w)
		}
		return
	}
	logger.Info("new user registered", slog.Int("user_id", id))

	if err := response.WriteJSON(w, http.StatusCreated, response.Envelope{"user_id": id}); err != nil {
		response.SendServerError(w)
	}
}

func (app *Application) activateUser(w http.ResponseWriter, r *http.Request) {
	const op = "app.routes.activateUser"
	logger := app.logger.With(slog.String("op", op))

	input := struct {
		Token string `json:"token"`
	}{}
	err := request.ReadJSON(r.Body, &input)
	if err != nil {
		logger.Error("failed to decode JSON user input", slog.String("error", err.Error()))
		response.SendError(w, http.StatusBadRequest, err.Error())
	}

	v := validator.New()
	v.Check(input.Token != "", "token", "must not be empty")
	if !v.IsValid() {
		response.SendError(w, http.StatusBadRequest, v.Errors)
		return
	}

	if err := app.userService.Activate(input.Token); err != nil {
		response.SendServerError(w)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, response.Envelope{"activated": true}); err != nil {
		response.SendServerError(w)
	}
}
