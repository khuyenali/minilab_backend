package service

import (
	"context"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"
	"strings"
	"errors"
	"fmt"
)

// TaskTypeService defines the interface for task type business operations
type TaskTypeService interface {
	GetTaskTypes(ctx context.Context) ([]*models.TaskType, error)
	GetTaskType(ctx context.Context, id int32) (*models.TaskType, error)
	CreateTaskType(ctx context.Context, req models.CreateTaskTypeRequest) (*models.TaskType, error)
	UpdateTaskType(ctx context.Context, id int32, req models.UpdateTaskTypeRequest) (*models.TaskType, error)
	DeleteTaskType(ctx context.Context, id int32) error
}

// taskTypeService implements TaskTypeService
type taskTypeService struct {
	taskTypeRepo repository.TaskTypeRepository
	userRepo     repository.UserRepository
}

// NewTaskTypeService creates a new task type service
func NewTaskTypeService(taskTypeRepo repository.TaskTypeRepository, userRepo repository.UserRepository) TaskTypeService {
	return &taskTypeService{
		taskTypeRepo: taskTypeRepo,
		userRepo:     userRepo,
	}
}

// GetTaskTypes retrieves all task types
func (s *taskTypeService) GetTaskTypes(ctx context.Context) ([]*models.TaskType, error) {
	return s.taskTypeRepo.List(ctx)
}

// GetTaskType retrieves a task type by ID
func (s *taskTypeService) GetTaskType(ctx context.Context, id int32) (*models.TaskType, error) {
	if id <= 0 {
		return nil, ErrInvalidTaskTypeID
	}
	
	taskType, err := s.taskTypeRepo.GetWithMachines(ctx, id)
	if err != nil {
		if err == repository.ErrTaskTypeNotFound {
			return nil, ErrTaskTypeNotFound
		}
		return nil, err
	}
	
	return taskType, nil
}

// CreateTaskType creates a new task type
func (s *taskTypeService) CreateTaskType(ctx context.Context, req models.CreateTaskTypeRequest) (*models.TaskType, error) {
	// Validate input
	if err := s.validateCreateTaskTypeRequest(req); err != nil {
		return nil, err
	}
	
	// Validate user IDs if provided
	if len(req.UserIDs) > 0 {
		if err := s.validateUserIDs(ctx, req.UserIDs); err != nil {
			return nil, err
		}
	}
	
	taskType, err := s.taskTypeRepo.Create(ctx, req)
	if err != nil {
		// Check for machine validation errors
		if errors.Is(err, repository.ErrInvalidMachineIDs) {
			return nil, ErrInvalidMachineIDs
		}
		return nil, err
	}
	
	return taskType, nil
}

// UpdateTaskType updates an existing task type
func (s *taskTypeService) UpdateTaskType(ctx context.Context, id int32, req models.UpdateTaskTypeRequest) (*models.TaskType, error) {
	if id <= 0 {
		return nil, ErrInvalidTaskTypeID
	}
	
	// Validate input
	if err := s.validateUpdateTaskTypeRequest(req); err != nil {
		return nil, err
	}
	
	// Validate user IDs if provided
	if len(req.UserIDs) > 0 {
		if err := s.validateUserIDs(ctx, req.UserIDs); err != nil {
			return nil, err
		}
	}
	
	taskType, err := s.taskTypeRepo.Update(ctx, id, req)
	if err != nil {
		if err == repository.ErrTaskTypeNotFound {
			return nil, ErrTaskTypeNotFound
		}
		// Check for machine validation errors
		if errors.Is(err, repository.ErrInvalidMachineIDs) {
			return nil, ErrInvalidMachineIDs
		}
		return nil, err
	}
	
	return taskType, nil
}

// DeleteTaskType deletes a task type by ID
func (s *taskTypeService) DeleteTaskType(ctx context.Context, id int32) error {
	if id <= 0 {
		return ErrInvalidTaskTypeID
	}
	
	err := s.taskTypeRepo.Delete(ctx, id)
	if err != nil {
		if err == repository.ErrTaskTypeNotFound {
			return ErrTaskTypeNotFound
		}
		return err
	}
	
	return nil
}

// validateCreateTaskTypeRequest validates the create task type request
func (s *taskTypeService) validateCreateTaskTypeRequest(req models.CreateTaskTypeRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidTaskTypeName
	}
	
	if len(req.Name) > 255 {
		return ErrTaskTypeNameTooLong
	}
	
	if req.Description != nil && len(*req.Description) > 1000 {
		return ErrTaskTypeDescriptionTooLong
	}
	
	return nil
}

// validateUpdateTaskTypeRequest validates the update task type request
func (s *taskTypeService) validateUpdateTaskTypeRequest(req models.UpdateTaskTypeRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidTaskTypeName
	}
	
	if len(req.Name) > 255 {
		return ErrTaskTypeNameTooLong
	}
	
	if req.Description != nil && len(*req.Description) > 1000 {
		return ErrTaskTypeDescriptionTooLong
	}
	
	return nil
}

// validateUserIDs checks if the provided user IDs are valid and are members
func (s *taskTypeService) validateUserIDs(ctx context.Context, userIDs []int32) error {
	if len(userIDs) == 0 {
		return nil
	}
	
	// Check each user ID exists and is a member
	var invalidIDs []int32
	var nonMemberIDs []int32
	
	for _, userID := range userIDs {
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			if err == repository.ErrUserNotFound {
				invalidIDs = append(invalidIDs, userID)
				continue
			}
			return err
		}
		
		// Check if user is a member (role_id = 3)
		if user.RoleID != 3 {
			nonMemberIDs = append(nonMemberIDs, userID)
		}
	}
	
	// Return specific error messages
	if len(invalidIDs) > 0 {
		return errors.New("invalid user IDs: " + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(invalidIDs)), ", "), "[]"))
	}
	
	if len(nonMemberIDs) > 0 {
		return errors.New("only members can be assigned to task types, invalid users: " + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(nonMemberIDs)), ", "), "[]"))
	}
	
	return nil
} 