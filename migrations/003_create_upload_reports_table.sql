-- Create table for upload completion reports
CREATE TABLE IF NOT EXISTS upload_reports (
    id SERIAL PRIMARY KEY,
    file_id INT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    owner_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    report_type VARCHAR(50) NOT NULL,
    message TEXT,
    error_code VARCHAR(50),
    error_message TEXT,
    file_checksum TEXT,
    file_size_actual BIGINT,
    upload_duration_ms INT,
    bandwidth_kbps DECIMAL(10, 2),
    reported_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_upload_reports_file_id ON upload_reports(file_id);
CREATE INDEX IF NOT EXISTS idx_upload_reports_owner_user_id ON upload_reports(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_upload_reports_status ON upload_reports(status);
CREATE INDEX IF NOT EXISTS idx_upload_reports_reported_at ON upload_reports(reported_at DESC);

-- Add unique constraint to prevent duplicate reports for the same file
CREATE UNIQUE INDEX IF NOT EXISTS idx_upload_reports_file_id_unique ON upload_reports(file_id) WHERE status = 'completed';
