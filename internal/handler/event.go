package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/OkanUysal/world-cup-predict-be/internal/httputil"
	"github.com/OkanUysal/world-cup-predict-be/internal/middleware"
	"github.com/OkanUysal/world-cup-predict-be/internal/repository"
	"github.com/OkanUysal/world-cup-predict-be/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AdminEventHandler struct {
	events *service.EventService
}

func NewAdminEventHandler(events *service.EventService) *AdminEventHandler {
	return &AdminEventHandler{events: events}
}

type createEventRequest struct {
	Type     domain.EventType `json:"type"`
	Title    string           `json:"title"`
	Metadata json.RawMessage  `json:"metadata"`
	Deadline time.Time        `json:"deadline"`
}

type updateEventRequest struct {
	Title    *string          `json:"title"`
	Metadata json.RawMessage  `json:"metadata"`
	Deadline *time.Time       `json:"deadline"`
}

type setResultRequest struct {
	Result json.RawMessage `json:"result"`
}

func (h *AdminEventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createEventRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event, err := h.events.Create(r.Context(), req.Type, req.Title, req.Metadata, req.Deadline)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, event)
}

func (h *AdminEventHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid event id")
		return
	}

	var req updateEventRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event, err := h.events.Update(r.Context(), id, req.Title, req.Metadata, req.Deadline)
	if errors.Is(err, repository.ErrNotFound) {
		httputil.WriteError(w, http.StatusNotFound, "event not found")
		return
	}
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, event)
}

func (h *AdminEventHandler) SetResult(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid event id")
		return
	}

	var req setResultRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event, err := h.events.SetResult(r.Context(), id, req.Result)
	if errors.Is(err, repository.ErrNotFound) {
		httputil.WriteError(w, http.StatusNotFound, "event not found")
		return
	}
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, event)
}

func (h *AdminEventHandler) CalculateAllScores(w http.ResponseWriter, r *http.Request) {
	result, err := h.events.CalculateAllScores(r.Context())
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, result)
}

func (h *AdminEventHandler) CalculateScores(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid event id")
		return
	}

	event, err := h.events.CalculateScores(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		httputil.WriteError(w, http.StatusNotFound, "event not found")
		return
	}
	if errors.Is(err, service.ErrResultRequired) {
		httputil.WriteError(w, http.StatusBadRequest, "result is required before scoring")
		return
	}
	if errors.Is(err, service.ErrAlreadyScored) {
		httputil.WriteError(w, http.StatusConflict, "scores already calculated")
		return
	}
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, event)
}

type EventHandler struct {
	events *service.EventService
	scores *service.ScoreService
}

func NewEventHandler(events *service.EventService, scores *service.ScoreService) *EventHandler {
	return &EventHandler{events: events, scores: scores}
}

type predictionRequest struct {
	Choice json.RawMessage `json:"choice"`
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	filter := domain.EventListFilter(r.URL.Query().Get("status"))
	if filter == "" {
		filter = domain.EventFilterOpen
	}

	switch filter {
	case domain.EventFilterOpen, domain.EventFilterLocked, domain.EventFilterPending, domain.EventFilterCompleted:
	default:
		httputil.WriteError(w, http.StatusBadRequest, "invalid status filter")
		return
	}

	events, err := h.events.List(r.Context(), filter, claims.UserID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to list events")
		return
	}
	if events == nil {
		events = []service.EventWithPrediction{}
	}

	httputil.WriteJSON(w, http.StatusOK, events)
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid event id")
		return
	}

	event, prediction, err := h.events.Get(r.Context(), id, claims.UserID)
	if errors.Is(err, repository.ErrNotFound) {
		httputil.WriteError(w, http.StatusNotFound, "event not found")
		return
	}
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to get event")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]any{
		"event":         event,
		"my_prediction": prediction,
	})
}

func (h *EventHandler) UpsertPrediction(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid event id")
		return
	}

	var req predictionRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	prediction, err := h.events.UpsertPrediction(r.Context(), id, claims.UserID, req.Choice)
	if errors.Is(err, repository.ErrNotFound) {
		httputil.WriteError(w, http.StatusNotFound, "event not found")
		return
	}
	if errors.Is(err, service.ErrEventClosed) {
		httputil.WriteError(w, http.StatusForbidden, "event is closed for predictions")
		return
	}
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, prediction)
}

func (h *EventHandler) ListPredictions(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if claims.ChannelID == uuid.Nil {
		httputil.WriteError(w, http.StatusForbidden, "channel membership required")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid event id")
		return
	}

	predictions, err := h.events.ListPredictions(r.Context(), id, claims.ChannelID)
	if errors.Is(err, repository.ErrNotFound) {
		httputil.WriteError(w, http.StatusNotFound, "event not found")
		return
	}
	if errors.Is(err, service.ErrPredictionsHidden) {
		httputil.WriteError(w, http.StatusForbidden, "predictions are not visible until deadline passes")
		return
	}
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to list predictions")
		return
	}
	if predictions == nil {
		predictions = []domain.Prediction{}
	}

	httputil.WriteJSON(w, http.StatusOK, predictions)
}

func (h *EventHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if claims.ChannelID == uuid.Nil {
		httputil.WriteError(w, http.StatusForbidden, "channel membership required")
		return
	}

	scores, err := h.scores.Leaderboard(r.Context(), claims.ChannelID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to load leaderboard")
		return
	}
	if scores == nil {
		scores = []domain.UserScore{}
	}

	httputil.WriteJSON(w, http.StatusOK, scores)
}

func (h *EventHandler) ListUserPredictions(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if claims.ChannelID == uuid.Nil {
		httputil.WriteError(w, http.StatusForbidden, "channel membership required")
		return
	}

	targetUserID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	targetUser, err := h.scores.GetUser(r.Context(), targetUserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			httputil.WriteError(w, http.StatusNotFound, "user not found")
		} else {
			httputil.WriteError(w, http.StatusInternalServerError, "failed to fetch user details")
		}
		return
	}

	if targetUser.ChannelID != claims.ChannelID {
		httputil.WriteError(w, http.StatusForbidden, "user belongs to a different channel")
		return
	}

	predictions, err := h.events.ListUserPredictions(r.Context(), targetUserID, claims.UserID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to load user predictions")
		return
	}

	if predictions == nil {
		predictions = []service.PredictionComparison{}
	}

	httputil.WriteJSON(w, http.StatusOK, predictions)
}
