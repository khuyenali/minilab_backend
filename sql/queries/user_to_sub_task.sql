-- name: CreateUserSubTaskAssignment :one
INSERT INTO user_to_sub_task (user_id, sub_task_id, report, status)
VALUES ($1, $2, $3, $4)
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

-- name: GetActiveAssignmentsByUserID :many
SELECT uts.*, st.type_id FROM user_to_sub_task uts
JOIN sub_tasks st ON uts.sub_task_id = st.id
WHERE uts.user_id = $1 AND uts.status IN ('pending', 'process')
ORDER BY uts.assigned_at DESC;

-- name: UpdateUserSubTaskAssignment :one
UPDATE user_to_sub_task
SET report = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateAssignmentStatus :one
UPDATE user_to_sub_task
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateAssignmentStatusAndReport :one
UPDATE user_to_sub_task
SET status = $2, report = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetAssignmentsByTaskID :many
SELECT uts.* FROM user_to_sub_task uts
JOIN sub_tasks st ON uts.sub_task_id = st.id
WHERE st.task_id = $1
ORDER BY uts.assigned_at ASC;

-- name: DeleteUserSubTaskAssignment :exec
DELETE FROM user_to_sub_task
WHERE id = $1;

-- name: DeleteUserSubTaskAssignmentByUserAndSubTask :exec
DELETE FROM user_to_sub_task
WHERE user_id = $1 AND sub_task_id = $2; 