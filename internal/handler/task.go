package handler

import (
	"net/http"

	"task-flow/internal/httpx"
	taskservice "task-flow/internal/service"
)

type TaskHandler struct {
	Service *taskservice.Service
}

func NewTaskHandler(task *taskservice.Service) *TaskHandler {
	return &TaskHandler{
		Service: task,
	}
}

type taskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Service.GetTasks(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.JSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	var req taskRequest

	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	err := h.Service.AddTask(r.Context(), req.Title, req.Description)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	httpx.JSON(w, http.StatusCreated, map[string]string{"message": " Add new task successfully"})
}
