-- name: GetMachine :one
SELECT id, machine_name, quantity, estimate_time, type_id, created_at, updated_at 
FROM machines 
WHERE id = $1 LIMIT 1;

-- name: GetMachineWithTaskType :one
SELECT m.id, m.machine_name, m.quantity, m.estimate_time, m.type_id, m.created_at, m.updated_at,
       tt.type_name as task_type_name
FROM machines m
LEFT JOIN task_types tt ON m.type_id = tt.id
WHERE m.id = $1 LIMIT 1;

-- name: ListMachines :many
SELECT id, machine_name, quantity, estimate_time, type_id, created_at, updated_at 
FROM machines 
ORDER BY id;

-- name: ListMachinesWithTaskType :many
SELECT m.id, m.machine_name, m.quantity, m.estimate_time, m.type_id, m.created_at, m.updated_at,
       tt.type_name as task_type_name
FROM machines m
LEFT JOIN task_types tt ON m.type_id = tt.id
ORDER BY m.id;

-- name: CreateMachine :one
INSERT INTO machines (machine_name, quantity, estimate_time, type_id) 
VALUES ($1, $2, $3, $4) 
RETURNING id, machine_name, quantity, estimate_time, type_id, created_at, updated_at;

-- name: UpdateMachine :one
UPDATE machines 
SET machine_name = $2, quantity = $3, estimate_time = $4, updated_at = NOW() 
WHERE id = $1 
RETURNING id, machine_name, quantity, estimate_time, type_id, created_at, updated_at;

-- name: UpdateMachineWithTaskType :one
UPDATE machines 
SET machine_name = $2, quantity = $3, estimate_time = $4, type_id = $5, updated_at = NOW() 
WHERE id = $1 
RETURNING id, machine_name, quantity, estimate_time, type_id, created_at, updated_at;

-- name: DeleteMachine :exec
DELETE FROM machines 
WHERE id = $1;

-- name: GetMachinesByIDs :many
SELECT id, machine_name, quantity, estimate_time, type_id, created_at, updated_at 
FROM machines 
WHERE id = ANY($1::int[])
ORDER BY id;

-- name: UpdateMachineTaskType :exec
UPDATE machines 
SET type_id = $2, updated_at = NOW() 
WHERE id = $1;

-- name: ClearMachineTaskType :exec
UPDATE machines 
SET type_id = NULL, updated_at = NOW() 
WHERE type_id = $1; 