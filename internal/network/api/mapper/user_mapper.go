package mapper

import (
	"fmt"

	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	"github.com/google/uuid"
)

func UserResponseToUser(resp *dto.UserResponse) (*domain.User, error) {
	userID, err := uuid.Parse(resp.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return &domain.User{
		ID:        userID,
		Email:     resp.Email,
		Username:  resp.Username,
		UserType:  resp.UserType,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}, nil
}

func UserToUserResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
