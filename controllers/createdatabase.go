package controllers

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	table = `CREATE TABLE IF NOT EXISTS account_user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		mail TEXT,
		password TEXT,
		auth_token TEXT
	);

	CREATE TABLE IF NOT EXISTS data_post (
		id	INTEGER,
		user_id	INTEGER,
		title	TEXT,
		content	TEXT,
		date	TEXT,
		num_com	INTEGER,
		category	TEXT,
		like	INTEGER,
		dislike	INTEGER,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("user_id") REFERENCES "account_user"("id")
	);

CREATE TABLE IF NOT EXISTS data_comments (
	id	INTEGER,
	post_id	INTEGER,
	user_id	INTEGER,
	comment	TEXT,
	date	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("post_id") REFERENCES "data_post"("id"),
	FOREIGN KEY("user_id") REFERENCES "data_post"("user_id")
);

CREATE TABLE IF NOT EXISTS user_likeordislike (
	id	INTEGER,
	user_id	INTEGER,
	post_id	INTEGER,
	likeordislike	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id") REFERENCES "account_user"("id"),
	FOREIGN KEY("post_id") REFERENCES "data_post"("id")
);`
)

func CreateDatabase(path string) error {
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(table)
	fmt.Println("database created")
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
