package http

import (
	"errors"
	"fmt"
)

// ValidationError represents validation errors with field information
type ValidationError struct {
	Field   string
	Message string
}

// ValidationErrors is a collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError
}

func (ve *ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation passed"
	}
	msg := "validation errors:\n"
	for _, err := range ve.Errors {
		msg += fmt.Sprintf("  - %s: %s\n", err.Field, err.Message)
	}
	return msg
}

func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// ValidateReportUploadCompleteRequest validates the upload completion report request
func ValidateReportUploadCompleteRequest(req *ReportUploadCompleteRequest) error {
	validationErrors := &ValidationErrors{}

	// Validate status
	if req.Status == "" {
		validationErrors.Add("status", "status is required")
	} else if req.Status != "completed" && req.Status != "failed" {
		validationErrors.Add("status", "status must be 'completed' or 'failed'")
	}

	// Validate report_type
	if req.ReportType == "" {
		validationErrors.Add("report_type", "report_type is required")
	} else if req.ReportType != "success" && req.ReportType != "error" && req.ReportType != "warning" {
		validationErrors.Add("report_type", "report_type must be 'success', 'error', or 'warning'")
	}

	// Validate message length
	if len(req.Message) > 1000 {
		validationErrors.Add("message", "message must not exceed 1000 characters")
	}

	// Validate error_code length
	if len(req.ErrorCode) > 50 {
		validationErrors.Add("error_code", "error_code must not exceed 50 characters")
	}

	// Validate error_message length
	if len(req.ErrorMessage) > 500 {
		validationErrors.Add("error_message", "error_message must not exceed 500 characters")
	}

	// Validate file_checksum format and length
	if req.FileChecksum != "" && (len(req.FileChecksum) < 8 || len(req.FileChecksum) > 128) {
		validationErrors.Add("file_checksum", "file_checksum must be between 8 and 128 characters")
	}

	// Validate file_size_actual
	if req.FileSizeActual < 0 {
		validationErrors.Add("file_size_actual", "file_size_actual must be non-negative")
	}

	// Validate upload_duration_ms
	if req.UploadDurationMs < 0 {
		validationErrors.Add("upload_duration_ms", "upload_duration_ms must be non-negative")
	}

	// Validate bandwidth_kbps
	if req.BandwidthKbps < 0 {
		validationErrors.Add("bandwidth_kbps", "bandwidth_kbps must be non-negative")
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

// ValidateTelegramHeaders validates required Telegram headers
func ValidateTelegramHeaders(telegramID, username string) error {
	if telegramID == "" {
		return errors.New("X-Telegram-User-Id header is required")
	}
	if username == "" {
		return errors.New("X-Telegram-Username header is required")
	}
	if len(username) > 255 {
		return errors.New("X-Telegram-Username exceeds maximum length")
	}
	return nil
}

// ValidateFileID validates file ID format
func ValidateFileID(fileID int64) error {
	if fileID <= 0 {
		return errors.New("file_id must be a positive integer")
	}
	return nil
}

// ValidatePaginationParams validates pagination parameters
func ValidatePaginationParams(limit, offset int) (int, int, error) {
	// Default limit
	if limit == 0 {
		limit = 20
	}

	// Validate limit
	if limit < 1 {
		return 0, 0, errors.New("limit must be at least 1")
	}
	if limit > 100 {
		limit = 100
	}

	// Validate offset
	if offset < 0 {
		offset = 0
	}
	if offset > 1000000 { // Prevent extremely large offsets
		offset = 1000000
	}

	return limit, offset, nil
}
