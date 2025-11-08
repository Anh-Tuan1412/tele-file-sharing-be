package storage

import (
	"context"
	"log"
	"github.com/jmoiron/sqlx"
	"file-sharing/internal/model"
)

type UserRepository interface {
	// UpsertByTelegram a user. If the user with the ginven Telegram ID exists,
	// update the username. Otherwise, create a new user.
	UpsertByTelegram(ctx context.Context, telegramID int64, username string) (*model.User, error)
}

// postgresUserRepository is the PostgreSQL implementation of UserRepository.
type postgresUserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) UpsertByTelegram(ctx context.Context, telegramID int64, username string) (*model.User, error) {
	const query = `
		INSERT INTO users (telegram_user_id, username)
		VALUES ($1, $2)
		ON CONFLICT (telegram_user_id) 
		DO UPDATE SET
			username = EXCLUDED.username,
			updated_at = NOW()
		RETURNING id, telegram_user_id, username, created_at, updated_at;
	`

	var user model.User
	err := r.db.GetContext(ctx, &user, query, telegramID, username)
	if err != nil {
		log.Printf("Failed to upsert user with telegramID %v: %v", telegramID, err)
		return nil, err
	}	
	
	return &user, nil
}