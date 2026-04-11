package dto

import (
	"time"
)

type RegisterRequest struct {
	Email    string `json:"email" binging:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binging:"required,min=6,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" binging:"required,email"`
	Password string `json:"password" binging:"required"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	UserType  string    `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type RegisterResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
