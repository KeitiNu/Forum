package data

import (
	"database/sql"
	"time"
)
// Fields we extract from database to use and render on webpages
type User struct {
	ID              string
	Username        string
	Hashed_password string
	CreatedAt       time.Time
}

type UserModel struct {
	DB *sql.DB
}
