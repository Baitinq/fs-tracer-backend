package lib

import "time"

type File struct {
	Id            string    `db:"id" json:"id"`
	User_id       string    `db:"user_id" json:"user_id"`
	Absolute_path string    `db:"absolute_path" json:"absolute_path"`
	Contents      string    `db:"contents" json:"contents"`
	Timestamp     time.Time `db:"timestamp" json:"timestamp"`
}
