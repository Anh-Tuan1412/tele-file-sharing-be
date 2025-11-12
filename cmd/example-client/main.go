package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Config holds the API configuration
type Config struct {
	BaseURL     string
	TelegramID  string
	Username    string
	HTTPClient  *http.Client
}

// ReportRequest represents the upload completion report request
type ReportRequest struct {
	Status           string  `json:"status"`
	ReportType       string  `json:"report_type"`
	Message          string  `json:"message,omitempty"`
	ErrorCode        string  `json:"error_code,omitempty"`
	ErrorMessage     string  `json:"error_message,omitempty"`
	FileChecksum     string  `json:"file_checksum,omitempty"`
	FileSizeActual   int64   `json:"file_size_actual,omitempty"`
	UploadDurationMs int     `json:"upload_duration_ms,omitempty"`
	BandwidthKbps    float64 `json:"bandwidth_kbps,omitempty"`
}

// ReportResponse represents the API response
type ReportResponse struct {
	ReportID   int64     `json:"report_id"`
	FileID     int64     `json:"file_id"`
	Status     string    `json:"status"`
	ReportType string    `json:"report_type"`
	Message    string    `json:"message"`
	ReportedAt time.Time `json:"reported_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewConfig creates a new API configuration
func NewConfig(baseURL, telegramID, username string) *Config {
	return &Config{
		BaseURL:    baseURL,
		TelegramID: telegramID,
		Username:   username,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// reportUploadComplete sends an upload completion report
func (c *Config) reportUploadComplete(fileID int64, req *ReportRequest) (*ReportResponse, error) {
	url := fmt.Sprintf("%s/v1/files/%d/report-complete", c.BaseURL, fileID)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Telegram-User-Id", c.TelegramID)
	httpReq.Header.Set("X-Telegram-Username", c.Username)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("request failed: %s - %s", errResp.Code, errResp.Message)
	}

	var reportResp ReportResponse
	if err := json.Unmarshal(respBody, &reportResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &reportResp, nil
}

// getUploadReport retrieves the latest upload report for a file
func (c *Config) getUploadReport(fileID int64) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v1/files/%d/report", c.BaseURL, fileID)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("X-Telegram-User-Id", c.TelegramID)
	httpReq.Header.Set("X-Telegram-Username", c.Username)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("request failed: %s - %s", errResp.Code, errResp.Message)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(respBody, &report); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return report, nil
}

// listUploadReports lists all upload reports for the current user
func (c *Config) listUploadReports(limit, offset int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v1/upload-reports?limit=%d&offset=%d", c.BaseURL, limit, offset)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("X-Telegram-User-Id", c.TelegramID)
	httpReq.Header.Set("X-Telegram-Username", c.Username)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("request failed: %s - %s", errResp.Code, errResp.Message)
	}

	var reports []map[string]interface{}
	if err := json.Unmarshal(respBody, &reports); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return reports, nil
}

func main() {
	// Command-line flags
	baseURL := flag.String("url", "http://localhost:8080", "Base URL of the API")
	telegramID := flag.String("telegram-id", "123456789", "Telegram user ID")
	username := flag.String("username", "testuser", "Telegram username")
	command := flag.String("cmd", "report", "Command: report, get, list")
	fileID := flag.Int64("file-id", 1, "File ID")
	status := flag.String("status", "completed", "Report status: completed, failed")
	reportType := flag.String("type", "success", "Report type: success, error, warning")
	message := flag.String("message", "Upload completed successfully", "Report message")
	fileSize := flag.Int64("file-size", 1024, "Actual file size in bytes")
	duration := flag.Int("duration", 5000, "Upload duration in milliseconds")
	bandwidth := flag.Float64("bandwidth", 512.5, "Upload bandwidth in Kbps")

	flag.Parse()

	// Create API config
	config := NewConfig(*baseURL, *telegramID, *username)

	switch *command {
	case "report":
		fmt.Println("📤 Reporting upload completion...")
		req := &ReportRequest{
			Status:           *status,
			ReportType:       *reportType,
			Message:          *message,
			FileSizeActual:   *fileSize,
			UploadDurationMs: *duration,
			BandwidthKbps:    *bandwidth,
		}

		resp, err := config.reportUploadComplete(*fileID, req)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Printf("✅ Report created successfully!\n")
		fmt.Printf("Report ID: %d\n", resp.ReportID)
		fmt.Printf("File ID: %d\n", resp.FileID)
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Report Type: %s\n", resp.ReportType)
		fmt.Printf("Message: %s\n", resp.Message)
		fmt.Printf("Reported At: %s\n", resp.ReportedAt)

	case "get":
		fmt.Println("📖 Retrieving upload report...")
		report, err := config.getUploadReport(*fileID)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Printf("✅ Report retrieved successfully!\n")
		respJSON, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(respJSON))

	case "list":
		fmt.Println("📋 Listing all upload reports...")
		reports, err := config.listUploadReports(20, 0)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		if len(reports) == 0 {
			fmt.Println("No reports found")
			return
		}

		fmt.Printf("✅ Found %d reports:\n", len(reports))
		for i, report := range reports {
			fmt.Printf("\n[%d] Report:\n", i+1)
			respJSON, _ := json.MarshalIndent(report, "    ", "  ")
			fmt.Println(string(respJSON))
		}

	default:
		fmt.Printf("❌ Unknown command: %s\n", *command)
		fmt.Println("Available commands: report, get, list")
	}
}
