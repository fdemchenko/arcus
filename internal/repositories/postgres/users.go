package postgres

import (
	"database/sql"
	"errors"

	"github.com/fdemchenko/arcus/internal/models"
	"github.com/fdemchenko/arcus/internal/repositories"
	"github.com/lib/pq"
)

type UsersRepository struct {
	DB *sql.DB
}

func (ur *UsersRepository) Insert(user models.User) (int, error) {
	query := `INSERT INTO users (name, email, password_hash)
					VALUES ($1, $2, $3) RETURNING id`
	row := ur.DB.QueryRow(query, user.Name, user.Email, user.Password.Hash)

	var id int
	err := row.Scan(&id)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) && pgError.Code == pq.ErrorCode(repositories.UniqueViolationErrorCode) {
			return 0, repositories.ErrEmailAlreadyExists
		}
		return 0, nil
	}

	return id, nil
}
