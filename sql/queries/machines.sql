-- name: GetMachine :one
SELECT id, machine_name, quantity, estimate_time, type_id, created_at, updated_at 
FROM machines 
WHERE id = $1 LIMIT 1;

-- name: ListMachines :many
SELECT id, machine_name, quantity, estimate_time, type_id, created_at, updated_at 
FROM machines 
ORDER BY machine_name;

-- name: CreateMachine :one
INSERT INTO machines (machine_name, quantity, estimate_time, type_id) 
VALUES ($1, $2, $3, $4) 
RETURNING id, machine_name, quantity, estimate_time, type_id, created_at, updated_at;

-- name: UpdateMachine :one
UPDATE machines 
SET machine_name = $2, quantity = $3, estimate_time = $4, updated_at = NOW() 
WHERE id = $1 
RETURNING id, machine_name, quantity, estimate_time, type_id, created_at, updated_at;

-- name: DeleteMachine :exec
DELETE FROM machines 
WHERE id = $1; 