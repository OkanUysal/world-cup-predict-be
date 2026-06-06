package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type EventType string

const (
	EventTypeMatchScore EventType = "match_score"
	EventTypeChampion   EventType = "champion"
	EventTypeRunnerUp   EventType = "runner_up"
	EventTypeThirdPlace EventType = "third_place"
)

type EventStatus string

const (
	EventStatusOpen      EventStatus = "open"
	EventStatusLocked    EventStatus = "locked"
	EventStatusCompleted EventStatus = "completed"
)

type Channel struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	ChannelID    uuid.UUID `json:"channel_id"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Event struct {
	ID        uuid.UUID       `json:"id"`
	Type      EventType       `json:"type"`
	Title     string          `json:"title"`
	Metadata  json.RawMessage `json:"metadata"`
	Deadline  time.Time       `json:"deadline"`
	Status    EventStatus     `json:"status"`
	Result    json.RawMessage `json:"result,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

type Prediction struct {
	ID             uuid.UUID       `json:"id"`
	EventID        uuid.UUID       `json:"event_id"`
	UserID         uuid.UUID       `json:"user_id"`
	UserName       string          `json:"user_name,omitempty"`
	Choice         json.RawMessage `json:"choice"`
	PointsAwarded  int             `json:"points_awarded"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type UserScore struct {
	UserID      uuid.UUID `json:"user_id"`
	UserName    string    `json:"user_name"`
	ChannelID   uuid.UUID `json:"channel_id"`
	TotalPoints int       `json:"total_points"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EventListFilter string

const (
	EventFilterOpen      EventListFilter = "open"
	EventFilterLocked    EventListFilter = "locked"
	EventFilterPending   EventListFilter = "pending"
	EventFilterCompleted EventListFilter = "completed"
)
