package repository

import (
	"context"

	"task-flow/internal/model"
)

type TaskRepo interface {
	AddTask(ctx context.Context, task model.Task) error
	GetTasks(ctx context.Context) ([]model.Task, error)
	FindByID(ctx context.Context, id string) (model.Task, error)
	DeleteTask(ctx context.Context, id string) error
}
