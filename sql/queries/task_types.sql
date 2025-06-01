-- name: GetTaskType :one
SELECT id, type_name, description, created_at, updated_at 
FROM task_types 
WHERE id = $1 LIMIT 1;

-- name: ListTaskTypes :many
SELECT id, type_name, description, created_at, updated_at 
FROM task_types 
ORDER BY id;

-- name: CreateTaskType :one
INSERT INTO task_types (type_name, description) 
VALUES ($1, $2) 
RETURNING id, type_name, description, created_at, updated_at;

-- name: UpdateTaskType :one
UPDATE task_types 
SET type_name = $2, description = $3, updated_at = NOW() 
WHERE id = $1 
RETURNING id, type_name, description, created_at, updated_at;

-- name: DeleteTaskType :exec
DELETE FROM task_types 
WHERE id = $1;

-- name: GetTaskTypeWithMachines :one
SELECT t.id, t.type_name, t.description, t.created_at, t.updated_at 
FROM task_types t 
WHERE t.id = $1 LIMIT 1;

-- name: GetMachinesByTaskType :many
SELECT m.id, m.machine_name, m.quantity, m.estimate_time, m.type_id, m.created_at, m.updated_at 
FROM machines m 
WHERE m.type_id = $1 
ORDER BY m.id; 