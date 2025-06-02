-- Revert user_to_sub_task table to use VARCHAR with CHECK constraint
ALTER TABLE user_to_sub_task 
ALTER COLUMN status TYPE VARCHAR(20) USING status::text;

ALTER TABLE user_to_sub_task 
ADD CONSTRAINT user_to_sub_task_status_check CHECK (status IN ('pending', 'process', 'finish'));

-- Drop assignment status enum type
DROP TYPE assignment_status; 