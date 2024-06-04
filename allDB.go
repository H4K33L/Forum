package authentification
/*
import (
	"database/sql"
	"log"
)

func OpenDb(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("OpenDb 1:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("OpenDb 2:", err)
	}
	return db
}

func InitDb(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS user (
				uuid VARCHAR(80) NOT NULL UNIQUE,
				email VARCHAR(80) NOT NULL UNIQUE,
				username VARCHAR(10) NOT NULL UNIQUE,
				pwd VARCHAR(255) NOT NULL
			);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal("InitDb :", (dberr.Error()))
	}
}

func InitDbpost(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS post
	(
	id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
	uuid VARCHAR(80) NOT NULL,
	username VARCHAR(80) NOT NULL,
	message LONG VARCHAR,
	image VARCHAR(80),
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

func InitDbProfile(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS profile (
				uuid VARCHAR(80) NOT NULL UNIQUE,
				username VARCHAR(10) NOT NULL UNIQUE,
				profilepicture VARCHAR(80),
				follow VARCHAR(80),
				follower VARCHAR(80),
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
*/