package handler

import (
	"net/http"

	"task-flow/internal/httpx"
	"task-flow/internal/middleware"
	"task-flow/internal/repository"
)

type UserHandler struct {
	Repo repository.UserRepo
}

func NewUserHandler(repo repository.UserRepo) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	u, ok, err := h.Repo.FindByID(r.Context(), userID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	if !ok {
		httpx.Error(w, http.StatusNotFound, "user not found")
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]string{
		"id":    u.ID,
		"email": u.Email,
	})
}
