package data

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

type Message struct {
	Sender    string
	Recipient string
	Content   string
	SentAt    string
}

type MessageModel struct {
	DB *sql.DB
}

func (msg MessageModel) Insert(message *Message) error {
	query := `INSERT INTO messages (sender_id, receiver_id, content, sent_at)
	VALUES(?, ? ,?, datetime('now'))`

	args := []interface{}{message.Sender, message.Recipient, message.Content}

	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "users_email_key"
	// constraint that we set up in the previous chapter. We check for this error
	// specifically, and return custom ErrDuplicateEmail error instead.
	_, err := msg.DB.Exec(query, args...)
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

func (msg MessageModel) GetMessages(rec_id string, sender_id string, offset int) ([]*Message, error) {

	// query := `INSERT INTO messages (sender_id, receiver_id, content, sent_at)
	// VALUES(?, ? ,?, datetime('now'))`

	query := `SELECT sender_id, receiver_id, content, sent_at FROM messages
	WHERE sender_id = ? AND receiver_id = ? OR receiver_id = ? AND sender_id = ?
    ORDER BY sent_at DESC LIMIT 10 OFFSET ?`

	args := []interface{}{rec_id, sender_id, rec_id, sender_id, offset}

	// rows, err := msg.DB.Exec(query, args...)

	rows, err := msg.DB.Query(query, args...)
	fmt.Println(rows)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := []*Message{}

	for rows.Next() {
		s := &Message{}

		err := rows.Scan(&s.Sender, &s.Recipient, &s.Content, &s.SentAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}
