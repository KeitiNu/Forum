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

func (s *PostModel) Insert(title, content, user string) error {
	stmt := `INSERT INTO posts (user, title, content, created)
			VALUES(?, ?, ?,datetime('now'))`
	_, err := s.DB.Exec(stmt, user, title, content)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostModel) Latest(category string) ([]*Post, error) {
	stmt := `SELECT user, title, content, created FROM posts
    ORDER BY created DESC LIMIT 15`

	rows, err := m.DB.Query(stmt, category)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*Post{}

	for rows.Next() {
		s := &Post{}

		err := rows.Scan(&s.Title, &s.Content, &s.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
