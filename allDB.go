package client

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
	// Open a connection to the SQLite database with foreign key constraints enabled.
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	// Check if there was an error opening the database connection.
	if err != nil {
		log.Fatal(" allDB OpenDb 1:", err)
	}
	// Ping the database to ensure the connection is valid.
	if err = db.Ping(); err != nil {
		log.Fatal("allDB OpenDb 2:", err)
	}
	// Return the database connection object.
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
	// SQL statement to create a 'user' table if it doesn't exist already.
	table := `CREATE TABLE IF NOT EXISTS user (
				uuid VARCHAR(80) NOT NULL UNIQUE,
				email VARCHAR(80) NOT NULL UNIQUE,
				username VARCHAR(10) NOT NULL UNIQUE,
				pwd VARCHAR(255) NOT NULL,
				PRIMARY KEY (uuid)
			);`

	// Execute the SQL statement to create the table in the database.
	_, dberr := db.Exec(table)
	// Check if there was an error during the execution of the SQL statement.
	if dberr != nil {
		log.Fatal("InitDb sql :", dberr.Error())
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
	// SQL statement to create a 'post' table if it doesn't exist already.
	table := `CREATE TABLE IF NOT EXISTS post (
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
		like INTEGER,
		dislike INTEGER,
		FOREIGN KEY(uuid) 
			REFERENCES user(uuid)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);`
	// Execute the SQL statement to create the 'post' table in the database.
	_, dberr := db.Exec(table)
	// Check if there was an error during the execution of the SQL statement.
	if dberr != nil {
		log.Fatal("allDB InitDbpost sql :", dberr.Error())
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
	// SQL statement to create a 'profile' table if it doesn't exist already.
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
	// Execute the SQL statement to create the 'profile' table in the database.
	_, dberr := db.Exec(table)
	// Check if there was an error during the execution of the SQL statement.
	if dberr != nil {
		log.Fatal("allDB InitDbProfile :", dberr.Error())
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
	// SQL statement to create a 'like' table if it doesn't exist already.
	table := `CREATE TABLE IF NOT EXISTS like (
		id INTEGER NOT NULL,
		uuid VARCHAR(80) NOT NULL,
		liked BOOLEAN,
		disliked BOOLEAN,
		PRIMARY KEY (id , uuid)
	);`
	// Execute the SQL statement to create the 'like' table in the database.
	_, dberr := db.Exec(table)
	// Check if there was an error during the execution of the SQL statement.
	if dberr != nil {
		log.Fatal("allDB InitDbLike :", dberr.Error())
	}
}
