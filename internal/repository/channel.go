package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")

type ChannelRepository struct {
	db *sql.DB
}

func NewChannelRepository(db *sql.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) Create(ctx context.Context, code, name string) (*domain.Channel, error) {
	ch := &domain.Channel{}
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO channels (code, name)
		VALUES ($1, $2)
		RETURNING id, code, name, created_at
	`, code, name).Scan(&ch.ID, &ch.Code, &ch.Name, &ch.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, ErrConflict
		}
		return nil, fmt.Errorf("create channel: %w", err)
	}
	return ch, nil
}

func (r *ChannelRepository) List(ctx context.Context) ([]domain.Channel, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, code, name, created_at FROM channels ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list channels: %w", err)
	}
	defer rows.Close()

	var channels []domain.Channel
	for rows.Next() {
		var ch domain.Channel
		if err := rows.Scan(&ch.ID, &ch.Code, &ch.Name, &ch.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan channel: %w", err)
		}
		channels = append(channels, ch)
	}
	return channels, rows.Err()
}

func (r *ChannelRepository) GetByCode(ctx context.Context, code string) (*domain.Channel, error) {
	ch := &domain.Channel{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, code, name, created_at FROM channels WHERE code = $1
	`, code).Scan(&ch.ID, &ch.Code, &ch.Name, &ch.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get channel by code: %w", err)
	}
	return ch, nil
}

func (r *ChannelRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Channel, error) {
	ch := &domain.Channel{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, code, name, created_at FROM channels WHERE id = $1
	`, id).Scan(&ch.ID, &ch.Code, &ch.Name, &ch.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get channel by id: %w", err)
	}
	return ch, nil
}
