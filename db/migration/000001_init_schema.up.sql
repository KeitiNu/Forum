CREATE TABLE IF NOT EXISTS user_type(
	id integer PRIMARY KEY NOT NULL,
	name text UNIQUE NOT NULL,
	description text
);

CREATE TABLE IF NOT EXISTS categories (
	title TEXT PRIMARY KEY NOT NULL UNIQUE,
	description TEXT NOT NULL,
	created DATETIME
);

CREATE TABLE IF NOT EXISTS gender_type (
	id integer PRIMARY KEY NOT NULL,
	name text UNIQUE NOT NULL,
	description text
);

INSERT INTO gender_type (id, name)
VALUES (0, "Female"),
(1, "Male"),
(2, "Non-Binary"),
(3, "Unknown"),
(4, "Undefined");

CREATE TABLE IF NOT EXISTS users (
	username PRIMARY KEY NOT NULL UNIQUE,
	forname TEXT NOT NULL DEFAULT "forname",
	surname TEXT NOT NULL DEFAULT "surname",
	age integer NOT NULL DEFAULT 0,
	gender_id integer NOT NULL DEFAULT 4,	
	email TEXT NOT NULL UNIQUE,
	hashed_password BLOB NOT NULL,
	token BLOB NOT NULL,
	created DATETIME,
	updated DATETIME,
	online integer NOT NULL DEFAULT 0, 
	FOREIGN KEY(gender_id)
		REFERENCES gender_type(id) 
			ON DELETE SET DEFAULT -- kui user_type kustutatse, sii
);

CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	votes INTEGER,
	created DATETIME,
	imagesrc TEXT,
	user TEXT NOT NULL,
	FOREIGN KEY (user) REFERENCES users (username)
);

CREATE TABLE IF NOT EXISTS post_category (
	id INTEGER PRIMARY KEY,
	post_id INTEGER,
	category_id TEXT,
	FOREIGN KEY (category_id) REFERENCES categories (title)
	FOREIGN KEY (post_id) REFERENCES posts (id)
);


CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY,
	post_id INTEGER NOT NULL,
	user_id TEXT NOT NULL,
	content TEXT NOT NULL,
	votes INTEGER,
	created DATETIME,
	FOREIGN KEY (user_id) REFERENCES users (username)
	FOREIGN KEY (post_id) REFERENCES posts (id)
);

CREATE TABLE IF NOT EXISTS vote(
	id integer PRIMARY KEY NOT NULL,
	post_id integer DEFAULT 0,
	comment_id integer DEFAULT 0,
	type bool,
	created datetime NOT NULL,
	user_id text NOT NULL,
	FOREIGN KEY (user_id)
		REFERENCES users(username)
	FOREIGN KEY(comment_id)
		REFERENCES comment(id)
				ON DELETE CASCADE
	FOREIGN KEY(post_id)
		REFERENCES post(id)
				ON DELETE CASCADE
);



CREATE TABLE IF NOT EXISTS messages(
id integer PRIMARY KEY NOT NULL,
sender_id integer NULL,
receiver_id integer NULL,
content text,
sent_at CURRENT_TIMESTAMP NOT NULL,
read_at datetime,
CONSTRAINT fk_sender
FOREIGN KEY (sender_id)
    REFERENCES users(id)
            ON DELETE SET NULL,
CONSTRAINT fk_receievr
FOREIGN KEY (receiver_id)
    REFERENCES users(id)
            ON DELETE SET NULL
);

-- Seed DATABASE

INSERT INTO categories (title, description, created)
VALUES ("Javascript", "For questions regarding programming in ECMAScript (JavaScript/JS) and its various dialects/implementations (excluding ActionScript).", CURRENT_TIMESTAMP),
("RUST", "About Rust programming language", CURRENT_TIMESTAMP),
("Python", "Python is a multi-paradigm, dynamically typed, multi-purpose programming language. ", CURRENT_TIMESTAMP),
("Golang", "Go is an open source programming language supported by Google", CURRENT_TIMESTAMP);
