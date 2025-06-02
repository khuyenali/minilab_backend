-- Add status column to user_to_sub_task table
ALTER TABLE user_to_sub_task 
ADD COLUMN status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'process', 'finish'));

-- Update existing records to have pending status
UPDATE user_to_sub_task SET status = 'pending' WHERE status IS NULL; 