package repository

import (
	"context"
	"database/sql"
	"mini-lab-api/internal/db"
	"mini-lab-api/internal/models"
)

// RoleRepository defines the interface for role data operations
type RoleRepository interface {
	GetAll(ctx context.Context) ([]*models.Role, error)
	GetByID(ctx context.Context, id int32) (*models.Role, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)
}

// roleRepository implements RoleRepository using sqlc generated code
type roleRepository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(database *sql.DB) RoleRepository {
	return &roleRepository{
		db:      database,
		queries: db.New(database),
	}
}

// GetAll retrieves all roles
func (r *roleRepository) GetAll(ctx context.Context) ([]*models.Role, error) {
	dbRoles, err := r.queries.GetRoles(ctx)
	if err != nil {
		return nil, err
	}
	
	return models.FromDBRoles(dbRoles), nil
}

// GetByID retrieves a role by ID
func (r *roleRepository) GetByID(ctx context.Context, id int32) (*models.Role, error) {
	dbRole, err := r.queries.GetRoleByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	
	return models.FromDBRole(dbRole), nil
}

// GetByName retrieves a role by name
func (r *roleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	dbRole, err := r.queries.GetRoleByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	
	return models.FromDBRole(dbRole), nil
} 