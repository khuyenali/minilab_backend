package service

import (
	"context"
	"database/sql"
	"mini-lab-api/internal/db"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"
)

// AssignmentService defines the interface for assignment business operations
type AssignmentService interface {
	UpdateAssignmentToProcessing(ctx context.Context, assignmentID int32) (*models.UserSubTaskAssignment, error)
	FinishAssignment(ctx context.Context, assignmentID int32, report string) (*models.UserSubTaskAssignment, error)
	GetAssignment(ctx context.Context, assignmentID int32) (*models.UserSubTaskAssignment, error)
}

// assignmentService implements AssignmentService
type assignmentService struct {
	queries  *db.Queries
	taskRepo repository.TaskRepository
	db       *sql.DB
}

// NewAssignmentService creates a new assignment service
func NewAssignmentService(database *sql.DB, taskRepo repository.TaskRepository) AssignmentService {
	return &assignmentService{
		queries:  db.New(database),
		taskRepo: taskRepo,
		db:       database,
	}
}

// GetAssignment retrieves an assignment by ID
func (s *assignmentService) GetAssignment(ctx context.Context, assignmentID int32) (*models.UserSubTaskAssignment, error) {
	if assignmentID <= 0 {
		return nil, ErrInvalidAssignmentID
	}

	dbAssignment, err := s.queries.GetUserSubTaskAssignment(ctx, assignmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAssignmentNotFound
		}
		return nil, err
	}

	return models.FromDBUserSubTaskAssignment(dbAssignment), nil
}

// UpdateAssignmentToProcessing updates assignment status from pending to processing
func (s *assignmentService) UpdateAssignmentToProcessing(ctx context.Context, assignmentID int32) (*models.UserSubTaskAssignment, error) {
	if assignmentID <= 0 {
		return nil, ErrInvalidAssignmentID
	}

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	txQueries := s.queries.WithTx(tx)

	// Get current assignment to check status
	currentAssignment, err := txQueries.GetUserSubTaskAssignment(ctx, assignmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAssignmentNotFound
		}
		return nil, err
	}

	// Check if assignment is in pending status
	currentStatus := db.AssignmentStatusPending // default
	if currentAssignment.Status.Valid {
		currentStatus = currentAssignment.Status.AssignmentStatus
	}
	
	if currentStatus != db.AssignmentStatusPending {
		return nil, ErrAssignmentNotInPendingStatus
	}

	// Update assignment status to processing
	updatedAssignment, err := txQueries.UpdateAssignmentStatus(ctx, db.UpdateAssignmentStatusParams{
		ID:     assignmentID,
		Status: db.NullAssignmentStatus{AssignmentStatus: db.AssignmentStatusProcess, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// Get the task ID for this assignment
	subTask, err := txQueries.GetSubTask(ctx, currentAssignment.SubTaskID)
	if err != nil {
		return nil, err
	}

	// Update task status to processing if it's currently pending
	task, err := txQueries.GetTask(ctx, subTask.TaskID)
	if err != nil {
		return nil, err
	}

	if task.Status == db.TaskStatusPending {
		_, err = txQueries.UpdateTask(ctx, db.UpdateTaskParams{
			ID:        task.ID,
			TaskName:  task.TaskName,
			Status:    db.TaskStatusProcessing,
			Priority:  task.Priority,
			StartTime: task.StartTime,
			EndTime:   task.EndTime,
			Note:      task.Note,
		})
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return models.FromDBUserSubTaskAssignment(updatedAssignment), nil
}

// FinishAssignment updates assignment status from processing to finish with a report
func (s *assignmentService) FinishAssignment(ctx context.Context, assignmentID int32, report string) (*models.UserSubTaskAssignment, error) {
	if assignmentID <= 0 {
		return nil, ErrInvalidAssignmentID
	}

	if report == "" {
		return nil, ErrAssignmentReportRequired
	}

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	txQueries := s.queries.WithTx(tx)

	// Get current assignment to check status
	currentAssignment, err := txQueries.GetUserSubTaskAssignment(ctx, assignmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAssignmentNotFound
		}
		return nil, err
	}

	// Check if assignment is in processing status
	currentStatus := db.AssignmentStatusPending // default
	if currentAssignment.Status.Valid {
		currentStatus = currentAssignment.Status.AssignmentStatus
	}
	
	if currentStatus != db.AssignmentStatusProcess {
		return nil, ErrAssignmentNotInProcessingStatus
	}

	// Update assignment status to finish with report
	updatedAssignment, err := txQueries.UpdateAssignmentStatusAndReport(ctx, db.UpdateAssignmentStatusAndReportParams{
		ID:     assignmentID,
		Status: db.NullAssignmentStatus{AssignmentStatus: db.AssignmentStatusFinish, Valid: true},
		Report: sql.NullString{String: report, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// Get the task ID for this assignment
	subTask, err := txQueries.GetSubTask(ctx, currentAssignment.SubTaskID)
	if err != nil {
		return nil, err
	}

	// Check if all assignments for this task are finished
	allAssignments, err := txQueries.GetAssignmentsByTaskID(ctx, subTask.TaskID)
	if err != nil {
		return nil, err
	}

	allFinished := true
	for _, assignment := range allAssignments {
		status := db.AssignmentStatusPending // default
		if assignment.Status.Valid {
			status = assignment.Status.AssignmentStatus
		}
		if status != db.AssignmentStatusFinish {
			allFinished = false
			break
		}
	}

	// If all assignments are finished, update task status to finish
	if allFinished {
		task, err := txQueries.GetTask(ctx, subTask.TaskID)
		if err != nil {
			return nil, err
		}

		_, err = txQueries.UpdateTask(ctx, db.UpdateTaskParams{
			ID:        task.ID,
			TaskName:  task.TaskName,
			Status:    db.TaskStatusFinish,
			Priority:  task.Priority,
			StartTime: task.StartTime,
			EndTime:   task.EndTime,
			Note:      task.Note,
		})
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return models.FromDBUserSubTaskAssignment(updatedAssignment), nil
} 