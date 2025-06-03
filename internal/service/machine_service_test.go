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

/* ------------------------------------------------------------------
   Mock MachineRepository
-------------------------------------------------------------------*/
type MockMachineRepository struct{ mock.Mock }

func (m *MockMachineRepository) GetByID(ctx context.Context, id int32) (*models.Machine, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Machine), args.Error(1)
}

func (m *MockMachineRepository) List(ctx context.Context) ([]*models.Machine, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Machine), args.Error(1)
}

func (m *MockMachineRepository) Create(ctx context.Context, req models.CreateMachineRequest) (*models.Machine, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Machine), args.Error(1)
}

func (m *MockMachineRepository) Update(ctx context.Context, id int32, req models.UpdateMachineRequest) (*models.Machine, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Machine), args.Error(1)
}

func (m *MockMachineRepository) Delete(ctx context.Context, id int32) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockMachineRepository) GetByIDs(ctx context.Context, ids []int32) ([]*models.Machine, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Machine), args.Error(1)
}

func (m *MockMachineRepository) UpdateTaskType(ctx context.Context, machineID, taskTypeID int32) error {
	return m.Called(ctx, machineID, taskTypeID).Error(0)
}

func (m *MockMachineRepository) ClearTaskTypeAssignments(ctx context.Context, taskTypeID int32) error {
	return m.Called(ctx, taskTypeID).Error(0)
}

/* ------------------------------------------------------------------
   Tests: GetMachines
-------------------------------------------------------------------*/
func TestMachineService_GetMachines(t *testing.T) {
	mockMachineRepo := new(MockMachineRepository)
	mockTaskTypeRepo := new(MockTaskTypeRepository)
	svc := NewMachineService(mockMachineRepo, mockTaskTypeRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		want := []*models.Machine{{ID: 1, MachineName: "A"}, {ID: 2, MachineName: "B"}}
		mockMachineRepo.On("List", ctx).Return(want, nil).Once()

		got, err := svc.GetMachines(ctx)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		mockMachineRepo.AssertExpectations(t)
	})

	t.Run("repo-error", func(t *testing.T) {
		repoErr := errors.New("db")
		mockMachineRepo.On("List", ctx).Return(nil, repoErr).Once()

		got, err := svc.GetMachines(ctx)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, got)
	})
}

/* ------------------------------------------------------------------
   Tests: GetMachine
-------------------------------------------------------------------*/
func TestMachineService_GetMachine(t *testing.T) {
	ctx := context.Background()
	mockMachineRepo := new(MockMachineRepository)
	svc := NewMachineService(mockMachineRepo, new(MockTaskTypeRepository))
	id := int32(1)

	t.Run("success", func(t *testing.T) {
		want := &models.Machine{ID: id, MachineName: "P"}
		mockMachineRepo.On("GetByID", ctx, id).Return(want, nil).Once()

		got, err := svc.GetMachine(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("invalid-id", func(t *testing.T) {
		got, err := svc.GetMachine(ctx, 0)
		assert.ErrorIs(t, err, ErrInvalidMachineID)
		assert.Nil(t, got)
	})

	t.Run("not-found", func(t *testing.T) {
		mockMachineRepo.On("GetByID", ctx, id).Return(nil, repository.ErrMachineNotFound).Once()
		_, err := svc.GetMachine(ctx, id)
		assert.ErrorIs(t, err, ErrMachineNotFound)
	})
}

/* ------------------------------------------------------------------
   Tests: CreateMachine
-------------------------------------------------------------------*/
func TestMachineService_CreateMachine(t *testing.T) {
    ctx := context.Background()
    mockMachineRepo := new(MockMachineRepository)
    mockTaskTypeRepo := new(MockTaskTypeRepository)
    svc := NewMachineService(mockMachineRepo, mockTaskTypeRepo)

    taskTypeID := int32(10)         // ← 變數保留

    /* ---------- Success - Basic ---------- */
    // ...（保持原樣）

    /* ---------- Success - With TaskTypeID ---------- */
    t.Run("success-with-tasktype", func(t *testing.T) {
        req := models.CreateMachineRequest{
            Name: "Printer TT", Quantity: 1, EstimateTime: 120, TaskTypeID: &taskTypeID,
        }
        want := &models.Machine{
            ID: 2, MachineName: "Printer TT", Quantity: 1, EstimateTime: 120, TaskTypeID: &taskTypeID,
        }

        mockTaskTypeRepo.On("GetByID", ctx, taskTypeID).
            Return(&models.TaskType{ID: taskTypeID}, nil).Once()
        mockMachineRepo.On("Create", ctx, req).
            Return(want, nil).Once()

        got, err := svc.CreateMachine(ctx, req)
        assert.NoError(t, err)
        assert.Equal(t, want, got)
        mockTaskTypeRepo.AssertExpectations(t)
        mockMachineRepo.AssertExpectations(t)
    })

    /* ---------- 其餘測試案例保持不變 ---------- */
}

/* ------------------------------------------------------------------
   Tests: UpdateMachine  (僅測幾個代表案例)
-------------------------------------------------------------------*/
func TestMachineService_UpdateMachine(t *testing.T) {
	ctx := context.Background()
	mockMachineRepo := new(MockMachineRepository)
	mockTaskTypeRepo := new(MockTaskTypeRepository)
	svc := NewMachineService(mockMachineRepo, mockTaskTypeRepo)

	id := int32(1)
	exist := &models.Machine{ID: id, MachineName: "Old"}

	t.Run("invalid-id", func(t *testing.T) {
		_, err := svc.UpdateMachine(ctx, 0, models.UpdateMachineRequest{Name: "x", Quantity: 1, EstimateTime: 1})
		assert.ErrorIs(t, err, ErrInvalidMachineID)
	})

	t.Run("not-found-current", func(t *testing.T) {
		mockMachineRepo.On("GetByID", ctx, id).Return(nil, repository.ErrMachineNotFound).Once()
		_, err := svc.UpdateMachine(ctx, id, models.UpdateMachineRequest{Name: "x", Quantity: 1, EstimateTime: 1})
		assert.ErrorIs(t, err, ErrMachineNotFound)
	})

	t.Run("success-basic", func(t *testing.T) {
		req := models.UpdateMachineRequest{Name: "Upd", Quantity: 2, EstimateTime: 3}
		want := &models.Machine{ID: id, MachineName: "Upd", Quantity: 2, EstimateTime: 3}
		mockMachineRepo.On("GetByID", ctx, id).Return(exist, nil).Once()
		mockMachineRepo.On("Update", ctx, id, req).Return(want, nil).Once()

		got, err := svc.UpdateMachine(ctx, id, req)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

/* ------------------------------------------------------------------
   Tests: DeleteMachine
-------------------------------------------------------------------*/
func TestMachineService_DeleteMachine(t *testing.T) {
	ctx := context.Background()
	mockMachineRepo := new(MockMachineRepository)
	svc := NewMachineService(mockMachineRepo, new(MockTaskTypeRepository))
	id := int32(1)

	t.Run("success", func(t *testing.T) {
		mockMachineRepo.On("Delete", ctx, id).Return(nil).Once()
		assert.NoError(t, svc.DeleteMachine(ctx, id))
	})

	t.Run("invalid-id", func(t *testing.T) {
		err := svc.DeleteMachine(ctx, 0)
		assert.ErrorIs(t, err, ErrInvalidMachineID)
	})

	t.Run("not-found", func(t *testing.T) {
		mockMachineRepo.On("Delete", ctx, id).Return(repository.ErrMachineNotFound).Once()
		err := svc.DeleteMachine(ctx, id)
		assert.ErrorIs(t, err, ErrMachineNotFound)
	})
}
