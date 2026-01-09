package repository

import (
	"context"

	"task-flow/internal/model"
)

type UserRepo interface {
	Create(ctx context.Context, user model.User) error
	FindByID(ctx context.Context, id string) (model.User, bool, error)
	FindByEmail(ctx context.Context, email string) (model.User, bool, error)
}
