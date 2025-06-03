-- Re-create assignment status enum type
CREATE TYPE assignment_status AS ENUM ('pending', 'process', 'finish');

-- Add status column back to user_to_sub_task table
ALTER TABLE user_to_sub_task ADD COLUMN status assignment_status DEFAULT 'pending';

-- Add report column back to user_to_sub_task table
ALTER TABLE user_to_sub_task ADD COLUMN report TEXT; 