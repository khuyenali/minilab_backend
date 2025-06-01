-- Drop indexes
DROP INDEX IF EXISTS idx_user_to_sub_task_sub_task_id;
DROP INDEX IF EXISTS idx_user_to_sub_task_user_id;

-- Drop table
DROP TABLE IF EXISTS user_to_sub_task;
