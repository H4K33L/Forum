package authentification

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type post struct {
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
	id INT NOT NULL UNIQUE,
	username VARCHAR(80) NOT NULL,
	message VARCHAR(80),
	image VARCHAR(80),
	date VARCHAR(80),
	chanel VARCHAR(80),
	target VARCHAR(80),
	answers VARCHAR(80),
	like INT,
	dislike INT,
	PRIMARY KEY(id AUTOINCREMENT)
	);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal(dberr.Error())
	}
}

func UserPost(w http.ResponseWriter, r *http.Request) {
	var post post
	if r.Method == "POST" {
		username, err := r.Cookie("Username")
		if err != nil {
			if err == http.ErrNoCookie {
				// Si le cookie n'existe pas
				log.Fatal("Cookie username not found")
			}
			// Autre erreur
			log.Fatal("Error retrieving cookie:", err)
		}
		post.username = username.Value
		post.message = r.FormValue("message")
		post.image = r.FormValue("image")
		then := time.Now()
		post.date = strconv.Itoa(then.Year()) + "/" + then.Month().String() + "/" + strconv.Itoa(then.Day()) + "/" + strconv.Itoa(then.Hour()) + "/" + strconv.Itoa(then.Minute()) + "/" + strconv.Itoa(then.Second())
		post.chanel = strings.Split(r.FormValue("chanel"), "R/")
		post.target = strings.Split(r.FormValue("target"), "\\\\")
		AddPost(OpenDb("../DATA/User_data.db"), post)
	}
}

func AddPost(db *sql.DB, post post) bool {
	statement, err := db.Prepare("INSERT INTO post(username, message, image, date, chanel, target, answers, like, dislike) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(post.username)
	chanel := convertToString(post.chanel)
	target := convertToString(post.target)
	answers := convertToString(post.answers)
	statement.Exec(post.username, post.message, post.image, post.date, chanel, target, answers, post.like, post.dislike)
	defer db.Close()
	return true
}

func convertToString(array []string) string {
	return strings.Join(array, "|\\/|-_-|\\/|+{}")
}
