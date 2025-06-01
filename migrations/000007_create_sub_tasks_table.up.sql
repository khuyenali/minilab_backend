-- Sub task table schema
CREATE TABLE sub_tasks (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    type_id INTEGER NOT NULL REFERENCES task_types(id) ON DELETE CASCADE,
    sub_task_name VARCHAR(255),
    description TEXT,
    estimate_effort INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_sub_tasks_task_id ON sub_tasks(task_id);
CREATE INDEX idx_sub_tasks_type_id ON sub_tasks(type_id);
