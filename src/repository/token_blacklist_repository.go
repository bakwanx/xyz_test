package repository

import (
	"database/sql"
	"time"
)

type TokenBlacklistRepository interface {
	Blacklist(jti string, expiresAt time.Time) error
	IsBlacklisted(jti string) (bool, error)
}

type tokenBlacklistRepo struct{ db *sql.DB }

func NewTokenBlacklistRepo(db *sql.DB) TokenBlacklistRepository {
	return &tokenBlacklistRepo{db}
}

func (r *tokenBlacklistRepo) Blacklist(jti string, expiresAt time.Time) error {
	_, err := r.db.Exec(`INSERT INTO token_blacklist (jti,expires_at) VALUES (?,?)`, jti, expiresAt)
	return err
}

func (r *tokenBlacklistRepo) IsBlacklisted(jti string) (bool, error) {
	var cnt int
	err := r.db.QueryRow(`SELECT COUNT(1) FROM token_blacklist WHERE jti=? AND expires_at > NOW()`, jti).Scan(&cnt)
	return cnt > 0, err
}
