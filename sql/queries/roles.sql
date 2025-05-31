-- name: GetRoles :many
SELECT id, role_name, created_at, updated_at FROM roles
ORDER BY id;

-- name: GetRoleByID :one
SELECT id, role_name, created_at, updated_at FROM roles
WHERE id = $1;

-- name: GetRoleByName :one
SELECT id, role_name, created_at, updated_at FROM roles
WHERE role_name = $1; 