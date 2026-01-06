package handler

import (
	"net/http"

	"task-flow/internal/httpx"
)

type TaskHandler struct {
	// Add task service here when ready
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	httpx.JSON(w, http.StatusOK, map[string]string{"message": "list tasks"})
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	httpx.JSON(w, http.StatusCreated, map[string]string{"message": "task created"})
}
