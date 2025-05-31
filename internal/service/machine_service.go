package service

import (
	"context"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"
	"strings"
)

// MachineService defines the interface for machine business operations
type MachineService interface {
	GetMachines(ctx context.Context) ([]*models.Machine, error)
	GetMachine(ctx context.Context, id int32) (*models.Machine, error)
	CreateMachine(ctx context.Context, req models.CreateMachineRequest) (*models.Machine, error)
	UpdateMachine(ctx context.Context, id int32, req models.UpdateMachineRequest) (*models.Machine, error)
	DeleteMachine(ctx context.Context, id int32) error
}

// machineService implements MachineService
type machineService struct {
	machineRepo repository.MachineRepository
}

// NewMachineService creates a new machine service
func NewMachineService(machineRepo repository.MachineRepository) MachineService {
	return &machineService{
		machineRepo: machineRepo,
	}
}

// GetMachines retrieves all machines
func (s *machineService) GetMachines(ctx context.Context) ([]*models.Machine, error) {
	return s.machineRepo.List(ctx)
}

// GetMachine retrieves a machine by ID
func (s *machineService) GetMachine(ctx context.Context, id int32) (*models.Machine, error) {
	if id <= 0 {
		return nil, ErrInvalidMachineID
	}
	
	machine, err := s.machineRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrMachineNotFound {
			return nil, ErrMachineNotFound
		}
		return nil, err
	}
	
	return machine, nil
}

// CreateMachine creates a new machine
func (s *machineService) CreateMachine(ctx context.Context, req models.CreateMachineRequest) (*models.Machine, error) {
	// Validate input
	if err := s.validateCreateMachineRequest(req); err != nil {
		return nil, err
	}
	
	machine, err := s.machineRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return machine, nil
}

// UpdateMachine updates an existing machine
func (s *machineService) UpdateMachine(ctx context.Context, id int32, req models.UpdateMachineRequest) (*models.Machine, error) {
	if id <= 0 {
		return nil, ErrInvalidMachineID
	}
	
	// Validate input
	if err := s.validateUpdateMachineRequest(req); err != nil {
		return nil, err
	}
	
	machine, err := s.machineRepo.Update(ctx, id, req)
	if err != nil {
		if err == repository.ErrMachineNotFound {
			return nil, ErrMachineNotFound
		}
		return nil, err
	}
	
	return machine, nil
}

// DeleteMachine deletes a machine by ID
func (s *machineService) DeleteMachine(ctx context.Context, id int32) error {
	if id <= 0 {
		return ErrInvalidMachineID
	}
	
	err := s.machineRepo.Delete(ctx, id)
	if err != nil {
		if err == repository.ErrMachineNotFound {
			return ErrMachineNotFound
		}
		return err
	}
	
	return nil
}

// validateCreateMachineRequest validates the create machine request
func (s *machineService) validateCreateMachineRequest(req models.CreateMachineRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidMachineName
	}
	
	if len(req.Name) > 255 {
		return ErrMachineNameTooLong
	}
	
	if req.Quantity <= 0 {
		return ErrInvalidMachineQuantity
	}
	
	if req.EstimateTime <= 0 {
		return ErrInvalidMachineEstimateTime
	}
	
	return nil
}

// validateUpdateMachineRequest validates the update machine request
func (s *machineService) validateUpdateMachineRequest(req models.UpdateMachineRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidMachineName
	}
	
	if len(req.Name) > 255 {
		return ErrMachineNameTooLong
	}
	
	if req.Quantity <= 0 {
		return ErrInvalidMachineQuantity
	}
	
	if req.EstimateTime <= 0 {
		return ErrInvalidMachineEstimateTime
	}
	
	return nil
} 