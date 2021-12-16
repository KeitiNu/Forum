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

func (p *PostModel) Insert(title, content, user, imagesrc string, category []string) (int, error) {
	stmt := `INSERT INTO posts (user, title, content, imagesrc, created)
			VALUES(?, ?, ?, ?, datetime('now'))`
	result, err := p.DB.Exec(stmt, user, title, content, imagesrc)
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

		_, err = p.DB.Exec(stmt, id, v)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}
	return int(id), nil
}

func (p *PostModel) Get(id int) (*Post, error) {

	stmt := `SELECT a.id, title, content, created, user, imagesrc, b.category_id From posts a
	LEFT JOIN post_category b ON a.id = b.post_id
    WHERE a.id = ?`

	rows, err := p.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	a := &Post{}
	for rows.Next() {
		s := &Post{}
		d := ""
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.User, &s.ImageSrc, &d)
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		} else if err != nil {
			fmt.Println(err)
			return nil, err
		}
		a.ID = s.ID
		a.Title = s.Title
		a.Content = s.Content
		a.Created = s.Created
		a.User = s.User
		if s.ImageSrc != "" {
			a.ImageSrc = s.ImageSrc[12:]
		}
		a.Category = append(a.Category, d)
	}
	return a, nil
}

func (p *PostModel) Latest(category string) ([]*Post, error) {
	stmt := `SELECT p.id, user, title, content, created FROM posts p
	LEFT JOIN post_category c ON p.id = c.post_id
	WHERE c.category_id = ?
    ORDER BY created DESC LIMIT 15`

	rows, err := p.DB.Query(stmt, category)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*Post{}

	for rows.Next() {
		s := &Post{}

		err := rows.Scan(&s.ID, &s.User, &s.Title, &s.Content, &s.Created)
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

func (p *PostModel) Update(title, content string, id int) error {
	stmt := `UPDATE posts SET title=?, content=?
			WHERE id = ?`
	_, err := p.DB.Exec(stmt, title, content, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (p *PostModel) Delete(id int) error {
	stmt := `DELETE FROM posts WHERE id = ?`
	_, err := p.DB.Exec(stmt, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
