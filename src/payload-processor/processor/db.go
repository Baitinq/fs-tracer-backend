package processor

import (
	"context"

	"github.com/Baitinq/fs-tracer-backend/lib"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=mock_$GOFILE
type DB interface {
	InsertFile(ctx context.Context, file lib.File) error
}

type DBImpl struct {
	db *sqlx.DB
}

var _ DB = (*DBImpl)(nil)

func NewDB(db *sqlx.DB) DB {
	return &DBImpl{db: db}
}

func (db DBImpl) InsertFile(ctx context.Context, file lib.File) error {
	_, err := db.db.NamedExecContext(ctx, "INSERT INTO private.file (user_id, absolute_path, contents, timestamp) VALUES (:user_id, :absolute_path, :contents, :timestamp)", file)
	return err
}
