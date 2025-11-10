CREATE TABLE IF NOT EXISTS files (
    id SERIAL PRIMARY KEY,
    owner_user_id INT NOT NULL REFERENCES users(id),
    object_key TEXT NOT NULL,
    filename TEXT NOT NULL,
    size BIGINT,
    mime TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_files_owner_user_id ON files(owner_user_id);
