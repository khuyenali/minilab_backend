package repository

import (
	"context"
	"database/sql"
	"mini-lab-api/internal/db"
	"mini-lab-api/internal/models"
)

// TaskRepository defines the interface for task data operations
type TaskRepository interface {
	GetByID(ctx context.Context, id int32) (*models.Task, error)
	List(ctx context.Context) ([]*models.Task, error)
	Create(ctx context.Context, req models.CreateTaskRequest) (*models.Task, error)
	Update(ctx context.Context, id int32, req models.UpdateTaskRequest) (*models.Task, error)
	Delete(ctx context.Context, id int32) error
	GetWithSubTasks(ctx context.Context, id int32) (*models.Task, error)
}

// taskRepository implements TaskRepository using sqlc generated code
type taskRepository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(database *sql.DB) TaskRepository {
	return &taskRepository{
		db:      database,
		queries: db.New(database),
	}
}

// GetByID retrieves a task by ID
func (r *taskRepository) GetByID(ctx context.Context, id int32) (*models.Task, error) {
	row, err := r.queries.GetTask(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	
	return models.FromDBTask(row), nil
}

// List retrieves all tasks
func (r *taskRepository) List(ctx context.Context) ([]*models.Task, error) {
	rows, err := r.queries.ListTasks(ctx)
	if err != nil {
		return nil, err
	}
	
	return models.FromDBTasks(rows), nil
}

// Create creates a new task with sub-tasks and assignments in a transaction
func (r *taskRepository) Create(ctx context.Context, req models.CreateTaskRequest) (*models.Task, error) {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create queries with transaction
	txQueries := r.queries.WithTx(tx)

	// Create the main task
	taskParams := req.ToCreateTaskParams()
	dbTask, err := txQueries.CreateTask(ctx, taskParams)
	if err != nil {
		return nil, err
	}

	// Convert to domain model
	task := models.FromDBTask(dbTask)

	// Create sub-tasks if provided
	if len(req.SubTasks) > 0 {
		subTasks := make([]*models.SubTask, 0, len(req.SubTasks))
		
		for _, subTaskReq := range req.SubTasks {
			// Create sub-task
			subTaskParams := db.CreateSubTaskParams{
				TaskID:         dbTask.ID,
				TypeID:         subTaskReq.TypeID,
				SubTaskName:    sql.NullString{}, // Will be set to null initially
				Description:    sql.NullString{}, // Will be set to null initially
				EstimateEffort: sql.NullInt32{},  // Will be set to null initially
			}

			dbSubTask, err := txQueries.CreateSubTask(ctx, subTaskParams)
			if err != nil {
				return nil, err
			}

			// Convert to domain model
			subTask := models.FromDBSubTask(dbSubTask)

			// Create user assignments if provided
			if len(subTaskReq.UserIDs) > 0 {
				assignments := make([]*models.UserSubTaskAssignment, 0, len(subTaskReq.UserIDs))
				
				for _, userID := range subTaskReq.UserIDs {
					assignmentParams := db.CreateUserSubTaskAssignmentParams{
						UserID:    userID,
						SubTaskID: dbSubTask.ID,
						Report:    sql.NullString{}, // Default to null
						Status:    db.NullAssignmentStatus{AssignmentStatus: db.AssignmentStatusPending, Valid: true}, // Default to pending
					}

					dbAssignment, err := txQueries.CreateUserSubTaskAssignment(ctx, assignmentParams)
					if err != nil {
						return nil, err
					}

					assignment := models.FromDBUserSubTaskAssignment(dbAssignment)
					assignments = append(assignments, assignment)
				}
				
				subTask.Assignments = assignments
			}

			subTasks = append(subTasks, subTask)
		}
		
		task.SubTasks = subTasks
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return task, nil
}

// Update updates an existing task
func (r *taskRepository) Update(ctx context.Context, id int32, req models.UpdateTaskRequest) (*models.Task, error) {
	// Check if task exists
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var startTime, endTime sql.NullTime
	var note sql.NullString

	if req.StartTime != nil {
		startTime = sql.NullTime{Time: *req.StartTime, Valid: true}
	}
	if req.EndTime != nil {
		endTime = sql.NullTime{Time: *req.EndTime, Valid: true}
	}
	if req.Note != nil {
		note = sql.NullString{String: *req.Note, Valid: true}
	}

	dbTask, err := r.queries.UpdateTask(ctx, db.UpdateTaskParams{
		ID:        id,
		TaskName:  req.TaskName,
		Status:    db.TaskStatus(req.Status),
		Priority:  req.Priority,
		StartTime: startTime,
		EndTime:   endTime,
		Note:      note,
	})
	if err != nil {
		return nil, err
	}

	return models.FromDBTask(dbTask), nil
}

// Delete deletes a task by ID
func (r *taskRepository) Delete(ctx context.Context, id int32) error {
	// Check if task exists
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = r.queries.DeleteTask(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// GetWithSubTasks retrieves a task with its sub-tasks and assignments
func (r *taskRepository) GetWithSubTasks(ctx context.Context, id int32) (*models.Task, error) {
	// Get the main task
	task, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get sub-tasks for this task
	dbSubTasks, err := r.queries.GetSubTasksByTaskID(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(dbSubTasks) > 0 {
		subTasks := make([]*models.SubTask, 0, len(dbSubTasks))
		
		for _, dbSubTask := range dbSubTasks {
			subTask := models.FromDBSubTask(dbSubTask)
			
			// Get assignments for this sub-task
			dbAssignments, err := r.queries.GetAssignmentsBySubTaskID(ctx, dbSubTask.ID)
			if err != nil {
				// Continue if we can't get assignments for this sub-task
				subTask.Assignments = []*models.UserSubTaskAssignment{}
			} else {
				subTask.Assignments = models.FromDBUserSubTaskAssignments(dbAssignments)
			}
			
			subTasks = append(subTasks, subTask)
		}
		
		task.SubTasks = subTasks
	}

	return task, nil
} 