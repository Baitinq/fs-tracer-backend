package handler

import (
	"context"

	"github.com/Baitinq/fs-tracer-backend/lib"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=mock_$GOFILE
type DB interface {
	GetLatestFileByPath(ctx context.Context, path string, user_id string) (*lib.File, error)
	GetUserIDByAPIKey(ctx context.Context, apiKey string) (string, error)
}

type DBImpl struct {
	db *sqlx.DB
}

var _ DB = (*DBImpl)(nil)

func NewDB(db *sqlx.DB) DB {
	return &DBImpl{db: db}
}

func (db DBImpl) GetLatestFileByPath(ctx context.Context, path string, user_id string) (*lib.File, error) {
	var file lib.File
	err := db.db.GetContext(ctx, &file, `
		SELECT * FROM private.file
		WHERE
			user_id = $1
			AND absolute_path = $2
		ORDER BY timestamp DESC
		LIMIT 1
	`, user_id, path)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// TODO: Add test
func (db DBImpl) GetUserIDByAPIKey(ctx context.Context, apiKey string) (string, error) {
	if len(apiKey) != 44 {
		return "", nil
	}

	var userID string
	err := db.db.GetContext(ctx, &userID, `
		SELECT id FROM private.api_keys
		WHERE api_key = $1
		LIMIT 1
	`, apiKey)
	if err != nil {
		return "", err
	}
	return userID, nil
}
