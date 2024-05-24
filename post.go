package authentification

import (
	"database/sql"
	"strings"

	//"html/template"
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
	id INTEGER NOT NULL UNIQUE,
	username VARCHAR(80) NOT NULL,
	message VARCHAR(80),
	image VARCHAR(80),
	date VARCHAR(80),
	chanel VARCHAR(80),
	target VARCHAR(80),
	answers VARCHAR(80),
	like INTIGER,
	dislike INTEGER,
	PRIMARY KEY(id AUTOINCREMENT)
	);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal(dberr.Error())
	}
}

func convertToString(array []string) string {
	return strings.Join(array, "|\\/|-_-|\\/|+{}")
}

func convertToArray(str string) []string {
	return strings.Split(str, "|\\/|-_-|\\/|+{}")
}

func UserPost(w http.ResponseWriter, r *http.Request) {
	var post post
	if r.Method == "POST" {
		post.username = r.FormValue("username")
		post.message = r.FormValue("message")
		post.image = r.FormValue("image")
		post.date = r.FormValue("date")
		post.chanel = strings.Split(r.FormValue("chanel"), "R/")
		post.target = strings.Split(r.FormValue("target"), " \\ ")
		InitDbpost(OpenDb("../DATA/User_data.db"))
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

/*
func getPost(db *sql.DB, username string, chanels []string) (bool, []post) {
	output := []post{}
	UserPost, err := db.Query("SELECT * FROM post WHERE username=?", username)
	if err != nil {
		fmt.Println("error in hash")
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var post post
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&post.username, &post.message, &post.image, &post.date, chanel, target, answers, &post.like, &post.dislike); err != nil {
			return false, output
		}
		post.chanel = convertToArray(chanel)
		post.target = convertToArray(target)
		post.answers = convertToArray(answers)
		for _,arg := range chanels {
			if !contains(post.chanel ,arg) {
				return false, output
			}
		}
		output = append(output, post)
	}
	if err = UserPost.Err(); err != nil {
		return false, output
	}
	return true, output
}

func contains(s []string, str string) bool {
    for _, v := range s {
        if v == str {
            return true
        }
    }

    return false
}
*/
