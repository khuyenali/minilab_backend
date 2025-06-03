package service

import (
	"context"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"
	"strings"
	"errors"
	"fmt"
	"database/sql"
	"mini-lab-api/internal/db"
)

// TaskService defines the interface for task business operations
type TaskService interface {
	GetTasks(ctx context.Context) ([]*models.Task, error)
	GetTask(ctx context.Context, id int32) (*models.Task, error)
	CreateTask(ctx context.Context, req models.CreateTaskRequest) (*models.Task, error)
	UpdateTask(ctx context.Context, id int32, req models.UpdateTaskRequest) (*models.Task, error)
	UpdateTaskStatusToPending(ctx context.Context, id int32) (*models.Task, error)
	DeleteTask(ctx context.Context, id int32) error
	GetAvailableTasks(ctx context.Context) ([]int32, error)
}

// taskService implements TaskService
type taskService struct {
	taskRepo     repository.TaskRepository
	taskTypeRepo repository.TaskTypeRepository
	userRepo     repository.UserRepository
	db           *sql.DB
	queries      *db.Queries
}

// NewTaskService creates a new task service
func NewTaskService(taskRepo repository.TaskRepository, taskTypeRepo repository.TaskTypeRepository, userRepo repository.UserRepository, database *sql.DB) TaskService {
	return &taskService{
		taskRepo:     taskRepo,
		taskTypeRepo: taskTypeRepo,
		userRepo:     userRepo,
		db:           database,
		queries:      db.New(database),
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

// GetAvailableTasks returns available task IDs based on priority and user loading constraints
func (s *taskService) GetAvailableTasks(ctx context.Context) ([]int32, error) {
	// Get draft tasks sorted by priority (1=high, 2=medium, 3=low)
	// Only draft tasks are available for auto assignment
	draftTasks, err := s.queries.GetDraftTasksSortedByPriority(ctx)
	if err != nil {
		return nil, err
	}
	
	// Get all member users (role_id = 3)
	memberUsers, err := s.queries.GetUsersByRoleID(ctx, 3)
	if err != nil {
		return nil, err
	}
	
	// Calculate current user loading (in minutes)
	userLoadings := make(map[int32]int32) // user_id -> total minutes
	
	for _, user := range memberUsers {
		// Get active assignments for this user
		activeAssignments, err := s.queries.GetActiveAssignmentsByUserID(ctx, user.ID)
		if err != nil {
			continue // Skip this user if we can't get their assignments
		}
		
		totalMinutes := int32(0)
		for _, assignment := range activeAssignments {
			// Get machines for this task type to calculate duration
			machines, err := s.queries.GetMachinesByTaskType(ctx, sql.NullInt32{Int32: assignment.TypeID, Valid: true})
			if err != nil {
				continue // Skip if we can't get machines
			}
			
			// Sum up estimate times from all machines for this task type
			taskTypeDuration := int32(0)
			for _, machine := range machines {
				taskTypeDuration += machine.EstimateTime
			}
			
			totalMinutes += taskTypeDuration
		}
		
		userLoadings[user.ID] = totalMinutes
	}
	
	// Calculate current machine usage for each task type
	// Only count pending/process assignments as using machines
	// Finished assignments have released their machines
	machineUsage := make(map[int32]int32) // task_type_id -> current usage count
	
	// Get all pending and processing tasks to count machine usage
	pendingTasks, err := s.queries.GetPendingTasksSortedByPriority(ctx)
	if err != nil {
		return nil, err
	}
	
	for _, task := range pendingTasks {
		subTasks, err := s.queries.GetSubTasksByTaskID(ctx, task.ID)
		if err != nil {
			continue
		}
		
		for _, subTask := range subTasks {
			// Get assignments for this sub-task
			assignments, err := s.queries.GetAssignmentsBySubTaskID(ctx, subTask.ID)
			if err != nil {
				continue
			}
			
			// Count assignments that are pending or processing (using machines)
			for _, assignment := range assignments {
				status := "pending" // default
				if assignment.Status.Valid {
					status = string(assignment.Status.AssignmentStatus)
				}
				
				if status == "pending" || status == "process" {
					machineUsage[subTask.TypeID]++
				}
			}
		}
	}
	
	// Filter draft tasks based on user loading constraints and machine availability
	var availableTaskIDs []int32
	const maxUserLoadingMinutes = 8 * 60 // 8 hours in minutes
	
	// Track how many machines would be needed if we include each task
	projectedMachineUsage := make(map[int32]int32)
	for k, v := range machineUsage {
		projectedMachineUsage[k] = v
	}
	
	for _, task := range draftTasks {
		// Get sub-tasks for this task
		subTasks, err := s.queries.GetSubTasksByTaskID(ctx, task.ID)
		if err != nil {
			continue // Skip this task if we can't get sub-tasks
		}
		
		canAssignTask := true
		var requiredMachines []int32 // Track which task types this task will use
		
		// Check each sub-task for both user loading and machine availability
		for _, subTask := range subTasks {
			// Calculate duration for this sub-task
			machines, err := s.queries.GetMachinesByTaskType(ctx, sql.NullInt32{Int32: subTask.TypeID, Valid: true})
			if err != nil {
				canAssignTask = false
				break
			}
			
			subTaskDuration := int32(0)
			totalMachineQuantity := int32(0)
			for _, machine := range machines {
				subTaskDuration += machine.EstimateTime
				totalMachineQuantity += machine.Quantity
			}
			
			// Check if there are available machines for this task type
			if projectedMachineUsage[subTask.TypeID] >= totalMachineQuantity {
				canAssignTask = false
				break
			}
			
			// Check if any member user can take this sub-task without exceeding 8 hours
			canAssignSubTask := false
			for _, user := range memberUsers {
				if userLoadings[user.ID] + subTaskDuration <= maxUserLoadingMinutes {
					canAssignSubTask = true
					break
				}
			}
			
			if !canAssignSubTask {
				canAssignTask = false
				break
			}
			
			// Track this task type for machine reservation
			requiredMachines = append(requiredMachines, subTask.TypeID)
		}
		
		if canAssignTask {
			availableTaskIDs = append(availableTaskIDs, task.ID)
			// Reserve machines for this task - increment usage for each sub-task
			for _, taskTypeID := range requiredMachines {
				projectedMachineUsage[taskTypeID]++
			}
		}
	}
	
	return availableTaskIDs, nil
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