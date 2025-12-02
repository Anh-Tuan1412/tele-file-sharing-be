package share

import (
	"context"
	"database/sql"
	"errors"
	"file-sharing/internal/model"
	"file-sharing/internal/storage"
	"log"
	"time"
	"crypto/rand"
    "encoding/hex"
    "golang.org/x/crypto/bcrypt"
)

type Service interface {
	RevokeShare(ctx context.Context, shareID int64, userID int64) error

	ListShares(ctx context.Context, userID int64, limit, offset int) ([]model.Share, error)

	// Tạo presigned URL cho download
	CreatePresignedURL(ctx context.Context, shareID int64, requesterUserID int64, expirySeconds int) (string, error)

	GetMetadata(ctx context.Context, id int64) (*model.ShareMetadataResponseDTO, error)
	CreateShare(ctx context.Context, userID int64, req *CreateShareRequest) (*model.Share, error)
}

// đối tượng DTO mới để nhận request từ FE 
// nên cải tiến lại trong docs cho giống với cấu trúc database hiện tại 
// là không có việc nhập tên các người được phép tải recipients và from_ts (lấy luôn thời điểm lúc mới tạo xong) - ứng với 
// trường created_at 
type CreateShareRequest struct {
	FileID    int64      `json:"file_id"`
    Password  string     `json:"password"`
    ExpiresAt *time.Time `json:"expires_at"`
}

type shareService struct {
	repo storage.ShareRepository
}

func NewShareService(repo storage.ShareRepository) Service {
	return &shareService{
		repo: repo,
	}
}

// Hàm sinh  chuỗi hash ngẫu nhiên 
func generateRandomHash(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func (s *shareService) RevokeShare(ctx context.Context, shareID int64, userID int64) error {
	return s.repo.RevokeShare(ctx, shareID, userID)
}

// ListShares lists all shares for a user
func (s *shareService) ListShares(ctx context.Context, userID int64, limit, offset int) ([]model.Share, error) {
	reports, err := s.repo.ListSharesByOwnerUserID(ctx, userID, limit, offset)
	if err != nil {
		log.Printf("Failed to list shares in service for user %d: %v", userID, err)
		return nil, err
	}

	if reports == nil {
		reports = []model.Share{} // Trả về mảng rỗng thay vì null
	}

	return reports, nil
}

// prepareDownload
// kiểm tra điều kiện: revoked/expired, check và cập nhật số lần download
// trả về: objectKey và filename
func (s *shareService) prepareDownload(ctx context.Context, shareID int64, requesterUserID int64) (string, string, error) {
	// load share meta
	shareRec, err := s.repo.GetShareByID(ctx, shareID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrShareNotFoundOrAccessDenied
		}
		return "", "", err
	}

	// check expired/revoked
	if shareRec != nil {
		if shareRec.Revoked {
			return "", "", ErrShareRevokedOrExpired
		}
		if shareRec.ExpiresAt != nil && shareRec.ExpiresAt.Before(time.Now()) {
			return "", "", ErrShareRevokedOrExpired
		}
	}

	// lấy object key + filename
	objectKey, filename, err := s.repo.GetActiveFilePathByShareID(ctx, shareID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrShareNotFoundOrAccessDenied
		}
		return "", "", err
	}

	// check + increment download counter (Hiện là no-op)
	err = s.repo.CheckAndIncrementDownload(ctx, shareID)
	if err != nil {
		return "", "", ErrMaxDownloadsExceeded
	}

	return objectKey, filename, nil
}

func (s *shareService) CreatePresignedURL(ctx context.Context, shareID int64, requesterUserID int64, expirySeconds int) (string, error) {
	// prepareDownload kiểm tra các điều kiện và trả về object key + filename + lỗi
	objectKey, _, err := s.prepareDownload(ctx, shareID, requesterUserID)
	if err != nil {
		return "", err
	}

	// presign via repo (MinIO)
	url, err := s.repo.PresignObject(ctx, objectKey, expirySeconds)
	if err != nil {
		return "", err
	}
	return url, nil
}

var (
	// trả về khi share không tồn tại hoặc không được phép truy cập
	ErrShareNotFoundOrAccessDenied = errors.New("share not found or access denied")
	// trả về khi share đã bị revoke hoặc đã hết hạn
	ErrShareRevokedOrExpired = errors.New("share revoked or expired")
	// trả khi đã đạt max_downloads
	ErrMaxDownloadsExceeded = errors.New("max downloads exceeded")
)

func (s *shareService) GetMetadata(ctx context.Context, id int64) (*model.ShareMetadataResponseDTO, error) {
	md, err := s.repo.GetShareMetadata(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.ShareMetadataResponseDTO{
		ID:        md.Share.ID,
		Hash:      md.Share.Hash,
		Revoked:   md.Share.Revoked,
		ExpiresAt: md.Share.ExpiresAt,
		CreatedAt: md.Share.CreatedAt,
		File: model.FileMetadataDTO{
			ID:        md.File.ID,
			Filename:  md.File.Filename,
			ObjectKey: md.File.ObjectKey,
			Size:      md.File.Size,
			Mime:      md.File.Mime,
			Status:    md.File.Status,
		},
		Owner: model.UserMinimalDTO{
			ID:         md.User.ID,
			Username:   md.User.Username,
			TelegramID: md.User.TelegramID,
		},
	}, nil
}

func (s *shareService) CreateShare(ctx context.Context, userID int64, req *CreateShareRequest) (*model.Share, error) {
    // Tạo chuỗi hash unique cho link (ví dụ 8 bytes -> 16 ký tự hex)
    linkHash, err := generateRandomHash(8) 
    if err != nil {
        return nil, err
    }

    // Xử lý mật khẩu (nếu có)
    var passwordHash string
    requirePassword := false
    
    if req.Password != "" {
        requirePassword = true
        // Hash mật khẩu bằng bcrypt
        bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
        if err != nil {
            return nil, err
        }
        passwordHash = string(bytes)
    }

    // Tạo model share
    newShare := &model.Share{
        FileID:          req.FileID,
        OwnerUserID:     userID,
        Hash:            linkHash,
        RequirePassword: requirePassword,
        HashPassword:    passwordHash,
        ExpiresAt:       req.ExpiresAt,
        Revoked:         false,
    }

    // Gọi Repo để lưu
    return s.repo.CreateShare(ctx, newShare)
}