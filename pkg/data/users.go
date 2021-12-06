package data

import (
	"database/sql"
	"time"

	"github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
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

// Insert user into database
func (u UserModel) Insert(username, password string) (int, error) {
	query := `INSERT INTO users (id, username, hashed_password, created)
	VALUES(?, ?, ?,  datetime('now'))`
	uuid := uuid.NewV4().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}
	result, err := u.DB.Exec(query, uuid, username, string(hashedPassword))
	if err != nil {
		if sqlErr, ok := err.(sqlite3.Error); ok {
			if string(sqlErr.Code.Error()) == "constraint failed" {
				return 0, ErrDuplicateUsername
			}
		}
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

// We'll use the Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(username, password string) (string, error) {
	// Retrieve the id and hashed password associated with the given email. If no
	// matching email exists, we return the ErrInvalidCredentials error.
	var id string
	var hashedPassword []byte
	row := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE username = ?", username)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return "", ErrInvalidCredentials
	} else if err != nil {
		return "", err
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return "", ErrInvalidCredentials
	} else if err != nil {
		return "", err
	}

	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}
