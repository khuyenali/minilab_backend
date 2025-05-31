package repository

import (
	"context"
	"database/sql"
	"mini-lab-api/internal/db"
)

// UserToTypeRepository defines the interface for user-to-type relationship operations
type UserToTypeRepository interface {
	AssignTaskTypes(ctx context.Context, userID int32, taskTypeIDs []int32) error
	GetUserTaskTypes(ctx context.Context, userID int32) ([]db.GetUserTaskTypesRow, error)
	RemoveUserTaskTypes(ctx context.Context, userID int32) error
}

// userToTypeRepository implements UserToTypeRepository using sqlc generated code
type userToTypeRepository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewUserToTypeRepository creates a new user to type repository
func NewUserToTypeRepository(database *sql.DB) UserToTypeRepository {
	return &userToTypeRepository{
		db:      database,
		queries: db.New(database),
	}
}

// AssignTaskTypes assigns task types to a user (replaces existing assignments)
func (r *userToTypeRepository) AssignTaskTypes(ctx context.Context, userID int32, taskTypeIDs []int32) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	queries := r.queries.WithTx(tx)
	
	// Remove all existing assignments for this user
	err = queries.RemoveAllUserTaskTypes(ctx, userID)
	if err != nil {
		return err
	}
	
	// Add new assignments
	for _, taskTypeID := range taskTypeIDs {
		err = queries.AddUserTaskType(ctx, db.AddUserTaskTypeParams{
			UserID: userID,
			TypeID: taskTypeID,
		})
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

// GetUserTaskTypes retrieves all task types assigned to a user
func (r *userToTypeRepository) GetUserTaskTypes(ctx context.Context, userID int32) ([]db.GetUserTaskTypesRow, error) {
	return r.queries.GetUserTaskTypes(ctx, userID)
}

// RemoveUserTaskTypes removes all task type assignments for a user
func (r *userToTypeRepository) RemoveUserTaskTypes(ctx context.Context, userID int32) error {
	return r.queries.RemoveAllUserTaskTypes(ctx, userID)
} 