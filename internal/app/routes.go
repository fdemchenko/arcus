package app

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", app.registerUser)
	mux.HandleFunc("PUT /auth/activate", app.activateUser)
	mux.HandleFunc("POST /auth/resend-activation-token", app.resendActivationToken)

	mux.HandleFunc("POST /posts", app.createPost)
	mux.HandleFunc("GET /posts", app.getPosts)
	mux.HandleFunc("GET /posts/{id}", app.getPostByID)
	mux.HandleFunc("DELETE /posts/{id}", app.deletePostByID)
	middlewares := alice.New(app.RecoveryMiddleware, app.LoggingMiddleware)

	return middlewares.Then(mux)
}
