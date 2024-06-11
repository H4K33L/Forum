package authentification

import (
	"database/sql"
	"log"
)

/*
OpenDb(path string) *sql.DB

This function opens a connection to an SQLite database located at the specified path.
It enables foreign key support for the database.

Input:

path : string, the path to the SQLite database file.

Output:

*sql.DB : a pointer to the database connection.
*/
func OpenDb(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	if err != nil {
		log.Fatal("OpenDb 1:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("OpenDb 2:", err)
	}
	return db
}

/*
InitDb(db *sql.DB)

This function initializes the database schema by creating the necessary tables if they do not exist already.

Input:

db : *sql.DB, a pointer to the database connection.

Output: none
*/
func InitDb(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS user (
				uuid VARCHAR(80) NOT NULL UNIQUE,
				email VARCHAR(80) NOT NULL UNIQUE,
				username VARCHAR(10) NOT NULL UNIQUE,
				pwd VARCHAR(255) NOT NULL,
				PRIMARY KEY (uuid)
			);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal("InitDb :", (dberr.Error()))
	}
}

/*
InitDbpost(db *sql.DB)

This function initializes the database schema by creating the necessary tables for posts if they do not exist already.

Input:

db : *sql.DB, a pointer to the database connection.

Output: none
*/

func InitDbpost(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS post
	(
	id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
	uuid VARCHAR(80) NOT NULL,
	postuuid VARCHAR(80) NOT NULL UNIQUE,
	username VARCHAR(80) NOT NULL,
	message LONG VARCHAR,
	document VARCHAR(80),
	ext VARCHAR(80),
	typedoc VARCHAR(80),
	date VARCHAR(80),
	chanel VARCHAR(80),
	target VARCHAR(80),
	answers LONG VARCHAR,
	like INTEGER,
	dislike INTEGER,
	FOREIGN KEY(uuid) 
		REFERENCES user(uuid)
		ON DELETE CASCADE
		ON UPDATE CASCADE
	);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal("InitDbpost :", dberr.Error())
	}
}

/*
InitDbProfile(db *sql.DB)

This function initializes the database schema by creating the necessary tables for profiles if they do not exist already.

Input:

db : *sql.DB, a pointer to the database connection.

Output: none
*/
func InitDbProfile(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS profile (
				uuid VARCHAR(80) NOT NULL UNIQUE,
				username VARCHAR(10) NOT NULL UNIQUE,
				profilepicture VARCHAR(80),
				PRIMARY KEY (uuid),
				FOREIGN KEY (uuid) 
					REFERENCES user(uuid) 
					ON DELETE CASCADE 
					ON UPDATE CASCADE
			);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal("InitDbProfile :", (dberr.Error()))
	}
}

/*
InitDbLike(db *sql.DB)

This function initializes the database schema by creating the necessary tables for likes if they do not exist already.

Input:

db : *sql.DB, a pointer to the database connection.

Output: none
*/
func InitDbLike(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS like
	(
		id INTEGER NOT NULL,
		uuid VARCHAR(80) NOT NULL,
		liked BOOLEAN,
		disliked BOOLEAN,
		PRIMARY KEY (id , uuid)
	);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal("InitDbLike :", dberr.Error())
	}
}
