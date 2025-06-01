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
	TaskType     *string  `json:"task_type" example:"3D Printing"`
	CreatedAt    time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// MachineBasic represents a machine without task_type info (for nested contexts)
type MachineBasic struct {
	ID           int32  `json:"id" example:"1"`
	MachineName  string `json:"name" example:"3D Printer Model X"`
	Quantity     int32  `json:"quantity" example:"2"`
	EstimateTime int32  `json:"estimate_time" example:"120"` // in minutes
}

// CreateMachineRequest represents the request to create a new machine
type CreateMachineRequest struct {
	Name         string `json:"name" binding:"required" example:"3D Printer Model X"`
	Quantity     int32  `json:"quantity" binding:"required,min=1" example:"2"`
	EstimateTime int32  `json:"estimate_time" binding:"required,min=1" example:"120"`
	TaskTypeID   *int32 `json:"task_type_id" example:"1"`
}

// UpdateMachineRequest represents the request to update a machine
type UpdateMachineRequest struct {
	Name         string `json:"name" example:"3D Printer Model Y"`
	Quantity     int32  `json:"quantity" binding:"min=1" example:"3"`
	EstimateTime int32  `json:"estimate_time" binding:"min=1" example:"90"`
	TaskTypeID   *int32 `json:"task_type_id" example:"2"`
	TaskTypeProvided bool `json:"-"` // Internal field to track if task_type_id was provided
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

// FromDBMachineBasic converts database model to basic domain model (without task_type)
func FromDBMachineBasic(dbMachine db.Machine) *MachineBasic {
	return &MachineBasic{
		ID:           dbMachine.ID,
		MachineName:  dbMachine.MachineName,
		Quantity:     dbMachine.Quantity,
		EstimateTime: dbMachine.EstimateTime,
	}
}

// FromGetMachinesByTaskTypeRows converts task type machine query results to basic domain models
func FromGetMachinesByTaskTypeRows(rows []db.Machine) []*MachineBasic {
	machines := make([]*MachineBasic, len(rows))
	for i, row := range rows {
		machines[i] = FromDBMachineBasic(row)
	}
	return machines
}

// FromGetMachineWithTaskTypeRow converts get machine with task type query result to domain model
func FromGetMachineWithTaskTypeRow(row db.GetMachineWithTaskTypeRow) *Machine {
	var createdAt, updatedAt time.Time
	var taskTypeID *int32
	var taskType *string
	
	if row.CreatedAt.Valid {
		createdAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		updatedAt = row.UpdatedAt.Time
	}
	if row.TypeID.Valid {
		taskTypeID = &row.TypeID.Int32
	}
	if row.TaskTypeName.Valid {
		taskType = &row.TaskTypeName.String
	}
	
	return &Machine{
		ID:           row.ID,
		MachineName:  row.MachineName,
		Quantity:     row.Quantity,
		EstimateTime: row.EstimateTime,
		TaskTypeID:   taskTypeID,
		TaskType:     taskType,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

// FromListMachinesWithTaskTypeRows converts list machines with task type query results to domain models
func FromListMachinesWithTaskTypeRows(rows []db.ListMachinesWithTaskTypeRow) []*Machine {
	machines := make([]*Machine, len(rows))
	for i, row := range rows {
		var createdAt, updatedAt time.Time
		var taskTypeID *int32
		var taskType *string
		
		if row.CreatedAt.Valid {
			createdAt = row.CreatedAt.Time
		}
		if row.UpdatedAt.Valid {
			updatedAt = row.UpdatedAt.Time
		}
		if row.TypeID.Valid {
			taskTypeID = &row.TypeID.Int32
		}
		if row.TaskTypeName.Valid {
			taskType = &row.TaskTypeName.String
		}
		
		machines[i] = &Machine{
			ID:           row.ID,
			MachineName:  row.MachineName,
			Quantity:     row.Quantity,
			EstimateTime: row.EstimateTime,
			TaskTypeID:   taskTypeID,
			TaskType:     taskType,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		}
	}
	return machines
} 