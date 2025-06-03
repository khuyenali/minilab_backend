package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"mini-lab-api/internal/db"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//
// ------------------------------ Mock -- UserRepository ------------------------------
//
type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) GetByID(ctx context.Context, id int32) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, filters models.UserFilters) ([]*models.User, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByRole(ctx context.Context, roleName string) ([]*models.User, error) {
	args := m.Called(ctx, roleName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByRoleID(ctx context.Context, roleID int32) ([]*models.User, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int32) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockUserRepository) GetByIDs(ctx context.Context, ids []int32) ([]*models.User, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) AddUserToTaskType(ctx context.Context, userID, taskTypeID int32) error {
	return m.Called(ctx, userID, taskTypeID).Error(0)
}

func (m *MockUserRepository) RemoveUserFromTaskType(ctx context.Context, userID, taskTypeID int32) error {
	return m.Called(ctx, userID, taskTypeID).Error(0)
}

func (m *MockUserRepository) ClearTaskTypeAssignments(ctx context.Context, taskTypeID int32) error {
	return m.Called(ctx, taskTypeID).Error(0)
}

//
// ------------------------------ Mock -- UserToTypeRepository ------------------------------
//
type MockUserToTypeRepository struct{ mock.Mock }

func (m *MockUserToTypeRepository) AssignTaskTypes(ctx context.Context, userID int32, taskTypeIDs []int32) error {
	return m.Called(ctx, userID, taskTypeIDs).Error(0)
}

func (m *MockUserToTypeRepository) GetUserTaskTypes(ctx context.Context, userID int32) ([]db.GetUserTaskTypesRow, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.GetUserTaskTypesRow), args.Error(1)
}

func (m *MockUserToTypeRepository) RemoveUserTaskTypes(ctx context.Context, userID int32) error {
	return m.Called(ctx, userID).Error(0)
}

//
// ------------------------------ Mock -- TaskTypeRepository ------------------------------
//
type MockTaskTypeRepository struct{ mock.Mock }

func (m *MockTaskTypeRepository) GetByID(ctx context.Context, id int32) (*models.TaskType, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskType), args.Error(1)
}

func (m *MockTaskTypeRepository) List(ctx context.Context) ([]*models.TaskType, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.TaskType), args.Error(1)
}

func (m *MockTaskTypeRepository) Create(ctx context.Context, req models.CreateTaskTypeRequest) (*models.TaskType, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskType), args.Error(1)
}

func (m *MockTaskTypeRepository) Update(ctx context.Context, id int32, req models.UpdateTaskTypeRequest) (*models.TaskType, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskType), args.Error(1)
}

func (m *MockTaskTypeRepository) Delete(ctx context.Context, id int32) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockTaskTypeRepository) GetWithMachines(ctx context.Context, id int32) (*models.TaskType, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskType), args.Error(1)
}

func (m *MockTaskTypeRepository) ValidateMachineIDs(ctx context.Context, machineIDs []int32) error {
	return m.Called(ctx, machineIDs).Error(0)
}

func (m *MockTaskTypeRepository) AssignMachinesToTaskType(ctx context.Context, machineIDs []int32, taskTypeID int32) error {
	return m.Called(ctx, machineIDs, taskTypeID).Error(0)
}

func (m *MockTaskTypeRepository) ClearMachineAssignments(ctx context.Context, taskTypeID int32) error {
	return m.Called(ctx, taskTypeID).Error(0)
}

func (m *MockTaskTypeRepository) ValidateUserIDs(ctx context.Context, userIDs []int32) error {
	return m.Called(ctx, userIDs).Error(0)
}

func (m *MockTaskTypeRepository) AssignUsersToTaskType(ctx context.Context, userIDs []int32, taskTypeID int32) error {
	return m.Called(ctx, userIDs, taskTypeID).Error(0)
}

func (m *MockTaskTypeRepository) ClearUserAssignments(ctx context.Context, taskTypeID int32) error {
	return m.Called(ctx, taskTypeID).Error(0)
}

//
// ------------------------------  UserService  Tests  ------------------------------
//

func TestUserService_GetUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockUserToTypeRepo := new(MockUserToTypeRepository)
	mockTaskTypeRepo := new(MockTaskTypeRepository)
	svc := NewUserService(mockUserRepo, mockUserToTypeRepo, mockTaskTypeRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		want := &models.User{ID: 1, Name: "Test User", Email: "test@example.com", Role: "member"}
		mockUserRepo.On("GetByID", ctx, int32(1)).Return(want, nil).Once()

		got, err := svc.GetUser(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("invalid-id", func(t *testing.T) {
		got, err := svc.GetUser(ctx, 0)
		assert.ErrorIs(t, err, ErrInvalidUserID)
		assert.Nil(t, got)
	})

	t.Run("not-found", func(t *testing.T) {
		mockUserRepo.On("GetByID", ctx, int32(2)).Return(nil, repository.ErrUserNotFound).Once()
		got, err := svc.GetUser(ctx, 2)
		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Nil(t, got)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("repo-error", func(t *testing.T) {
		repoErr := errors.New("db boom")
		mockUserRepo.On("GetByID", ctx, int32(3)).Return(nil, repoErr).Once()
		got, err := svc.GetUser(ctx, 3)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, got)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUsers(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	svc := NewUserService(mockUserRepo, new(MockUserToTypeRepository), new(MockTaskTypeRepository))
	ctx := context.Background()

	t.Run("no-filters", func(t *testing.T) {
		want := []*models.User{{ID: 1, Name: "u1"}, {ID: 2, Name: "u2"}}
		expectedFilters := models.UserFilters{Limit: 50}
		mockUserRepo.On("List", ctx, expectedFilters).Return(want, nil).Once()

		got, err := svc.GetUsers(ctx, models.UserFilters{})
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("limit-adjustment", func(t *testing.T) {
		// over-limit → 100
		mockUserRepo.On("List", ctx, models.UserFilters{Limit: 100}).Return([]*models.User{}, nil).Once()
		_, _ = svc.GetUsers(ctx, models.UserFilters{Limit: 200})

		// negative → 50
		mockUserRepo.On("List", ctx, models.UserFilters{Limit: 50}).Return([]*models.User{}, nil).Once()
		_, _ = svc.GetUsers(ctx, models.UserFilters{Limit: -10})

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("repo-error", func(t *testing.T) {
		repoErr := errors.New("list err")
		mockUserRepo.On("List", ctx, models.UserFilters{Limit: 50}).Return(nil, repoErr).Once()

		got, err := svc.GetUsers(ctx, models.UserFilters{})
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, got)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockUserToTypeRepo := new(MockUserToTypeRepository)
	mockTaskTypeRepo := new(MockTaskTypeRepository)
	svc := NewUserService(mockUserRepo, mockUserToTypeRepo, mockTaskTypeRepo)

	memberRoleID, adminRoleID := int32(3), int32(1)

	t.Run("basic-user", func(t *testing.T) {
		req := models.CreateUserRequest{Name: "New User", Email: "new@example.com"}
		want := &models.User{ID: 1, Name: "New User", Email: "new@example.com", RoleID: memberRoleID, Role: "member"}
		mockUserRepo.On("Create", ctx, req).Return(want, nil).Once()

		got, err := svc.CreateUser(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("specific-role-admin", func(t *testing.T) {
		req := models.CreateUserRequest{Name: "Admin", Email: "admin@example.com", RoleID: &adminRoleID}
		want := &models.User{ID: 2, Name: "Admin", Email: "admin@example.com", RoleID: adminRoleID, Role: "admin"}
		mockUserRepo.On("Create", ctx, req).Return(want, nil).Once()

		got, err := svc.CreateUser(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("member-with-taskTypes", func(t *testing.T) {
		taskIDs := []int32{10, 11}
		req := models.CreateUserRequest{Name: "Tasky", Email: "tasky@example.com", RoleID: &memberRoleID, TaskTypeIDs: taskIDs}
		created := &models.User{ID: 3, Name: "Tasky", Email: "tasky@example.com", RoleID: memberRoleID, Role: "member"}
		withTT := &models.User{ID: 3, Name: "Tasky", Email: "tasky@example.com", RoleID: memberRoleID, Role: "member",
			TaskTypes: []*models.TaskTypeMinimal{{ID: 10}, {ID: 11}}}

		mockTaskTypeRepo.On("GetByID", ctx, int32(10)).Return(&models.TaskType{ID: 10, TypeName: "TT10"}, nil).Once()
		mockTaskTypeRepo.On("GetByID", ctx, int32(11)).Return(&models.TaskType{ID: 11, TypeName: "TT11"}, nil).Once()
		mockUserRepo.On("Create", ctx, req).Return(created, nil).Once()
		mockUserToTypeRepo.On("AssignTaskTypes", ctx, created.ID, taskIDs).Return(nil).Once()
		mockUserRepo.On("GetByID", ctx, created.ID).Return(withTT, nil).Once()

		got, err := svc.CreateUser(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, withTT, got)
		mockUserRepo.AssertExpectations(t)
		mockUserToTypeRepo.AssertExpectations(t)
		mockTaskTypeRepo.AssertExpectations(t)
	})

	t.Run("invalid-task-type-ids", func(t *testing.T) {
		req := models.CreateUserRequest{Name: "BadTT", Email: "badtt@example.com", RoleID: &memberRoleID, TaskTypeIDs: []int32{10, 999}}
		mockTaskTypeRepo.On("GetByID", ctx, int32(10)).Return(&models.TaskType{ID: 10}, nil).Once()
		mockTaskTypeRepo.On("GetByID", ctx, int32(999)).Return(nil, repository.ErrTaskTypeNotFound).Once()

		got, err := svc.CreateUser(ctx, req)
		assert.Nil(t, got)
		assert.EqualError(t, err, fmt.Sprintf("invalid task type IDs: %v", []int32{999}))
		mockTaskTypeRepo.AssertExpectations(t)
	})

	t.Run("email-exists", func(t *testing.T) {
		req := models.CreateUserRequest{Name: "Dup", Email: "dup@example.com"}
		mockUserRepo.On("Create", ctx, req).Return(nil, repository.ErrUserEmailExists).Once()

		got, err := svc.CreateUser(ctx, req)
		assert.ErrorIs(t, err, ErrUserEmailExists)
		assert.Nil(t, got)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	svc := NewUserService(mockUserRepo, new(MockUserToTypeRepository), new(MockTaskTypeRepository))
	userID := int32(1)

	t.Run("success", func(t *testing.T) {
		req := models.UpdateUserRequest{Name: "Upd", Email: "upd@example.com"}
		want := &models.User{ID: userID, Name: "Upd", Email: "upd@example.com"}
		mockUserRepo.On("Update", ctx, userID, req).Return(want, nil).Once()

		got, err := svc.UpdateUser(ctx, userID, req)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("invalid-id", func(t *testing.T) {
		got, err := svc.UpdateUser(ctx, 0, models.UpdateUserRequest{Name: "x"})
		assert.ErrorIs(t, err, ErrInvalidUserID)
		assert.Nil(t, got)
	})

	t.Run("repo-not-found", func(t *testing.T) {
		req := models.UpdateUserRequest{Name: "x"}
		mockUserRepo.On("Update", ctx, userID, req).Return(nil, repository.ErrUserNotFound).Once()

		got, err := svc.UpdateUser(ctx, userID, req)
		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Nil(t, got)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	svc := NewUserService(mockUserRepo, new(MockUserToTypeRepository), new(MockTaskTypeRepository))
	userID := int32(1)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Delete", ctx, userID).Return(nil).Once()
		assert.NoError(t, svc.DeleteUser(ctx, userID))
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("repo-not-found", func(t *testing.T) {
		mockUserRepo.On("Delete", ctx, userID).Return(repository.ErrUserNotFound).Once()
		err := svc.DeleteUser(ctx, userID)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestUserService_AssignTaskTypes(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockUserToTypeRepo := new(MockUserToTypeRepository)
	mockTaskTypeRepo := new(MockTaskTypeRepository)
	svc := NewUserService(mockUserRepo, mockUserToTypeRepo, mockTaskTypeRepo)

	memberID, adminID := int32(1), int32(2)
	memberRoleID, adminRoleID := int32(3), int32(1)
	taskIDs := []int32{10, 11}

	t.Run("success", func(t *testing.T) {
		user := &models.User{ID: memberID, RoleID: memberRoleID}
		mockUserRepo.On("GetByID", ctx, memberID).Return(user, nil).Once()
		mockTaskTypeRepo.On("GetByID", ctx, int32(10)).Return(&models.TaskType{ID: 10}, nil).Once()
		mockTaskTypeRepo.On("GetByID", ctx, int32(11)).Return(&models.TaskType{ID: 11}, nil).Once()
		mockUserToTypeRepo.On("AssignTaskTypes", ctx, memberID, taskIDs).Return(nil).Once()

		err := svc.AssignTaskTypes(ctx, memberID, models.UserTaskAssignmentRequest{TaskTypeIDs: taskIDs})
		assert.NoError(t, err)
	})

	t.Run("invalid-ids", func(t *testing.T) {
		user := &models.User{ID: memberID, RoleID: memberRoleID}
		mockUserRepo.On("GetByID", ctx, memberID).Return(user, nil).Once()
		mockTaskTypeRepo.On("GetByID", ctx, int32(10)).Return(nil, repository.ErrTaskTypeNotFound).Once()

		err := svc.AssignTaskTypes(ctx, memberID, models.UserTaskAssignmentRequest{TaskTypeIDs: []int32{10}})
		assert.EqualError(t, err, fmt.Sprintf("invalid task type IDs: %v", []int32{10}))
	})

	t.Run("non-member", func(t *testing.T) {
		user := &models.User{ID: adminID, RoleID: adminRoleID}
		mockUserRepo.On("GetByID", ctx, adminID).Return(user, nil).Once()

		err := svc.AssignTaskTypes(ctx, adminID, models.UserTaskAssignmentRequest{TaskTypeIDs: taskIDs})
		assert.ErrorIs(t, err, ErrInvalidUserRole)
	})
}
