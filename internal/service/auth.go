package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/OkanUysal/world-cup-predict-be/internal/auth"
	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/OkanUysal/world-cup-predict-be/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrChannelNotFound    = errors.New("channel not found")
	ErrUserExists         = errors.New("user already exists")
)

type AuthService struct {
	users    *repository.UserRepository
	channels *repository.ChannelRepository
	tokens   *auth.TokenService
}

func NewAuthService(users *repository.UserRepository, channels *repository.ChannelRepository, tokens *auth.TokenService) *AuthService {
	return &AuthService{users: users, channels: channels, tokens: tokens}
}

type AuthResponse struct {
	AccessToken string      `json:"access_token"`
	ExpiresAt   string      `json:"expires_at"`
	User        UserProfile `json:"user"`
}

type UserProfile struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Nickname    string          `json:"nickname,omitempty"`
	DisplayName string          `json:"display_name"`
	Role        domain.Role     `json:"role"`
	ChannelID   *string         `json:"channel_id,omitempty"`
	Channel     *ChannelSummary `json:"channel,omitempty"`
	TotalPoints *int            `json:"total_points,omitempty"`
}

type ChannelSummary struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (s *AuthService) Register(ctx context.Context, name, password, channelCode string) (*AuthResponse, error) {
	name = strings.TrimSpace(name)
	channelCode = strings.TrimSpace(channelCode)
	if name == "" || password == "" || channelCode == "" {
		return nil, fmt.Errorf("name, password and channel_code are required")
	}
	if len(password) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters")
	}

	channel, err := s.channels.GetByCode(ctx, channelCode)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrChannelNotFound
	}
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.users.Create(ctx, channel.ID, name, string(hash), domain.RoleUser)
	if errors.Is(err, repository.ErrConflict) {
		return nil, ErrUserExists
	}
	if err != nil {
		return nil, err
	}

	if err := s.users.EnsureScoreRow(ctx, user.ID, channel.ID); err != nil {
		return nil, err
	}

	return s.buildAuthResponse(user)
}

func (s *AuthService) Login(ctx context.Context, name, password, channelCode string) (*AuthResponse, error) {
	name = strings.TrimSpace(name)
	channelCode = strings.TrimSpace(channelCode)
	if name == "" || password == "" {
		return nil, fmt.Errorf("name and password are required")
	}

	var user *domain.User
	var err error

	if channelCode == "" {
		user, err = s.users.GetAdminByName(ctx, name)
	} else {
		channel, chErr := s.channels.GetByCode(ctx, channelCode)
		if errors.Is(chErr, repository.ErrNotFound) {
			return nil, ErrChannelNotFound
		}
		if chErr != nil {
			return nil, chErr
		}
		user, err = s.users.GetByChannelAndName(ctx, channel.ID, name)
	}

	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.buildAuthResponse(user)
}

func (s *AuthService) BootstrapAdmin(ctx context.Context, name, password string) error {
	if name == "" || password == "" {
		return nil
	}

	count, err := s.users.CountAdmins(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	_, err = s.users.CreateAdminBootstrap(ctx, name, string(hash))
	return err
}

func (s *AuthService) GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.users.GetByID(ctx, userID)
}

func (s *AuthService) buildAuthResponse(user *domain.User) (*AuthResponse, error) {
	token, expiresAt, err := s.tokens.Generate(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken: token,
		ExpiresAt:   expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		User:        toUserProfile(user, nil),
	}, nil
}

func toUserProfile(user *domain.User, totalPoints *int) UserProfile {
	profile := UserProfile{
		ID:          user.ID.String(),
		Name:        user.Name,
		Nickname:    user.Nickname,
		DisplayName: domain.DisplayName(user.Name, user.Nickname),
		Role:        user.Role,
		TotalPoints: totalPoints,
	}
	if user.ChannelID != uuid.Nil {
		chID := user.ChannelID.String()
		profile.ChannelID = &chID
	}
	return profile
}
