package share

import (
	"context"
	"file-sharing/internal/storage"
)

type Service interface {
	RevokeShare(ctx context.Context, shareID int64, userID int64) error
}

type shareService struct {
	repo storage.ShareRepository
}

func NewShareService(repo storage.ShareRepository) Service {
	return &shareService{
		repo: repo,
	}
}

func (s *shareService) RevokeShare(ctx context.Context, shareID int64, userID int64) error {
	return s.repo.RevokeShare(ctx, shareID, userID)
}