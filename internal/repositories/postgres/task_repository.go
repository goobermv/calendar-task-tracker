package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/google/uuid"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *domain.Task) error {
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	query := `INSERT INTO tasks (id, user_id, title, description, status, priority, created_at, updated_at, due_date) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query, task.ID, task.UserID, task.Title, task.Description, task.Status, task.Priority, task.CreatedAt, task.UpdatedAt, task.DueDate)
	if err != nil {
		return fmt.Errorf("failed to execute create task in database: %w", err)
	}
	return nil
}

func (r *TaskRepository) FindByID(id uuid.UUID) (*domain.Task, error) {
	task := &domain.Task{}

	query := `SELECT id, user_id, title, description, status, priority, created_at, updated_at, due_date
			  FROM tasks
			  WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt, &task.DueDate,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find task by id: %w", err)
	}

	return task, nil
}

func (r *TaskRepository) FindByTitle(title string) ([]*domain.Task, error) {
	tasks := []*domain.Task{}

	query := `SELECT id, user_id, title, description, status, priority, created_at, updated_at, due_date
			  FROM tasks
			  WHERE title = $1`

	rows, err := r.db.Query(query, title)
	if err != nil {
		return nil, fmt.Errorf("failed to send query to database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &domain.Task{}

		err = rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt, &task.DueDate,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning rows: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) FindByDate(date time.Time) ([]*domain.Task, error) {
	tasks := []*domain.Task{}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `SELECT id, user_id, title, description, status, priority, created_at, updated_at, due_date
			  FROM tasks
			  WHERE due_date >= $1 AND due_date < $2`

	rows, err := r.db.Query(query, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to send query to database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &domain.Task{}

		err = rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt, &task.DueDate,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning rows: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) Update(task *domain.Task) error {
	task.UpdatedAt = time.Now()

	query := `UPDATE tasks
			  SET title = $2, description = $3, status = $4, priority = $5, updated_at = $6, due_date = $7
			  WHERE id = $1`

	result, err := r.db.Exec(query, task.ID, task.Title, task.Description, task.Status, task.Priority, task.UpdatedAt, task.DueDate)
	if err != nil {
		return fmt.Errorf("failed to execute update task in database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %s not found", task.ID)
	}
	return nil
} // and whether other users have access to updated others' tasks + how to update each field individually

func (r *TaskRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM tasks 
			  WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete task from database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %s not found", id)
	}

	return nil
}
