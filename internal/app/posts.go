package app

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/fdemchenko/arcus/internal/api/request"
	"github.com/fdemchenko/arcus/internal/api/response"
	"github.com/fdemchenko/arcus/internal/models"
	"github.com/fdemchenko/arcus/internal/repositories/postgres"
	"github.com/fdemchenko/arcus/internal/validator"
)

func (app *Application) createPost(w http.ResponseWriter, r *http.Request) {
	const op = "app.routes.createPost"
	logger := app.logger.With(slog.String("op", op))

	var input struct {
		Title   string   `json:"title"`
		Content *string  `json:"content"`
		Tags    []string `json:"tags"`
	}
	err := request.ReadJSON(r.Body, &input)
	if err != nil {
		logger.Error("failed to decode JSON user input", slog.String("error", err.Error()))
		response.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	post := models.Post{
		Title:   input.Title,
		Content: input.Content,
		Tags:    input.Tags,
	}

	v := validator.New()
	if post.Validate(v); !v.IsValid() {
		logger.Error("failed to validate incoming post", slog.Any("errors", v.Errors))
		response.SendError(w, http.StatusBadRequest, v.Errors)
		return
	}

	postID, err := app.postsService.Create(post)
	if err != nil {
		logger.Error("failed to create post", slog.String("err", err.Error()))
		response.SendServerError(w)
		return
	}

	if err := response.WriteJSON(w, http.StatusCreated, response.Envelope{"post_id": postID}); err != nil {
		response.SendServerError(w)
	}
}

func (app *Application) getPosts(w http.ResponseWriter, r *http.Request) {
	const op = "app.routes.getPosts"
	logger := app.logger.With(slog.String("op", op))

	posts, err := app.postsService.GetAll()
	if err != nil {
		logger.Error("failed to get all posts", slog.String("err", err.Error()))
		response.SendServerError(w)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, response.Envelope{"posts": posts}); err != nil {
		response.SendServerError(w)
	}
}

func (app *Application) getPostByID(w http.ResponseWriter, r *http.Request) {
	const op = "app.routes.getPostByID"
	logger := app.logger.With(slog.String("op", op))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.SendError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	post, err := app.postsService.GetByID(id)
	if err != nil {
		logger.Error("failed to retrieve post from DB", slog.String("err", err.Error()))
		if errors.Is(err, postgres.ErrPostDoesNotExist) {
			response.SendError(w, http.StatusNotFound, err.Error())
			return
		}
		response.SendServerError(w)
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{"post": post})
}

func (app *Application) deletePostByID(w http.ResponseWriter, r *http.Request) {
	const op = "app.routes.deletePostByID"
	logger := app.logger.With(slog.String("op", op))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.SendError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	deletedID, err := app.postsService.DeleteByID(id)
	if err != nil {
		logger.Error("failed to delete post from DB", slog.String("err", err.Error()))
		if errors.Is(err, postgres.ErrPostDoesNotExist) {
			response.SendError(w, http.StatusNotFound, err.Error())
			return
		}
		response.SendServerError(w)
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{"post_id": deletedID})
}

func (app *Application) updatePost(w http.ResponseWriter, r *http.Request) {
	const op = "app.routes.updatePost"
	logger := app.logger.With(slog.String("op", op))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.SendError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	_, err = app.postsService.GetByID(id)
	if err != nil {
		logger.Error("failed to get post from DB", slog.String("err", err.Error()))
		if errors.Is(err, postgres.ErrPostDoesNotExist) {
			response.SendError(w, http.StatusNotFound, err.Error())
			return
		}
		response.SendServerError(w)
		return
	}

	var input struct {
		Title   string   `json:"title"`
		Content *string  `json:"content"`
		Tags    []string `json:"tags"`
	}
	err = request.ReadJSON(r.Body, &input)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	post := models.Post{
		Title:   input.Title,
		Content: input.Content,
		Tags:    input.Tags,
	}
	post.ID = id

	v := validator.New()
	if post.Validate(v); !v.IsValid() {
		response.SendError(w, http.StatusBadRequest, v.Errors)
		return
	}

	err = app.postsService.UpdateByID(post)
	if err != nil {
		logger.Error("failed to update post", slog.String("err", err.Error()))
		response.SendServerError(w)
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{"post_id": id})
}

func (app *Application) partialPostUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "app.routes.updatePostPartial"
	logger := app.logger.With(slog.String("op", op))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.SendError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	post, err := app.postsService.GetByID(id)
	if err != nil {
		logger.Error("failed to get post from DB", slog.String("err", err.Error()))
		if errors.Is(err, postgres.ErrPostDoesNotExist) {
			response.SendError(w, http.StatusNotFound, err.Error())
			return
		}
		response.SendServerError(w)
		return
	}

	input := struct {
		Title   *string  `json:"title"`
		Content *string  `json:"content"`
		Tags    []string `json:"tags"`
	}{}

	err = request.ReadJSON(r.Body, &input)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if input.Title != nil {
		post.Title = *input.Title
	}

	if input.Content != nil {
		post.Content = input.Content
		if *input.Content == "" {
			post.Content = nil
		}
	}

	if input.Tags != nil {
		post.Tags = input.Tags
	}

	v := validator.New()
	if post.Validate(v); !v.IsValid() {
		logger.Error("failed to validate incoming post", slog.Any("err", v.Errors))
		response.SendError(w, http.StatusBadRequest, v.Errors)
		return
	}

	err = app.postsService.UpdateByID(*post)
	if err != nil {
		if errors.Is(err, postgres.ErrEditConflict) {
			response.SendError(w, http.StatusConflict, err.Error())
			return
		}
		response.SendServerError(w)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, response.Envelope{"post": *post}); err != nil {
		response.SendServerError(w)
	}
}
