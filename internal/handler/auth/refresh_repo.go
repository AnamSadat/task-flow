package auth

import "context"

type RefreshRepo interface {
	Insert(ctx context.Context, userID string, tokenHash []byte, expUnix int64) error
	FindUserIDByHash(ctx context.Context, tokenHash []byte) (string, bool, error)
	RevokeByHash(ctx context.Context, tokenHash []byte) error
}