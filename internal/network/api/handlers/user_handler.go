package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/middleware"
	userUsecase "github.com/goobermv/calendar-task-tracker/internal/usescases/user"
)

type UserHandler struct {
	userService *userUsecase.Service
}

func NewUserHandler(userService *userUsecase.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		domain.HandleValidationError(c, err)
		return
	}

	response, err := h.userService.Register(req)
	if err != nil {
		domain.HandleAlreadyExistsError(c, err)
		return
	}

	userResponse := dto.UserResponse{
		ID:        response.User.ID,
		Email:     response.User.Email,
		Username:  response.User.Username,
		UserType:  string(response.User.UserType),
		CreatedAt: response.User.CreatedAt,
		UpdatedAt: response.User.UpdatedAt,
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{
		User:  userResponse,
		Token: response.Token,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		domain.HandleValidationError(c, err)
		return
	}
	response, err := h.userService.Login(req)
	if err != nil {
		domain.HandleIncorrectInfoError(c, err)
		return
	}

	userResponse := dto.UserResponse{
		ID:        response.User.ID,
		Email:     response.User.Email,
		Username:  response.User.Username,
		UserType:  string(response.User.UserType),
		CreatedAt: response.User.CreatedAt,
		UpdatedAt: response.User.UpdatedAt,
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		User:  userResponse,
		Token: response.Token,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		domain.HandleUnathorizedError(c)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		domain.HandleUserNotFoundError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		UserType:  string(user.UserType),
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now(),
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		domain.HandleUnathorizedError(c)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		domain.HandleValidationError(c, err)
		return
	}

	user, err := h.userService.UpdateUser(userID, req)
	if err != nil {
		domain.HandleUpdateUserError(c, err)
		return
	}

	response := dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		UserType:  string(user.UserType),
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		domain.HandleUnathorizedError(c)
		return
	}

	err := h.userService.DeleteUser(userID)
	if err != nil {
		domain.HandleDeleteUserError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
