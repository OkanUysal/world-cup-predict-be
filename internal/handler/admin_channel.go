package handler

import (
	"errors"
	"net/http"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/OkanUysal/world-cup-predict-be/internal/httputil"
	"github.com/OkanUysal/world-cup-predict-be/internal/repository"
	"github.com/OkanUysal/world-cup-predict-be/internal/service"
)

type AdminChannelHandler struct {
	channels *service.ChannelService
}

func NewAdminChannelHandler(channels *service.ChannelService) *AdminChannelHandler {
	return &AdminChannelHandler{channels: channels}
}

type createChannelRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (h *AdminChannelHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createChannelRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	channel, err := h.channels.Create(r.Context(), req.Code, req.Name)
	if errors.Is(err, repository.ErrConflict) {
		httputil.WriteError(w, http.StatusConflict, "channel code already exists")
		return
	}
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, channel)
}

func (h *AdminChannelHandler) List(w http.ResponseWriter, r *http.Request) {
	channels, err := h.channels.List(r.Context())
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to list channels")
		return
	}
	if channels == nil {
		channels = []domain.Channel{}
	}
	httputil.WriteJSON(w, http.StatusOK, channels)
}
