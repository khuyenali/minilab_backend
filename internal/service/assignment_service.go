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
	CreateAssignment(ctx context.Context, req models.CreateAssignmentRequest) (*models.UserSubTaskAssignment, error)
	GetAssignment(ctx context.Context, assignmentID int32) (*models.UserSubTaskAssignment, error)
	DeleteAssignment(ctx context.Context, assignmentID int32) error
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

// CreateAssignment creates a new assignment
func (s *assignmentService) CreateAssignment(ctx context.Context, req models.CreateAssignmentRequest) (*models.UserSubTaskAssignment, error) {
	if req.SubTaskID <= 0 {
		return nil, ErrInvalidSubTaskForAssignment
	}
	if req.UserID <= 0 {
		return nil, ErrInvalidUserForAssignment
	}

	// Validate sub-task exists
	_, err := s.queries.GetSubTask(ctx, req.SubTaskID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidSubTaskForAssignment
		}
		return nil, err
	}

	// Validate user exists and is a member (role_id = 3)
	user, err := s.queries.GetUser(ctx, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidUserForAssignment
		}
		return nil, err
	}
	
	if user.RoleID != 3 {
		return nil, ErrInvalidUserForAssignment
	}

	// Check if assignment already exists for this user and sub-task
	existingAssignments, err := s.queries.GetAssignmentsBySubTaskID(ctx, req.SubTaskID)
	if err != nil {
		return nil, err
	}
	
	for _, assignment := range existingAssignments {
		if assignment.UserID == req.UserID {
			return nil, ErrAssignmentAlreadyExists
		}
	}

	// Create the assignment
	dbAssignment, err := s.queries.CreateUserSubTaskAssignment(ctx, db.CreateUserSubTaskAssignmentParams{
		UserID:    req.UserID,
		SubTaskID: req.SubTaskID,
	})
	if err != nil {
		return nil, err
	}

	return models.FromDBUserSubTaskAssignment(dbAssignment), nil
}

// DeleteAssignment deletes an assignment by ID
func (s *assignmentService) DeleteAssignment(ctx context.Context, assignmentID int32) error {
	if assignmentID <= 0 {
		return ErrInvalidAssignmentID
	}

	// Check if assignment exists
	_, err := s.queries.GetUserSubTaskAssignment(ctx, assignmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrAssignmentNotFound
		}
		return err
	}

	// Delete the assignment
	err = s.queries.DeleteUserSubTaskAssignment(ctx, assignmentID)
	if err != nil {
		return err
	}

	return nil
} 