package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OkanUysal/world-cup-predict-be/internal/auth"
	"github.com/OkanUysal/world-cup-predict-be/internal/config"
	"github.com/OkanUysal/world-cup-predict-be/internal/db"
	"github.com/OkanUysal/world-cup-predict-be/internal/handler"
	"github.com/OkanUysal/world-cup-predict-be/internal/repository"
	"github.com/OkanUysal/world-cup-predict-be/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	sqlDB, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer sqlDB.Close()

	if err := db.RunMigrations(sqlDB); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	channelRepo := repository.NewChannelRepository(sqlDB)
	userRepo := repository.NewUserRepository(sqlDB)
	eventRepo := repository.NewEventRepository(sqlDB)
	predictionRepo := repository.NewPredictionRepository(sqlDB)
	scoreRepo := repository.NewScoreRepository(sqlDB)

	tokenService := auth.NewTokenService(cfg.JWTSecret, cfg.JWTExpiration)

	authService := service.NewAuthService(userRepo, channelRepo, tokenService)
	channelService := service.NewChannelService(channelRepo)
	eventService := service.NewEventService(eventRepo, predictionRepo, scoreRepo)
	scoreService := service.NewScoreService(scoreRepo, userRepo)

	if err := authService.BootstrapAdmin(context.Background(), cfg.AdminBootstrapName, cfg.AdminBootstrapPassword); err != nil {
		log.Fatalf("bootstrap admin: %v", err)
	}

	handlers := &handler.Handlers{
		Auth:         handler.NewAuthHandler(authService),
		Me:           handler.NewMeHandler(scoreService),
		AdminChannel: handler.NewAdminChannelHandler(channelService),
		AdminEvent:   handler.NewAdminEventHandler(eventService),
		Event:        handler.NewEventHandler(eventService, scoreService),
	}

	router := handler.NewRouter(handlers, tokenService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go runStatusSync(ctx, eventService)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func runStatusSync(ctx context.Context, events *service.EventService) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := events.SyncStatuses(ctx); err != nil {
				log.Printf("sync event statuses: %v", err)
			}
		}
	}
}
