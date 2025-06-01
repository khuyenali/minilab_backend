-- name: CreateTask :one
INSERT INTO tasks (task_name, status, priority, start_time, end_time, note)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY created_at DESC;

-- name: UpdateTask :one
UPDATE tasks
SET task_name = $2, status = $3, priority = $4, start_time = $5, end_time = $6, note = $7, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1; 