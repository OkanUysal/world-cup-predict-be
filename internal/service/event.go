package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/OkanUysal/world-cup-predict-be/internal/repository"
	"github.com/OkanUysal/world-cup-predict-be/internal/scoring"
	"github.com/google/uuid"
)

var (
	ErrEventClosed       = errors.New("event is closed for predictions")
	ErrResultRequired    = errors.New("event result is required before scoring")
	ErrAlreadyScored     = errors.New("event scores already calculated")
	ErrPredictionsHidden = errors.New("predictions are not visible yet")
)

type EventService struct {
	events      *repository.EventRepository
	predictions *repository.PredictionRepository
	scores      *repository.ScoreRepository
}

func NewEventService(events *repository.EventRepository, predictions *repository.PredictionRepository, scores *repository.ScoreRepository) *EventService {
	return &EventService{events: events, predictions: predictions, scores: scores}
}

func (s *EventService) SyncStatuses(ctx context.Context) error {
	return s.events.LockExpiredOpenEvents(ctx)
}

func (s *EventService) Create(ctx context.Context, eventType domain.EventType, title string, metadata json.RawMessage, deadline time.Time) (*domain.Event, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if !deadline.After(time.Now()) {
		return nil, fmt.Errorf("deadline must be in the future")
	}
	if eventType != domain.EventTypeMatchScore && eventType != domain.EventTypeChampion &&
		eventType != domain.EventTypeRunnerUp && eventType != domain.EventTypeThirdPlace {
		return nil, fmt.Errorf("invalid event type")
	}
	return s.events.Create(ctx, eventType, title, metadata, deadline)
}

func (s *EventService) Update(ctx context.Context, id uuid.UUID, title *string, metadata json.RawMessage, deadline *time.Time) (*domain.Event, error) {
	if deadline != nil && !deadline.After(time.Now()) {
		return nil, fmt.Errorf("deadline must be in the future")
	}
	return s.events.UpdateOpen(ctx, id, title, metadata, deadline)
}

func (s *EventService) SetResult(ctx context.Context, id uuid.UUID, result json.RawMessage) (*domain.Event, error) {
	if len(result) == 0 {
		return nil, fmt.Errorf("result is required")
	}

	event, err := s.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.events.RefreshStatus(ctx, event)

	if event.Status == domain.EventStatusCompleted {
		return nil, fmt.Errorf("cannot change result on completed event")
	}

	if err := validateResult(event.Type, result); err != nil {
		return nil, err
	}

	return s.events.SetResult(ctx, id, result)
}

func (s *EventService) CalculateScores(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	event, err := s.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.scoreEvent(ctx, event); err != nil {
		return nil, err
	}
	if err := s.scores.RecalculateAll(ctx); err != nil {
		return nil, err
	}
	return s.events.GetByID(ctx, id)
}

type CalculateAllScoresResult struct {
	ProcessedCount int            `json:"processed_count"`
	Events         []domain.Event `json:"events"`
}

func (s *EventService) CalculateAllScores(ctx context.Context) (*CalculateAllScoresResult, error) {
	if err := s.SyncStatuses(ctx); err != nil {
		return nil, err
	}

	events, err := s.events.ListReadyForScoring(ctx)
	if err != nil {
		return nil, err
	}

	result := &CalculateAllScoresResult{
		Events: make([]domain.Event, 0, len(events)),
	}

	for i := range events {
		e := &events[i]
		s.events.RefreshStatus(ctx, e)
		if err := s.scoreEvent(ctx, e); err != nil {
			return nil, fmt.Errorf("score event %s (%s): %w", e.ID, e.Title, err)
		}
		completed, err := s.events.GetByID(ctx, e.ID)
		if err != nil {
			return nil, err
		}
		result.Events = append(result.Events, *completed)
		result.ProcessedCount++
	}

	if result.ProcessedCount > 0 {
		if err := s.scores.RecalculateAll(ctx); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *EventService) scoreEvent(ctx context.Context, event *domain.Event) error {
	s.events.RefreshStatus(ctx, event)

	if event.Status == domain.EventStatusCompleted {
		return ErrAlreadyScored
	}
	if len(event.Result) == 0 {
		return ErrResultRequired
	}

	predictions, err := s.predictions.ListByEventID(ctx, event.ID)
	if err != nil {
		return err
	}

	for _, p := range predictions {
		points, err := scoring.CalculatePoints(event.Type, p.Choice, event.Result)
		if err != nil {
			return fmt.Errorf("score prediction %s: %w", p.ID, err)
		}
		if err := s.predictions.UpdatePoints(ctx, p.ID, points); err != nil {
			return err
		}
	}

	return s.events.MarkCompleted(ctx, event.ID)
}

func (s *EventService) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Event, *domain.Prediction, error) {
	if err := s.SyncStatuses(ctx); err != nil {
		return nil, nil, err
	}

	event, err := s.events.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	s.events.RefreshStatus(ctx, event)

	var prediction *domain.Prediction
	p, err := s.predictions.GetByEventAndUser(ctx, id, userID)
	if err == nil {
		prediction = p
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, nil, err
	}

	return event, prediction, nil
}

func (s *EventService) List(ctx context.Context, filter domain.EventListFilter, userID uuid.UUID) ([]EventWithPrediction, error) {
	if err := s.SyncStatuses(ctx); err != nil {
		return nil, err
	}

	events, err := s.events.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make([]EventWithPrediction, 0, len(events))
	for _, e := range rangeEvents(events) {
		s.events.RefreshStatus(ctx, &e)
		item := EventWithPrediction{Event: e}
		p, err := s.predictions.GetByEventAndUser(ctx, e.ID, userID)
		if err == nil {
			item.MyPrediction = p
		}
		result = append(result, item)
	}
	return result, nil
}

type EventWithPrediction struct {
	Event        domain.Event        `json:"event"`
	MyPrediction *domain.Prediction  `json:"my_prediction,omitempty"`
}

func rangeEvents(events []domain.Event) []domain.Event {
	return events
}

func (s *EventService) UpsertPrediction(ctx context.Context, eventID, userID uuid.UUID, choice json.RawMessage) (*domain.Prediction, error) {
	if err := s.SyncStatuses(ctx); err != nil {
		return nil, err
	}

	event, err := s.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	s.events.RefreshStatus(ctx, event)

	if event.Status != domain.EventStatusOpen || !time.Now().Before(event.Deadline) {
		return nil, ErrEventClosed
	}

	if err := validateChoice(event.Type, choice); err != nil {
		return nil, err
	}

	return s.predictions.Upsert(ctx, eventID, userID, choice)
}

func (s *EventService) ListPredictions(ctx context.Context, eventID, channelID uuid.UUID) ([]domain.Prediction, error) {
	if err := s.SyncStatuses(ctx); err != nil {
		return nil, err
	}

	event, err := s.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	s.events.RefreshStatus(ctx, event)

	if event.Status == domain.EventStatusOpen && time.Now().Before(event.Deadline) {
		return nil, ErrPredictionsHidden
	}

	return s.predictions.ListByEventAndChannel(ctx, eventID, channelID)
}

func (s *EventService) ListUserPredictions(ctx context.Context, targetUserID, requestingUserID uuid.UUID) ([]EventWithPrediction, error) {
	if err := s.SyncStatuses(ctx); err != nil {
		return nil, err
	}

	events, err := s.events.List(ctx, "") // Get all events
	if err != nil {
		return nil, err
	}

	result := make([]EventWithPrediction, 0, len(events))
	for _, e := range rangeEvents(events) {
		s.events.RefreshStatus(ctx, &e)
		item := EventWithPrediction{Event: e}
		p, err := s.predictions.GetByEventAndUser(ctx, e.ID, targetUserID)
		if err == nil {
			// Redact prediction details if deadline hasn't passed and the requester is another user
			if requestingUserID != targetUserID && e.Status == domain.EventStatusOpen && time.Now().Before(e.Deadline) {
				p.Choice = json.RawMessage("null")
			}
			item.MyPrediction = p
		}
		result = append(result, item)
	}
	return result, nil
}

func validateChoice(eventType domain.EventType, choice json.RawMessage) error {
	switch eventType {
	case domain.EventTypeMatchScore:
		var c scoring.MatchScoreChoice
		if err := json.Unmarshal(choice, &c); err != nil {
			return fmt.Errorf("invalid match_score choice: expected {\"home_score\": N, \"away_score\": N}")
		}
		if c.HomeScore < 0 || c.AwayScore < 0 {
			return fmt.Errorf("scores must be non-negative")
		}
	case domain.EventTypeChampion, domain.EventTypeRunnerUp, domain.EventTypeThirdPlace:
		var c scoring.TeamChoice
		if err := json.Unmarshal(choice, &c); err != nil {
			return fmt.Errorf("invalid team choice: expected {\"team\": \"...\"}")
		}
		if c.Team == "" {
			return fmt.Errorf("team is required")
		}
	default:
		return fmt.Errorf("unknown event type")
	}
	return nil
}

func validateResult(eventType domain.EventType, result json.RawMessage) error {
	return validateChoice(eventType, result)
}
