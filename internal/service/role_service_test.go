package service

import (
	"context"
	"errors"
	"testing"

	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
   ----------------------------------------------------------------
   Mock infrastructure
   ----------------------------------------------------------------
*/

// MockRoleRepository 是 repository.RoleRepository 的 testify mock 實作
type MockRoleRepository struct {
	mock.Mock
}

// 確保 MockRoleRepository 仍符合 interface
var _ repository.RoleRepository = (*MockRoleRepository)(nil)

/* 實際在 roleService 中會被呼叫的方法 */
func (m *MockRoleRepository) GetAll(ctx context.Context) ([]*models.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id int32) (*models.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

/* 下面這些方法目前 service 用不到；先做空實作，未來 interface 變動仍能編譯 */
func (m *MockRoleRepository) Create(ctx context.Context, role *models.Role) (*models.Role, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) Update(ctx context.Context, role *models.Role) error {
	return m.Called(ctx, role).Error(0)
}

func (m *MockRoleRepository) Delete(ctx context.Context, id int32) error {
	return m.Called(ctx, id).Error(0)
}

/*
   ----------------------------------------------------------------
   Tests
   ----------------------------------------------------------------
*/

func TestRoleService_GetRoles(t *testing.T) {
	ctx := context.Background()
	expected := []*models.Role{
		{ID: 1, RoleName: "admin"},
		{ID: 2, RoleName: "user"},
	}

	mockRepo := new(MockRoleRepository)
	mockRepo.On("GetAll", ctx).Return(expected, nil)

	svc := NewRoleService(mockRepo)

	got, err := svc.GetRoles(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	mockRepo.AssertExpectations(t)
}

func TestRoleService_GetRoles_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repoErr := errors.New("db error")

	mockRepo := new(MockRoleRepository)
	mockRepo.On("GetAll", ctx).Return(nil, repoErr)

	svc := NewRoleService(mockRepo)

	got, err := svc.GetRoles(ctx)
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, got)

	mockRepo.AssertExpectations(t)
}

func TestRoleService_GetRole(t *testing.T) {
	ctx := context.Background()
	roleID := int32(1)
	expected := &models.Role{ID: roleID, RoleName: "admin"}

	mockRepo := new(MockRoleRepository)
	mockRepo.On("GetByID", ctx, roleID).Return(expected, nil)

	svc := NewRoleService(mockRepo)

	got, err := svc.GetRole(ctx, roleID)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	mockRepo.AssertExpectations(t)
}

func TestRoleService_GetRole_InvalidID(t *testing.T) {
	ctx := context.Background()
	svc := NewRoleService(new(MockRoleRepository)) // repo 不會被呼叫

	for _, id := range []int32{0, -10} {
		got, err := svc.GetRole(ctx, id)
		assert.ErrorIs(t, err, ErrInvalidRoleID)
		assert.Nil(t, got)
	}
}

func TestRoleService_GetRole_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepository)
	mockRepo.On("GetByID", ctx, int32(99)).Return(nil, repository.ErrRoleNotFound)

	svc := NewRoleService(mockRepo)

	got, err := svc.GetRole(ctx, 99)
	assert.ErrorIs(t, err, ErrRoleNotFound)
	assert.Nil(t, got)

	mockRepo.AssertExpectations(t)
}

func TestRoleService_GetRole_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repoErr := errors.New("network glitch")
	mockRepo := new(MockRoleRepository)
	mockRepo.On("GetByID", ctx, int32(7)).Return(nil, repoErr)

	svc := NewRoleService(mockRepo)

	got, err := svc.GetRole(ctx, 7)
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, got)

	mockRepo.AssertExpectations(t)
}

func TestRoleService_GetRoleByName(t *testing.T) {
	ctx := context.Background()
	roleName := "admin"
	expected := &models.Role{ID: 1, RoleName: roleName}

	mockRepo := new(MockRoleRepository)
	mockRepo.On("GetByName", ctx, roleName).Return(expected, nil)

	svc := NewRoleService(mockRepo)

	got, err := svc.GetRoleByName(ctx, roleName)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	mockRepo.AssertExpectations(t)
}

func TestRoleService_GetRoleByName_InvalidName(t *testing.T) {
	ctx := context.Background()
	svc := NewRoleService(new(MockRoleRepository)) // repo 不會被呼叫

	got, err := svc.GetRoleByName(ctx, "")
	assert.ErrorIs(t, err, ErrInvalidRoleName)
	assert.Nil(t, got)
}

func TestRoleService_GetRoleByName_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepository)
	mockRepo.On("GetByName", ctx, "ghost").Return(nil, repository.ErrRoleNotFound)

	svc := NewRoleService(mockRepo)

	got, err := svc.GetRoleByName(ctx, "ghost")
	assert.ErrorIs(t, err, ErrRoleNotFound)
	assert.Nil(t, got)

	mockRepo.AssertExpectations(t)
}

func TestRoleService_GetRoleByName_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repoErr := errors.New("timeout")

	mockRepo := new(MockRoleRepository)
	mockRepo.On("GetByName", ctx, "admin").Return(nil, repoErr)

	svc := NewRoleService(mockRepo)

	got, err := svc.GetRoleByName(ctx, "admin")
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, got)

	mockRepo.AssertExpectations(t)
}
