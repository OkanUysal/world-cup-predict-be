package handler

import (
	"errors"
	"net/http"

	"github.com/OkanUysal/world-cup-predict-be/internal/httputil"
	"github.com/OkanUysal/world-cup-predict-be/internal/middleware"
	"github.com/OkanUysal/world-cup-predict-be/internal/repository"
	"github.com/OkanUysal/world-cup-predict-be/internal/service"
)

type MeHandler struct {
	scores *service.ScoreService
}

func NewMeHandler(scores *service.ScoreService) *MeHandler {
	return &MeHandler{scores: scores}
}

type updateNicknameRequest struct {
	Nickname string `json:"nickname"`
}

func (h *MeHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.scores.GetUserProfile(r.Context(), claims.UserID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to load profile")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, profile)
}

func (h *MeHandler) UpdateNickname(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req updateNicknameRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	profile, err := h.scores.UpdateNickname(r.Context(), claims.UserID, req.Nickname)
	if errors.Is(err, repository.ErrNotFound) {
		httputil.WriteError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, profile)
}
