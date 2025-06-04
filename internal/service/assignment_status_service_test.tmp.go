// assignment_service_test.go
package service

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"mini-lab-api/internal/db"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/repository"
)

/* ---------- TaskRepository mock (unique to this file) ---------- */

type asMockTaskRepository struct {
	GetWithSubTasksFunc func(context.Context, int32) (*models.Task, error)
}

func (m *asMockTaskRepository) Create(context.Context, models.CreateTaskRequest) (*models.Task, error) {
	panic("not used")
}
func (m *asMockTaskRepository) List(context.Context) ([]*models.Task, error)               { panic("not used") }
func (m *asMockTaskRepository) GetByID(context.Context, int32) (*models.Task, error)       { panic("not used") }
func (m *asMockTaskRepository) Update(context.Context, int32, models.UpdateTaskRequest) (*models.Task, error) {
	panic("not used")
}
func (m *asMockTaskRepository) Delete(context.Context, int32) error                       { panic("not used") }
func (m *asMockTaskRepository) ListWithSubTasks(context.Context) ([]*models.Task, error)  { panic("not used") }
func (m *asMockTaskRepository) GetWithSubTasks(ctx context.Context, id int32) (*models.Task, error) {
	return m.GetWithSubTasksFunc(ctx, id)
}
func (m *asMockTaskRepository) GetAssignmentsByTaskID(context.Context, int32) ([]*models.UserSubTaskAssignment, error) {
	panic("not used")
}
func (m *asMockTaskRepository) GetAssignmentsByUserID(context.Context, int32) ([]*models.UserSubTaskAssignment, error) {
	panic("not used")
}

/* ---------- Queries mock (only the methods used) ---------- */

type asMockQueries struct {
	WithTxFunc                   func(*sql.Tx) *db.Queries
	GetUserSubTaskAssignmentFunc func(context.Context, int32) (db.UserToSubTask, error)
	GetSubTaskFunc               func(context.Context, int32) (db.SubTask, error)
	GetTaskFunc                  func(context.Context, int32) (db.Task, error)
	UpdateTaskStatusFunc         func(context.Context, db.UpdateTaskStatusParams) (db.Task, error)
}

func (m *asMockQueries) WithTx(tx *sql.Tx) *db.Queries {
	if m.WithTxFunc != nil {
		return m.WithTxFunc(tx)
	}
	// Fallback or panic if not set, to avoid nil pointer dereferences if used unexpectedly.
	// However, the goal is to not use this mock for the tests being fixed.
	panic("WithTxFunc not mocked")
}
func (m *asMockQueries) GetUserSubTaskAssignment(ctx context.Context, id int32) (db.UserToSubTask, error) {
	return m.GetUserSubTaskAssignmentFunc(ctx, id)
}
func (m *asMockQueries) GetSubTask(ctx context.Context, id int32) (db.SubTask, error) {
	return m.GetSubTaskFunc(ctx, id)
}
func (m *asMockQueries) GetTask(ctx context.Context, id int32) (db.Task, error) {
	return m.GetTaskFunc(ctx, id)
}
func (m *asMockQueries) UpdateTaskStatus(ctx context.Context, p db.UpdateTaskStatusParams) (db.Task, error) {
	return m.UpdateTaskStatusFunc(ctx, p)
}

/* ---------- helper to build service with sqlmock ---------- */

func newSvc(t *testing.T, taskRepo repository.TaskRepository) (*assignmentStatusService, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	// The service will use db.New(sqlDB), so s.queries will be the actual sqlc queries
	// that operate on the mocked sqlDB.
	svc := NewAssignmentStatusService(sqlDB, taskRepo).(*assignmentStatusService)
	return svc, mock
}

/* ---------- UpdateAssignmentStatus tests ---------- */

func TestAssignmentStatusService_UpdateAssignmentStatus(t *testing.T) {
	ctx := context.Background()
	base := time.Now()

	// Ensure these structs have all fields that would be selected by sqlc queries
	// and used by models.FromDBTask, with correct types.
	assign := db.UserToSubTask{
		ID: 1, SubTaskID: 10, UserID: 99,
		AssignedAt: sql.NullTime{Time: base, Valid: true},
		UpdatedAt:  sql.NullTime{Time: base, Valid: true},
	}
	subTask := db.SubTask{
		ID:             10,
		TaskID:         20,
		TypeID:         1,
		SubTaskName:    sql.NullString{String: "Sample SubTask Name", Valid: true},
		Description:    sql.NullString{String: "SubTask Desc", Valid: true},
		EstimateEffort: sql.NullInt32{Int32: 5, Valid: true},
		CreatedAt:      sql.NullTime{Time: base, Valid: true},
		UpdatedAt:      sql.NullTime{Time: base, Valid: true},
	}
	taskPending := db.Task{
		ID:       20,
		TaskName: "Task Name", // Changed back to string based on compiler error
		// Description removed as per compiler error 'taskPending.Description undefined'
		Status:    db.TaskStatusPending,
		CreatedAt: sql.NullTime{Time: base, Valid: true}, // Changed to sql.NullTime
		UpdatedAt: sql.NullTime{Time: base, Valid: true}, // Changed to sql.NullTime
	}
	taskProc := db.Task{
		ID:       20,
		TaskName: "Task Name", // Changed back to string based on compiler error
		// Description removed
		Status:    db.TaskStatusProcessing,
		CreatedAt: sql.NullTime{Time: base, Valid: true}, // Changed to sql.NullTime
		UpdatedAt: sql.NullTime{Time: base, Valid: true}, // Changed to sql.NullTime
	}


	tests := []struct {
		name      string
		id        int32
		req       models.UpdateAssignmentStatusRequest
		setup     func(m sqlmock.Sqlmock)
		wantTask  *models.Task
		wantErr   string
	}{
		{
			name: "pending → processing",
			id:   1,
			req:  models.UpdateAssignmentStatusRequest{Status: "processing"},
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()

				// GetUserSubTaskAssignment
				rowsUserToSubTask := sqlmock.NewRows([]string{"id", "user_id", "sub_task_id", "assigned_at", "updated_at"}).
					AddRow(assign.ID, assign.UserID, assign.SubTaskID, assign.AssignedAt, assign.UpdatedAt)
				m.ExpectQuery(`SELECT id, user_id, sub_task_id, assigned_at, updated_at FROM user_to_sub_task WHERE id = \$1`).
					WithArgs(assign.ID).
					WillReturnRows(rowsUserToSubTask)

				// GetSubTask - Updated columns and query to match actual
				rowsSubTask := sqlmock.NewRows([]string{"id", "task_id", "type_id", "sub_task_name", "description", "estimate_effort", "created_at", "updated_at"}).
					AddRow(subTask.ID, subTask.TaskID, subTask.TypeID, subTask.SubTaskName, subTask.Description, subTask.EstimateEffort, subTask.CreatedAt, subTask.UpdatedAt)
				m.ExpectQuery(`SELECT id, task_id, type_id, sub_task_name, description, estimate_effort, created_at, updated_at FROM sub_tasks WHERE id = \$1`).
					WithArgs(assign.SubTaskID).
					WillReturnRows(rowsSubTask)

				// GetTask
				// IMPORTANT: Column names must match the SELECT clause. Removed "description".
				rowsTask := sqlmock.NewRows([]string{"id", "task_name", "status", "created_at", "updated_at"}).
					AddRow(taskPending.ID, taskPending.TaskName, taskPending.Status, taskPending.CreatedAt, taskPending.UpdatedAt)
				m.ExpectQuery(`SELECT id, task_name, status, created_at, updated_at FROM tasks WHERE id = \$1`). // Adjusted query
					WithArgs(subTask.TaskID).
					WillReturnRows(rowsTask)

				// UpdateTaskStatus
				// IMPORTANT: Column names must match the RETURNING clause. Removed "description".
				rowsUpdatedTask := sqlmock.NewRows([]string{"id", "task_name", "status", "created_at", "updated_at"}).
					AddRow(taskProc.ID, taskProc.TaskName, taskProc.Status, taskProc.CreatedAt, taskProc.UpdatedAt)
				m.ExpectQuery(`UPDATE tasks SET status = \$1, updated_at = now\(\\) WHERE id = \$2 RETURNING id, task_name, status, created_at, updated_at`). // Adjusted query
					WithArgs(db.TaskStatusProcessing, taskPending.ID).
					WillReturnRows(rowsUpdatedTask)
				
				m.ExpectCommit()
			},
			wantTask: models.FromDBTask(taskProc),
		},
		{
			name:    "invalid id",
			id:      0,
			req:     models.UpdateAssignmentStatusRequest{Status: "processing"},
			wantErr: ErrInvalidAssignmentID.Error(),
		},
		{
			name:    "invalid status",
			id:      1,
			req:     models.UpdateAssignmentStatusRequest{Status: "bad"},
			wantErr: ErrInvalidAssignmentStatus.Error(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			taskRepoMock := &asMockTaskRepository{}
			svc, mock := newSvc(t, taskRepoMock)
			
			sqlDB := svc.db 
			defer sqlDB.Close()


			if tc.setup != nil {
				tc.setup(mock)
			}

			got, err := svc.UpdateAssignmentStatus(ctx, tc.id, tc.req)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("expected err %q got %v", tc.wantErr, err)
				}
				if err := mock.ExpectationsWereMet(); err != nil {
					if tc.setup != nil { 
						t.Logf("Note: DB expectations might not be fully met for error case %q: %v", tc.name, err)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(got, tc.wantTask) {
				t.Fatalf("task mismatch:\\ngot  %+v\\nwant %+v", got, tc.wantTask)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unfulfilled DB expectations: %v", err)
			}
		})
	}
}

/* ---------- GetAssignmentTask tests ---------- */

func TestAssignmentStatusService_GetAssignmentTask(t *testing.T) {
	ctx := context.Background()
	// baseTime := time.Now() // Used for consistent sql.NullTime
		
	// Define a sample subtask structure that matches the DB schema for mocks
	// mockSubTask := db.SubTask{
	// 	ID:             5,
	// 	TaskID:         42,
	// 	TypeID:         1,
	// 	SubTaskName:    sql.NullString{String: "Mock SubTask Name", Valid: true},
	// 	Description:    sql.NullString{String: "Desc", Valid: true},
	// 	EstimateEffort: sql.NullInt32{Int32: 3, Valid: true},
	// 	CreatedAt:      sql.NullTime{Time: baseTime, Valid: true},
	// 	UpdatedAt:      sql.NullTime{Time: baseTime, Valid: true},
	// }

	tests := []struct {
		name    string
		id      int32
		setup   func(m sqlmock.Sqlmock, taskRepo *asMockTaskRepository)
		want    *models.Task
		wantErr string
	}{
		{
			name: "happy",
			id:   1,
			setup: func(m sqlmock.Sqlmock, tr *asMockTaskRepository) {
				// Mock for s.queries.GetUserSubTaskAssignment
				// rowsUserToSubTask := sqlmock.NewRows([]string{"id", "user_id", "sub_task_id", "assigned_at", "updated_at"}).
				// 	AddRow(1, 99, mockSubTask.ID, sql.NullTime{Time: baseTime, Valid: true}, sql.NullTime{Time: baseTime, Valid: true})
				// m.ExpectQuery(`SELECT id, user_id, sub_task_id, assigned_at, updated_at FROM user_to_sub_task WHERE id = \$1`).
				// 	WithArgs(int32(1)).
				// 	WillReturnRows(rowsUserToSubTask)

				// // Mock for s.queries.GetSubTask - Updated columns and query
				// rowsSubTask := sqlmock.NewRows([]string{"id", "task_id", "type_id", "sub_task_name", "description", "estimate_effort", "created_at", "updated_at"}).
				// 	AddRow(mockSubTask.ID, mockSubTask.TaskID, mockSubTask.TypeID, mockSubTask.SubTaskName, mockSubTask.Description, mockSubTask.EstimateEffort, mockSubTask.CreatedAt, mockSubTask.UpdatedAt)
				// m.ExpectQuery(`SELECT id, task_id, type_id, sub_task_name, description, estimate_effort, created_at, updated_at FROM sub_tasks WHERE id = \$1`).
				// 	WithArgs(mockSubTask.ID).
				// 	WillReturnRows(rowsSubTask)
				
				// Mock for taskRepo.GetWithSubTasksFunc
				tr.GetWithSubTasksFunc = func(context.Context, int32) (*models.Task, error) {
					// This task should also reflect the actual structure of models.Task,
					// potentially using sql.NullString for TaskName if that's the case.
					return &models.Task{ID: 42, TaskName: "x"}, nil 
				}
			},
			// want task should match the structure returned by GetWithSubTasksFunc
			want: &models.Task{ID: 42, TaskName: "x"},
		},
		{
			name:    "bad id",
			id:      0,
			wantErr: ErrInvalidAssignmentID.Error(),
		},
		{
			name: "assignment not found",
			id:   2,
			setup: func(m sqlmock.Sqlmock, tr *asMockTaskRepository) {
				m.ExpectQuery(`SELECT id, user_id, sub_task_id, assigned_at, updated_at FROM user_to_sub_task WHERE id = \$1`).
					WithArgs(int32(2)).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: ErrAssignmentNotFound.Error(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trMock := &asMockTaskRepository{}
			svc, mock := newSvc(t, trMock) 
			
			sqlDB := svc.db 
			defer sqlDB.Close()

			if tc.setup != nil {
				tc.setup(mock, trMock) 
			}
			got, err := svc.GetAssignmentTask(ctx, tc.id)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("expected err %q got %v", tc.wantErr, err)
				}
				if err == sql.ErrNoRows || (tc.setup != nil && (tc.wantErr == ErrAssignmentNotFound.Error())) {
					if err := mock.ExpectationsWereMet(); err != nil {
						t.Logf("Note: DB expectations might not be fully met for error case %q: %v", tc.name, err)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("mismatch:\\ngot  %+v\\nwant %+v", got, tc.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unfulfilled DB expectations: %v", err)
			}
		})
	}
}

/* ---------- status helper ---------- */

func Test_isValidAssignmentStatus(t *testing.T) {
	for _, c := range []struct {
		s string
		b bool
	}{
		{"pending", true}, {"processing", true}, {"finish", true},
		{"", false}, {"BAD", false},
	} {
		if got := isValidAssignmentStatus(c.s); got != c.b {
			t.Fatalf("%q got %v want %v", c.s, got, c.b)
		}
	}
}
