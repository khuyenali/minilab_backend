-- name: CreateTask :one
INSERT INTO tasks (task_name, status, priority, start_time, end_time, note, report)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY created_at DESC;

-- name: GetPendingTasksSortedByPriority :many
SELECT * FROM tasks
WHERE status = 'pending'
ORDER BY priority ASC, created_at ASC;

-- name: GetDraftTasksSortedByPriority :many
SELECT * FROM tasks
WHERE status = 'draft'
ORDER BY priority ASC, created_at ASC;

-- name: GetProcessingTasksSortedByPriority :many
SELECT * FROM tasks
WHERE status = 'processing'
ORDER BY priority ASC, created_at ASC;

-- name: GetTasksByStatus :many
SELECT * FROM tasks
WHERE status = $1
ORDER BY priority ASC, created_at ASC;

-- name: UpdateTask :one
UPDATE tasks
SET task_name = $2, status = $3, priority = $4, start_time = $5, end_time = $6, note = $7, report = $8, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateTaskStatusWithReport :one
UPDATE tasks
SET status = $2, report = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1; 