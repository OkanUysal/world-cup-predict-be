package main

import (
	"context"
	"flag"
	"log"

	"github.com/OkanUysal/world-cup-predict-be/internal/config"
	"github.com/OkanUysal/world-cup-predict-be/internal/db"
	"github.com/OkanUysal/world-cup-predict-be/internal/repository"
	"github.com/OkanUysal/world-cup-predict-be/internal/seed/worldcup2026"
)

func main() {
	force := flag.Bool("force", false, "re-seed even if WC2026 events already exist")
	flag.Parse()

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

	eventRepo := repository.NewEventRepository(sqlDB)
	seeder := worldcup2026.NewSeeder(eventRepo)

	created, err := seeder.Run(context.Background(), *force)
	if err != nil {
		log.Fatalf("seed: %v", err)
	}

	log.Printf("seeded %d WC2026 events (3 placement + 72 group matches)", created)
	log.Printf("placement deadline (GMT): %s", worldcup2026.PredictionDeadline.Format("2006-01-02T15:04:05Z"))
}
