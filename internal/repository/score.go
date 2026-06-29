package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/google/uuid"
)

type ScoreRepository struct {
	db *sql.DB
}

func NewScoreRepository(db *sql.DB) *ScoreRepository {
	return &ScoreRepository{db: db}
}

func (r *ScoreRepository) RecalculateForChannel(ctx context.Context, channelID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE user_scores us SET
			total_points = COALESCE((
				SELECT SUM(p.points_awarded)
				FROM predictions p
				WHERE p.user_id = us.user_id
			), 0),
			updated_at = NOW()
		WHERE us.channel_id = $1
	`, channelID)
	if err != nil {
		return fmt.Errorf("recalculate channel scores: %w", err)
	}
	return nil
}

func (r *ScoreRepository) RecalculateAll(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE user_scores us SET
			total_points = COALESCE((
				SELECT SUM(p.points_awarded)
				FROM predictions p
				WHERE p.user_id = us.user_id
			), 0),
			updated_at = NOW()
	`)
	if err != nil {
		return fmt.Errorf("recalculate all scores: %w", err)
	}
	return nil
}

func (r *ScoreRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserScore, error) {
	s := &domain.UserScore{}
	err := r.db.QueryRowContext(ctx, `
		SELECT 
			us.user_id, 
			COALESCE(NULLIF(TRIM(u.nickname), ''), u.name), 
			us.channel_id, 
			us.total_points,
			COALESCE(COUNT(p.id) FILTER (WHERE e.type = 'match_score' AND p.points_awarded = 3), 0) AS exact_score_count,
			COALESCE(COUNT(p.id) FILTER (WHERE e.type = 'match_score' AND p.points_awarded = 1), 0) AS correct_outcome_count,
			us.updated_at
		FROM user_scores us
		JOIN world_cup_users u ON u.id = us.user_id
		LEFT JOIN predictions p ON p.user_id = us.user_id
		LEFT JOIN events e ON e.id = p.event_id
		WHERE us.user_id = $1
		GROUP BY us.user_id, u.nickname, u.name, us.channel_id, us.total_points, us.updated_at
	`, userID).Scan(&s.UserID, &s.UserName, &s.ChannelID, &s.TotalPoints, &s.ExactScoreCount, &s.CorrectOutcomeCount, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user score: %w", err)
	}
	return s, nil
}

func (r *ScoreRepository) Leaderboard(ctx context.Context, channelID uuid.UUID) ([]domain.UserScore, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT 
			us.user_id, 
			COALESCE(NULLIF(TRIM(u.nickname), ''), u.name) AS user_name, 
			us.channel_id, 
			us.total_points,
			COALESCE(COUNT(p.id) FILTER (WHERE e.type = 'match_score' AND p.points_awarded = 3), 0) AS exact_score_count,
			COALESCE(COUNT(p.id) FILTER (WHERE e.type = 'match_score' AND p.points_awarded = 1), 0) AS correct_outcome_count,
			us.updated_at
		FROM user_scores us
		JOIN world_cup_users u ON u.id = us.user_id
		LEFT JOIN predictions p ON p.user_id = us.user_id
		LEFT JOIN events e ON e.id = p.event_id
		WHERE us.channel_id = $1
		GROUP BY us.user_id, u.nickname, u.name, us.channel_id, us.total_points, us.updated_at
		ORDER BY us.total_points DESC, user_name ASC
	`, channelID)
	if err != nil {
		return nil, fmt.Errorf("leaderboard: %w", err)
	}
	defer rows.Close()

	var scores []domain.UserScore
	for rows.Next() {
		var s domain.UserScore
		if err := rows.Scan(&s.UserID, &s.UserName, &s.ChannelID, &s.TotalPoints, &s.ExactScoreCount, &s.CorrectOutcomeCount, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan score: %w", err)
		}
		scores = append(scores, s)
	}
	return scores, rows.Err()
}
