-- User to sub task assignment table schema
CREATE TABLE user_to_sub_task (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sub_task_id INTEGER NOT NULL REFERENCES sub_tasks(id) ON DELETE CASCADE,
    report TEXT,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, sub_task_id)
);

-- Indexes for performance
CREATE INDEX idx_user_to_sub_task_user_id ON user_to_sub_task(user_id);
CREATE INDEX idx_user_to_sub_task_sub_task_id ON user_to_sub_task(sub_task_id);
