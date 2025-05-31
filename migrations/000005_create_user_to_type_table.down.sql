-- Drop indexes first
DROP INDEX IF EXISTS idx_user_to_type_type_id;
DROP INDEX IF EXISTS idx_user_to_type_user_id;

-- Drop user_to_type table
DROP TABLE IF EXISTS user_to_type; 