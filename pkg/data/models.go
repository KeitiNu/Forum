package data

import "database/sql"

// Encapsulate the code for working with sqlite3 in a separate package to the rest of our application.

// Storage for database data from tables
type Models struct {
	Users    UserModel
	Posts    PostModel
	Comments CommentsModel
}

// All database data has been divided into categories and is ready to be used. Actually it links the whole database to every category.
func NewModels(db *sql.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Posts:    PostModel{DB: db},
		Comments: CommentsModel{DB: db},
	}
}
