package lib

import "time"

type File struct {
	Id            string    `db:"id"`
	User_id       string    `db:"user_id"`
	Absolute_path string    `db:"absolute_path"`
	Contents      string    `db:"contents"`
	Timestamp     time.Time `db:"timestamp"`
}
