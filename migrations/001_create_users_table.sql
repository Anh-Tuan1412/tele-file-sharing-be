--- migrations/001_create_users_table.sql

--- table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegrame_user_id BIGINT NOT NULL UNIQUE,
    username varchar(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()   
);

--- index
CREATE INDEX IF NOT EXISTS idx_users_telegrame_user_id ON users(telegrame_user_id);
