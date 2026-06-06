package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/OkanUysal/world-cup-predict-be/internal/domain"
	"github.com/google/uuid"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func scanEvent(row interface {
	Scan(dest ...any) error
}) (*domain.Event, error) {
	e := &domain.Event{}
	var result sql.NullString
	err := row.Scan(
		&e.ID, &e.Type, &e.Title, &e.Metadata, &e.Deadline, &e.Status, &result, &e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if result.Valid {
		e.Result = json.RawMessage(result.String)
	}
	return e, nil
}

func (r *EventRepository) Create(ctx context.Context, eventType domain.EventType, title string, metadata json.RawMessage, deadline time.Time) (*domain.Event, error) {
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO events (type, title, metadata, deadline, status)
		VALUES ($1, $2, $3, $4, 'open')
		RETURNING id, type, title, metadata, deadline, status, result, created_at
	`, eventType, title, metadata, deadline)

	e, err := scanEvent(row)
	if err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}
	return e, nil
}

func (r *EventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, type, title, metadata, deadline, status, result, created_at
		FROM events WHERE id = $1
	`, id)

	e, err := scanEvent(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}
	return e, nil
}

func (r *EventRepository) UpdateOpen(ctx context.Context, id uuid.UUID, title *string, metadata json.RawMessage, deadline *time.Time) (*domain.Event, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current.Status != domain.EventStatusOpen {
		return nil, fmt.Errorf("event is not open")
	}

	newTitle := current.Title
	if title != nil {
		newTitle = *title
	}
	newMetadata := current.Metadata
	if len(metadata) > 0 {
		newMetadata = metadata
	}
	newDeadline := current.Deadline
	if deadline != nil {
		newDeadline = *deadline
	}

	row := r.db.QueryRowContext(ctx, `
		UPDATE events SET title = $2, metadata = $3, deadline = $4
		WHERE id = $1 AND status = 'open'
		RETURNING id, type, title, metadata, deadline, status, result, created_at
	`, id, newTitle, newMetadata, newDeadline)

	e, err := scanEvent(row)
	if err != nil {
		return nil, fmt.Errorf("update event: %w", err)
	}
	return e, nil
}

func (r *EventRepository) SetResult(ctx context.Context, id uuid.UUID, result json.RawMessage) (*domain.Event, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE events SET result = $2
		WHERE id = $1 AND status IN ('open', 'locked')
		RETURNING id, type, title, metadata, deadline, status, result, created_at
	`, id, result)

	e, err := scanEvent(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("set result: %w", err)
	}
	return e, nil
}

func (r *EventRepository) MarkCompleted(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE events SET status = 'completed' WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("mark completed: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *EventRepository) LockExpiredOpenEvents(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE events SET status = 'locked'
		WHERE status = 'open' AND deadline <= NOW()
	`)
	if err != nil {
		return fmt.Errorf("lock expired events: %w", err)
	}
	return nil
}

func (r *EventRepository) List(ctx context.Context, filter domain.EventListFilter) ([]domain.Event, error) {
	query := `
		SELECT id, type, title, metadata, deadline, status, result, created_at
		FROM events
	`
	var args []any

	switch filter {
	case domain.EventFilterOpen:
		query += ` WHERE status = 'open' AND deadline > NOW()`
	case domain.EventFilterLocked:
		query += ` WHERE status = 'locked'`
	case domain.EventFilterPending:
		query += ` WHERE status = 'locked' AND result IS NULL`
	case domain.EventFilterCompleted:
		query += ` WHERE status = 'completed'`
	}

	query += ` ORDER BY deadline ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		e, err := scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, *e)
	}
	return events, rows.Err()
}

func (r *EventRepository) ListReadyForScoring(ctx context.Context) ([]domain.Event, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, type, title, metadata, deadline, status, result, created_at
		FROM events
		WHERE result IS NOT NULL AND status != 'completed'
		ORDER BY deadline ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list events ready for scoring: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		e, err := scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, *e)
	}
	return events, rows.Err()
}

func (r *EventRepository) CountByTitlePrefix(ctx context.Context, prefix string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM events WHERE title LIKE $1
	`, prefix+"%").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count events by prefix: %w", err)
	}
	return count, nil
}

func (r *EventRepository) DeleteByTitlePrefix(ctx context.Context, prefix string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM events WHERE title LIKE $1
	`, prefix+"%")
	if err != nil {
		return fmt.Errorf("delete events by prefix: %w", err)
	}
	return nil
}

func (r *EventRepository) RefreshStatus(ctx context.Context, e *domain.Event) {
	if e.Status == domain.EventStatusOpen && !time.Now().Before(e.Deadline) {
		e.Status = domain.EventStatusLocked
	}
}
