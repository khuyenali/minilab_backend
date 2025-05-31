-- name: GetUser :one
SELECT u.id, u.name, u.email, u.role_id, r.role_name, u.created_at, u.updated_at 
FROM users u
JOIN roles r ON u.role_id = r.id
WHERE u.id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT u.id, u.name, u.email, u.role_id, r.role_name, u.created_at, u.updated_at 
FROM users u
JOIN roles r ON u.role_id = r.id
ORDER BY u.created_at DESC;

-- name: CreateUser :one
INSERT INTO users (name, email, role_id) 
VALUES ($1, $2, $3) 
RETURNING id, name, email, role_id, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users 
SET name = $2, email = $3, role_id = COALESCE($4, role_id), updated_at = NOW() 
WHERE id = $1 
RETURNING id, name, email, role_id, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT u.id, u.name, u.email, u.role_id, r.role_name, u.created_at, u.updated_at 
FROM users u
JOIN roles r ON u.role_id = r.id
WHERE u.email = $1 LIMIT 1;

-- name: GetUsersByRole :many
SELECT u.id, u.name, u.email, u.role_id, r.role_name, u.created_at, u.updated_at 
FROM users u
JOIN roles r ON u.role_id = r.id
WHERE r.role_name = $1
ORDER BY u.created_at DESC; 