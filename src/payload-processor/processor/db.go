package processor

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=mock_$GOFILE
type DB interface {
	TestInsert(ctx context.Context, message string) error
}

type DBImpl struct {
	db *sqlx.DB
}

var _ DB = (*DBImpl)(nil)

func NewDB(db *sqlx.DB) DB {
	return &DBImpl{db: db}
}

func (db DBImpl) TestInsert(ctx context.Context, message string) error {
	_, err := db.db.ExecContext(ctx, "INSERT INTO test (created_at, test) VALUES ($1, $2)", time.Now(), message)
	return err
}
