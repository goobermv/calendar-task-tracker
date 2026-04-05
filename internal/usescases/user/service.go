package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/auth"
	repositories "github.com/goobermv/calendar-task-tracker/internal/repositories/interfaces"

	"github.com/google/uuid"
)

type Service struct {
	userRepo        repositories.UserRepository
	passwordService *auth.PasswordService
	jwtService      *auth.JWTService
}

func NewService(userRepo repositories.UserRepository, passwordService *auth.PasswordService, jwtService *auth.JWTService) *Service {
	return &Service{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

type RegisterRequest struct {
	Email    string
	Username string
	Password string
}

type RegisterResponce struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

func (s *Service) Register(req RegisterRequest) (*RegisterResponce, error) {
	existing_email, _ := s.userRepo.FindByEmail(req.Email)
	if existing_email != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	existing_user, _ := s.userRepo.FindByUsername(req.Username)
	if existing_user != nil {
		return nil, fmt.Errorf("user with this username already exists")
	}

	hashedPassword, err := s.passwordService.Hash(req.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	user := *&domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		UserType:     "regular",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	user.PasswordHash = ""

	return &RegisterResponce{
		User:  &user,
		Token: token,
	}, nil
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if user == nil {
		return nil, errors.New("invalid credentails")
	}

	if !s.passwordService.Verify(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	user.PasswordHash = ""

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *Service) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	user.PasswordHash = ""
	return user, nil
}
