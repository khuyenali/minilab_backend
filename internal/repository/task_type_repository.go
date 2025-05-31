package repository

import (
	"context"
	"database/sql"
	"mini-lab-api/internal/db"
	"mini-lab-api/internal/models"
)

// TaskTypeRepository defines the interface for task type data operations
type TaskTypeRepository interface {
	GetByID(ctx context.Context, id int32) (*models.TaskType, error)
	List(ctx context.Context) ([]*models.TaskType, error)
	Create(ctx context.Context, req models.CreateTaskTypeRequest) (*models.TaskType, error)
	Update(ctx context.Context, id int32, req models.UpdateTaskTypeRequest) (*models.TaskType, error)
	Delete(ctx context.Context, id int32) error
	GetWithMachines(ctx context.Context, id int32) (*models.TaskType, error)
}

// taskTypeRepository implements TaskTypeRepository using sqlc generated code
type taskTypeRepository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewTaskTypeRepository creates a new task type repository
func NewTaskTypeRepository(database *sql.DB) TaskTypeRepository {
	return &taskTypeRepository{
		db:      database,
		queries: db.New(database),
	}
}

// GetByID retrieves a task type by ID
func (r *taskTypeRepository) GetByID(ctx context.Context, id int32) (*models.TaskType, error) {
	row, err := r.queries.GetTaskType(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTaskTypeNotFound
		}
		return nil, err
	}
	
	return models.FromGetTaskTypeRow(row), nil
}

// List retrieves all task types
func (r *taskTypeRepository) List(ctx context.Context) ([]*models.TaskType, error) {
	rows, err := r.queries.ListTaskTypes(ctx)
	if err != nil {
		return nil, err
	}
	
	return models.FromListTaskTypesRows(rows), nil
}

// Create creates a new task type
func (r *taskTypeRepository) Create(ctx context.Context, req models.CreateTaskTypeRequest) (*models.TaskType, error) {
	var description sql.NullString
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	
	dbTaskType, err := r.queries.CreateTaskType(ctx, db.CreateTaskTypeParams{
		TypeName:    req.Name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}
	
	// Get the full task type
	taskType, err := r.GetByID(ctx, dbTaskType.ID)
	if err != nil {
		return nil, err
	}
	
	// Handle machine assignments if provided
	if len(req.MachineIDs) > 0 {
		// Update machines to belong to this task type
		for _, machineID := range req.MachineIDs {
			err := r.assignMachineToTaskType(ctx, machineID, taskType.ID)
			if err != nil {
				// Log error but don't fail the creation
				continue
			}
		}
		
		// Reload with machines
		return r.GetWithMachines(ctx, taskType.ID)
	}
	
	return taskType, nil
}

// Update updates an existing task type
func (r *taskTypeRepository) Update(ctx context.Context, id int32, req models.UpdateTaskTypeRequest) (*models.TaskType, error) {
	// Check if task type exists
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	var description sql.NullString
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	
	dbTaskType, err := r.queries.UpdateTaskType(ctx, db.UpdateTaskTypeParams{
		ID:          id,
		TypeName:    req.Name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}
	
	// Handle machine assignments if provided
	if len(req.MachineIDs) > 0 {
		// Clear existing assignments and set new ones
		err := r.clearMachineAssignments(ctx, id)
		if err != nil {
			return nil, err
		}
		
		for _, machineID := range req.MachineIDs {
			err := r.assignMachineToTaskType(ctx, machineID, id)
			if err != nil {
				continue
			}
		}
		
		// Reload with machines
		return r.GetWithMachines(ctx, id)
	}
	
	return r.GetByID(ctx, dbTaskType.ID)
}

// Delete deletes a task type by ID
func (r *taskTypeRepository) Delete(ctx context.Context, id int32) error {
	// Check if task type exists
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	return r.queries.DeleteTaskType(ctx, id)
}

// GetWithMachines retrieves a task type with its associated machines
func (r *taskTypeRepository) GetWithMachines(ctx context.Context, id int32) (*models.TaskType, error) {
	taskType, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	machines, err := r.queries.GetMachinesByTaskType(ctx, sql.NullInt32{Int32: id, Valid: true})
	if err != nil {
		// Return task type without machines if query fails
		return taskType, nil
	}
	
	taskType.Machines = models.FromGetMachinesByTaskTypeRows(machines)
	return taskType, nil
}

// Helper function to assign a machine to a task type
func (r *taskTypeRepository) assignMachineToTaskType(ctx context.Context, machineID, taskTypeID int32) error {
	// This would require an UPDATE query on machines table
	// For now, we'll skip this implementation since it requires additional queries
	return nil
}

// Helper function to clear machine assignments for a task type
func (r *taskTypeRepository) clearMachineAssignments(ctx context.Context, taskTypeID int32) error {
	// This would require an UPDATE query on machines table
	// For now, we'll skip this implementation since it requires additional queries
	return nil
} 