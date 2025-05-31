package service

import (
	"context"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"
)

// RoleService defines the interface for role business operations
type RoleService interface {
	GetRoles(ctx context.Context) ([]*models.Role, error)
	GetRole(ctx context.Context, id int32) (*models.Role, error)
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
}

// roleService implements RoleService
type roleService struct {
	roleRepo repository.RoleRepository
}

// NewRoleService creates a new role service
func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}

// GetRoles retrieves all roles
func (s *roleService) GetRoles(ctx context.Context) ([]*models.Role, error) {
	return s.roleRepo.GetAll(ctx)
}

// GetRole retrieves a role by ID
func (s *roleService) GetRole(ctx context.Context, id int32) (*models.Role, error) {
	if id <= 0 {
		return nil, ErrInvalidRoleID
	}
	
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrRoleNotFound {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	
	return role, nil
}

// GetRoleByName retrieves a role by name
func (s *roleService) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	if name == "" {
		return nil, ErrInvalidRoleName
	}
	
	role, err := s.roleRepo.GetByName(ctx, name)
	if err != nil {
		if err == repository.ErrRoleNotFound {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	
	return role, nil
} 