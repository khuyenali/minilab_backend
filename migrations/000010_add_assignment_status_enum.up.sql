-- Create assignment status enum type
CREATE TYPE assignment_status AS ENUM ('pending', 'process', 'finish');

-- Update user_to_sub_task table to use the enum
ALTER TABLE user_to_sub_task 
DROP CONSTRAINT IF EXISTS user_to_sub_task_status_check;

-- Remove the default temporarily
ALTER TABLE user_to_sub_task 
ALTER COLUMN status DROP DEFAULT;

-- Convert the column type
ALTER TABLE user_to_sub_task 
ALTER COLUMN status TYPE assignment_status USING status::assignment_status;

-- Add the default back
ALTER TABLE user_to_sub_task 
ALTER COLUMN status SET DEFAULT 'pending'::assignment_status; 