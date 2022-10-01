package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"git.01.kood.tech/roosarula/forum/pkg/forms"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// Fields we extract from database to use and render on webpages
type User struct {
	Name            string
	Forname         string
	Surname         string
	Email           string
	Age             int
	Password        password
	Gender          Gender
	Hashed_password string
	CreatedAt       time.Time
	Online          int
}

type password struct {
	plaintext *string
	hash      []byte
}

type Gender int64

const (
	Male Gender = iota
	Female
	NonBinary
	Undefined
	Unknown
)

type UserModel struct {
	DB *sql.DB
}

// The Set() method calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

// The Matches() method checks whether the provided plaintext password matches the
// hashed password stored in the struct, returning true if it matches and false
// otherwise.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, ErrInvalidCredentials
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *forms.Validator, email string) {
	v.Check(email != "", "email", "Please enter a valid email address!")
	v.Check(forms.Matches(email, forms.EmailRX), "email", "Please enter a valid email address!")
}

func ValidatePasswordPlaintext(v *forms.Validator, password string) {
	v.Check(password != "", "password", "Please enter your password!")
	v.Check(len(password) >= 3, "password", "Password must be atleast 8 characters long!")
	v.Check(len(password) <= 72, "password", "Password must not be over 72 characters long!")
}

func ValidateUser(v *forms.Validator, user *User) {
	v.Check(user.Name != "", "username", "Please enter your username!")
	v.Check(len(user.Name) >= 5, "username", "Username must be atleast 5 characters long!")
	v.Check(len(user.Name) <= 30, "username", "Username must not be over 30 characters long!")

	//AGE
	v.Check(user.Age >= 5, "age", "User cannot be younger than 5 years!")
	v.Check(user.Age <= 120, "age", "Are you really that old? I dont´t think so!")
	//FORNAME
	v.Check(user.Forname != "", "forname", "Please enter your forname!")
	v.Check(len(user.Forname) >= 2, "forname", "Forname must be atleast 2 characters long!")
	//SURNAME
	v.Check(user.Surname != "", "surname", "Please enter your surname!")
	v.Check(len(user.Surname) >= 2, "surname", "Username must be atleast 2 characters long!")
	//GENDER
	// v.Check(user.Gender == Unknown, "gender", "Please enter gender!")

	// Call the standalone ValidateEmail() helper.
	ValidateEmail(v, user.Email)

	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func ValidateLogin(v *forms.Validator, user *User) {
	v.Check(user.Name != "", "username", "Please enter your username!")
	// v.Check(len(user.Name) >= 8, "username", "Username must be atleast 8 characters long!")
	// v.Check(len(user.Name) <= 16, "username", "Username must not be over 16 characters long!")

	// ValidatePasswordPlaintext() helper.
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

// Insert user into database
func (u UserModel) Insert(user *User, token string) error {
	query := `INSERT INTO users (username, forname, surname, email, age, gender_id, hashed_password, token, created)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?,  datetime('now'))`

	args := []interface{}{user.Name, user.Forname, user.Surname, user.Email, user.Age, user.Gender, user.Password.hash, token}

	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "users_email_key"
	// constraint that we set up in the previous chapter. We check for this error
	// specifically, and return custom ErrDuplicateEmail error instead.
	_, err := u.DB.Exec(query, args...)
	if err != nil {
		if sqlErr, ok := err.(sqlite3.Error); ok {
			if sqlErr.Error() == "UNIQUE constraint failed: users.username" {
				return ErrDuplicateUsername
			}
			if sqlErr.Error() == "UNIQUE constraint failed: users.email" {
				return ErrDuplicateEmail
			}
			fmt.Println(sqlErr.Error())
		}
	}

	return nil
}

// We'll use the Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (u *UserModel) Authenticate(username, password string) error {
	// Retrieve the id and hashed password associated with the given email. If no
	// matching email exists, we return the ErrInvalidCredentials error.

	fmt.Println(username)
	var hashedPassword []byte
	row := u.DB.QueryRow("SELECT hashed_password FROM users WHERE email = ? OR username = ?", username, username)

	err := row.Scan(&hashedPassword)

	if err == sql.ErrNoRows {
		return ErrInvalidCredentials
	} else if err != nil {
		return err
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidCredentials
	} else if err != nil {
		return err
	}

	// Otherwise, the password is correct..
	return nil
}

func (u *UserModel) GetByToken(token string) (*User, error) {
	row := u.DB.QueryRow("SELECT username, email FROM users WHERE token = ?", token)
	user := &User{}
	err := row.Scan(&user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return user, nil
}

func (u *UserModel) GetByUserCredentials(credentials string) (*User, error) {
	row := u.DB.QueryRow("SELECT username, forname, surname, email, age FROM users WHERE username = ? or email = ?", credentials, credentials)

	user := &User{}
	err := row.Scan(&user.Name, &user.Forname, &user.Surname, &user.Email, &user.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return user, nil
}

func (u *UserModel) UpdateByToken(token, username string) error {
	_, err := u.DB.Exec("UPDATE users SET token = ? WHERE username = ?", token, username)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// _, err = u.DB.Exec("UPDATE users SET online = ? WHERE username = ?", 1, username)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	return nil
}

func (u *UserModel) EmailExist(email string) (bool, string, error) {
	row := u.DB.QueryRow("SELECT username, email FROM users WHERE email = ?", email)
	user := &User{}
	err := row.Scan(&user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, user.Name, nil
		}
		return false, user.Name, err
	}
	return true, user.Name, nil
}

func (u *UserModel) GetAllUsers() ([]*User, error) {
	stmt := `SELECT username, forname, surname, email FROM users ORDER BY username`

	rows, err := u.DB.Query(stmt)
	fmt.Println(rows)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		s := &User{}

		err := rows.Scan(&s.Name, &s.Forname, &s.Surname, &s.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
