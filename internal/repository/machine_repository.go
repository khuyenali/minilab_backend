package repository

import (
	"context"
	"database/sql"
	"mini-lab-api/internal/db"
	"mini-lab-api/internal/models"
)

// MachineRepository defines the interface for machine data operations
type MachineRepository interface {
	GetByID(ctx context.Context, id int32) (*models.Machine, error)
	List(ctx context.Context) ([]*models.Machine, error)
	Create(ctx context.Context, req models.CreateMachineRequest) (*models.Machine, error)
	Update(ctx context.Context, id int32, req models.UpdateMachineRequest) (*models.Machine, error)
	Delete(ctx context.Context, id int32) error
}

// machineRepository implements MachineRepository using sqlc generated code
type machineRepository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewMachineRepository creates a new machine repository
func NewMachineRepository(database *sql.DB) MachineRepository {
	return &machineRepository{
		db:      database,
		queries: db.New(database),
	}
}

// GetByID retrieves a machine by ID
func (r *machineRepository) GetByID(ctx context.Context, id int32) (*models.Machine, error) {
	row, err := r.queries.GetMachine(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrMachineNotFound
		}
		return nil, err
	}
	
	return models.FromGetMachineRow(row), nil
}

// List retrieves all machines
func (r *machineRepository) List(ctx context.Context) ([]*models.Machine, error) {
	rows, err := r.queries.ListMachines(ctx)
	if err != nil {
		return nil, err
	}
	
	return models.FromListMachinesRows(rows), nil
}

// Create creates a new machine
func (r *machineRepository) Create(ctx context.Context, req models.CreateMachineRequest) (*models.Machine, error) {
	var taskTypeID sql.NullInt32
	if req.TaskTypeID != nil {
		taskTypeID = sql.NullInt32{Int32: *req.TaskTypeID, Valid: true}
	}
	
	dbMachine, err := r.queries.CreateMachine(ctx, db.CreateMachineParams{
		MachineName:  req.Name,
		Quantity:     req.Quantity,
		EstimateTime: req.EstimateTime,
		TypeID:       taskTypeID,
	})
	if err != nil {
		return nil, err
	}
	
	return r.GetByID(ctx, dbMachine.ID)
}

// Update updates an existing machine
func (r *machineRepository) Update(ctx context.Context, id int32, req models.UpdateMachineRequest) (*models.Machine, error) {
	// Check if machine exists
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	dbMachine, err := r.queries.UpdateMachine(ctx, db.UpdateMachineParams{
		ID:           id,
		MachineName:  req.Name,
		Quantity:     req.Quantity,
		EstimateTime: req.EstimateTime,
	})
	if err != nil {
		return nil, err
	}
	
	return r.GetByID(ctx, dbMachine.ID)
}

// Delete deletes a machine by ID
func (r *machineRepository) Delete(ctx context.Context, id int32) error {
	// Check if machine exists
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	return r.queries.DeleteMachine(ctx, id)
} 