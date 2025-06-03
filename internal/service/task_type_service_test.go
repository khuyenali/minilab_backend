package service

import (
	"context"
	"errors"
	"testing"

	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

/*
   ----------------------------------------------------------------
   GetTaskTypes
   ----------------------------------------------------------------
*/
func TestTaskTypeService_GetTaskTypes(t *testing.T) {
	mockTaskTypeRepo := new(MockTaskTypeRepository)
	mockUserRepo := new(MockUserRepository)
	svc := NewTaskTypeService(mockTaskTypeRepo, mockUserRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		desc := "desc"
		want := []*models.TaskType{
			{ID: 1, TypeName: "TT1", Description: &desc},
			{ID: 2, TypeName: "TT2"},
		}
		mockTaskTypeRepo.On("List", ctx).Return(want, nil).Once()

		got, err := svc.GetTaskTypes(ctx)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		mockTaskTypeRepo.AssertExpectations(t)
	})

	t.Run("repo-error", func(t *testing.T) {
		repoErr := errors.New("db error")
		mockTaskTypeRepo.On("List", ctx).Return(nil, repoErr).Once()

		got, err := svc.GetTaskTypes(ctx)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, got)
		mockTaskTypeRepo.AssertExpectations(t)
	})
}

/*
   ----------------------------------------------------------------
   GetTaskType
   ----------------------------------------------------------------
*/
func TestTaskTypeService_GetTaskType(t *testing.T) {
	mockTaskTypeRepo := new(MockTaskTypeRepository)
	svc := NewTaskTypeService(mockTaskTypeRepo, new(MockUserRepository))
	ctx := context.Background()
	id := int32(1)

	t.Run("success", func(t *testing.T) {
		desc := "d"
		want := &models.TaskType{ID: id, TypeName: "TT", Description: &desc}
		mockTaskTypeRepo.On("GetWithMachines", ctx, id).Return(want, nil).Once()

		got, err := svc.GetTaskType(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		mockTaskTypeRepo.AssertExpectations(t)
	})

	t.Run("invalid-id", func(t *testing.T) {
		got, err := svc.GetTaskType(ctx, 0)
		assert.ErrorIs(t, err, ErrInvalidTaskTypeID)
		assert.Nil(t, got)
	})

	t.Run("not-found", func(t *testing.T) {
		mockTaskTypeRepo.On("GetWithMachines", ctx, id).
			Return(nil, repository.ErrTaskTypeNotFound).Once()

		got, err := svc.GetTaskType(ctx, id)
		assert.ErrorIs(t, err, ErrTaskTypeNotFound)
		assert.Nil(t, got)
		mockTaskTypeRepo.AssertExpectations(t)
	})
}

/*
   ----------------------------------------------------------------
   CreateTaskType
   ----------------------------------------------------------------
*/
func TestTaskTypeService_CreateTaskType(t *testing.T) {
	ctx := context.Background()
	mockTaskTypeRepo := new(MockTaskTypeRepository)
	mockUserRepo := new(MockUserRepository)
	svc := NewTaskTypeService(mockTaskTypeRepo, mockUserRepo)

	desc := "Test Description"
	memberRoleID := int32(3)

	t.Run("success-basic", func(t *testing.T) {
		req := models.CreateTaskTypeRequest{Name: "New"}
		want := &models.TaskType{ID: 1, TypeName: "New"}

		mockTaskTypeRepo.On("Create", ctx, req).Return(want, nil).Once()

		got, err := svc.CreateTaskType(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		mockTaskTypeRepo.AssertExpectations(t)
	})

	t.Run("success-with-users", func(t *testing.T) {
		userIDs := []int32{1, 2}
		req := models.CreateTaskTypeRequest{Name: "WithUsers", Description: &desc, UserIDs: userIDs}
		want := &models.TaskType{ID: 2, TypeName: "WithUsers", Description: &desc}

		for _, uid := range userIDs {
			mockUserRepo.On("GetByID", ctx, uid).
				Return(&models.User{ID: uid, RoleID: memberRoleID}, nil).Once()
		}
		mockTaskTypeRepo.On("Create", ctx, req).Return(want, nil).Once()

		got, err := svc.CreateTaskType(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		mockUserRepo.AssertExpectations(t)
		mockTaskTypeRepo.AssertExpectations(t)
	})

	t.Run("invalid-user-ids", func(t *testing.T) {
		req := models.CreateTaskTypeRequest{Name: "BadUsers", UserIDs: []int32{1, 999}}

		mockUserRepo.On("GetByID", ctx, int32(1)).
			Return(&models.User{ID: 1, RoleID: memberRoleID}, nil).Once()
		mockUserRepo.On("GetByID", ctx, int32(999)).
			Return(nil, repository.ErrUserNotFound).Once()

		got, err := svc.CreateTaskType(ctx, req)
		assert.Nil(t, got)
		assert.EqualError(t, err, "invalid user IDs: 999")
	})
}

/*
   ----------------------------------------------------------------
   UpdateTaskType
   ----------------------------------------------------------------
*/
func TestTaskTypeService_UpdateTaskType(t *testing.T) {
	ctx := context.Background()
	mockTaskTypeRepo := new(MockTaskTypeRepository)
	mockUserRepo := new(MockUserRepository)
	svc := NewTaskTypeService(mockTaskTypeRepo, mockUserRepo)

	id := int32(1)
	memberRoleID := int32(3)
	desc := "upd"

	t.Run("success", func(t *testing.T) {
		req := models.UpdateTaskTypeRequest{Name: "Upd", Description: &desc, UserIDs: []int32{1}}
		want := &models.TaskType{ID: id, TypeName: "Upd", Description: &desc}

		mockUserRepo.On("GetByID", ctx, int32(1)).
			Return(&models.User{ID: 1, RoleID: memberRoleID}, nil).Once()
		mockTaskTypeRepo.On("Update", ctx, id, req).Return(want, nil).Once()

		got, err := svc.UpdateTaskType(ctx, id, req)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

/*
   ----------------------------------------------------------------
   DeleteTaskType
   ----------------------------------------------------------------
*/
func TestTaskTypeService_DeleteTaskType(t *testing.T) {
	ctx := context.Background()
	mockTaskTypeRepo := new(MockTaskTypeRepository)
	svc := NewTaskTypeService(mockTaskTypeRepo, new(MockUserRepository))
	id := int32(1)

	t.Run("success", func(t *testing.T) {
		mockTaskTypeRepo.On("Delete", ctx, id).Return(nil).Once()
		assert.NoError(t, svc.DeleteTaskType(ctx, id))
	})

	t.Run("invalid-id", func(t *testing.T) {
		err := svc.DeleteTaskType(ctx, 0)
		assert.ErrorIs(t, err, ErrInvalidTaskTypeID)
	})

	t.Run("not-found", func(t *testing.T) {
		mockTaskTypeRepo.On("Delete", ctx, id).
			Return(repository.ErrTaskTypeNotFound).Once()
		err := svc.DeleteTaskType(ctx, id)
		assert.ErrorIs(t, err, ErrTaskTypeNotFound)
	})
}
