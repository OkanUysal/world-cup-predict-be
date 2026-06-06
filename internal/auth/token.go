package auth

import (
	"fmt"
	"time"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID    uuid.UUID   `json:"sub"`
	ChannelID uuid.UUID   `json:"channel_id"`
	Role      domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret     []byte
	expiration time.Duration
}

func NewTokenService(secret string, expiration time.Duration) *TokenService {
	return &TokenService{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

func (s *TokenService) Generate(user *domain.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.expiration)

	claims := Claims{
		UserID:    user.ID,
		ChannelID: user.ChannelID,
		Role:      user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign token: %w", err)
	}

	return signed, expiresAt, nil
}

func (s *TokenService) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
