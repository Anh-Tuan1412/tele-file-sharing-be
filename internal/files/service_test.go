package files

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"file-sharing/internal/model"
	"file-sharing/internal/storage"
)

// MockFileRepository is a mock implementation of FileRepository
type MockFileRepository struct {
	GetFileByIDAndUserIDFunc func(ctx context.Context, fileID, userID int64) (*model.FileWithOwner, error)
	CreateUploadReportFunc   func(ctx context.Context, report *model.UploadReport) (*model.UploadReport, error)
	UpdateFileStatusFunc     func(ctx context.Context, fileID int64, status string) error
	GetUploadReportFunc      func(ctx context.Context, fileID int64) (*model.UploadReport, error)
	GetUploadReportsByUserIDFunc func(ctx context.Context, userID int64, limit, offset int) ([]model.UploadReport, error)
}

func (m *MockFileRepository) GetFileByID(ctx context.Context, fileID int64) (*model.FileWithOwner, error) {
	return nil, nil
}

func (m *MockFileRepository) GetFileByIDAndUserID(ctx context.Context, fileID, userID int64) (*model.FileWithOwner, error) {
	if m.GetFileByIDAndUserIDFunc != nil {
		return m.GetFileByIDAndUserIDFunc(ctx, fileID, userID)
	}
	return nil, sql.ErrNoRows
}

func (m *MockFileRepository) UpdateFileStatus(ctx context.Context, fileID int64, status string) error {
	if m.UpdateFileStatusFunc != nil {
		return m.UpdateFileStatusFunc(ctx, fileID, status)
	}
	return nil
}

func (m *MockFileRepository) CreateUploadReport(ctx context.Context, report *model.UploadReport) (*model.UploadReport, error) {
	if m.CreateUploadReportFunc != nil {
		return m.CreateUploadReportFunc(ctx, report)
	}
	report.ID = 1
	report.ReportedAt = time.Now()
	report.UpdatedAt = time.Now()
	return report, nil
}

func (m *MockFileRepository) GetUploadReport(ctx context.Context, fileID int64) (*model.UploadReport, error) {
	if m.GetUploadReportFunc != nil {
		return m.GetUploadReportFunc(ctx, fileID)
	}
	return nil, sql.ErrNoRows
}

func (m *MockFileRepository) GetUploadReportsByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.UploadReport, error) {
	if m.GetUploadReportsByUserIDFunc != nil {
		return m.GetUploadReportsByUserIDFunc(ctx, userID, limit, offset)
	}
	return []model.UploadReport{}, nil
}

func (m *MockFileRepository) UpdateUploadReportStatus(ctx context.Context, reportID int64, status string) error {
	return nil
}

// TestReportUploadComplete tests the ReportUploadComplete method
func TestReportUploadComplete_Success(t *testing.T) {
	mockRepo := &MockFileRepository{
		GetFileByIDAndUserIDFunc: func(ctx context.Context, fileID, userID int64) (*model.FileWithOwner, error) {
			return &model.FileWithOwner{
				ID:          fileID,
				OwnerUserID: userID,
				Filename:    "test.txt",
				Size:        1024,
				Status:      "pending",
			}, nil
		},
		CreateUploadReportFunc: func(ctx context.Context, report *model.UploadReport) (*model.UploadReport, error) {
			report.ID = 1
			report.ReportedAt = time.Now()
			report.UpdatedAt = time.Now()
			return report, nil
		},
		UpdateFileStatusFunc: func(ctx context.Context, fileID int64, status string) error {
			return nil
		},
	}

	service := NewFileService(nil, mockRepo)

	req := &ReportUploadCompleteRequest{
		Status:           "completed",
		ReportType:       "success",
		Message:          "Upload successful",
		FileSizeActual:   1024,
		UploadDurationMs: 5000,
		BandwidthKbps:    512.5,
	}

	resp, err := service.ReportUploadComplete(context.Background(), 1, 1, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", resp.Status)
	}

	if resp.ReportType != "success" {
		t.Errorf("Expected report_type 'success', got %s", resp.ReportType)
	}
}

// TestReportUploadComplete_FileNotFound tests when file is not found
func TestReportUploadComplete_FileNotFound(t *testing.T) {
	mockRepo := &MockFileRepository{
		GetFileByIDAndUserIDFunc: func(ctx context.Context, fileID, userID int64) (*model.FileWithOwner, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := NewFileService(nil, mockRepo)

	req := &ReportUploadCompleteRequest{
		Status:     "completed",
		ReportType: "success",
	}

	resp, err := service.ReportUploadComplete(context.Background(), 1, 999, req)

	if resp != nil {
		t.Fatal("Expected nil response")
	}

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err != ErrFileNotFound {
		t.Errorf("Expected ErrFileNotFound, got %v", err)
	}
}

// TestGetUploadReport_Success tests retrieving an upload report
func TestGetUploadReport_Success(t *testing.T) {
	now := time.Now()
	mockRepo := &MockFileRepository{
		GetFileByIDAndUserIDFunc: func(ctx context.Context, fileID, userID int64) (*model.FileWithOwner, error) {
			return &model.FileWithOwner{
				ID:          fileID,
				OwnerUserID: userID,
				Filename:    "test.txt",
			}, nil
		},
		GetUploadReportFunc: func(ctx context.Context, fileID int64) (*model.UploadReport, error) {
			return &model.UploadReport{
				ID:       1,
				FileID:   fileID,
				Status:   "completed",
				ReportedAt: now,
			}, nil
		},
	}

	service := NewFileService(nil, mockRepo)

	report, err := service.GetUploadReport(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if report == nil {
		t.Fatal("Expected report, got nil")
	}

	if report.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", report.Status)
	}
}

// TestListUploadReports tests listing upload reports
func TestListUploadReports_Success(t *testing.T) {
	mockRepo := &MockFileRepository{
		GetUploadReportsByUserIDFunc: func(ctx context.Context, userID int64, limit, offset int) ([]model.UploadReport, error) {
			return []model.UploadReport{
				{
					ID:       1,
					FileID:   1,
					Status:   "completed",
					ReportType: "success",
				},
				{
					ID:       2,
					FileID:   2,
					Status:   "completed",
					ReportType: "success",
				},
			}, nil
		},
	}

	service := NewFileService(nil, mockRepo)

	reports, err := service.ListUploadReports(context.Background(), 1, 20, 0)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(reports) != 2 {
		t.Errorf("Expected 2 reports, got %d", len(reports))
	}
}
