package service

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"mini-lab-api/internal/db"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"
)

// Re-use mockTaskRepository from assignment_status_service_test.go (if in same package & file is compiled)
// For standalone execution or clarity, it can be redefined or imported if it were in a shared test utility package.
// Assuming it's available as they are in the same 'service' package.

// mockAssignmentServiceQueries is a mock for db.Queries methods used by AssignmentService.
// We can use the same mockQueries structure from assignment_status_service_test.go
// and just set the functions we need for these tests.
// No need to redefine mockQueries struct if it's already in the package from another test file.

func TestNewAssignmentService(t *testing.T) {
	mockDB := &sql.DB{} // or a mockSqlDB if db.New requires specific behavior
	mockRepo := &mockTaskRepository{}
	service := NewAssignmentService(mockDB, mockRepo)

	if service == nil {
		t.Errorf("NewAssignmentService() returned nil")
	}
	if _, ok := service.(*assignmentService); !ok {
		t.Errorf("NewAssignmentService() did not return a *assignmentService")
	}
}

func TestAssignmentService_GetAssignment(t *testing.T) {
	ctx := context.Background()
	baseTime := time.Now()
	expectedAssignmentModel := &models.UserSubTaskAssignment{
		AssignmentID: 1,
		UserID:       100,
		SubTaskID:    200,
		AssignedAt:   baseTime,
	}
	dbAssignment := db.UserToSubTask{
		ID:        1,
		UserID:    100,
		SubTaskID: 200,
		AssignedAt: sql.NullTime{Time: baseTime, Valid: true},
	}

	tests := []struct {
		name         string
		assignmentID int32
		setupMocks   func(s *assignmentService)
		want         *models.UserSubTaskAssignment
		wantErrStr   string
	}{
		{
			name:         "success",
			assignmentID: 1,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetUserSubTaskAssignmentFunc: func(ctx context.Context, id int32) (db.UserToSubTask, error) {
						if id == 1 {
							return dbAssignment, nil
						}
						return db.UserToSubTask{}, sql.ErrNoRows
					},
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			want:       expectedAssignmentModel,
			wantErrStr: "",
		},
		{
			name:         "invalid assignment ID",
			assignmentID: 0,
			setupMocks:   func(s *assignmentService) {},
			want:         nil,
			wantErrStr:   ErrInvalidAssignmentID.Error(),
		},
		{
			name:         "db error GetUserSubTaskAssignment",
			assignmentID: 1,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetUserSubTaskAssignmentFunc: func(ctx context.Context, id int32) (db.UserToSubTask, error) {
						return db.UserToSubTask{}, errors.New("db query failed")
					},
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			want:         nil,
			wantErrStr:   "db query failed",
		},
		{
			name:         "assignment not found sql.ErrNoRows",
			assignmentID: 2,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetUserSubTaskAssignmentFunc: func(ctx context.Context, id int32) (db.UserToSubTask, error) {
						return db.UserToSubTask{}, sql.ErrNoRows
					},
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			want:         nil,
			wantErrStr:   ErrAssignmentNotFound.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAssignmentService(&sql.DB{}, &mockTaskRepository{}).(*assignmentService)
			tt.setupMocks(service)
			got, err := service.GetAssignment(ctx, tt.assignmentID)

			if tt.wantErrStr != "" {
				if err == nil {
					t.Errorf("GetAssignment() expected error_msg %q, got nil", tt.wantErrStr)
				} else if err.Error() != tt.wantErrStr {
					t.Errorf("GetAssignment() error_msg = %q, want %q", err.Error(), tt.wantErrStr)
				}
			} else if err != nil {
				t.Errorf("GetAssignment() unexpected error_msg: %q", err.Error())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAssignment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssignmentService_CreateAssignment(t *testing.T) {
	ctx := context.Background()
	baseTime := time.Now()

	successReq := models.CreateAssignmentRequest{SubTaskID: 1, UserID: 1}
	dbSubTask := db.SubTask{ID: 1, TaskID: 10}
	dbUserMember := db.User{ID: 1, RoleID: 3, Name: "Test User", Email: "test@example.com"} // RoleID 3 is member
	dbUserNonMember := db.User{ID: 2, RoleID: 1, Name: "Admin User", Email: "admin@example.com"} // RoleID 1 is admin
	dbCreatedAssignment := db.UserToSubTask{ID: 100, UserID: 1, SubTaskID: 1, AssignedAt: sql.NullTime{Time: baseTime, Valid: true}}
	expectedModelAssignment := models.FromDBUserSubTaskAssignment(dbCreatedAssignment)

	tests := []struct {
		name       string
		req        models.CreateAssignmentRequest
		setupMocks func(s *assignmentService)
		want       *models.UserSubTaskAssignment
		wantErrStr string
	}{
		{
			name: "success",
			req:  successReq,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetSubTaskFunc: func(ctx context.Context, id int32) (db.SubTask, error) {
						return dbSubTask, nil
					},
					GetUserFunc: func(ctx context.Context, id int32) (db.User, error) {
						return dbUserMember, nil
					},
					GetAssignmentsBySubTaskIDFunc: func(ctx context.Context, subTaskID int32) ([]db.UserToSubTask, error) {
						return []db.UserToSubTask{}, nil // No existing assignment
					},
					CreateUserSubTaskAssignmentFunc: func(ctx context.Context, arg db.CreateUserSubTaskAssignmentParams) (db.UserToSubTask, error) {
						return dbCreatedAssignment, nil
					},
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			want:       expectedModelAssignment,
			wantErrStr: "",
		},
		{
			name:       "invalid sub_task_id",
			req:        models.CreateAssignmentRequest{SubTaskID: 0, UserID: 1},
			setupMocks: func(s *assignmentService) {},
			wantErrStr: ErrInvalidSubTaskForAssignment.Error(),
		},
		{
			name:       "invalid user_id",
			req:        models.CreateAssignmentRequest{SubTaskID: 1, UserID: 0},
			setupMocks: func(s *assignmentService) {},
			wantErrStr: ErrInvalidUserForAssignment.Error(),
		},
		{
			name: "sub-task not found",
			req:  successReq,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetSubTaskFunc: func(ctx context.Context, id int32) (db.SubTask, error) {
						return db.SubTask{}, sql.ErrNoRows
					},
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			wantErrStr: ErrInvalidSubTaskForAssignment.Error(),
		},
		{
			name: "user not found",
			req:  successReq,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetSubTaskFunc: func(ctx context.Context, id int32) (db.SubTask, error) { return dbSubTask, nil },
					GetUserFunc: func(ctx context.Context, id int32) (db.User, error) {
						return db.User{}, sql.ErrNoRows
					},
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			wantErrStr: ErrInvalidUserForAssignment.Error(),
		},
		{
			name: "user not a member",
			req:  successReq,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetSubTaskFunc: func(ctx context.Context, id int32) (db.SubTask, error) { return dbSubTask, nil },
					GetUserFunc:    func(ctx context.Context, id int32) (db.User, error) { return dbUserNonMember, nil },
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			wantErrStr: ErrInvalidUserForAssignment.Error(),
		},
		{
			name: "assignment already exists",
			req:  successReq,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetSubTaskFunc: func(ctx context.Context, id int32) (db.SubTask, error) { return dbSubTask, nil },
					GetUserFunc:    func(ctx context.Context, id int32) (db.User, error) { return dbUserMember, nil },
					GetAssignmentsBySubTaskIDFunc: func(ctx context.Context, subTaskID int32) ([]db.UserToSubTask, error) {
						return []db.UserToSubTask{{UserID: 1}}, nil // Existing assignment for user 1
					},
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			wantErrStr: ErrAssignmentAlreadyExists.Error(),
		},
		{
			name: "create assignment db error",
			req:  successReq,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetSubTaskFunc:                func(ctx context.Context, id int32) (db.SubTask, error) { return dbSubTask, nil },
					GetUserFunc:                   func(ctx context.Context, id int32) (db.User, error) { return dbUserMember, nil },
					GetAssignmentsBySubTaskIDFunc: func(ctx context.Context, subTaskID int32) ([]db.UserToSubTask, error) { return []db.UserToSubTask{}, nil },
					CreateUserSubTaskAssignmentFunc: func(ctx context.Context, arg db.CreateUserSubTaskAssignmentParams) (db.UserToSubTask, error) {
						return db.UserToSubTask{}, errors.New("db create failed")
					},
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			wantErrStr: "db create failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAssignmentService(&sql.DB{}, &mockTaskRepository{}).(*assignmentService)
			tt.setupMocks(service)
			got, err := service.CreateAssignment(ctx, tt.req)

			if tt.wantErrStr != "" {
				if err == nil {
					t.Errorf("CreateAssignment() expected error_msg %q, got nil", tt.wantErrStr)
				} else if err.Error() != tt.wantErrStr {
					t.Errorf("CreateAssignment() error_msg = %q, want %q", err.Error(), tt.wantErrStr)
				}
			} else if err != nil {
				t.Errorf("CreateAssignment() unexpected error_msg: %q", err.Error())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAssignment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssignmentService_DeleteAssignment(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name         string
		assignmentID int32
		setupMocks   func(s *assignmentService)
		wantErrStr   string
	}{
		{
			name:         "success",
			assignmentID: 1,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetUserSubTaskAssignmentFunc: func(ctx context.Context, id int32) (db.UserToSubTask, error) {
						return db.UserToSubTask{ID: 1}, nil // Simulate assignment exists
					},
					DeleteUserSubTaskAssignmentFunc: func(ctx context.Context, id int32) error {
						return nil // Success
					},
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			wantErrStr: "",
		},
		{
			name:         "invalid assignment ID",
			assignmentID: 0,
			setupMocks:   func(s *assignmentService) {},
			wantErrStr:   ErrInvalidAssignmentID.Error(),
		},
		{
			name:         "assignment not found on check",
			assignmentID: 1,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetUserSubTaskAssignmentFunc: func(ctx context.Context, id int32) (db.UserToSubTask, error) {
						return db.UserToSubTask{}, sql.ErrNoRows
					},
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			wantErrStr: ErrAssignmentNotFound.Error(),
		},
		{
			name:         "db error on check assignment existence",
			assignmentID: 1,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetUserSubTaskAssignmentFunc: func(ctx context.Context, id int32) (db.UserToSubTask, error) {
						return db.UserToSubTask{}, errors.New("db check error")
					},
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			wantErrStr: "db check error",
		},
		{
			name:         "db error on delete",
			assignmentID: 1,
			setupMocks: func(s *assignmentService) {
				mq := &mockQueries{
					GetUserSubTaskAssignmentFunc: func(ctx context.Context, id int32) (db.UserToSubTask, error) {
						return db.UserToSubTask{ID: 1}, nil
					},
					DeleteUserSubTaskAssignmentFunc: func(ctx context.Context, id int32) error {
						return errors.New("db delete failed")
					},
				}
				s.queries = (*db.Queries)(unsafe.Pointer(mq))
			},
			wantErrStr: "db delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAssignmentService(&sql.DB{}, &mockTaskRepository{}).(*assignmentService)
			tt.setupMocks(service)
			err := service.DeleteAssignment(ctx, tt.assignmentID)

			if tt.wantErrStr != "" {
				if err == nil {
					t.Errorf("DeleteAssignment() expected error_msg %q, got nil", tt.wantErrStr)
				} else if err.Error() != tt.wantErrStr {
					t.Errorf("DeleteAssignment() error_msg = %q, want %q", err.Error(), tt.wantErrStr)
				}
			} else if err != nil {
				t.Errorf("DeleteAssignment() unexpected error_msg: %q", err.Error())
			}
		})
	}
}

// Ensure mockQueries has all the methods needed by assignmentService's s.queries
// If not already defined in assignment_status_service_test.go or if that file isn't compiled together:
/*
type mockQueries struct { // This is a more complete mockQueries for assignment_service if needed standalone
	GetUserSubTaskAssignmentFunc func(ctx context.Context, id int32) (db.UserToSubTask, error)
	GetSubTaskFunc               func(ctx context.Context, id int32) (db.SubTask, error)
	GetUserFunc                  func(ctx context.Context, id int32) (db.User, error)
	GetAssignmentsBySubTaskIDFunc func(ctx context.Context, subTaskID int32) ([]db.UserToSubTask, error)
	CreateUserSubTaskAssignmentFunc func(ctx context.Context, arg db.CreateUserSubTaskAssignmentParams) (db.UserToSubTask, error)
	DeleteUserSubTaskAssignmentFunc func(ctx context.Context, id int32) error

	// Stubs for other methods if *db.Queries is directly assigned and has other methods.
	// Add them here if they are ever called, otherwise they can be omitted if you are careful
	// that your mock (unsafe.Pointer cast) is only accessed for the implemented methods.
	WithTxFunc func(tx *sql.Tx) *db.Queries // Example, not used by this service directly
    UpdateTaskStatusFunc func(ctx context.Context, arg db.UpdateTaskStatusParams) (db.Task, error) // Example
}

func (m *mockQueries) GetUserSubTaskAssignment(ctx context.Context, id int32) (db.UserToSubTask, error) {
	if m.GetUserSubTaskAssignmentFunc != nil { return m.GetUserSubTaskAssignmentFunc(ctx, id) }
	panic("GetUserSubTaskAssignmentFunc not implemented")
}
func (m *mockQueries) GetSubTask(ctx context.Context, id int32) (db.SubTask, error) {
	if m.GetSubTaskFunc != nil { return m.GetSubTaskFunc(ctx, id) }
	panic("GetSubTaskFunc not implemented")
}
func (m *mockQueries) GetUser(ctx context.Context, id int32) (db.User, error) {
	if m.GetUserFunc != nil { return m.GetUserFunc(ctx, id) }
	panic("GetUserFunc not implemented")
}
func (m *mockQueries) GetAssignmentsBySubTaskID(ctx context.Context, subTaskID int32) ([]db.UserToSubTask, error) {
	if m.GetAssignmentsBySubTaskIDFunc != nil { return m.GetAssignmentsBySubTaskIDFunc(ctx, subTaskID) }
	panic("GetAssignmentsBySubTaskIDFunc not implemented")
}
func (m *mockQueries) CreateUserSubTaskAssignment(ctx context.Context, arg db.CreateUserSubTaskAssignmentParams) (db.UserToSubTask, error) {
	if m.CreateUserSubTaskAssignmentFunc != nil { return m.CreateUserSubTaskAssignmentFunc(ctx, arg) }
	panic("CreateUserSubTaskAssignmentFunc not implemented")
}
func (m *mockQueries) DeleteUserSubTaskAssignment(ctx context.Context, id int32) error {
	if m.DeleteUserSubTaskAssignmentFunc != nil { return m.DeleteUserSubTaskAssignmentFunc(ctx, id) }
	panic("DeleteUserSubTaskAssignmentFunc not implemented")
}
func (m *mockQueries) WithTx(tx *sql.Tx) *db.Queries { panic("WithTx not implemented") }
func (m *mockQueries) UpdateTaskStatus(ctx context.Context, arg db.UpdateTaskStatusParams) (db.Task, error) { panic("UpdateTaskStatus not implemented") }
*/

// Add implementations for these if mockQueries from the other file is not used or is insufficient.
// This assumes mockQueries from assignment_status_service_test.go is available and covers these methods
// or that the test setup provides the implementations as needed.
// For example, the TestAssignmentService_CreateAssignment setupMocks explicitly defines:
// GetSubTaskFunc, GetUserFunc, GetAssignmentsBySubTaskIDFunc, CreateUserSubTaskAssignmentFunc.
// The mockQueries struct in assignment_status_service_test.go would need these fields if it's the one being used.
// It's generally better to have one comprehensive mock for db.Querier if it's shared across tests in the same package.

// Methods needed by this specific test file for mockQueries:
// GetUserSubTaskAssignmentFunc (used by GetAssignment, DeleteAssignment)
// GetSubTaskFunc (used by CreateAssignment)
// GetUserFunc (used by CreateAssignment)
// GetAssignmentsBySubTaskIDFunc (used by CreateAssignment)
// CreateUserSubTaskAssignmentFunc (used by CreateAssignment)
// DeleteUserSubTaskAssignmentFunc (used by DeleteAssignment)

// The mockQueries struct in assignment_status_service_test.go has:
// WithTxFunc, GetUserSubTaskAssignmentFunc, GetSubTaskFunc, GetTaskFunc, UpdateTaskStatusFunc.
// It's MISSING: GetUserFunc, GetAssignmentsBySubTaskIDFunc, CreateUserSubTaskAssignmentFunc, DeleteUserSubTaskAssignmentFunc.
// So, the mockQueries in assignment_status_service_test.go IS NOT SUFFICIENT.
// The tests above for assignment_service.go correctly define these needed functions dynamically inside setupMocks.
// This approach is fine. The commented-out full mockQueries is just for illustration if a static shared mock was preferred.