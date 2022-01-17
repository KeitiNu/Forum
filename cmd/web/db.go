package main

import "database/sql"

func Migrate(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS user_type(
		id integer PRIMARY KEY NOT NULL,
		name text UNIQUE NOT NULL,
		description text
		);

		CREATE TABLE IF NOT EXISTS categories (
			title TEXT PRIMARY KEY NOT NULL UNIQUE,
			description TEXT NOT NULL,
			created DATETIME
		);

		CREATE TABLE IF NOT EXISTS users (
			username PRIMARY KEY NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			hashed_password BLOB NOT NULL,
			token BLOB NOT NULL,
			created DATETIME,
			updated DATETIME
		);

		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			votes INTEGER,
			created DATETIME,
			imagesrc TEXT,
			user TEXT NOT NULL,
			FOREIGN KEY (user) REFERENCES users (username)
		);

		CREATE TABLE IF NOT EXISTS post_category (
			id INTEGER PRIMARY KEY,
			post_id INTEGER,
			category_id TEXT,
			FOREIGN KEY (category_id) REFERENCES categories (title)
			FOREIGN KEY (post_id) REFERENCES posts (id)
		);


		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY,
			post_id INTEGER NOT NULL,
			user_id TEXT NOT NULL,
			content TEXT NOT NULL,
			votes INTEGER,
			created DATETIME,
			FOREIGN KEY (user_id) REFERENCES users (username)
			FOREIGN KEY (post_id) REFERENCES posts (id)
		);


		CREATE TABLE IF NOT EXISTS vote(
			id integer PRIMARY KEY NOT NULL,
			post_id integer DEFAULT 0,
			comment_id integer DEFAULT 0,
			type bool,
			created datetime NOT NULL,
			user_id text NOT NULL,
			FOREIGN KEY (user_id)
				REFERENCES users(username)
			FOREIGN KEY(comment_id)
				REFERENCES comment(id)
						ON DELETE CASCADE
			FOREIGN KEY(post_id)
				REFERENCES post(id)
						ON DELETE CASCADE
			);

`

	_, err := db.Exec(stmt)
	return err
}
