-- name: CreateSubTask :one
INSERT INTO sub_tasks (task_id, type_id, sub_task_name, description, estimate_effort)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSubTask :one
SELECT * FROM sub_tasks
WHERE id = $1;

-- name: GetSubTasksByTaskID :many
SELECT * FROM sub_tasks
WHERE task_id = $1
ORDER BY created_at ASC;

-- name: UpdateSubTask :one
UPDATE sub_tasks
SET sub_task_name = $2, description = $3, estimate_effort = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteSubTask :exec
DELETE FROM sub_tasks
WHERE id = $1; 