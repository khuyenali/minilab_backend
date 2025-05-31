package models

import (
	"time"
	"mini-lab-api/internal/db"
)

// Machine represents a machine in the system (domain model)
type Machine struct {
	ID           int32    `json:"id" example:"1"`
	MachineName  string   `json:"name" example:"3D Printer Model X"`
	Quantity     int32    `json:"quantity" example:"2"`
	EstimateTime int32    `json:"estimate_time" example:"120"` // in minutes
	TaskTypeID   *int32   `json:"-"` // Internal use only
	TaskType     *string  `json:"task_type,omitempty" example:"3D Printing"`
	CreatedAt    time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CreateMachineRequest represents the request to create a new machine
type CreateMachineRequest struct {
	Name         string `json:"name" binding:"required" example:"3D Printer Model X"`
	Quantity     int32  `json:"quantity" binding:"required,min=1" example:"2"`
	EstimateTime int32  `json:"estimate_time" binding:"required,min=1" example:"120"`
	TaskTypeID   *int32 `json:"task_type,omitempty" example:"1"`
}

// UpdateMachineRequest represents the request to update a machine
type UpdateMachineRequest struct {
	Name         string `json:"name" example:"3D Printer Model Y"`
	Quantity     int32  `json:"quantity" binding:"min=1" example:"3"`
	EstimateTime int32  `json:"estimate_time" binding:"min=1" example:"90"`
}

// FromDBMachine converts database model to domain model
func FromDBMachine(dbMachine db.Machine) *Machine {
	var createdAt, updatedAt time.Time
	var taskTypeID *int32
	
	if dbMachine.CreatedAt.Valid {
		createdAt = dbMachine.CreatedAt.Time
	}
	if dbMachine.UpdatedAt.Valid {
		updatedAt = dbMachine.UpdatedAt.Time
	}
	if dbMachine.TypeID.Valid {
		taskTypeID = &dbMachine.TypeID.Int32
	}
	
	return &Machine{
		ID:           dbMachine.ID,
		MachineName:  dbMachine.MachineName,
		Quantity:     dbMachine.Quantity,
		EstimateTime: dbMachine.EstimateTime,
		TaskTypeID:   taskTypeID,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

// FromDBMachines converts slice of database models to domain models
func FromDBMachines(dbMachines []db.Machine) []*Machine {
	machines := make([]*Machine, len(dbMachines))
	for i, dbMachine := range dbMachines {
		machines[i] = FromDBMachine(dbMachine)
	}
	return machines
}

// FromListMachinesRows converts list query results to domain models
func FromListMachinesRows(rows []db.Machine) []*Machine {
	machines := make([]*Machine, len(rows))
	for i, row := range rows {
		machines[i] = FromDBMachine(row)
	}
	return machines
}

// FromGetMachineRow converts get query result to domain model
func FromGetMachineRow(row db.Machine) *Machine {
	return FromDBMachine(row)
}

// FromGetMachinesByTaskTypeRows converts task type machine query results to domain models
func FromGetMachinesByTaskTypeRows(rows []db.Machine) []*Machine {
	machines := make([]*Machine, len(rows))
	for i, row := range rows {
		machines[i] = FromDBMachine(row)
	}
	return machines
} 