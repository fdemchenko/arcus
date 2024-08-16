package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

const (
	ScopeActivation = "activation"
)

const TokenBytesLength = 18

type Token struct {
	ID        int
	UserID    int
	ExpiresAt time.Time
	Hash      []byte
	PlainText string
	Scope     string
}

func GenerateToken(scope string, ttl time.Duration, userID int) (*Token, error) {
	tokenBytes := make([]byte, TokenBytesLength)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return nil, err
	}

	plaintextToken := hex.EncodeToString(tokenBytes)
	tokenHash := sha256.Sum256(tokenBytes)

	token := Token{
		UserID:    userID,
		ExpiresAt: time.Now().UTC().Add(ttl),
		Hash:      tokenHash[:],
		PlainText: plaintextToken,
		Scope:     scope,
	}
	return &token, nil
}
