package service

import (
	"context"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"
	"strings"
	"errors"
	"fmt"
)

// TaskService defines the interface for task business operations
type TaskService interface {
	GetTasks(ctx context.Context) ([]*models.Task, error)
	GetTask(ctx context.Context, id int32) (*models.Task, error)
	CreateTask(ctx context.Context, req models.CreateTaskRequest) (*models.Task, error)
	UpdateTask(ctx context.Context, id int32, req models.UpdateTaskRequest) (*models.Task, error)
	UpdateTaskStatusToPending(ctx context.Context, id int32) (*models.Task, error)
	DeleteTask(ctx context.Context, id int32) error
}

// taskService implements TaskService
type taskService struct {
	taskRepo     repository.TaskRepository
	taskTypeRepo repository.TaskTypeRepository
	userRepo     repository.UserRepository
}

// NewTaskService creates a new task service
func NewTaskService(taskRepo repository.TaskRepository, taskTypeRepo repository.TaskTypeRepository, userRepo repository.UserRepository) TaskService {
	return &taskService{
		taskRepo:     taskRepo,
		taskTypeRepo: taskTypeRepo,
		userRepo:     userRepo,
	}
}

// GetTasks retrieves all tasks
func (s *taskService) GetTasks(ctx context.Context) ([]*models.Task, error) {
	return s.taskRepo.ListWithSubTasks(ctx)
}

// GetTask retrieves a task by ID with sub-tasks and assignments
func (s *taskService) GetTask(ctx context.Context, id int32) (*models.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidTaskID
	}
	
	task, err := s.taskRepo.GetWithSubTasks(ctx, id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	
	return task, nil
}

// CreateTask creates a new task with sub-tasks and assignments
func (s *taskService) CreateTask(ctx context.Context, req models.CreateTaskRequest) (*models.Task, error) {
	// Validate input
	if err := s.validateCreateTaskRequest(req); err != nil {
		return nil, err
	}
	
	// Validate sub-tasks if provided
	if len(req.SubTasks) > 0 {
		if err := s.validateSubTaskRequests(ctx, req.SubTasks); err != nil {
			return nil, err
		}
	}
	
	task, err := s.taskRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return task, nil
}

// UpdateTask updates an existing task
func (s *taskService) UpdateTask(ctx context.Context, id int32, req models.UpdateTaskRequest) (*models.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidTaskID
	}
	
	// Validate input
	if err := s.validateUpdateTaskRequest(req); err != nil {
		return nil, err
	}
	
	task, err := s.taskRepo.Update(ctx, id, req)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	
	return task, nil
}

// UpdateTaskStatusToPending updates a task's status from draft to pending
func (s *taskService) UpdateTaskStatusToPending(ctx context.Context, id int32) (*models.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidTaskID
	}
	
	// Get current task to check status
	currentTask, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	
	// Check if task is currently in draft status
	if currentTask.Status != "draft" {
		return nil, ErrInvalidTaskStatusTransition
	}
	
	// Create update request to change status to pending
	updateReq := models.UpdateTaskRequest{
		TaskName:  currentTask.TaskName,
		Status:    "pending",
		Priority:  currentTask.Priority,
		StartTime: currentTask.StartTime,
		EndTime:   currentTask.EndTime,
		Note:      currentTask.Note,
	}
	
	// Update the task
	task, err := s.taskRepo.Update(ctx, id, updateReq)
	if err != nil {
		return nil, err
	}
	
	return task, nil
}

// DeleteTask deletes a task by ID
func (s *taskService) DeleteTask(ctx context.Context, id int32) error {
	if id <= 0 {
		return ErrInvalidTaskID
	}
	
	err := s.taskRepo.Delete(ctx, id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			return ErrTaskNotFound
		}
		return err
	}
	
	return nil
}

// validateCreateTaskRequest validates the create task request
func (s *taskService) validateCreateTaskRequest(req models.CreateTaskRequest) error {
	if strings.TrimSpace(req.TaskName) == "" {
		return ErrInvalidTaskName
	}
	
	if len(req.TaskName) > 255 {
		return ErrTaskNameTooLong
	}
	
	if req.Note != nil && len(*req.Note) > 1000 {
		return ErrTaskNoteTooLong
	}
	
	// Validate status if provided
	if req.Status != "" && !isValidTaskStatus(req.Status) {
		return ErrInvalidTaskStatus
	}
	
	// Validate priority if provided
	if req.Priority != 0 && !isValidTaskPriority(req.Priority) {
		return ErrInvalidTaskPriority
	}
	
	return nil
}

// validateUpdateTaskRequest validates the update task request
func (s *taskService) validateUpdateTaskRequest(req models.UpdateTaskRequest) error {
	if strings.TrimSpace(req.TaskName) == "" {
		return ErrInvalidTaskName
	}
	
	if len(req.TaskName) > 255 {
		return ErrTaskNameTooLong
	}
	
	if req.Note != nil && len(*req.Note) > 1000 {
		return ErrTaskNoteTooLong
	}
	
	// Validate status
	if !isValidTaskStatus(req.Status) {
		return ErrInvalidTaskStatus
	}
	
	// Validate priority
	if !isValidTaskPriority(req.Priority) {
		return ErrInvalidTaskPriority
	}
	
	return nil
}

// validateSubTaskRequests validates the sub-task requests
func (s *taskService) validateSubTaskRequests(ctx context.Context, subTasks []models.SubTaskRequest) error {
	for i, subTask := range subTasks {
		// Validate type_id exists
		_, err := s.taskTypeRepo.GetByID(ctx, subTask.TypeID)
		if err != nil {
			if err == repository.ErrTaskTypeNotFound {
				return errors.New(fmt.Sprintf("invalid type_id %d in sub-task %d: task type not found", subTask.TypeID, i+1))
			}
			return err
		}
		
		// Validate user IDs if provided
		if len(subTask.UserIDs) > 0 {
			if err := s.validateUserIDs(ctx, subTask.UserIDs, i+1); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// validateUserIDs checks if the provided user IDs are valid and are members
func (s *taskService) validateUserIDs(ctx context.Context, userIDs []int32, subTaskIndex int) error {
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
		return errors.New(fmt.Sprintf("invalid user IDs in sub-task %d: %v", subTaskIndex, invalidIDs))
	}
	
	if len(nonMemberIDs) > 0 {
		return errors.New(fmt.Sprintf("only members can be assigned to sub-tasks, invalid users in sub-task %d: %v", subTaskIndex, nonMemberIDs))
	}
	
	return nil
}

// isValidTaskStatus checks if the status is valid
func isValidTaskStatus(status string) bool {
	validStatuses := []string{"draft", "pending", "processing", "finish"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// isValidTaskPriority checks if the priority is valid
func isValidTaskPriority(priority int32) bool {
	return priority >= 1 && priority <= 3 // 1=high, 2=medium, 3=low
} 