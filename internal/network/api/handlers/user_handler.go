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
		HandleError(c, err)
		return
	}

	response, err := h.userService.Register(req)
	if err != nil {
		HandleError(c, err)
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
		HandleError(c, err)
		return
	}
	response, err := h.userService.Login(req)
	if err != nil {
		HandleError(c, err)
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
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		UserType:  string(user.UserType),
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now(),
	})
}

func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	var req dto.UpdateUserInfoRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	user, err := h.userService.UpdateUserInfo(userID, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		UserType:  string(user.UserType),
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now(),
	})
}

func (h *UserHandler) UpdateUserPassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	var req dto.UpdateUserPasswordRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.userService.UpdateUserPassword(userID, req); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	err := h.userService.DeleteUser(userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *UserHandler) PromoteUserToAdmin(c *gin.Context) {
	adminID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	targetUserID, exists := middleware.GetID(c)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing user ID in path"})
		return
	}

	if adminID == targetUserID {
		HandleError(c, domain.ErrCannotPromoteYourself)
		return
	}

	err := h.userService.PromoteUserToAdmin(targetUserID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User promoted to admin successfully",
		"user_id": targetUserID.String(),
	})
}

func (h *UserHandler) DemoteAdminToUser(c *gin.Context) {
	adminID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	targetUserID, exists := middleware.GetID(c)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing user ID in path"})
		return
	}

	if adminID == targetUserID {
		HandleError(c, domain.ErrCannotDemoteYourself)
		return
	}

	err := h.userService.DemoteAdminToUser(targetUserID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin demoted to user successfully",
		"user_id": targetUserID.String(),
	})
}

func (h *UserHandler) AdminGetAllUsers(c *gin.Context) {
	users, err := h.userService.AdminGetAllUsers()
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}
