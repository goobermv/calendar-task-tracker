package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `INSERT INTO users (id, email, username, password, usertype, created_at, updated_at) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query, user.ID, user.Email, user.Username, user.PasswordHash, user.UserType, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute create user in database: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	user := &domain.User{}

	query := `SELECT id, email, username, password, usertype, created_at, updated_at
			  FROM users
			  WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.UserType, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return user, nil
}

func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	user := &domain.User{}

	query := `SELECT id, email, username, password, usertype, created_at, updated_at
			  FROM users 
			  WHERE username = $1`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.UserType, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	query := `SELECT id, email, username, password, usertype, created_at, updated_at
			  FROM users 
			  WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.UserType, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	user.UpdatedAt = time.Now()

	query := `UPDATE users
			  SET email = $2, username = $3, password = $4, usertype = $5, updated_at = $6
			  WHERE id = $1`

	result, err := r.db.Exec(query, user.ID, user.Email, user.Username, user.PasswordHash, user.UserType, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute update user in database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %s not found", user.ID)
	}

	return nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users 
			  WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete user from database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	return nil
}
