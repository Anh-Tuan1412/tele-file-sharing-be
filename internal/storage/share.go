package storage

import (
	"context"
	"errors"
	"log"

	"file-sharing/internal/model"
	"github.com/jmoiron/sqlx"
)

type ShareRepository interface {
	// RevokeShare cập nhật trạng thái revoked = true
	// Check ownerUserID để đảm bảo chính chủ mới được revoke
	RevokeShare(ctx context.Context, shareID int64, ownerUserID int64) error
	
    // GetShareByID lấy thông tin share (nếu cần dùng sau này)
    GetShareByID(ctx context.Context, shareID int64) (*model.Share, error)
}

type postgresShareRepository struct {
	db *sqlx.DB
}

func NewShareRepository(db *sqlx.DB) ShareRepository {
	return &postgresShareRepository{db: db}
}

func (r *postgresShareRepository) RevokeShare(ctx context.Context, shareID int64, ownerUserID int64) error {
	const query = `
		UPDATE shares
		SET revoked = true, updated_at = NOW()
		WHERE id = $1 AND owner_user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, shareID, ownerUserID)
	if err != nil {
		log.Printf("Failed to revoke share ID %d: %v", shareID, err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("share not found or access denied")
	}

	return nil
}

func (r *postgresShareRepository) GetShareByID(ctx context.Context, shareID int64) (*model.Share, error) {
    const query = `SELECT * FROM shares WHERE id = $1`
    var share model.Share
    err := r.db.GetContext(ctx, &share, query, shareID)
    return &share, err
}

//Các hàm truy xuất DB liên quan đến authorize password 
func (r *postgresShareRepository) GetPasswordHash(ctx context.Context, shareID int64) (string, error) {
	const query = `SELECT hash_password FROM shares WHERE id = $1`
	var passwordHash string
	err := r.db.GetContext(ctx, &passwordHash, query, shareID)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}
