package postgres

import (
	"database/sql"

	"github.com/fdemchenko/arcus/internal/models"
)

type TokensRepository struct {
	DB *sql.DB
}

func (tr *TokensRepository) Insert(token models.Token) error {
	query := `INSERT INTO tokens (scope, expires, token_hash, user_id)
					VALUES ($1, $2, $3, $4)`
	_, err := tr.DB.Exec(query, token.Scope, token.ExpiresAt, token.Hash, token.UserID)
	return err
}
