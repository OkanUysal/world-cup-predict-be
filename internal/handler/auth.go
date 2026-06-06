package handler

import (
	"errors"
	"net/http"

	"github.com/OkanUysal/world-cup-predict-be/internal/httputil"
	"github.com/OkanUysal/world-cup-predict-be/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type registerRequest struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	ChannelCode string `json:"channel_code"`
}

type loginRequest struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	ChannelCode string `json:"channel_code"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.auth.Register(r.Context(), req.Name, req.Password, req.ChannelCode)
	if errors.Is(err, service.ErrChannelNotFound) {
		httputil.WriteError(w, http.StatusNotFound, "channel not found")
		return
	}
	if errors.Is(err, service.ErrUserExists) {
		httputil.WriteError(w, http.StatusConflict, "user already exists in this channel")
		return
	}
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.auth.Login(r.Context(), req.Name, req.Password, req.ChannelCode)
	if errors.Is(err, service.ErrChannelNotFound) {
		httputil.WriteError(w, http.StatusNotFound, "channel not found")
		return
	}
	if errors.Is(err, service.ErrInvalidCredentials) {
		httputil.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, resp)
}
