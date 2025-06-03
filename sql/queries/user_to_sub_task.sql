-- name: CreateUserSubTaskAssignment :one
INSERT INTO user_to_sub_task (user_id, sub_task_id)
VALUES ($1, $2)
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
WHERE uts.user_id = $1
ORDER BY uts.assigned_at DESC;

-- name: UpdateUserSubTaskAssignment :one
UPDATE user_to_sub_task
SET updated_at = CURRENT_TIMESTAMP
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

-- name: DeleteAssignmentsForDraftTasks :exec
DELETE FROM user_to_sub_task
USING sub_tasks, tasks
WHERE user_to_sub_task.sub_task_id = sub_tasks.id 
  AND sub_tasks.task_id = tasks.id 
  AND tasks.status = 'draft';

-- name: CheckUserHasActiveTasks :one
SELECT EXISTS(
    SELECT 1 FROM user_to_sub_task uts
    JOIN sub_tasks st ON uts.sub_task_id = st.id  
    JOIN tasks t ON st.task_id = t.id
    WHERE uts.user_id = $1 AND t.status IN ('draft', 'pending', 'processing')
) AS has_active_tasks; 