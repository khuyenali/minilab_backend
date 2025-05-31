CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add index on email for faster lookups
CREATE INDEX idx_users_email ON users(email);

-- Add index on created_at for sorting
CREATE INDEX idx_users_created_at ON users(created_at); 