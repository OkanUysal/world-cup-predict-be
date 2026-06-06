package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/google/uuid"
)

type PredictionRepository struct {
	db *sql.DB
}

func NewPredictionRepository(db *sql.DB) *PredictionRepository {
	return &PredictionRepository{db: db}
}

func (r *PredictionRepository) Upsert(ctx context.Context, eventID, userID uuid.UUID, choice json.RawMessage) (*domain.Prediction, error) {
	p := &domain.Prediction{}
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO predictions (event_id, user_id, choice)
		VALUES ($1, $2, $3)
		ON CONFLICT (event_id, user_id) DO UPDATE SET
			choice = EXCLUDED.choice,
			updated_at = NOW()
		RETURNING id, event_id, user_id, choice, points_awarded, created_at, updated_at
	`, eventID, userID, choice).Scan(
		&p.ID, &p.EventID, &p.UserID, &p.Choice, &p.PointsAwarded, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert prediction: %w", err)
	}
	return p, nil
}

func (r *PredictionRepository) GetByEventAndUser(ctx context.Context, eventID, userID uuid.UUID) (*domain.Prediction, error) {
	p := &domain.Prediction{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, event_id, user_id, choice, points_awarded, created_at, updated_at
		FROM predictions WHERE event_id = $1 AND user_id = $2
	`, eventID, userID).Scan(
		&p.ID, &p.EventID, &p.UserID, &p.Choice, &p.PointsAwarded, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get prediction: %w", err)
	}
	return p, nil
}

func (r *PredictionRepository) ListByEventAndChannel(ctx context.Context, eventID, channelID uuid.UUID) ([]domain.Prediction, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.id, p.event_id, p.user_id, u.name, p.choice, p.points_awarded, p.created_at, p.updated_at
		FROM predictions p
		JOIN world_cup_users u ON u.id = p.user_id
		WHERE p.event_id = $1 AND u.channel_id = $2
		ORDER BY u.name ASC
	`, eventID, channelID)
	if err != nil {
		return nil, fmt.Errorf("list predictions: %w", err)
	}
	defer rows.Close()

	var predictions []domain.Prediction
	for rows.Next() {
		var p domain.Prediction
		if err := rows.Scan(
			&p.ID, &p.EventID, &p.UserID, &p.UserName, &p.Choice, &p.PointsAwarded, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan prediction: %w", err)
		}
		predictions = append(predictions, p)
	}
	return predictions, rows.Err()
}

func (r *PredictionRepository) ListByEventID(ctx context.Context, eventID uuid.UUID) ([]domain.Prediction, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, event_id, user_id, choice, points_awarded, created_at, updated_at
		FROM predictions WHERE event_id = $1
	`, eventID)
	if err != nil {
		return nil, fmt.Errorf("list predictions by event: %w", err)
	}
	defer rows.Close()

	var predictions []domain.Prediction
	for rows.Next() {
		var p domain.Prediction
		if err := rows.Scan(
			&p.ID, &p.EventID, &p.UserID, &p.Choice, &p.PointsAwarded, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan prediction: %w", err)
		}
		predictions = append(predictions, p)
	}
	return predictions, rows.Err()
}

func (r *PredictionRepository) UpdatePoints(ctx context.Context, predictionID uuid.UUID, points int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE predictions SET points_awarded = $2 WHERE id = $1
	`, predictionID, points)
	if err != nil {
		return fmt.Errorf("update points: %w", err)
	}
	return nil
}
