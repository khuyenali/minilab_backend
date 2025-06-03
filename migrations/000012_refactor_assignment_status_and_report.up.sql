-- Remove status column from user_to_sub_task table
ALTER TABLE user_to_sub_task DROP COLUMN IF EXISTS status;

-- Remove report column from user_to_sub_task table (moving it to tasks only)
ALTER TABLE user_to_sub_task DROP COLUMN IF EXISTS report;

-- Drop assignment status enum type as it's no longer needed
DROP TYPE IF EXISTS assignment_status; 