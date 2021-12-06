package data

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID      int
	PostID  int
	User    string
	Content string
	Created time.Time
	Votes   int
}

type CommentsModel struct {
	DB *sql.DB
}
