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

const userSelectColumns = `id, channel_id, name, nickname, password_hash, role, created_at`

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

type scannable interface {
	Scan(dest ...any) error
}

func scanUser(row scannable) (*domain.User, error) {
	u := &domain.User{}
	var channelID, nickname sql.NullString
	err := row.Scan(&u.ID, &channelID, &u.Name, &nickname, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	if channelID.Valid {
		u.ChannelID, err = uuid.Parse(channelID.String)
		if err != nil {
			return nil, fmt.Errorf("parse channel id: %w", err)
		}
	}
	if nickname.Valid {
		u.Nickname = nickname.String
	}
	return u, nil
}

func (r *UserRepository) Create(ctx context.Context, channelID uuid.UUID, name, passwordHash string, role domain.Role) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO world_cup_users (channel_id, name, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING `+userSelectColumns+`
	`, channelID, name, passwordHash, role)

	u, err := scanUser(row)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, ErrConflict
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetByChannelAndName(ctx context.Context, channelID uuid.UUID, name string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT `+userSelectColumns+`
		FROM world_cup_users
		WHERE channel_id = $1 AND name = $2
	`, channelID, name)

	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT `+userSelectColumns+`
		FROM world_cup_users WHERE id = $1
	`, id)

	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil
}

func (r *UserRepository) UpdateNickname(ctx context.Context, userID uuid.UUID, nickname string) (*domain.User, error) {
	var nickParam any
	if nickname == "" {
		nickParam = nil
	} else {
		nickParam = nickname
	}

	row := r.db.QueryRowContext(ctx, `
		UPDATE world_cup_users SET nickname = $2
		WHERE id = $1
		RETURNING `+userSelectColumns+`
	`, userID, nickParam)

	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update nickname: %w", err)
	}
	return u, nil
}

func (r *UserRepository) CountAdmins(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM world_cup_users WHERE role = 'admin'
	`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count admins: %w", err)
	}
	return count, nil
}

func (r *UserRepository) CreateAdminBootstrap(ctx context.Context, name, passwordHash string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO world_cup_users (channel_id, name, password_hash, role)
		VALUES (NULL, $1, $2, 'admin')
		RETURNING `+userSelectColumns+`
	`, name, passwordHash)

	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("create admin bootstrap: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetAdminByName(ctx context.Context, name string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT `+userSelectColumns+`
		FROM world_cup_users
		WHERE role = 'admin' AND LOWER(name) = LOWER($1)
	`, name)

	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get admin: %w", err)
	}
	return u, nil
}

func (r *UserRepository) EnsureScoreRow(ctx context.Context, userID, channelID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_scores (user_id, channel_id, total_points)
		VALUES ($1, $2, 0)
		ON CONFLICT (user_id) DO NOTHING
	`, userID, channelID)
	if err != nil {
		return fmt.Errorf("ensure score row: %w", err)
	}
	return nil
}
