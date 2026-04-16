package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/auth"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/mapper"
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

func (s *Service) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
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

	userResponse := mapper.UserToUserResponse(&user)

	return &dto.RegisterResponse{
		User:  *userResponse,
		Token: token,
	}, nil
}

func (s *Service) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
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

	userResonse := mapper.UserToUserResponse(user)

	return &dto.LoginResponse{
		User:  *userResonse,
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

func (s *Service) UpdateUserInfo(userID uuid.UUID, req dto.UpdateUserInfoRequest) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.ID != userID {
		return nil, errors.New("you do not have the permissions to update this user")
	}

	if req.Email != nil {
		existing_email, _ := s.userRepo.FindByEmail(*req.Email)
		if existing_email != nil {
			return nil, fmt.Errorf("user with this email already exists")
		}
		user.Email = *req.Email
	}

	if req.Username != nil {
		existing_user, _ := s.userRepo.FindByUsername(*req.Username)
		if existing_user != nil {
			return nil, fmt.Errorf("user with this username already exists")
		}
		user.Username = *req.Username
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *Service) UpdateUserPassword(userID uuid.UUID, req dto.UpdateUserPasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to find user by ID: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	if user.ID != userID {
		return errors.New("you do not have the permissions to update this user")
	}

	if req.Password != nil {
		hashedPassword, err := s.passwordService.Hash(*req.Password)
		if err != nil {
			return errors.New("failed to process password")
		}

		if user.PasswordHash == hashedPassword {
			return errors.New("new password cannot be the same as old password")
		}

		user.PasswordHash = hashedPassword
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *Service) DeleteUser(userID uuid.UUID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to find user by ID: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	if err := s.userRepo.Delete(userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
