-- Drop indexes first
DROP INDEX IF EXISTS idx_users_role_id;
DROP INDEX IF EXISTS idx_users_email;

-- Drop users table
DROP TABLE IF EXISTS users; 