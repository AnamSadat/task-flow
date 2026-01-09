package service

import (
	"context"

	"task-flow/internal/model"
	"task-flow/internal/repository"
	"task-flow/internal/utils"
)

type Service struct {
	TaskRepo repository.TaskRepo
}

func NewServiceTask(task repository.TaskRepo) *Service {
	return &Service{
		TaskRepo: task,
	}
}

func (s *Service) AddTask(ctx context.Context, title, description string) error {
	id, err := utils.GenerateID()
	if err != nil {
		return err
	}

	task := model.Task{
		ID:          id,
		Title:       title,
		Description: description,
	}

	return s.TaskRepo.AddTask(ctx, task)
}

func (s *Service) GetTasks(ctx context.Context) ([]model.Task, error) {
	tasks, err := s.TaskRepo.GetTasks(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
