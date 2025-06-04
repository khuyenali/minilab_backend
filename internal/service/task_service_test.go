package service

import (
    "context"
    "database/sql"
    "errors"
    "reflect"
    "testing"
    "time"

    "mini-lab-api/internal/models"
    "mini-lab-api/internal/repository"
)

// -----------------------------------------------------------------------------
// universal no-op helpers – used by every stub so we don't repeat ourselves
// -----------------------------------------------------------------------------
var (
    noopTask    = &models.Task{}
    noopTasks   = []*models.Task{}
    noopTT      = &models.TaskType{}
    noopTTs     = []*models.TaskType{}
    noopUsers   = []*models.User{}
    noopUser    = &models.User{}
    noopErr     error
)

// -----------------------------------------------------------------------------
// mock TaskRepository (must satisfy every method in the real interface)
// -----------------------------------------------------------------------------
type mockTaskRepository struct {
    CreateFunc          func(context.Context, models.CreateTaskRequest) (*models.Task, error)
    GetWithSubTasksFunc func(context.Context, int32) (*models.Task, error)
    ListWithSubTasksFunc func(context.Context) ([]*models.Task, error)
    GetByIDFunc         func(context.Context, int32) (*models.Task, error)
    UpdateFunc          func(context.Context, int32, models.UpdateTaskRequest) (*models.Task, error)
    DeleteFunc          func(context.Context, int32) error
    ListFunc            func(context.Context) ([]*models.Task, error)  // Add this missing method
}

func (m *mockTaskRepository) Create(ctx context.Context, r models.CreateTaskRequest) (*models.Task, error) {
    if m.CreateFunc != nil { return m.CreateFunc(ctx, r) }
    return noopTask, noopErr
}
func (m *mockTaskRepository) GetWithSubTasks(ctx context.Context, id int32) (*models.Task, error) {
    if m.GetWithSubTasksFunc != nil { return m.GetWithSubTasksFunc(ctx, id) }
    return noopTask, noopErr
}
func (m *mockTaskRepository) ListWithSubTasks(ctx context.Context) ([]*models.Task, error) {
    if m.ListWithSubTasksFunc != nil { return m.ListWithSubTasksFunc(ctx) }
    return noopTasks, noopErr
}
func (m *mockTaskRepository) GetByID(ctx context.Context, id int32) (*models.Task, error) {
    if m.GetByIDFunc != nil { return m.GetByIDFunc(ctx, id) }
    return noopTask, noopErr
}
func (m *mockTaskRepository) Update(ctx context.Context, id int32, r models.UpdateTaskRequest) (*models.Task, error) {
    if m.UpdateFunc != nil { return m.UpdateFunc(ctx, id, r) }
    return noopTask, noopErr
}
func (m *mockTaskRepository) Delete(ctx context.Context, id int32) error {
    if m.DeleteFunc != nil { return m.DeleteFunc(ctx, id) }
    return noopErr
}
func (m *mockTaskRepository) List(ctx context.Context) ([]*models.Task, error) {
    if m.ListFunc != nil { return m.ListFunc(ctx) }
    return noopTasks, noopErr
}

// -----------------------------------------------------------------------------
// mock TaskTypeRepository
// -----------------------------------------------------------------------------
type mockTaskTypeRepository struct {
    CreateFunc                    func(context.Context, models.CreateTaskTypeRequest) (*models.TaskType, error)
    ListFunc                      func(context.Context) ([]*models.TaskType, error)
    GetByIDFunc                   func(context.Context, int32) (*models.TaskType, error)
    UpdateFunc                    func(context.Context, int32, models.UpdateTaskTypeRequest) (*models.TaskType, error)
    DeleteFunc                    func(context.Context, int32) error
    AssignMachinesToTaskTypeFunc  func(context.Context, []int32, int32) error
    AssignUsersToTaskTypeFunc     func(context.Context, []int32, int32) error
    GetWithMachinesFunc           func(context.Context, int32) (*models.TaskType, error)
    ClearMachineAssignmentsFunc   func(context.Context, int32) error
    ClearUserAssignmentsFunc      func(context.Context, int32) error
    ValidateMachineIDsFunc        func(context.Context, []int32) error
    ValidateUserIDsFunc           func(context.Context, []int32) error
}

func (m *mockTaskTypeRepository) Create(ctx context.Context, r models.CreateTaskTypeRequest) (*models.TaskType, error) {
    if m.CreateFunc != nil { return m.CreateFunc(ctx, r) }
    return noopTT, noopErr
}
func (m *mockTaskTypeRepository) List(ctx context.Context) ([]*models.TaskType, error) {
    if m.ListFunc != nil { return m.ListFunc(ctx) }
    return noopTTs, noopErr
}
func (m *mockTaskTypeRepository) GetByID(ctx context.Context, id int32) (*models.TaskType, error) {
    if m.GetByIDFunc != nil { return m.GetByIDFunc(ctx, id) }
    return noopTT, noopErr
}
func (m *mockTaskTypeRepository) Update(ctx context.Context, id int32, r models.UpdateTaskTypeRequest) (*models.TaskType, error) {
    if m.UpdateFunc != nil { return m.UpdateFunc(ctx, id, r) }
    return noopTT, noopErr
}
func (m *mockTaskTypeRepository) Delete(ctx context.Context, id int32) error {
    if m.DeleteFunc != nil { return m.DeleteFunc(ctx, id) }
    return noopErr
}
func (m *mockTaskTypeRepository) AssignMachinesToTaskType(ctx context.Context, mids []int32, tid int32) error {
    if m.AssignMachinesToTaskTypeFunc != nil { return m.AssignMachinesToTaskTypeFunc(ctx, mids, tid) }
    return noopErr
}
func (m *mockTaskTypeRepository) AssignUsersToTaskType(ctx context.Context, uids []int32, tid int32) error {
    if m.AssignUsersToTaskTypeFunc != nil { return m.AssignUsersToTaskTypeFunc(ctx, uids, tid) }
    return noopErr
}
func (m *mockTaskTypeRepository) GetWithMachines(ctx context.Context, id int32) (*models.TaskType, error) {
    if m.GetWithMachinesFunc != nil { return m.GetWithMachinesFunc(ctx, id) }
    return noopTT, noopErr
}
func (m *mockTaskTypeRepository) ClearMachineAssignments(ctx context.Context, taskTypeID int32) error {
    if m.ClearMachineAssignmentsFunc != nil { return m.ClearMachineAssignmentsFunc(ctx, taskTypeID) }
    return noopErr
}
func (m *mockTaskTypeRepository) ClearUserAssignments(ctx context.Context, taskTypeID int32) error {
    if m.ClearUserAssignmentsFunc != nil { return m.ClearUserAssignmentsFunc(ctx, taskTypeID) }
    return noopErr
}
func (m *mockTaskTypeRepository) ValidateMachineIDs(ctx context.Context, machineIDs []int32) error {
	if m.ValidateMachineIDsFunc != nil {
		return m.ValidateMachineIDsFunc(ctx, machineIDs)
	}
	return noopErr
}
func (m *mockTaskTypeRepository) ValidateUserIDs(ctx context.Context, userIDs []int32) error {
	if m.ValidateUserIDsFunc != nil {
		return m.ValidateUserIDsFunc(ctx, userIDs)
	}
	return noopErr
}

// -----------------------------------------------------------------------------
// mock UserRepository
// -----------------------------------------------------------------------------
type mockUserRepository struct {
    CreateFunc                   func(context.Context, models.CreateUserRequest) (*models.User, error)
    ListFunc                     func(context.Context, models.UserFilters) ([]*models.User, error)
    GetByIDFunc                  func(context.Context, int32) (*models.User, error)
    GetByEmailFunc               func(context.Context, string) (*models.User, error)
    UpdateFunc                   func(context.Context, int32, models.UpdateUserRequest) (*models.User, error)
    DeleteFunc                   func(context.Context, int32) error
    GetTaskTypesFunc             func(context.Context, int32) ([]*models.TaskTypeBasic, error)
    AddTaskTypesFunc             func(context.Context, int32, []int32) error
    RemoveTaskTypesFunc          func(context.Context, int32, []int32) error
    GetMultiByIDsFunc            func(context.Context, []int32) ([]*models.User, error)
    AddUserToTaskTypeFunc        func(context.Context, int32, int32) error
    RemoveUserFromTaskTypeFunc   func(context.Context, int32, int32) error
    ClearTaskTypeAssignmentsFunc func(context.Context, int32) error
    GetByIDsFunc                 func(context.Context, []int32) ([]*models.User, error)  // Add this missing method here
    GetByRoleFunc                func(context.Context, string) ([]*models.User, error)
    GetByRoleIDFunc              func(context.Context, int32) ([]*models.User, error)
}

func (m *mockUserRepository) Create(ctx context.Context, r models.CreateUserRequest) (*models.User, error) {
    if m.CreateFunc != nil { return m.CreateFunc(ctx, r) }
    return noopUser, noopErr
}
func (m *mockUserRepository) List(ctx context.Context, f models.UserFilters) ([]*models.User, error) {
    if m.ListFunc != nil { return m.ListFunc(ctx, f) }
    return noopUsers, noopErr
}
func (m *mockUserRepository) GetByID(ctx context.Context, id int32) (*models.User, error) {
    if m.GetByIDFunc != nil { return m.GetByIDFunc(ctx, id) }
    return noopUser, noopErr
}
func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    if m.GetByEmailFunc != nil { return m.GetByEmailFunc(ctx, email) }
    return noopUser, noopErr
}
func (m *mockUserRepository) Update(ctx context.Context, id int32, r models.UpdateUserRequest) (*models.User, error) {
    if m.UpdateFunc != nil { return m.UpdateFunc(ctx, id, r) }
    return noopUser, noopErr
}
func (m *mockUserRepository) Delete(ctx context.Context, id int32) error {
    if m.DeleteFunc != nil { return m.DeleteFunc(ctx, id) }
    return noopErr
}
func (m *mockUserRepository) GetTaskTypes(ctx context.Context, uid int32) ([]*models.TaskTypeBasic, error) {
    if m.GetTaskTypesFunc != nil { return m.GetTaskTypesFunc(ctx, uid) }
    return nil, noopErr
}
func (m *mockUserRepository) AddTaskTypes(ctx context.Context, uid int32, tids []int32) error {
    if m.AddTaskTypesFunc != nil { return m.AddTaskTypesFunc(ctx, uid, tids) }
    return noopErr
}
func (m *mockUserRepository) RemoveTaskTypes(ctx context.Context, uid int32, tids []int32) error {
    if m.RemoveTaskTypesFunc != nil { return m.RemoveTaskTypesFunc(ctx, uid, tids) }
    return noopErr
}
func (m *mockUserRepository) GetMultiByIDs(ctx context.Context, ids []int32) ([]*models.User, error) {
    if m.GetMultiByIDsFunc != nil { return m.GetMultiByIDsFunc(ctx, ids) }
    return noopUsers, noopErr
}
func (m *mockUserRepository) AddUserToTaskType(ctx context.Context, uid int32, tid int32) error {
    if m.AddUserToTaskTypeFunc != nil { return m.AddUserToTaskTypeFunc(ctx, uid, tid) }
    return noopErr
}
func (m *mockUserRepository) RemoveUserFromTaskType(ctx context.Context, uid int32, tid int32) error {
    if m.RemoveUserFromTaskTypeFunc != nil { return m.RemoveUserFromTaskTypeFunc(ctx, uid, tid) }
    return noopErr
}
func (m *mockUserRepository) ClearTaskTypeAssignments(ctx context.Context, tid int32) error {
    if m.ClearTaskTypeAssignmentsFunc != nil { return m.ClearTaskTypeAssignmentsFunc(ctx, tid) }
    return noopErr
}
func (m *mockUserRepository) GetByIDs(ctx context.Context, ids []int32) ([]*models.User, error) {
    if m.GetByIDsFunc != nil { return m.GetByIDsFunc(ctx, ids) }
    return noopUsers, noopErr
}
func (m *mockUserRepository) GetByRole(ctx context.Context, role string) ([]*models.User, error) {
	if m.GetByRoleFunc != nil {
		return m.GetByRoleFunc(ctx, role)
	}
	return noopUsers, noopErr
}
func (m *mockUserRepository) GetByRoleID(ctx context.Context, roleID int32) ([]*models.User, error) {
	if m.GetByRoleIDFunc != nil {
		return m.GetByRoleIDFunc(ctx, roleID)
	}
	return noopUsers, noopErr
}

// TestTaskService_GetTask tests the GetTask method
func TestTaskService_GetTask(t *testing.T) {
    ctx := context.Background()
    expectedTask := &models.Task{
        ID:        1,
        TaskName:  "Test Task",
        Status:    "draft",
        Priority:  1,
        SubTasks:  []*models.SubTask{},
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    tests := []struct {
        name       string
        taskID     int32
        setupMocks func(mtr *mockTaskRepository)
        want       *models.Task
        wantErrStr string
    }{
        {
            name:   "success",
            taskID: 1,
            setupMocks: func(mtr *mockTaskRepository) {
                mtr.GetWithSubTasksFunc = func(ctx context.Context, id int32) (*models.Task, error) {
                    if id == 1 {
                        return expectedTask, nil
                    }
                    return nil, repository.ErrTaskNotFound
                }
            },
            want:       expectedTask,
            wantErrStr: "",
        },
        {
            name:       "invalid task ID",
            taskID:     0,
            setupMocks: func(mtr *mockTaskRepository) {},
            want:       nil,
            wantErrStr: ErrInvalidTaskID.Error(),
        },
        {
            name:   "task not found",
            taskID: 2,
            setupMocks: func(mtr *mockTaskRepository) {
                mtr.GetWithSubTasksFunc = func(ctx context.Context, id int32) (*models.Task, error) {
                    return nil, repository.ErrTaskNotFound
                }
            },
            want:       nil,
            wantErrStr: ErrTaskNotFound.Error(),
        },
        {
            name:   "other repository error",
            taskID: 1,
            setupMocks: func(mtr *mockTaskRepository) {
                mtr.GetWithSubTasksFunc = func(ctx context.Context, id int32) (*models.Task, error) {
                    return nil, errors.New("repo db error")
                }
            },
            want:       nil,
            wantErrStr: "repo db error",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockTaskRepo := &mockTaskRepository{}
            service := NewTaskService(mockTaskRepo, &mockTaskTypeRepository{}, &mockUserRepository{}, &sql.DB{}).(*taskService)
            tt.setupMocks(mockTaskRepo)

            got, err := service.GetTask(ctx, tt.taskID)

            if tt.wantErrStr != "" {
                if err == nil {
                    t.Errorf("GetTask() expected error_msg %q, got nil", tt.wantErrStr)
                } else if err.Error() != tt.wantErrStr {
                    t.Errorf("GetTask() error_msg = %q, want %q", err.Error(), tt.wantErrStr)
                }
            } else if err != nil {
                t.Errorf("GetTask() unexpected error_msg: %q", err.Error())
            }

            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetTask() got = %v, want %v", got, tt.want)
            }
        })
    }
}

// TestTaskService_GetTasks tests the GetTasks method
func TestTaskService_GetTasks(t *testing.T) {
	ctx := context.Background()
	expectedTasks := []*models.Task{
		{
			ID:        1,
			TaskName:  "Test Task 1",
			Status:    "draft",
			Priority:  1,
			SubTasks:  []*models.SubTask{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			TaskName:  "Test Task 2",
			Status:    "pending",
			Priority:  2,
			SubTasks:  []*models.SubTask{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	tests := []struct {
		name       string
		setupMocks func(mtr *mockTaskRepository)
		want       []*models.Task
		wantErrStr string
	}{
		{
			name: "success",
			setupMocks: func(mtr *mockTaskRepository) {
				mtr.ListWithSubTasksFunc = func(ctx context.Context) ([]*models.Task, error) {
					return expectedTasks, nil
				}
			},
			want:       expectedTasks,
			wantErrStr: "",
		},
		{
			name: "repository error",
			setupMocks: func(mtr *mockTaskRepository) {
				mtr.ListWithSubTasksFunc = func(ctx context.Context) ([]*models.Task, error) {
					return nil, errors.New("repo db error")
				}
			},
			want:       nil,
			wantErrStr: "repo db error",
		},
		{
			name: "no tasks found",
			setupMocks: func(mtr *mockTaskRepository) {
				mtr.ListWithSubTasksFunc = func(ctx context.Context) ([]*models.Task, error) {
					return []*models.Task{}, nil // Empty slice
				}
			},
			want:       []*models.Task{},
			wantErrStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTaskRepo := &mockTaskRepository{}
			service := NewTaskService(mockTaskRepo, &mockTaskTypeRepository{}, &mockUserRepository{}, &sql.DB{}).(*taskService)
			tt.setupMocks(mockTaskRepo)

			got, err := service.GetTasks(ctx)

			if tt.wantErrStr != "" {
				if err == nil {
					t.Errorf("GetTasks() expected error_msg %q, got nil", tt.wantErrStr)
				} else if err.Error() != tt.wantErrStr {
					t.Errorf("GetTasks() error_msg = %q, want %q", err.Error(), tt.wantErrStr)
				}
			} else if err != nil {
				t.Errorf("GetTasks() unexpected error_msg: %q", err.Error())
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTasks() got = %v, want %v", got, tt.want)
			}
		})
	}
}
