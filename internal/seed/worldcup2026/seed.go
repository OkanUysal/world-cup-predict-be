package worldcup2026

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/OkanUysal/world-cup-predict-be/internal/repository"
)

type Seeder struct {
	events *repository.EventRepository
}

func NewSeeder(events *repository.EventRepository) *Seeder {
	return &Seeder{events: events}
}

func (s *Seeder) Run(ctx context.Context, force bool) (int, error) {
	count, err := s.events.CountByTitlePrefix(ctx, SeedTitlePrefix)
	if err != nil {
		return 0, err
	}
	if count > 0 && !force {
		return 0, fmt.Errorf("WC2026 events already exist (%d). use -force to re-seed", count)
	}
	if count > 0 && force {
		if err := s.events.DeleteByTitlePrefix(ctx, SeedTitlePrefix); err != nil {
			return 0, fmt.Errorf("delete existing WC2026 events: %w", err)
		}
	}

	teamsMetadata, err := json.Marshal(map[string]any{"teams": Teams})
	if err != nil {
		return 0, err
	}

	placementEvents := []struct {
		eventType domain.EventType
		title     string
	}{
		{domain.EventTypeChampion, SeedTitlePrefix + " Şampiyon"},
		{domain.EventTypeRunnerUp, SeedTitlePrefix + " İkinci"},
		{domain.EventTypeThirdPlace, SeedTitlePrefix + " Üçüncü"},
	}

	created := 0
	for _, pe := range placementEvents {
		if _, err := s.events.Create(ctx, pe.eventType, pe.title, teamsMetadata, PredictionDeadline); err != nil {
			return created, fmt.Errorf("create %s: %w", pe.title, err)
		}
		created++
	}

	for _, match := range GroupMatches {
		metadata, err := json.Marshal(map[string]any{
			"home_team":   match.HomeTeam,
			"away_team":   match.AwayTeam,
			"group":       match.Group,
			"match_no":    match.MatchNo,
			"venue":       match.Venue,
			"kickoff_gmt": match.Kickoff.UTC().Format("2006-01-02T15:04:05Z"),
		})
		if err != nil {
			return created, err
		}

		if _, err := s.events.Create(ctx, domain.EventTypeMatchScore, match.Title(), metadata, match.Deadline()); err != nil {
			return created, fmt.Errorf("create match %d: %w", match.MatchNo, err)
		}
		created++
	}

	return created, nil
}
