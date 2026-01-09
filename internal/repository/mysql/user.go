package mysql

import (
	"context"
	"database/sql"

	"task-flow/internal/model"
	"task-flow/internal/repository"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) repository.UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user model.User) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (id, email, pass_hash) VALUES (?, ?, ?)",
		user.ID, user.Email, user.PassHash,
	)
	return err
}

func (r *userRepo) FindByID(ctx context.Context, id string) (model.User, bool, error) {
	var u model.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, email, pass_hash FROM users WHERE id = ?", id,
	).Scan(&u.ID, &u.Email, &u.PassHash)

	if err == sql.ErrNoRows {
		return model.User{}, false, nil
	}
	if err != nil {
		return model.User{}, false, err
	}
	return u, true, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (model.User, bool, error) {
	var u model.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, email, pass_hash FROM users WHERE email = ?", email,
	).Scan(&u.ID, &u.Email, &u.PassHash)

	if err == sql.ErrNoRows {
		return model.User{}, false, nil
	}
	if err != nil {
		return model.User{}, false, err
	}
	return u, true, nil
}
