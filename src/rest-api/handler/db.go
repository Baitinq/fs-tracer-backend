package handler

import (
	"context"

	"github.com/Baitinq/fs-tracer-backend/lib"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=mock_$GOFILE
type DB interface {
	GetLatestFileByPath(ctx context.Context, path string) (*lib.File, error)
}

type DBImpl struct {
	db *sqlx.DB
}

var _ DB = (*DBImpl)(nil)

func NewDB(db *sqlx.DB) DB {
	return &DBImpl{db: db}
}

func (db DBImpl) GetLatestFileByPath(ctx context.Context, path string) (*lib.File, error) {
	var file lib.File
	err := db.db.GetContext(ctx, &file, `
		SELECT * FROM private.file
		WHERE absolute_path = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`, path)
	if err != nil {
		return nil, err
	}
	return &file, nil
}
