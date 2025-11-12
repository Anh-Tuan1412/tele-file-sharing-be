package model

import "time"

type UploadReport struct {
	ID                   int64     `json:"id" db:"id"`
	FileID               int64     `json:"file_id" db:"file_id"`
	OwnerUserID          int64     `json:"owner_user_id" db:"owner_user_id"`
	Status               string    `json:"status" db:"status"`
	ReportType           string    `json:"report_type" db:"report_type"`
	Message              string    `json:"message,omitempty" db:"message"`
	ErrorCode            string    `json:"error_code,omitempty" db:"error_code"`
	ErrorMessage         string    `json:"error_message,omitempty" db:"error_message"`
	FileChecksum         string    `json:"file_checksum,omitempty" db:"file_checksum"`
	FileSizeActual       int64     `json:"file_size_actual,omitempty" db:"file_size_actual"`
	UploadDurationMs     int       `json:"upload_duration_ms,omitempty" db:"upload_duration_ms"`
	BandwidthKbps        float64   `json:"bandwidth_kbps,omitempty" db:"bandwidth_kbps"`
	ReportedAt           time.Time `json:"reported_at" db:"reported_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// Status constants for upload reports
const (
	ReportStatusPending   = "pending"
	ReportStatusCompleted = "completed"
	ReportStatusFailed    = "failed"
	ReportStatusProcessing = "processing"
)

// ReportType constants
const (
	ReportTypeSuccess = "success"
	ReportTypeError   = "error"
	ReportTypeWarning = "warning"
)
