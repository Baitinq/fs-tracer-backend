package processor

import (
	"context"
	"log"

	"github.com/Baitinq/fs-tracer-backend/lib"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=mock_$GOFILE
type DB interface {
	InsertFiles(ctx context.Context, files []lib.File, user_id string) error
}

type DBImpl struct {
	db *sqlx.DB
}

var _ DB = (*DBImpl)(nil)

func NewDB(db *sqlx.DB) DB {
	return &DBImpl{db: db}
}

func (db DBImpl) InsertFiles(ctx context.Context, files []lib.File, user_id string) error {
	for i := range files {
		files[i].User_id = user_id
	}
	_, err := db.db.NamedExecContext(ctx, "INSERT INTO public.file (user_id, absolute_path, contents, timestamp) VALUES (:user_id, :absolute_path, :contents, :timestamp)", files)
	if err != nil {
		return err
	}
	log.Println("Inserted files: ", files)
	return nil
}
