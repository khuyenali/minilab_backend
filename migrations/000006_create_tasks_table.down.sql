-- Drop indexes
DROP INDEX IF EXISTS idx_tasks_created_at;
DROP INDEX IF EXISTS idx_tasks_priority;
DROP INDEX IF EXISTS idx_tasks_status;

-- Drop table
DROP TABLE IF EXISTS tasks;

-- Drop enum type
DROP TYPE IF EXISTS task_status;
