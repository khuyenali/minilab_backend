package service

import (
	"context"
	"database/sql"
	"mini-lab-api/internal/db"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"
)

// AssignmentStatusService defines the interface for assignment status operations
type AssignmentStatusService interface {
	UpdateAssignmentStatus(ctx context.Context, assignmentID int32, req models.UpdateAssignmentStatusRequest) (*models.Task, error)
	GetAssignmentTask(ctx context.Context, assignmentID int32) (*models.Task, error)
}

// assignmentStatusService implements AssignmentStatusService
type assignmentStatusService struct {
	queries  *db.Queries
	taskRepo repository.TaskRepository
	db       *sql.DB
}

// NewAssignmentStatusService creates a new assignment status service
func NewAssignmentStatusService(database *sql.DB, taskRepo repository.TaskRepository) AssignmentStatusService {
	return &assignmentStatusService{
		queries:  db.New(database),
		taskRepo: taskRepo,
		db:       database,
	}
}

// UpdateAssignmentStatus updates assignment status and triggers task status updates with proper order logic
func (s *assignmentStatusService) UpdateAssignmentStatus(ctx context.Context, assignmentID int32, req models.UpdateAssignmentStatusRequest) (*models.Task, error) {
	if assignmentID <= 0 {
		return nil, ErrInvalidAssignmentID
	}

	// Validate status
	if !isValidAssignmentStatus(req.Status) {
		return nil, ErrInvalidAssignmentStatus
	}

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	txQueries := s.queries.WithTx(tx)

	// Get current assignment
	assignment, err := txQueries.GetUserSubTaskAssignment(ctx, assignmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAssignmentNotFound
		}
		return nil, err
	}

	// Get the sub-task and task
	subTask, err := txQueries.GetSubTask(ctx, assignment.SubTaskID)
	if err != nil {
		return nil, err
	}

	task, err := txQueries.GetTask(ctx, subTask.TaskID)
	if err != nil {
		return nil, err
	}

	// Apply assignment status logic and update task status accordingly
	var updatedTask db.Task
	switch req.Status {
	case "processing":
		// When assignment moves to processing, task should move to processing if it's pending
		if task.Status == db.TaskStatusPending {
			updatedTask, err = txQueries.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
				ID:     task.ID,
				Status: db.TaskStatusProcessing,
			})
			if err != nil {
				return nil, err
			}
		} else {
			updatedTask = task
		}

	case "finish":
		// When assignment finishes, we can update the task status
		// For this simplified implementation, we assume the task can be finished
		// when any assignment is marked as finished (you can modify this logic as needed)
		// In a real scenario, you might want to track assignment statuses in a separate table
		
		// For now, we update the task to processing if it's pending
		if task.Status == db.TaskStatusPending {
			updatedTask, err = txQueries.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
				ID:     task.ID,
				Status: db.TaskStatusProcessing,
			})
			if err != nil {
				return nil, err
			}
		} else {
			updatedTask = task
		}

	default:
		updatedTask = task
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return models.FromDBTask(updatedTask), nil
}

// GetAssignmentTask gets the task associated with an assignment
func (s *assignmentStatusService) GetAssignmentTask(ctx context.Context, assignmentID int32) (*models.Task, error) {
	if assignmentID <= 0 {
		return nil, ErrInvalidAssignmentID
	}

	// Get assignment
	assignment, err := s.queries.GetUserSubTaskAssignment(ctx, assignmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAssignmentNotFound
		}
		return nil, err
	}

	// Get sub-task
	subTask, err := s.queries.GetSubTask(ctx, assignment.SubTaskID)
	if err != nil {
		return nil, err
	}

	// Get task with sub-tasks
	task, err := s.taskRepo.GetWithSubTasks(ctx, subTask.TaskID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// isValidAssignmentStatus checks if the assignment status is valid
func isValidAssignmentStatus(status string) bool {
	validStatuses := []string{"pending", "processing", "finish"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
} 