package postgres

import (
	"database/sql"
	"errors"

	"github.com/fdemchenko/arcus/internal/models"
	"github.com/lib/pq"
)

var ErrPostDoesNotExist = errors.New("post does not exist")
var ErrEditConflict = errors.New("entity was edited concurrently")

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

func (pr *PostsRepository) GetAll() ([]models.Post, error) {
	query := `SELECT id, title, content, created_at, updated_at, tags, version FROM posts`

	rows, err := pr.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]models.Post, 0)

	for rows.Next() {
		var post models.Post
		tags := pq.StringArray{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt, &tags, &post.Version)
		if err != nil {
			return nil, err
		}
		post.Tags = tags
		posts = append(posts, post)
	}
	return posts, nil
}

func (pr *PostsRepository) GetByID(id int) (*models.Post, error) {
	query := `SELECT id, title, content, created_at, updated_at, tags, version FROM posts
					WHERE id = $1`

	var post models.Post

	tags := pq.StringArray{}
	err := pr.DB.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt, &tags, &post.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostDoesNotExist
		}
		return nil, err
	}
	post.Tags = tags
	return &post, nil
}

func (pr *PostsRepository) DeleteByID(id int) (int, error) {
	query := `DELETE FROM posts WHERE id = $1 RETURNING id`
	var deletedID int
	err := pr.DB.QueryRow(query, id).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrPostDoesNotExist
		}
		return 0, err
	}

	return deletedID, nil
}

func (pr *PostsRepository) UpdateByID(post models.Post) error {
	query := `UPDATE posts 
					 SET title = $1,
					 content = $2,
					 version = version + 1,
					 tags = $3 WHERE id = $4 AND version = $5`
	result, err := pr.DB.Exec(query, post.Title, post.Content, pq.StringArray(post.Tags), post.ID, post.Version)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrEditConflict
	}
	return nil
}
