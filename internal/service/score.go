package service

import (
	"context"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/OkanUysal/world-cup-predict-be/internal/repository"
	"github.com/google/uuid"
)

type ScoreService struct {
	scores *repository.ScoreRepository
	users  *repository.UserRepository
}

func NewScoreService(scores *repository.ScoreRepository, users *repository.UserRepository) *ScoreService {
	return &ScoreService{scores: scores, users: users}
}

func (s *ScoreService) Leaderboard(ctx context.Context, channelID uuid.UUID) ([]domain.UserScore, error) {
	return s.scores.Leaderboard(ctx, channelID)
}

func (s *ScoreService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var totalPoints *int
	if user.ChannelID != uuid.Nil {
		score, err := s.scores.GetByUserID(ctx, userID)
		if err == nil {
			totalPoints = &score.TotalPoints
		} else {
			zero := 0
			totalPoints = &zero
		}
	}

	profile := toUserProfile(user, totalPoints)
	return &profile, nil
}
