package handler

import (
	"net/http"

	"github.com/OkanUysal/world-cup-predict-be/internal/httputil"
	"github.com/OkanUysal/world-cup-predict-be/internal/middleware"
	"github.com/OkanUysal/world-cup-predict-be/internal/service"
)

type MeHandler struct {
	scores *service.ScoreService
}

func NewMeHandler(scores *service.ScoreService) *MeHandler {
	return &MeHandler{scores: scores}
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
