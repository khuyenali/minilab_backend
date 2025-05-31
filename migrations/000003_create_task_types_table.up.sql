-- Create task_types table
CREATE TABLE IF NOT EXISTS task_types (
    id SERIAL PRIMARY KEY,
    type_name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_task_types_type_name ON task_types(type_name); 