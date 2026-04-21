package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/auth"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/logger"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/mapper"
	repositories "github.com/goobermv/calendar-task-tracker/internal/repositories/interfaces"

	"github.com/google/uuid"
)

type Service struct {
	userRepo        repositories.UserRepository
	passwordService *auth.PasswordService
	jwtService      *auth.JWTService
	logger          logger.Logger
}

func NewService(userRepo repositories.UserRepository, passwordService *auth.PasswordService, jwtService *auth.JWTService, logger logger.Logger) *Service {
	return &Service{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
		logger:          logger,
	}
}

func (s *Service) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	s.logger.Infof("Attempting to register user with email: %s", req.Email)

	existing_email, _ := s.userRepo.FindByEmail(req.Email)
	if existing_email != nil {
		s.logger.Warnf("Registration failed: email already exists - %s", req.Email)
		return nil, domain.ErrEmailExists
	}

	existing_user, _ := s.userRepo.FindByUsername(req.Username)
	if existing_user != nil {
		s.logger.Warnf("Registration failed: username already exists - %s", req.Username)
		return nil, domain.ErrUsernameExists
	}

	hashedPassword, err := s.passwordService.Hash(req.Password)
	if err != nil {
		s.logger.Errorf(err, "Failed to hash password for user: %s", req.Email)
		return nil, errors.New("failed to process password")
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		UserType:     "regular",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		s.logger.Errorf(err, "Failed to create user in database: %s", req.Email)
		return nil, err
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		s.logger.Errorf(err, "Failed to generate token for user: %s", user.ID)
		return nil, errors.New("failed to generate token")
	}

	user.PasswordHash = ""

	userResponse := mapper.UserToUserResponse(user)

	s.logger.Successf("User registered successfully - ID: %s, Email: %s", user.ID, user.Email)

	return &dto.RegisterResponse{
		User:  *userResponse,
		Token: token,
	}, nil
}

func (s *Service) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	s.logger.Infof("Login attempt for email: %s", req.Email)

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		s.logger.Errorf(err, "Database error during login for email: %s", req.Email)
		return nil, domain.ErrInvalidCredentials
	}
	if user == nil {
		s.logger.Errorf(err, "Login failed: user not found - %s", req.Email)
		return nil, domain.ErrInvalidCredentials
	}

	if !s.passwordService.Verify(req.Password, user.PasswordHash) {
		s.logger.Errorf(err, "Login failed: invalid password for user - %s", req.Email)
		return nil, domain.ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		s.logger.Errorf(err, "Failed to generate token for user: %s", user.ID)
		return nil, errors.New("failed to generate token")
	}

	user.PasswordHash = ""

	userResonse := mapper.UserToUserResponse(user)

	s.logger.Successf("User successfully logged in - ID: %s, Email: %s", user.ID, user.Email)

	return &dto.LoginResponse{
		User:  *userResonse,
		Token: token,
	}, nil
}

func (s *Service) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	s.logger.Infof("Attempting to get user by ID: %s", userID)

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Errorf(err, "Failed to find user: %s", userID)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		s.logger.Errorf(err, "User not found: %s", userID)
		return nil, domain.ErrUserNotFound
	}

	user.PasswordHash = ""

	s.logger.Successf("User retrieved successfully: %s", userID)
	return user, nil
}

func (s *Service) UpdateUserInfo(userID uuid.UUID, req dto.UpdateUserInfoRequest) (*domain.User, error) {
	s.logger.Infof("Updating user info for ID: %s", userID)

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Errorf(err, "Failed to find user for update: %s", userID)
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	if user == nil {
		s.logger.Warnf("User not found for update: %s", userID)
		return nil, domain.ErrUserNotFound
	}

	if user.ID != userID {
		s.logger.Warnf("User ID mismatch - expected: %s, got: %s", userID, user.ID)
		return nil, domain.ErrNoPermissionUpdateUser
	}

	if req.Email != nil {
		s.logger.Debugf("Checking email availability: %s", *req.Email)
		existing_email, _ := s.userRepo.FindByEmail(*req.Email)
		if existing_email != nil {
			s.logger.Warnf("Email already taken: %s", *req.Email)
			return nil, domain.ErrEmailExists
		}
		user.Email = *req.Email
		s.logger.Infof("Email updated to: %s for user: %s", *req.Email, userID)
	}

	if req.Username != nil {
		s.logger.Debugf("Checking username availability: %s", *req.Username)
		existing_user, _ := s.userRepo.FindByUsername(*req.Username)
		if existing_user != nil {
			s.logger.Warnf("Username already taken: %s", *req.Username)
			return nil, domain.ErrUsernameExists
		}
		user.Username = *req.Username
		s.logger.Infof("Username updated to: %s for user: %s", *req.Username, userID)
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		s.logger.Errorf(err, "Failed to update user in database: %s", userID)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Successf("User info updated successfully: %s", userID)
	return user, nil
}

func (s *Service) UpdateUserPassword(userID uuid.UUID, req dto.UpdateUserPasswordRequest) error {
	s.logger.Infof("Updating password for user: %s", userID)

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Errorf(err, "Failed to find user for password update: %s", userID)
		return fmt.Errorf("failed to find user by ID: %w", err)
	}
	if user == nil {
		s.logger.Warnf("User not found for password update: %s", userID)
		return domain.ErrUserNotFound
	}

	if user.ID != userID {
		s.logger.Warnf("User ID mismatch for password update - expected: %s, got: %s", userID, user.ID)
		return domain.ErrNoPermissionUpdateUser
	}

	if req.Password != nil {
		s.logger.Debugf("Hashing new password for user: %s", userID)
		hashedPassword, err := s.passwordService.Hash(*req.Password)
		if err != nil {
			s.logger.Errorf(err, "Failed to hash new password for user: %s", userID)
			return errors.New("failed to process password")
		}

		if user.PasswordHash == hashedPassword {
			s.logger.Warnf("New password is same as old password for user: %s", userID)
			return errors.New("new password cannot be the same as old password")
		}

		user.PasswordHash = hashedPassword
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		s.logger.Errorf(err, "Failed to update password in database for user: %s", userID)
		return fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Successf("Password updated successfully for user: %s", userID)
	return nil
}

func (s *Service) DeleteUser(userID uuid.UUID) error {
	s.logger.Infof("Deleting user: %s", userID)

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Errorf(err, "Failed to find user for deletion: %s", userID)
		return fmt.Errorf("failed to find user by ID: %w", err)
	}
	if user == nil {
		s.logger.Warnf("User not found for deletion: %s", userID)
		return domain.ErrUserNotFound
	}

	if err := s.userRepo.Delete(userID); err != nil {
		s.logger.Errorf(err, "Failed to delete user from database: %s", userID)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Successf("User deleted successfully: %s (Email: %s)", userID, user.Email)
	return nil
}

func (s *Service) PromoteUserToAdmin(userID uuid.UUID) error {
	s.logger.Infof("Promoting user to admin: %s", userID)

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Errorf(err, "Failed to find user for promotion: %s", userID)
		return fmt.Errorf("failed to find user by ID: %w", err)
	}
	if user == nil {
		s.logger.Warnf("User not found for promotion: %s", userID)
		return domain.ErrUserNotFound
	}

	if user.UserType == domain.UserTypeAdmin {
		s.logger.Warnf("User is already an admin: %s", userID)
		return errors.New("user is already an admin")
	}

	user.UserType = domain.UserTypeAdmin
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		s.logger.Errorf(err, "Failed to update user role to admin: %s", userID)
		return fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Successf("User promoted to admin successfully: %s (Email: %s)", userID, user.Email)
	return nil
}

func (s *Service) DemoteAdminToUser(userID uuid.UUID) error {
	s.logger.Infof("Demoting admin to user: %s", userID)

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Errorf(err, "Failed to find user for demotion: %s", userID)
		return fmt.Errorf("failed to find user by ID: %w", err)
	}
	if user == nil {
		s.logger.Warnf("User not found for demotion: %s", userID)
		return domain.ErrUserNotFound
	}

	if user.UserType != domain.UserTypeAdmin {
		s.logger.Warnf("User is not an admin: %s (Current role: %s)", userID, user.UserType)
		return errors.New("user is not an admin")
	}

	user.UserType = domain.UserTypeRegular
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		s.logger.Errorf(err, "Failed to update user role to regular: %s", userID)
		return fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Successf("Admin demoted to user successfully: %s (Email: %s)", userID, user.Email)
	return nil
}
