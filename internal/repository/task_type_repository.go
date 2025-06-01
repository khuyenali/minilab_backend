package repository

import (
	"context"
	"database/sql"
	"fmt"
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
	ValidateMachineIDs(ctx context.Context, machineIDs []int32) error
	AssignMachinesToTaskType(ctx context.Context, machineIDs []int32, taskTypeID int32) error
	ClearMachineAssignments(ctx context.Context, taskTypeID int32) error
}

// taskTypeRepository implements TaskTypeRepository using sqlc generated code
type taskTypeRepository struct {
	db        *sql.DB
	queries   *db.Queries
	machineRepo MachineRepository
}

// NewTaskTypeRepository creates a new task type repository
func NewTaskTypeRepository(database *sql.DB, machineRepo MachineRepository) TaskTypeRepository {
	return &taskTypeRepository{
		db:          database,
		queries:     db.New(database),
		machineRepo: machineRepo,
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
	
	taskTypes := models.FromListTaskTypesRows(rows)
	
	// Get machines for each task type
	for _, taskType := range taskTypes {
		machines, err := r.queries.GetMachinesByTaskType(ctx, sql.NullInt32{Int32: taskType.ID, Valid: true})
		if err != nil {
			// Continue if we can't get machines for this task type
			continue
		}
		taskType.Machines = models.FromGetMachinesByTaskTypeRows(machines)
	}
	
	return taskTypes, nil
}

// Create creates a new task type
func (r *taskTypeRepository) Create(ctx context.Context, req models.CreateTaskTypeRequest) (*models.TaskType, error) {
	var description sql.NullString
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	
	// Validate machine IDs if provided
	if len(req.MachineIDs) > 0 {
		if err := r.ValidateMachineIDs(ctx, req.MachineIDs); err != nil {
			return nil, err
		}
	}
	
	dbTaskType, err := r.queries.CreateTaskType(ctx, db.CreateTaskTypeParams{
		TypeName:    req.Name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}
	
	// Assign machines if provided
	if len(req.MachineIDs) > 0 {
		if err := r.AssignMachinesToTaskType(ctx, req.MachineIDs, dbTaskType.ID); err != nil {
			// If machine assignment fails, we should probably rollback the task type creation
			// For now, we'll continue but this should be in a transaction
			return nil, err
		}
	}
	
	// Get the full task type with machines
	return r.GetWithMachines(ctx, dbTaskType.ID)
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
	
	// Validate machine IDs if provided
	if len(req.MachineIDs) > 0 {
		if err := r.ValidateMachineIDs(ctx, req.MachineIDs); err != nil {
			return nil, err
		}
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
		// Clear existing assignments first
		if err := r.ClearMachineAssignments(ctx, id); err != nil {
			return nil, err
		}
		
		// Assign new machines
		if err := r.AssignMachinesToTaskType(ctx, req.MachineIDs, id); err != nil {
			return nil, err
		}
	}
	
	// Get the full task type with machines
	return r.GetWithMachines(ctx, dbTaskType.ID)
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

// ValidateMachineIDs checks if the provided machine IDs are valid
func (r *taskTypeRepository) ValidateMachineIDs(ctx context.Context, machineIDs []int32) error {
	if len(machineIDs) == 0 {
		return nil
	}
	
	machines, err := r.machineRepo.GetByIDs(ctx, machineIDs)
	if err != nil {
		return err
	}
	
	// Check if all requested machines were found
	if len(machines) != len(machineIDs) {
		// Find which machine IDs don't exist
		foundIDs := make(map[int32]bool)
		for _, machine := range machines {
			foundIDs[machine.ID] = true
		}
		
		var missingIDs []int32
		for _, id := range machineIDs {
			if !foundIDs[id] {
				missingIDs = append(missingIDs, id)
			}
		}
		
		return fmt.Errorf("%w: IDs %v not found", ErrInvalidMachineIDs, missingIDs)
	}
	
	return nil
}

// AssignMachinesToTaskType assigns machines to a task type
func (r *taskTypeRepository) AssignMachinesToTaskType(ctx context.Context, machineIDs []int32, taskTypeID int32) error {
	for _, machineID := range machineIDs {
		if err := r.machineRepo.UpdateTaskType(ctx, machineID, taskTypeID); err != nil {
			return fmt.Errorf("failed to assign machine %d to task type %d: %w", machineID, taskTypeID, err)
		}
	}
	return nil
}

// ClearMachineAssignments clears machine assignments for a task type
func (r *taskTypeRepository) ClearMachineAssignments(ctx context.Context, taskTypeID int32) error {
	return r.machineRepo.ClearTaskTypeAssignments(ctx, taskTypeID)
} 