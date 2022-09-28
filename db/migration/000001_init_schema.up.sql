-- CREATE TABLE IF NOT EXISTS user_type(
-- 		id integer PRIMARY KEY NOT NULL,
-- 		name text UNIQUE NOT NULL,
-- 		description text
-- 		);

-- CREATE TABLE IF NOT EXISTS user (
-- id integer PRIMARY KEY NOT NULL,
-- username text UNIQUE NOT NULL, 
-- email text UNIQUE NOT NULL,  
-- hashed_password BLOB NOT NULL, 
-- token BLOB NOT NULL, 
-- auth_token BLOB,
-- provider text, --kui üks on täidetud, peab teine ka olema! Mis constraint panna??
-- created datetime NOT NULL,
-- updated datetime,
-- type_id integer NOT NULL DEFAULT 1,
-- -- 1 peaks siis olema nö tavauser - 
-- --st, et kui ise inputi ei anna, siis on value automaatselt "1"
-- FOREIGN KEY(type_id)
--     REFERENCES user_type(id) 
--         ON DELETE SET DEFAULT -- kui user_type kustutatse, siis muudetakse default value peale
-- );


-- CREATE TABLE IF NOT EXISTS post(
-- id integer PRIMARY KEY NOT NULL,
-- user_id integer NOT NULL,
-- title text NOT NULL,
-- content text NOT NULL,
-- likes integer,
-- dislikes integer,
-- created datetime NOT NULL,
-- FOREIGN KEY(user_id) 
--     REFERENCES user(id)
--             ON DELETE CASCADE -- kui kasutaja kustutatakse, kustuvad ka temaga seotud postitused
-- );

-- -- CREATE TABLE post_category(
-- -- id integer PRIMARY KEY NOT NULL,
-- -- name text NOT NULL
-- -- )


-- -- CREATE TABLE post_in_category(
-- -- id integer PRIMARY KEY NOT NULL,
-- -- post_id integer NOT NULL,
-- -- category_id integer NOT NULL,
-- -- FOREIGN KEY(post_id) 
-- --     REFERENCES post(id)
-- --             ON DELETE CASCADE,
-- -- FOREIGN KEY(category_id) 
-- --     REFERENCES post_category(id)
-- --             ON DELETE CASCADE
-- -- )

-- -- CREATE TABLE comment(
-- -- id integer PRIMARY KEY NOT NULL,
-- -- post_id integer NOT NULL,
-- -- user_id integer NOT NULL,
-- -- content text NOT NULL,
-- -- likes integer,
-- -- dislikes integer,
-- -- created datetime NOT NULL,
-- -- FOREIGN KEY(user_id) 
-- --     REFERENCES user(id)
-- --             ON DELETE CASCADE, -- kui kasuataja kustutatakse, kustuvad ka temaga seotud kommentaarid
-- -- FOREIGN KEY(post_id) 
-- --     REFERENCES post(id)
-- --             ON DELETE CASCADE -- postituse kustutamisel kustutatakse ka seotud kommentaar
-- -- )


-- -- CREATE TABLE vote(
-- -- id integer PRIMARY KEY NOT NULL,
-- -- user_id integer NOT NULL,
-- -- post_id integer,
-- -- comment_id integer,
-- -- type bool NOT NULL, -- false(0) for dislike and true(1) for a like
-- -- created datetime NOT NULL,
-- -- FOREIGN KEY(comment_id) 
-- --     REFERENCES comment(id)
-- --             ON DELETE CASCADE, -- kommentaari kustutamisel kustuvad ka seotud hääled
-- -- FOREIGN KEY(post_id)
-- --     REFERENCES post(id)
-- --             ON DELETE CASCADE, -- postituse kustutamisel kustutatakse ka seotud hääled
-- -- FOREIGN KEY(user_id) 
-- --     REFERENCES user(id)
-- --             ON DELETE CASCADE
-- -- CHECK(post_id NOT NULL AND comment_id is NULL 
-- -- OR post_id is NULL AND comment_id NOT NULL) 
-- -- -- post_id and comment_id: only one of these colmuns must be filled
-- -- )


-- -- CREATE TABLE report(
-- -- id integer PRIMARY KEY NOT NULL,
-- -- post_id integer NOT NULL,
-- -- reporter_id NOT NULL,
-- -- created datetime NOT NULL,
-- -- admin_id integer,
-- -- accepted bool,
-- -- reacted datetime,
-- -- FOREIGN KEY(reporter_id) 
-- --     REFERENCES user(id)
-- --             ON DELETE CASCADE, -- kui kasuataja kustutatakse, kustuvad ka temaga saadetud raportid
-- -- FOREIGN KEY(post_id) 
-- --     REFERENCES post(id)
-- --             ON DELETE CASCADE, -- postituse kustutamisel kustutatakse ka seotud raportid
-- -- FOREIGN KEY(admin_id) 
-- --     REFERENCES user(id)
-- --             ON DELETE SET NULL -- admini kustutamisel reporti ei kustutata
-- -- )


-- -- CREATE TABLE upgrade_request(
-- -- id integer PRIMARY KEY NOT NULL,
-- -- requester_id UNIQUE NOT NULL, -- ei saa olla mitut upgrade_requesti
-- -- created datetime NOT NULL,
-- -- admin_id integer,
-- -- accepted bool,
-- -- reacted datetime,
-- -- FOREIGN KEY(requester_id) 
-- --     REFERENCES user(id)
-- --             ON DELETE CASCADE, -- kui kasuataja kustutatakse, kustub ka upgrade_request
-- -- FOREIGN KEY(admin_id) 
-- --     REFERENCES user(id)
-- --             ON DELETE SET NULL -- admini kustutamisel requesti ei kustutata 
-- -- )

-- -- CREATE TABLE notification(
-- -- id integer PRIMARY KEY NOT NULL,
-- -- user_id integer NOT NULL,
-- -- vote_id integer,
-- -- comment_id integer,
-- -- report_id integer,
-- -- upgrade_id integer,
-- -- seen bool NOT NULL DEFAULT false,
-- -- FOREIGN KEY(comment_id) 
-- --     REFERENCES comment(id)
-- --             ON DELETE CASCADE,
-- -- FOREIGN KEY(vote_id) 
-- --     REFERENCES vote(id)
-- --             ON DELETE CASCADE,
-- -- FOREIGN KEY(report_id) 
-- --     REFERENCES report(id)
-- --             ON DELETE CASCADE,
-- -- FOREIGN KEY(upgrade_id) 
-- --     REFERENCES upgrade(id)
-- --             ON DELETE CASCADE,
-- -- CHECK(vote_id NOT NULL AND comment_id IS NULL AND report_id IS NULL AND upgrade_id IS NULL
-- -- OR vote_id IS NULL AND comment_id NOT NULL AND report_id IS NULL AND upgrade_id IS NULL
-- -- OR vote_id IS NULL AND comment_id IS NULL AND report_id NOT NULL AND upgrade_id IS NULL
-- -- OR vote_id IS NULL AND comment_id IS NULL AND report_id IS NULL AND upgrade_id NOT NULL) 
-- -- )


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
	online integer DEFAULT 0,
	token BLOB NOT NULL,
	created DATETIME,
	updated DATETIME,
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
