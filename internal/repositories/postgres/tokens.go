package postgres

import (
	"database/sql"
	"errors"

	"github.com/fdemchenko/arcus/internal/models"
)

var ErrTokenNotFound = errors.New("token not found")

type TokensRepository struct {
	DB *sql.DB
}

func (tr *TokensRepository) Insert(token models.Token) error {
	query := `INSERT INTO tokens (scope, expires, token_hash, user_id)
					VALUES ($1, $2, $3, $4)`
	_, err := tr.DB.Exec(query, token.Scope, token.ExpiresAt, token.Hash, token.UserID)
	return err
}

func (tr *TokensRepository) GetByTokenHash(hash []byte, scope string) (*models.Token, error) {
	query := `SELECT id, user_id, expires, token_hash, scope FROM tokens
					WHERE scope = $1 AND token_hash = $2`

	var token models.Token
	err := tr.DB.QueryRow(query, scope, hash).
		Scan(
			&token.ID,
			&token.UserID,
			&token.ExpiresAt,
			&token.Hash,
			&token.Scope,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}
	return &token, nil
}

func (tr *TokensRepository) DeleteAllForUser(userID int, scope string) error {
	query := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`
	_, err := tr.DB.Exec(query, userID, scope)
	return err
}
