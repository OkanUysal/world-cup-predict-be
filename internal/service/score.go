package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/OkanUysal/world-cup-predict-be/internal/repository"
	"github.com/google/uuid"
)

type ScoreService struct {
	scores   *repository.ScoreRepository
	users    *repository.UserRepository
	channels *repository.ChannelRepository
}

func NewScoreService(scores *repository.ScoreRepository, users *repository.UserRepository, channels *repository.ChannelRepository) *ScoreService {
	return &ScoreService{scores: scores, users: users, channels: channels}
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
	var channelSummary *ChannelSummary

	if user.ChannelID != uuid.Nil {
		score, err := s.scores.GetByUserID(ctx, userID)
		if err == nil {
			totalPoints = &score.TotalPoints
		} else {
			zero := 0
			totalPoints = &zero
		}

		channel, err := s.channels.GetByID(ctx, user.ChannelID)
		if err == nil {
			channelSummary = &ChannelSummary{
				ID:   channel.ID.String(),
				Code: channel.Code,
				Name: channel.Name,
			}
		}
	}

	profile := toUserProfile(user, totalPoints)
	profile.Channel = channelSummary
	return &profile, nil
}

func (s *ScoreService) UpdateNickname(ctx context.Context, userID uuid.UUID, nickname string) (*UserProfile, error) {
	nickname = strings.TrimSpace(nickname)
	if err := validateNickname(nickname); err != nil {
		return nil, err
	}

	user, err := s.users.UpdateNickname(ctx, userID, nickname)
	if err != nil {
		return nil, err
	}

	return s.GetUserProfile(ctx, user.ID)
}

func validateNickname(nickname string) error {
	if nickname == "" {
		return nil
	}
	if len(nickname) < 2 {
		return fmt.Errorf("nickname must be at least 2 characters")
	}
	if len(nickname) > 64 {
		return fmt.Errorf("nickname must be at most 64 characters")
	}
	return nil
}
