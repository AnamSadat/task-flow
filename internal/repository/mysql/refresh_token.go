package mysql

import (
	"context"
	"database/sql"

	"task-flow/internal/repository"
)

type refreshTokenRepo struct {
	db *sql.DB
}

func NewRefreshTokenRepo(db *sql.DB) repository.RefreshTokenRepo {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Insert(ctx context.Context, userID string, tokenHash []byte, expUnix int64) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES (?, ?, FROM_UNIXTIME(?))",
		userID, tokenHash, expUnix,
	)
	return err
}

func (r *refreshTokenRepo) FindUserIDByHash(ctx context.Context, tokenHash []byte) (string, bool, error) {
	var userID string
	err := r.db.QueryRowContext(ctx,
		"SELECT user_id FROM refresh_tokens WHERE token_hash = ? AND expires_at > NOW() AND revoked = FALSE",
		tokenHash,
	).Scan(&userID)

	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return userID, true, nil
}

func (r *refreshTokenRepo) RevokeByHash(ctx context.Context, tokenHash []byte) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE refresh_tokens SET revoked = TRUE WHERE token_hash = ?",
		tokenHash,
	)
	return err
}
