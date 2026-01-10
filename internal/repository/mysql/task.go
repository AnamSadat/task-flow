package mysql

import (
	"context"
	"database/sql"

	"task-flow/internal/model"
	"task-flow/internal/repository"
)

type taskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) repository.TaskRepo {
	return &taskRepo{db: db}
}

func (r *taskRepo) AddTask(ctx context.Context, task model.Task) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO tasks (id, title, description) VALUES (?, ?, ?)",
		task.ID, task.Title, task.Description,
	)

	return err
}

func (r *taskRepo) GetTasks(ctx context.Context) ([]model.Task, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Created_At); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *taskRepo) DeleteTask(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM tasks WHERE id = ?", id)

	return err
}

func (r *taskRepo) FindByID(ctx context.Context, id string) (model.Task, error) {
	var t model.Task
	err := r.db.QueryRowContext(ctx,
		"SELECT id, title, description FROM tasks WHERE id = ?",
		id,
	).Scan(&t.ID, &t.Title, &t.Description)

	if err == sql.ErrNoRows {
		return model.Task{}, nil
	}

	if err != nil {
		return model.Task{}, err
	}

	return t, nil
}
