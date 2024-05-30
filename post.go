package authentification

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
	
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type post struct {
	uuid     string
	username string
	message  string
	image    string
	date     string
	like     int
	dislike  int
	chanel   []string
	target   []string
	answers  []string
}

func InitDbpost(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS post
	(
	id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
	uuid VARCHAR(80) NOT NULL UNIQUE,
	username VARCHAR(80) NOT NULL,
	message VARCHAR(80),
	image VARCHAR(80),
	date VARCHAR(80),
	chanel VARCHAR(80),
	target VARCHAR(80),
	answers VARCHAR(80),
	like INTEGER,
	dislike INTEGER,
	UNIQUE( uuid, username),
	FOREIGN KEY(uuid, username) REFERENCES user(uuid, username)
	);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal("InitDbpost :", dberr.Error())
	}
}

func UserPost(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("../DATA/User_data.db")
	var post post
	if r.Method == "POST" {
		uid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				// Si le cookie n'existe pas
				log.Fatal("Cookie uuid not found")
			}
			// Autre erreur
			log.Fatal("Error retrieving cookie uuid :", err)
		}
		var username string
		err1 := db.QueryRow("SELECT username FROM user WHERE uuid=?", uid.Value).Scan(&username)
		if err1 != nil {
			if err1 == sql.ErrNoRows {
				log.Fatal("sql user post :", err1)
			}
			log.Fatal(err1)
		}
		post.uuid = uid.Value
		post.username = username
		post.message = r.FormValue("message")
		post.image = r.FormValue("image")
		then := time.Now()
		post.date = strconv.Itoa(then.Year()) + "/" + then.Month().String() + "/" + strconv.Itoa(then.Day()) + "/" + strconv.Itoa(then.Hour()) + "/" + strconv.Itoa(then.Minute()) + "/" + strconv.Itoa(then.Second())
		post.chanel = strings.Split(r.FormValue("chanel"), "R/")
		post.target = strings.Split(r.FormValue("target"), "\\\\")
		AddPost(OpenDb("../DATA/User_data.db"), post)
	}
}

func AddPost(db *sql.DB, post post) {
	statement, err := db.Prepare("INSERT INTO post(uuid, username, message, image, date, chanel, target, answers, like, dislike) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("sql add post", err)
	}
	chanel := convertToString(post.chanel)
	target := convertToString(post.target)
	answers := convertToString(post.answers)
	statement.Exec(post.uuid, post.username, post.message, post.image, post.date, chanel, target, answers, post.like, post.dislike)
	defer db.Close()
}

func convertToString(array []string) string {
	return strings.Join(array, "|\\/|-_-|\\/|+{}")
}
