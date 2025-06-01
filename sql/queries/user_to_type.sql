-- name: GetUserTaskTypes :many
SELECT ut.user_id, ut.type_id, t.type_name, t.description, t.created_at, t.updated_at, ut.created_at as assignment_created_at
FROM user_to_type ut
JOIN task_types t ON ut.type_id = t.id
WHERE ut.user_id = $1
ORDER BY t.type_name;

-- name: AddUserTaskType :exec
INSERT INTO user_to_type (user_id, type_id)
VALUES ($1, $2)
ON CONFLICT (user_id, type_id) DO NOTHING;

-- name: RemoveUserTaskType :exec
DELETE FROM user_to_type
WHERE user_id = $1 AND type_id = $2;

-- name: RemoveAllUserTaskTypes :exec
DELETE FROM user_to_type
WHERE user_id = $1;

-- name: RemoveAllTaskTypeUsers :exec
DELETE FROM user_to_type
WHERE type_id = $1;

-- name: CountTaskTypeUsers :one
SELECT COUNT(*) as user_count
FROM user_to_type
WHERE type_id = $1;

-- name: GetTaskTypeUsers :many
SELECT ut.user_id, ut.type_id, u.name, u.email, u.role_id, r.role_name, ut.created_at
FROM user_to_type ut
JOIN users u ON ut.user_id = u.id
JOIN roles r ON u.role_id = r.id
WHERE ut.type_id = $1
ORDER BY u.name; 