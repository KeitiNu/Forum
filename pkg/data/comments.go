package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Comment struct {
	ID      int
	PostID  int
	User    string
	Content string
	Votes   int
	Created time.Time
}

type CommentsModel struct {
	DB *sql.DB
}

func (c *CommentsModel) Insert(co *Comment) (int, error) {
	stmt := `INSERT INTO comments (user_id, post_id, content, votes, created)
			VALUES(?, ?, ?, 0, datetime('now'))`
	result, err := c.DB.Exec(stmt, co.User, co.PostID, co.Content)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	return int(id), nil
}

func (c *CommentsModel) Get(id int) (*Comment, error) {

	row := c.DB.QueryRow("SELECT user_id FROM comments WHERE id = ?", id)
	comment := &Comment{}
	err := row.Scan(&comment.User)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return comment, nil

}

func (c *CommentsModel) Latest(id int) ([]*Comment, error) {
	stmt := `SELECT id, user_id, content, created, votes FROM comments p
	WHERE c.post_id = ?
    ORDER BY created DESC LIMIT 15`

	rows, err := c.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []*Comment{}

	for rows.Next() {
		s := &Comment{}

		err := rows.Scan(&s.ID, &s.User, &s.Content, &s.Created, &s.Votes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *CommentsModel) Update(content string, id int) error {
	stmt := `UPDATE comments SET content=?
			WHERE id = ?`
	_, err := c.DB.Exec(stmt, content, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (c *CommentsModel) Delete(id int) error {
	stmt := `DELETE FROM comments WHERE id = ?`
	_, err := c.DB.Exec(stmt, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (c *CommentsModel) GetUserComments(username string) ([]*Comment, error) {
	stmt := `SELECT id, user_id, content, created FROM comments p
	WHERE c.user_id = ?
    ORDER BY created DESC LIMIT 15`

	rows, err := c.DB.Query(stmt, username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []*Comment{}

	for rows.Next() {
		s := &Comment{}

		err := rows.Scan(&s.ID, &s.User, &s.Content, &s.Created, &s.Votes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

// Check single vote
func (c *CommentsModel) GetVote(id, vote, username string) (string, error) {
	var s string
	stmt := `SELECT type FROM vote WHERE user_id = ? AND post_id = ?`
	res := c.DB.QueryRow(stmt, username, id)
	err := res.Scan(&s)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
	}

	return s, nil
}

func (c *CommentsModel) AddVote(id, vote, username string) error {
	var stmt string
	i, err := c.GetVote(id, vote, username)
	if err != nil {
		if err == sql.ErrNoRows {
			switch vote {
			case "up":
				stmt = `INSERT INTO vote (type, comment_id, created, user_id) VALUES
										(?, ?, datetime('now'), ?)`
				_, err := c.DB.Exec(stmt, true, id, username)
				if err != nil {
					fmt.Println(err, "1")
					return err
				}
				stmt = `UPDATE comments SET votes = votes + 1 WHERE id = ?`
				_, err = c.DB.Exec(stmt, id)
				if err != nil {
					fmt.Println(err, "3")
					return err
				}
				return nil
			case "down":
				stmt = `INSERT INTO vote (type, comment_id, created, user_id) VALUES
										(?, ?, datetime('now'), ?)`
				_, err := c.DB.Exec(stmt, false, id, username)
				if err != nil {
					fmt.Println(err, "2")
					return err
				}
				stmt = `UPDATE comments SET votes = votes - 1 WHERE id = ?`
				_, err = c.DB.Exec(stmt, id)
				if err != nil {
					fmt.Println(err, "3")
					return err
				}
				return nil
			}
		}
	}
	stmt = `DELETE FROM vote WHERE comment_id = ? AND user_id = ?`
	_, err = c.DB.Exec(stmt, id, username)
	if err != nil {
		fmt.Println("ERROR")
		return err
	}
	switch vote {
	case "up":
		if i == "0" {
			stmt = `INSERT INTO vote (type, comment_id, created, user_id) VALUES
								(?, ?, datetime('now'), ?)`
			_, err := c.DB.Exec(stmt, true, id, username)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
			stmt = `UPDATE comments SET votes = votes + 2 WHERE id = ?`
			_, err = c.DB.Exec(stmt, id)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}

		} else if i == "" {
			stmt = `UPDATE comments SET votes = votes + 1 WHERE id = ?`
			_, err := c.DB.Exec(stmt, id)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
			stmt = `INSERT INTO vote (type, comment_id, created, user_id) VALUES
								(?, ?, datetime('now'), ?)`
			_, err = c.DB.Exec(stmt, true, id, username)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
		} else {
			stmt = `UPDATE comments SET votes = votes - 1 WHERE id = ?`
			_, err := c.DB.Exec(stmt, id)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
			stmt = `INSERT INTO vote (type, comment_id, created, user_id) VALUES
								(?, ?, datetime('now'), ?)`
			_, err = c.DB.Exec(stmt, nil, id, username)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
		}
	case "down":
		if i == "1" {
			stmt = `INSERT INTO vote (type, comment_id, created, user_id) VALUES
								(?, ?, datetime('now'), ?)`
			_, err := c.DB.Exec(stmt, false, id, username)
			if err != nil {
				fmt.Println(err, "4")
				return err
			}
			stmt = `UPDATE comments SET votes = votes - 2 WHERE id = ?`
			_, err = c.DB.Exec(stmt, id)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
		} else if i == "" {
			stmt = `UPDATE comments SET votes = votes - 1 WHERE id = ?`
			_, err := c.DB.Exec(stmt, id)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
			stmt = `INSERT INTO vote (type, comment_id, created, user_id) VALUES
								(?, ?, datetime('now'), ?)`
			_, err = c.DB.Exec(stmt, false, id, username)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
		} else {
			stmt = `UPDATE comments SET votes = votes + 1 WHERE id = ?`
			_, err := c.DB.Exec(stmt, id)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
			stmt = `INSERT INTO vote (type, comment_id, created, user_id) VALUES
								(?, ?, datetime('now'), ?)`
			_, err = c.DB.Exec(stmt, nil, id, username)
			if err != nil {
				fmt.Println(err, "3")
				return err
			}
		}
	}
	return nil
}
