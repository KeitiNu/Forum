package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Comment struct {
	ID       int
	PostID   int
	User     string
	Content  string
	Likes    int
	Dislikes int
	Created  time.Time
}

type CommentsModel struct {
	DB *sql.DB
}

func (c *CommentsModel) Insert(co *Comment) (int, error) {
	stmt := `INSERT INTO comments (user_id, post_id, content, created)
			VALUES(?, ?, ?, datetime('now'))`
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
	stmt := `SELECT id, user_id, content, created FROM comments p
	WHERE p.post_id = ?
    ORDER BY created DESC LIMIT 15`

	rows, err := c.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []*Comment{}

	for rows.Next() {
		s := &Comment{}

		err := rows.Scan(&s.ID, &s.User, &s.Content, &s.Created)
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
	WHERE p.user_id = ?
    ORDER BY created DESC LIMIT 15`

	rows, err := c.DB.Query(stmt, username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []*Comment{}

	for rows.Next() {
		s := &Comment{}

		err := rows.Scan(&s.ID, &s.User, &s.Content, &s.Created)
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
