package app

import (
	"errors"
	"net/http"
	"strings"

	"github.com/fdemchenko/arcus/internal/api/request"
	"github.com/fdemchenko/arcus/internal/api/response"
	"github.com/fdemchenko/arcus/internal/models"
	"github.com/fdemchenko/arcus/internal/repositories"
	"github.com/fdemchenko/arcus/internal/validator"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", app.registerUser)

	return mux
}

func (app *Application) registerUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := request.ReadJSON(r.Body, &input)
	if err != nil {
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
		response.SendError(w, http.StatusBadRequest, v.Errors)
		return
	}

	id, err := app.userService.Register(newUser)
	if err != nil {
		if errors.Is(err, repositories.ErrEmailAlreadyExists) {
			response.SendError(w, http.StatusConflict, err.Error())
		} else {
			response.SendServerError(w)
		}
		return
	}

	if err := response.WriteJSON(w, http.StatusCreated, response.Envelope{"user_id": id}); err != nil {
		response.SendServerError(w)
	}
}
