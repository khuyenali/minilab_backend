-- Create machines table
CREATE TABLE IF NOT EXISTS machines (
    id SERIAL PRIMARY KEY,
    machine_name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    estimate_time INTEGER NOT NULL, -- in minutes
    type_id INTEGER REFERENCES task_types(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_machines_machine_name ON machines(machine_name);
CREATE INDEX IF NOT EXISTS idx_machines_type_id ON machines(type_id); 