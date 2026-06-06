package handler

import (
	"net/http"

	"github.com/OkanUysal/world-cup-predict-be/internal/auth"
	"github.com/OkanUysal/world-cup-predict-be/internal/middleware"
	"github.com/OkanUysal/world-cup-predict-be/internal/swagger"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Handlers struct {
	Auth         *AuthHandler
	Me           *MeHandler
	AdminChannel *AdminChannelHandler
	AdminEvent   *AdminEventHandler
	Event        *EventHandler
}

func NewRouter(h *Handlers, tokenService *auth.TokenService) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
	r.Get("/swagger/", swagger.UI)
	r.Get("/swagger/index.html", swagger.UI)
	r.Get("/swagger/openapi.json", swagger.OpenAPI)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(tokenService))

			r.Get("/me", h.Me.Get)
			r.Get("/leaderboard", h.Event.Leaderboard)

			r.Route("/events", func(r chi.Router) {
				r.Get("/", h.Event.List)
				r.Get("/{id}", h.Event.Get)
				r.Put("/{id}/prediction", h.Event.UpsertPrediction)
				r.Get("/{id}/predictions", h.Event.ListPredictions)
			})

			r.Route("/admin", func(r chi.Router) {
				r.Use(middleware.Admin)

				r.Route("/channels", func(r chi.Router) {
					r.Post("/", h.AdminChannel.Create)
					r.Get("/", h.AdminChannel.List)
				})

				r.Route("/events", func(r chi.Router) {
					r.Post("/", h.AdminEvent.Create)
					r.Post("/calculate-scores", h.AdminEvent.CalculateAllScores)
					r.Patch("/{id}", h.AdminEvent.Update)
					r.Post("/{id}/result", h.AdminEvent.SetResult)
					r.Post("/{id}/calculate-scores", h.AdminEvent.CalculateScores)
				})
			})
		})
	})

	return r
}
