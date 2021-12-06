package data

import (
	"database/sql"
	"time"
)

// Fields we extract from database to use and render on webpages
type Post struct {
	ID       int
	Category string
	User     string
	Title    string
	Content  string
	Created  time.Time
	Votes    int
	Comments int
}

type PostModel struct {
	DB *sql.DB
}
