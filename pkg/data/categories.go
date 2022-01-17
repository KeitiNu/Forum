package data

import (
	"database/sql"
	"time"
)

// Fields we extract from database to use and render on webpages
type Category struct {
	ID          int
	Title       string
	Description string
	Created     time.Time
}

type CategoryModel struct {
	DB *sql.DB
}

func (s *CategoryModel) Insert(title, description string) error {
	stmt := `INSERT INTO categories (title, description, created)
			VALUES(?, ?, datetime('now'))`
	_, err := s.DB.Exec(stmt, title, description)
	if err != nil {
		return err
	}

	return nil
}

func (m *CategoryModel) Latest() ([]*Category, error) {
	stmt := `SELECT title, description, created FROM categories
    ORDER BY created DESC LIMIT 15`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := []*Category{}

	for rows.Next() {
		s := &Category{}

		err := rows.Scan(&s.Title, &s.Description, &s.Created)
		if err != nil {
			return nil, err
		}
		categories = append(categories, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (m *CategoryModel) GetOne(title string) ([]*Category, error) {
	stmt := `SELECT title, description FROM categories WHERE title = ?`

	rows, err := m.DB.Query(stmt, title)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	category := []*Category{}

	for rows.Next() {
		s := &Category{}

		err := rows.Scan(&s.Title, &s.Description)
		if err != nil {
			return nil, err
		}
		category= append(category, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return category, nil

}
