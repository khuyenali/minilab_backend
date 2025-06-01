package service

import (
	"context"
	"fmt"
	"strings"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"
)

// UserService defines the interface for user business operations
type UserService interface {
	GetUser(ctx context.Context, id int32) (*models.User, error)
	GetUsers(ctx context.Context, filters models.UserFilters) ([]*models.User, error)
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error)
	UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id int32) error
	AssignTaskTypes(ctx context.Context, userID int32, req models.UserTaskAssignmentRequest) error
}

// userService implements UserService
type userService struct {
	userRepo       repository.UserRepository
	userToTypeRepo repository.UserToTypeRepository
	taskTypeRepo   repository.TaskTypeRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, userToTypeRepo repository.UserToTypeRepository, taskTypeRepo repository.TaskTypeRepository) UserService {
	return &userService{
		userRepo:       userRepo,
		userToTypeRepo: userToTypeRepo,
		taskTypeRepo:   taskTypeRepo,
	}
}

// GetUser retrieves a user by ID
func (s *userService) GetUser(ctx context.Context, id int32) (*models.User, error) {
	if id <= 0 {
		return nil, ErrInvalidUserID
	}
	
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return user, nil
}

// GetUsers retrieves all users with optional filtering
func (s *userService) GetUsers(ctx context.Context, filters models.UserFilters) ([]*models.User, error) {
	// Apply business rules for filtering
	if filters.Limit <= 0 {
		filters.Limit = 50 // Default limit
	}
	if filters.Limit > 100 {
		filters.Limit = 100 // Max limit
	}
	
	users, err := s.userRepo.List(ctx, filters)
	if err != nil {
		return nil, err
	}
	
	return users, nil
}

// CreateUser creates a new user with business validation
func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	// Business validation
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, err
	}
	
	// Clean and normalize data
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	
	user, err := s.userRepo.Create(ctx, req)
	if err != nil {
		if err == repository.ErrUserEmailExists {
			return nil, ErrUserEmailExists
		}
		return nil, err
	}
	
	return user, nil
}

// UpdateUser updates an existing user with business validation
func (s *userService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.User, error) {
	if id <= 0 {
		return nil, ErrInvalidUserID
	}
	
	// Business validation
	if err := s.validateUpdateUserRequest(req); err != nil {
		return nil, err
	}
	
	// Clean and normalize data
	if req.Name != "" {
		req.Name = strings.TrimSpace(req.Name)
	}
	if req.Email != "" {
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	}
	
	user, err := s.userRepo.Update(ctx, id, req)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		if err == repository.ErrUserEmailExists {
			return nil, ErrUserEmailExists
		}
		return nil, err
	}
	
	return user, nil
}

// DeleteUser deletes a user by ID
func (s *userService) DeleteUser(ctx context.Context, id int32) error {
	if id <= 0 {
		return ErrInvalidUserID
	}
	
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return ErrUserNotFound
		}
		return err
	}
	
	return nil
}

// validateCreateUserRequest validates create user request
func (s *userService) validateCreateUserRequest(req models.CreateUserRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidUserName
	}
	
	if len(req.Name) > 255 {
		return ErrUserNameTooLong
	}
	
	if strings.TrimSpace(req.Email) == "" {
		return ErrInvalidUserEmail
	}
	
	if len(req.Email) > 255 {
		return ErrUserEmailTooLong
	}
	
	// Basic email validation (Gin's email validator is more comprehensive)
	if !strings.Contains(req.Email, "@") {
		return ErrInvalidUserEmail
	}
	
	return nil
}

// validateUpdateUserRequest validates update user request
func (s *userService) validateUpdateUserRequest(req models.UpdateUserRequest) error {
	if req.Name != "" {
		if strings.TrimSpace(req.Name) == "" {
			return ErrInvalidUserName
		}
		if len(req.Name) > 255 {
			return ErrUserNameTooLong
		}
	}
	
	if req.Email != "" {
		if strings.TrimSpace(req.Email) == "" {
			return ErrInvalidUserEmail
		}
		if len(req.Email) > 255 {
			return ErrUserEmailTooLong
		}
		if !strings.Contains(req.Email, "@") {
			return ErrInvalidUserEmail
		}
	}
	
	return nil
}

// AssignTaskTypes assigns task types to a user (only for members)
func (s *userService) AssignTaskTypes(ctx context.Context, userID int32, req models.UserTaskAssignmentRequest) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}
	
	// Get user to check role
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return ErrUserNotFound
		}
		return err
	}
	
	// Only members (role_id = 3) can be assigned task types
	if user.RoleID != 3 {
		return ErrInvalidUserRole
	}
	
	// Validate task type IDs
	if len(req.TaskTypeIDs) == 0 {
		return ErrInvalidTaskTypeAssignment
	}
	
	// Validate that all task type IDs exist
	var invalidIDs []int32
	for _, taskTypeID := range req.TaskTypeIDs {
		_, err := s.taskTypeRepo.GetByID(ctx, taskTypeID)
		if err != nil {
			if err == repository.ErrTaskTypeNotFound {
				invalidIDs = append(invalidIDs, taskTypeID)
			} else {
				return err
			}
		}
	}
	
	if len(invalidIDs) > 0 {
		return fmt.Errorf("invalid task type IDs: %v", invalidIDs)
	}
	
	// Assign task types
	err = s.userToTypeRepo.AssignTaskTypes(ctx, userID, req.TaskTypeIDs)
	if err != nil {
		return err
	}
	
	return nil
} 
