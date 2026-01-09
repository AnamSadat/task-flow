package handler

import (
	"net/http"

	"task-flow/internal/httpx"
	authservice "task-flow/internal/service/auth"
)

type AuthHandler struct {
	service *authservice.Service
}

func NewAuthHandler(service *authservice.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	access, refresh, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	httpx.JSON(w, http.StatusOK, tokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	access, refresh, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	httpx.JSON(w, http.StatusOK, tokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := h.service.Logout(r.Context(), req.RefreshToken); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "logout failed")
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := h.service.Register(r.Context(), req.Email, req.Password); err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	httpx.JSON(w, http.StatusCreated, map[string]string{"message": "registered"})
}
