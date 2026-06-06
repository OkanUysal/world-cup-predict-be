package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/OkanUysal/world-cup-predict-be/internal/repository"
)

type ChannelService struct {
	channels *repository.ChannelRepository
}

func NewChannelService(channels *repository.ChannelRepository) *ChannelService {
	return &ChannelService{channels: channels}
}

func (s *ChannelService) Create(ctx context.Context, code, name string) (*domain.Channel, error) {
	code = strings.TrimSpace(strings.ToUpper(code))
	name = strings.TrimSpace(name)
	if code == "" || name == "" {
		return nil, fmt.Errorf("code and name are required")
	}
	if len(code) < 3 {
		return nil, fmt.Errorf("code must be at least 3 characters")
	}
	return s.channels.Create(ctx, code, name)
}

func (s *ChannelService) List(ctx context.Context) ([]domain.Channel, error) {
	return s.channels.List(ctx)
}
