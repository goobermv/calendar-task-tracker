package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/middleware"
	userUsecase "github.com/goobermv/calendar-task-tracker/internal/usescases/user"
	"github.com/google/uuid"
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
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation failed",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	userCaseReq := dto.RegisterRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	response, err := h.userService.Register(userCaseReq)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "user with this email already exists" || err.Error() == "user with this username already exists" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   err.Error(),
			Code:    statusCode,
			Details: "Please check your input and try again",
		})
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
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation failed",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	userUseReq := dto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := h.userService.Login(userUseReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   err.Error(),
			Code:    http.StatusUnauthorized,
			Details: "Email or password is incorrect",
		})
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
	userIDstr, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
			Details: "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "User not found",
			Code:    http.StatusNotFound,
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		UserType:  string(user.UserType),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
