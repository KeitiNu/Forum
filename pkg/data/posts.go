package data

import (
	"database/sql"
	"fmt"
	"time"
)

// Fields we extract from database to use and render on webpages
type Post struct {
	ID       int
	Title    string
	Content  string
	Likes    int
	Dislikes int
	Created  time.Time
	ImageSrc string
	User     string
	Category []string
}

type PostModel struct {
	DB *sql.DB
}

func (s *PostModel) Insert(title, content, user, imagesrc string, category []string) (int, error) {
	stmt := `INSERT INTO posts (user, title, content, imagesrc, created)
			VALUES(?, ?, ?, ?, datetime('now'))`
	result, err := s.DB.Exec(stmt, user, title, content, imagesrc)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("INSERT", err)
	}

	for _, v := range category {
		stmt = `INSERT INTO post_category (post_id, category_id)
			VALUES(?,?)`

		_, err = s.DB.Exec(stmt, id, v)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}
	return int(id), nil
}

func (m *PostModel) Get(id int) (*Post, error) {

	stmt := `SELECT title, content, created, user, imagesrc, b.category_id From posts a
	LEFT JOIN post_category b ON a.id = b.post_id
    WHERE a.id = ?`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	a := &Post{}
	for rows.Next() {
		s := &Post{}
		d := ""
		err := rows.Scan(&s.Title, &s.Content, &s.Created, &s.User, &s.ImageSrc, &d)
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		} else if err != nil {
			fmt.Println(err)
			return nil, err
		}
		a.Title = s.Title
		a.Content = s.Content
		a.Created = s.Created
		a.User = s.User
		if s.ImageSrc != "" {
			a.ImageSrc = s.ImageSrc[12:]
		}
		a.Category = append(a.Category, d)
	}
	fmt.Println(a)
	return a, nil
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
