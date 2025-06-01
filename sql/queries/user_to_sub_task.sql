-- name: CreateUserSubTaskAssignment :one
INSERT INTO user_to_sub_task (user_id, sub_task_id, report)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserSubTaskAssignment :one
SELECT * FROM user_to_sub_task
WHERE id = $1;

-- name: GetAssignmentsBySubTaskID :many
SELECT * FROM user_to_sub_task
WHERE sub_task_id = $1
ORDER BY assigned_at ASC;

-- name: GetAssignmentsByUserID :many
SELECT * FROM user_to_sub_task
WHERE user_id = $1
ORDER BY assigned_at DESC;

-- name: UpdateUserSubTaskAssignment :one
UPDATE user_to_sub_task
SET report = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUserSubTaskAssignment :exec
DELETE FROM user_to_sub_task
WHERE id = $1;

-- name: DeleteUserSubTaskAssignmentByUserAndSubTask :exec
DELETE FROM user_to_sub_task
WHERE user_id = $1 AND sub_task_id = $2; 