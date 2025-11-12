package files

import (
	"context"
	"database/sql"
	"log"
	"time"

	"file-sharing/internal/model"
	"file-sharing/internal/storage"
)

// Service interface for file operations
type Service interface {
	// ReportUploadComplete creates and processes an upload completion report
	ReportUploadComplete(ctx context.Context, userID, fileID int64, req *ReportUploadCompleteRequest) (*ReportUploadCompleteResponse, error)

	// GetUploadReport retrieves the latest upload report for a file
	GetUploadReport(ctx context.Context, userID, fileID int64) (*model.UploadReport, error)

	// ListUploadReports lists all upload reports for a user with pagination
	ListUploadReports(ctx context.Context, userID int64, limit, offset int) ([]model.UploadReport, error)

	// GetUploadReportCount gets the total count of upload reports for a user
	GetUploadReportCount(ctx context.Context, userID int64) (int, error)
}

// ReportUploadCompleteRequest represents request data
type ReportUploadCompleteRequest struct {
	Status           string
	ReportType       string
	Message          string
	ErrorCode        string
	ErrorMessage     string
	FileChecksum     string
	FileSizeActual   int64
	UploadDurationMs int
	BandwidthKbps    float64
}

// ReportUploadCompleteResponse represents response data
type ReportUploadCompleteResponse struct {
	ReportID   int64
	FileID     int64
	Status     string
	ReportType string
	Message    string
	ReportedAt time.Time
}

// fileService implements Service interface
type fileService struct {
	db   *sql.DB
	repo storage.FileRepository
}

// NewFileService creates a new file service
func NewFileService(db *sql.DB, repo storage.FileRepository) Service {
	return &fileService{
		db:   db,
		repo: repo,
	}
}

// ReportUploadComplete creates and processes an upload completion report
func (s *fileService) ReportUploadComplete(
	ctx context.Context,
	userID, fileID int64,
	req *ReportUploadCompleteRequest,
) (*ReportUploadCompleteResponse, error) {

	// Verify file exists and belongs to user
	file, err := s.repo.GetFileByIDAndUserID(ctx, fileID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewFileNotFoundError()
		}
		log.Printf("Failed to verify file ownership: %v", err)
		return nil, NewDatabaseError()
	}

	if file == nil {
		return nil, NewFileNotFoundError()
	}

	// Create upload report
	report := &model.UploadReport{
		FileID:               fileID,
		OwnerUserID:          userID,
		Status:               req.Status,
		ReportType:           req.ReportType,
		Message:              req.Message,
		ErrorCode:            req.ErrorCode,
		ErrorMessage:         req.ErrorMessage,
		FileChecksum:         req.FileChecksum,
		FileSizeActual:       req.FileSizeActual,
		UploadDurationMs:     req.UploadDurationMs,
		BandwidthKbps:        req.BandwidthKbps,
	}

	createdReport, err := s.repo.CreateUploadReport(ctx, report)
	if err != nil {
		log.Printf("Failed to create upload report: %v", err)
		return nil, NewDatabaseError()
	}

	// Update file status
	newStatus := "completed"
	if req.Status == "failed" {
		newStatus = "failed"
	}

	err = s.repo.UpdateFileStatus(ctx, fileID, newStatus)
	if err != nil {
		log.Printf("Failed to update file status: %v", err)
		return nil, NewDatabaseError()
	}

	resp := &ReportUploadCompleteResponse{
		ReportID:   createdReport.ID,
		FileID:     fileID,
		Status:     req.Status,
		ReportType: req.ReportType,
		Message:    req.Message,
		ReportedAt: createdReport.ReportedAt,
	}

	return resp, nil
}

// GetUploadReport retrieves the latest upload report for a file
func (s *fileService) GetUploadReport(ctx context.Context, userID, fileID int64) (*model.UploadReport, error) {
	// Verify file exists and belongs to user
	file, err := s.repo.GetFileByIDAndUserID(ctx, fileID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewFileNotFoundError()
		}
		log.Printf("Failed to verify file ownership: %v", err)
		return nil, NewDatabaseError()
	}

	if file == nil {
		return nil, NewFileNotFoundError()
	}

	// Get upload report
	report, err := s.repo.GetUploadReport(ctx, fileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewUploadReportNotFoundError()
		}
		log.Printf("Failed to get upload report: %v", err)
		return nil, NewDatabaseError()
	}

	return report, nil
}

// ListUploadReports lists all upload reports for a user with pagination
func (s *fileService) ListUploadReports(ctx context.Context, userID int64, limit, offset int) ([]model.UploadReport, error) {
	reports, err := s.repo.GetUploadReportsByUserID(ctx, userID, limit, offset)
	if err != nil {
		log.Printf("Failed to list upload reports: %v", err)
		return nil, NewDatabaseError()
	}

	if reports == nil {
		reports = []model.UploadReport{}
	}

	return reports, nil
}

// GetUploadReportCount gets the total count of upload reports for a user
func (s *fileService) GetUploadReportCount(ctx context.Context, userID int64) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM upload_reports
		WHERE owner_user_id = $1
	`

	err := s.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		log.Printf("Failed to get upload report count: %v", err)
		return 0, NewDatabaseError()
	}

	return count, nil
}
