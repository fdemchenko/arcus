package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func (app *Application) RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error("panic recovered", slog.String("error", fmt.Sprintf("%v", err)))
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Connection", "close")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *Application) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			timeElapsed := time.Since(start)
			app.logger.Info("request processed",
				slog.String("method", r.Method),
				slog.String("url", r.URL.Path),
				slog.Duration("time_elapsed", timeElapsed),
				slog.String("remote_addr", r.RemoteAddr),
			)
		}()

		next.ServeHTTP(w, r)
	})
}
