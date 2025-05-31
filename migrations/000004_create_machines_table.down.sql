-- Drop indexes first
DROP INDEX IF EXISTS idx_machines_type_id;
DROP INDEX IF EXISTS idx_machines_machine_name;

-- Drop machines table
DROP TABLE IF EXISTS machines; 