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
	UpdateTaskStatus(ctx context.Context, id int32, req models.UpdateTaskStatusRequest) (*models.Task, error)
	FinishTask(ctx context.Context, id int32, req models.FinishTaskRequest) (*models.Task, error)
	DeleteTask(ctx context.Context, id int32) error
	GetAvailableTasks(ctx context.Context) ([]int32, error)
	CleanDraftTaskAssignments(ctx context.Context) error
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

// UpdateTaskStatus updates a task's status
func (s *taskService) UpdateTaskStatus(ctx context.Context, id int32, req models.UpdateTaskStatusRequest) (*models.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidTaskID
	}
	
	// Validate status
	if !isValidTaskStatus(req.Status) {
		return nil, ErrInvalidTaskStatus
	}
	
	// Get current task to check status
	currentTask, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	
	// Validate status transition
	if err := s.validateStatusTransition(currentTask.Status, req.Status); err != nil {
		return nil, err
	}
	
	var dbTask db.Task
	
	// If status is finish and no report is provided, return error
	if req.Status == "finish" && (req.Report == nil || *req.Report == "") {
		return nil, ErrTaskReportRequired
	}
	
	// Update with or without report
	if req.Status == "finish" && req.Report != nil {
		dbTask, err = s.queries.UpdateTaskStatusWithReport(ctx, db.UpdateTaskStatusWithReportParams{
			ID:     id,
			Status: db.TaskStatus(req.Status),
			Report: sql.NullString{String: *req.Report, Valid: true},
		})
	} else {
		dbTask, err = s.queries.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
			ID:     id,
			Status: db.TaskStatus(req.Status),
		})
	}
	
	if err != nil {
		return nil, err
	}
	
	return models.FromDBTask(dbTask), nil
}

// FinishTask finishes a task with a required report
func (s *taskService) FinishTask(ctx context.Context, id int32, req models.FinishTaskRequest) (*models.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidTaskID
	}
	
	if req.Report == "" {
		return nil, ErrTaskReportRequired
	}
	
	// Get current task to check status
	currentTask, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	
	// Validate status transition to finish
	if err := s.validateStatusTransition(currentTask.Status, "finish"); err != nil {
		return nil, err
	}
	
	// Update task status to finish with report
	dbTask, err := s.queries.UpdateTaskStatusWithReport(ctx, db.UpdateTaskStatusWithReportParams{
		ID:     id,
		Status: db.TaskStatusFinish,
		Report: sql.NullString{String: req.Report, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	
	return models.FromDBTask(dbTask), nil
}

// validateStatusTransition validates if the status transition is allowed
func (s *taskService) validateStatusTransition(currentStatus, newStatus string) error {
	// Define allowed transitions
	allowedTransitions := map[string][]string{
		"draft":      {"pending"},
		"pending":    {"processing"},
		"processing": {"finish"},
		"finish":     {}, // No transitions allowed from finish
	}
	
	validTransitions, exists := allowedTransitions[currentStatus]
	if !exists {
		return ErrInvalidTaskStatusTransition
	}
	
	for _, validStatus := range validTransitions {
		if validStatus == newStatus {
			return nil
		}
	}
	
	return ErrInvalidTaskStatusTransition
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

// GetAvailableTasks returns available task IDs with auto-assignment based on priority, machine availability, and user task_type matching
func (s *taskService) GetAvailableTasks(ctx context.Context) ([]int32, error) {
	// Get draft tasks sorted by priority (1=high, 2=medium, 3=low)
	draftTasks, err := s.queries.GetDraftTasksSortedByPriority(ctx)
	if err != nil {
		return nil, err
	}
	
	// Get all member users (role_id = 3) with their task types
	memberUsers, err := s.queries.GetUsersByRoleID(ctx, 3)
	if err != nil {
		return nil, err
	}
	
	// Build user task types mapping
	userTaskTypes := make(map[int32][]int32) // user_id -> []task_type_ids
	for _, user := range memberUsers {
		userTaskTypeRows, err := s.queries.GetUserTaskTypes(ctx, user.ID)
		if err != nil {
			continue // Skip user if we can't get their task types
		}
		
		var taskTypes []int32
		for _, row := range userTaskTypeRows {
			taskTypes = append(taskTypes, row.TypeID)
		}
		userTaskTypes[user.ID] = taskTypes
	}
	
	// Find users who are already occupied (have assignments in pending/processing tasks)
	occupiedUsers := make(map[int32]bool)
	for _, user := range memberUsers {
		hasActiveTasks, err := s.queries.CheckUserHasActiveTasks(ctx, user.ID)
		if err != nil {
			continue // Skip if we can't check user status
		}
		if hasActiveTasks {
			occupiedUsers[user.ID] = true
		}
	}
	
	// Calculate current machine usage for each task type
	// Count pending and processing tasks that use machines
	machineUsage := make(map[int32]int32) // task_type_id -> current usage count
	
	// Get pending tasks
	pendingTasks, err := s.queries.GetPendingTasksSortedByPriority(ctx)
	if err != nil {
		return nil, err
	}
	
	// Get processing tasks
	processingTasks, err := s.queries.GetProcessingTasksSortedByPriority(ctx)
	if err != nil {
		return nil, err
	}
	
	// Count machine usage from pending and processing tasks
	allActiveTasks := append(pendingTasks, processingTasks...)
	for _, task := range allActiveTasks {
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
			
			// Each assignment consumes 1 machine usage for this task type
			machineUsage[subTask.TypeID] += int32(len(assignments))
		}
	}
	
	// Get machine quantities for each task type
	taskTypeMachines := make(map[int32]int32) // task_type_id -> total_quantity
	allTaskTypes, err := s.taskTypeRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	
	for _, taskType := range allTaskTypes {
		taskTypeWithMachines, err := s.taskTypeRepo.GetWithMachines(ctx, taskType.ID)
		if err != nil {
			continue
		}
		
		totalQuantity := int32(0)
		for _, machine := range taskTypeWithMachines.Machines {
			totalQuantity += machine.Quantity
		}
		taskTypeMachines[taskType.ID] = totalQuantity
	}
	
	// Process each draft task and check if it can be auto-assigned
	var availableTaskIDs []int32
	
	for _, task := range draftTasks {
		// Get sub-tasks for this task
		subTasks, err := s.queries.GetSubTasksByTaskID(ctx, task.ID)
		if err != nil {
			continue // Skip task if we can't get sub-tasks
		}
		
		canAssignTask := true
		taskAssignments := make(map[int32][]int32) // sub_task_id -> []user_ids to assign
		projectedMachineUsage := make(map[int32]int32) // Copy current usage for projection
		for k, v := range machineUsage {
			projectedMachineUsage[k] = v
		}
		
		hasAssignableSubTasks := false // Track if there are sub-tasks that can be assigned
		
		// Check each sub-task
		for _, subTask := range subTasks {
			// Check if this sub-task already has assignments
			existingAssignments, err := s.queries.GetAssignmentsBySubTaskID(ctx, subTask.ID)
			if err != nil {
				continue // Skip if we can't check assignments
			}
			
			// Skip sub-task if it already has assignments
			if len(existingAssignments) > 0 {
				continue
			}
			
			hasAssignableSubTasks = true // Found at least one unassigned sub-task
			
			// Check if there are enough machines available for this task type
			totalMachineQuantity := taskTypeMachines[subTask.TypeID]
			if projectedMachineUsage[subTask.TypeID] >= totalMachineQuantity {
				canAssignTask = false
				break
			}
			
			// Find available users who can handle this task type
			var availableUsers []int32
			for _, user := range memberUsers {
				// Skip if user is already occupied
				if occupiedUsers[user.ID] {
					continue
				}
				
				// Check if user has the required task type
				userTypes := userTaskTypes[user.ID]
				hasTaskType := false
				for _, userTaskType := range userTypes {
					if userTaskType == subTask.TypeID {
						hasTaskType = true
						break
					}
				}
				
				if hasTaskType {
					availableUsers = append(availableUsers, user.ID)
				}
			}
			
			// Check if we have at least one available user for this sub-task
			if len(availableUsers) == 0 {
				canAssignTask = false
				break
			}
			
			// For auto-assignment, assign the first available user
			assignedUser := availableUsers[0]
			taskAssignments[subTask.ID] = []int32{assignedUser}
			
			// Mark this user as occupied for subsequent sub-tasks in this task
			occupiedUsers[assignedUser] = true
			
			// Reserve machine usage
			projectedMachineUsage[subTask.TypeID]++
		}
		
		// Only consider task assignable if it has unassigned sub-tasks that can be assigned
		if canAssignTask && hasAssignableSubTasks && len(taskAssignments) > 0 {
			// Task can be assigned - create the assignments in the database
			for subTaskID, userIDs := range taskAssignments {
				for _, userID := range userIDs {
					_, err := s.queries.CreateUserSubTaskAssignment(ctx, db.CreateUserSubTaskAssignmentParams{
						UserID:    userID,
						SubTaskID: subTaskID,
					})
					if err != nil {
						// Log error but continue with other assignments
						continue
					}
				}
			}
			
			// Add to available list
			availableTaskIDs = append(availableTaskIDs, task.ID)
			
			// Update machine usage projection for next tasks
			for k, v := range projectedMachineUsage {
				machineUsage[k] = v
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

// CleanDraftTaskAssignments cleans all assignments from draft tasks
func (s *taskService) CleanDraftTaskAssignments(ctx context.Context) error {
	// First, delete all existing assignments for draft tasks to ensure clean slate
	// Get all draft tasks first
	allDraftTasks, err := s.queries.GetDraftTasksSortedByPriority(ctx)
	if err != nil {
		return fmt.Errorf("failed to get draft tasks for cleanup: %w", err)
	}
	
	// For each draft task, get its assignments and delete them
	for _, draftTask := range allDraftTasks {
		assignments, err := s.queries.GetAssignmentsByTaskID(ctx, draftTask.ID)
		if err != nil {
			continue // Skip if can't get assignments
		}
		
		// Delete each assignment
		for _, assignment := range assignments {
			err := s.queries.DeleteUserSubTaskAssignment(ctx, assignment.ID)
			if err != nil {
				continue // Skip failed deletions but continue
			}
		}
	}
	
	return nil
} 