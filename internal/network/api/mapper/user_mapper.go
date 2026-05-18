package mapper

import (
	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
)

func UserToUserResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
