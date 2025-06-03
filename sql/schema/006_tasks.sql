-- Task table schema
CREATE TYPE task_status AS ENUM ('draft', 'pending', 'processing', 'finish');

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    task_name VARCHAR(255) NOT NULL,
    status task_status NOT NULL DEFAULT 'draft',
    priority INTEGER NOT NULL DEFAULT 2,  -- 1=high, 2=medium, 3=low
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    note TEXT,
    report TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_created_at ON tasks(created_at); 