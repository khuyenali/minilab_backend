-- Create user_to_type many-to-many table
CREATE TABLE IF NOT EXISTS user_to_type (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type_id INTEGER NOT NULL REFERENCES task_types(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, type_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_to_type_user_id ON user_to_type(user_id);
CREATE INDEX IF NOT EXISTS idx_user_to_type_type_id ON user_to_type(type_id); 