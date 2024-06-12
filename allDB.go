package client

import (
	"database/sql"
	"fmt"
	"net/http"
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
func OpenDb(path string, w http.ResponseWriter, r *http.Request) *sql.DB {
	// Open a connection to the SQLite database with foreign key constraints enabled.
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	// Check if there was an error opening the database connection.
	if err != nil {
		fmt.Println("allDB OpenDb 1 :", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return nil
	}
	// Ping the database to ensure the connection is valid.
	if err = db.Ping(); err != nil {
		fmt.Println("allDB OpenDb 2 :", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return nil
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
func InitDb(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("alldb InitDb sql :", dberr.Error())
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
}

/*
InitDbpost(db *sql.DB)

This function initializes the database schema by creating the necessary tables for posts if they do not exist already.

Input:

db : *sql.DB, a pointer to the database connection.

Output: none
*/
func InitDbpost(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("allDB InitDbpost sql :", dberr.Error())
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
}

/*
InitDbProfile(db *sql.DB)

This function initializes the database schema by creating the necessary tables for profiles if they do not exist already.

Input:

db : *sql.DB, a pointer to the database connection.

Output: none
*/
func InitDbProfile(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("allDB InitDbProfile :", dberr.Error())
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
}

/*
InitDbLike(db *sql.DB)

This function initializes the database schema by creating the necessary tables for likes if they do not exist already.

Input:

db : *sql.DB, a pointer to the database connection.

Output: none
*/
func InitDbLike(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("allDB InitDbLike :", dberr.Error())
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
}
