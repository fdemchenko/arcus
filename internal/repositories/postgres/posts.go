package postgres

import (
	"database/sql"

	"github.com/fdemchenko/arcus/internal/models"
	"github.com/lib/pq"
)

type PostsRepository struct {
	DB *sql.DB
}

func (pr *PostsRepository) Insert(post models.Post) (int, error) {
	query := `INSERT INTO posts (title, content, tags)
					VALUES ($1, $2, $3) RETURNING id`
	var newPostID int
	err := pr.DB.QueryRow(query, post.Title, post.Content, pq.StringArray(post.Tags)).Scan(&newPostID)
	return newPostID, err
}
