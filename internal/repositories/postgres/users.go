package postgres

import (
	"database/sql"
	"errors"

	"github.com/fdemchenko/arcus/internal/models"
	"github.com/fdemchenko/arcus/internal/repositories"
	"github.com/lib/pq"
)

var ErrUserDoesNotExists = errors.New("user does not exists")

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

func (ur *UsersRepository) GetByID(userID int) (*models.User, error) {
	query := `SELECT id, name, email, password_hash, activated, created_at, updated_at
					 FROM users WHERE id = $1`
	var user models.User
	row := ur.DB.QueryRow(query, userID)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash, &user.Activated, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserDoesNotExists
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UsersRepository) Activate(userID int) error {
	query := `UPDATE users SET activated = TRUE WHERE id = $1`
	_, err := ur.DB.Exec(query, userID)
	return err
}
