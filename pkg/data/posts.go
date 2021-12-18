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
	Votes    int
	Created  time.Time
	ImageSrc string
	User     string
	Category []string
}

type PostModel struct {
	DB *sql.DB
}

func (p *PostModel) Insert(title, content, user, imagesrc string, category []string) (int, error) {
	stmt := `INSERT INTO posts (user, title, content, imagesrc, created, votes)
			VALUES(?, ?, ?, ?, datetime('now'), 0)`
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
	stmt := `SELECT p.id, user, title, content, created, votes FROM posts p
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

		err := rows.Scan(&s.ID, &s.User, &s.Title, &s.Content, &s.Created, &s.Votes)
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

func (p *PostModel) AddVote(id, vote, username string) error {
	stmt := ``
	i, err := p.GetVote(id, vote, username)
	if err != nil {
		if err == sql.ErrNoRows {
			switch vote {
			case "up":
				stmt = `INSERT INTO vote (type, post_id, created, user_id) VALUES
										(?, ?, datetime('now'), ?)`
				_, err := p.DB.Exec(stmt, true, id, username)
				if err != nil {
					fmt.Println(err, "1")
					return err
				}
				return nil
			case "down":
				stmt = `INSERT INTO vote (type, post_id, created, user_id) VALUES
										(?, ?, datetime('now'), ?)`
				_, err := p.DB.Exec(stmt, false, id, username)
				if err != nil {
					fmt.Println(err, "2")
					return err
				}
				return nil
			}
		}
	}
	stmt = `DELETE FROM vote WHERE post_id = ? AND user_id = ?`
	_, err = p.DB.Exec(stmt, id, username)
	if err != nil {
		fmt.Println("ERROR")
		return err
	}
	switch vote {
	case "up":
		if i == "0" {
			stmt = `INSERT INTO vote (type, post_id, created, user_id) VALUES
								(?, ?, datetime('now'), ?)`
			_, err := p.DB.Exec(stmt, true, id, username)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
		} else {
			stmt = `UPDATE posts SET votes = votes - 1 WHERE id = ?`
			_, err := p.DB.Exec(stmt, id)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
		}
	case "down":
		if i == "1" {
			stmt = `INSERT INTO vote (type, post_id, created, user_id) VALUES
								(?, ?, datetime('now'), ?)`
			_, err := p.DB.Exec(stmt, false, id, username)
			if err != nil {
				fmt.Println(err, "4")
				return err
			}
		} else {
			stmt = `UPDATE posts SET votes = votes + 1 WHERE id = ?`
			_, err := p.DB.Exec(stmt, id)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
		}
	}
	return nil
}

// Check single vote
func (p *PostModel) GetVote(id, vote, username string) (string, error) {
	s := ""
	stmt := `SELECT type FROM vote WHERE user_id = ? AND post_id = ?`
	res := p.DB.QueryRow(stmt, username, id)
	err := res.Scan(&s)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
	}

	fmt.Println(s)
	return s, nil
}

// Check all votes for user
func (p *PostModel) GetUserVotes(username string) map[string][]int {
	votes := map[string][]int{}
	stmt := `SELECT post_id, type from vote WHERE user_id = ?`
	if username == "" {
		return nil
	}
	rows, err := p.DB.Query(stmt, username)
	if err != nil {
		return nil
	}

	defer rows.Close()

	vote := 0
	postid := 0

	for rows.Next() {

		err := rows.Scan(&postid, &vote)
		if err != nil {
			return nil
		}
		if vote == 1 {
			votes["upposts"] = append(votes["upposts"], postid)
		} else {
			votes["downposts"] = append(votes["downposts"], postid)
		}
	}

	if err = rows.Err(); err != nil {
		return nil
	}
	fmt.Println(votes)
	return votes
}
