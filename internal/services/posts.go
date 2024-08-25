package services

import (
	"log/slog"

	"github.com/fdemchenko/arcus/internal/models"
)

type PostsRepository interface {
	Insert(models.Post) (int, error)
}

type PostsService struct {
	logger          *slog.Logger
	postsRepository PostsRepository
}

func NewPostsService(logger *slog.Logger, postsRepository PostsRepository) *PostsService {
	return &PostsService{
		logger:          logger,
		postsRepository: postsRepository,
	}
}

func (ps *PostsService) Create(post models.Post) (int, error) {
	return ps.postsRepository.Insert(post)
}
