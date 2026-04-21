package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/google/uuid"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(event *domain.Event) error {
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	query := `INSERT INTO events (id, user_id, title, description, start_time, end_time, event_type, created_at, updated_at) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query, event.ID, event.UserID, event.Title, event.Description, event.StartTime, event.EndTime, event.EventType, event.CreatedAt, event.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute create event in database: %w", err)
	}
	return nil
}

func (r *EventRepository) FindByID(id uuid.UUID) (*domain.Event, error) {
	event := &domain.Event{}

	query := `SELECT id, user_id, title, description, start_time, end_time, event_type, created_at, updated_at
			  FROM events
			  WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&event.ID, &event.UserID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.EventType, &event.CreatedAt, &event.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find event by id: %w", err)
	}

	return event, nil
}

func (r *EventRepository) FindByTitle(title string) ([]*domain.Event, error) {
	events := []*domain.Event{}

	query := `SELECT id, user_id, title, description, start_time, end_time, event_type, created_at, updated_at
			  FROM events
			  WHERE title = $1`

	rows, err := r.db.Query(query, title)
	if err != nil {
		return nil, fmt.Errorf("failed to send query to database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		event := &domain.Event{}

		err = rows.Scan(
			&event.ID, &event.UserID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.EventType, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning rows: %w", err)
		}

		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return events, nil
}

func (r *EventRepository) FindByDate(start, end time.Time) ([]*domain.Event, error) {
	events := []*domain.Event{}

	startOfRange := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endOfRange := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location()).Add(24 * time.Hour)

	query := `SELECT id, user_id, title, description, start_time, end_time, event_type, created_at, updated_at
			  FROM events
			  WHERE starttime < $2 AND endtime > $1`

	rows, err := r.db.Query(query, startOfRange, endOfRange)
	if err != nil {
		return nil, fmt.Errorf("failed to send query to database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		event := &domain.Event{}

		err = rows.Scan(
			&event.ID, &event.UserID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.EventType, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning rows: %w", err)
		}

		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return events, nil
}

func (r *EventRepository) FindByUserID(id uuid.UUID) ([]*domain.Event, error) {
	events := []*domain.Event{}

	query := `SELECT id, user_id, title, description, start_time, end_time, event_type, created_at, updated_at
			  FROM events
			  WHERE user_id = $1
			  ORDER BY 
              	CASE 
                	WHEN start_time IS NULL THEN 1 
                	ELSE 0
            	END,
            	start_time ASC
			  `

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to send query to database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		event := &domain.Event{}

		err = rows.Scan(
			&event.ID, &event.UserID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.EventType, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning rows: %w", err)
		}

		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return events, nil
}

func (r *EventRepository) Update(event *domain.Event) error {
	event.UpdatedAt = time.Now()

	query := `UPDATE events
			  SET title = $2, description = $3, start_time = $4, end_time = $5, event_type = $6, updated_at = $7
			  WHERE id = $1`

	result, err := r.db.Exec(query, event.ID, event.Title, event.Description, event.StartTime, event.EndTime, event.EventType, event.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute update event in database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("event with id %s not found", event.ID)
	}
	return nil
}

func (r *EventRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM events 
			  WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete event from database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("event with id %s not found", id)
	}
	return nil
}

func (r *EventRepository) FindAll() ([]*domain.Event, error) {
	events := []*domain.Event{}

	query := `SELECT id, user_id, title, description, start_time, end_time, event_type, created_at, updated_at
			  FROM events
			  ORDER BY start_time DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to send query to database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		event := &domain.Event{}

		err = rows.Scan(
			&event.ID, &event.UserID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.EventType, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning rows: %w", err)
		}

		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return events, nil
}
