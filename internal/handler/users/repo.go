package users

import "context"

type User struct {
	ID      string
	Email   string
	PassHas []byte
}

type Repo interface {
	FindByID(ctx context.Context, id string) (User, bool, error)
	FindByEmail(ctx context.Context, email string) (User, bool, error)
}