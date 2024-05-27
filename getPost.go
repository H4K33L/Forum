package authentification

import (
	"database/sql"
	"strings"

	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

)

func GetPost(w http.ResponseWriter, r *http.Request) []post {
	var posts []post
	if r.Method == "GET" {
		usename := r.FormValue("username")
		chanels := r.FormValue("chanels")
		if usename != "" {
			return getPostByUser(OpenDb("../DATA/User_data.db"), usename)
		} else if chanels != "" {
			return getPostByChanel(OpenDb("../DATA/User_data.db"), chanels)
		}
	}
	return posts
}

func getPostByUser(db *sql.DB, username string) []post {
	output := []post{}
	UserPost, err := db.Query("SELECT * FROM post WHERE username=?", username)
	if err != nil {
		log.Fatal("error in hash")
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var post post
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&post.username, &post.message, &post.image, &post.date, chanel, target, answers, &post.like, &post.dislike); err != nil {
			log.Fatal()
		}
		post.chanel = convertToArray(chanel)
		post.target = convertToArray(target)
		post.answers = convertToArray(answers)
		output = append(output, post)
	}
	if err = UserPost.Err(); err != nil {
		log.Fatal()
	}
	return output
}

func convertToArray(str string) []string {
	return strings.Split(str, "|\\/|-_-|\\/|+{}")
}

func getPostByChanel(db *sql.DB, chanel string) []post {
	output := []post{}
	UserPost, err := db.Query("SELECT * FROM post WHERE chanel=?", chanel)
	if err != nil {
		log.Fatal("error in hash")
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var post post
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&post.username, &post.message, &post.image, &post.date, chanel, target, answers, &post.like, &post.dislike); err != nil {
			log.Fatal()
		}
		post.chanel = convertToArray(chanel)
		post.target = convertToArray(target)
		post.answers = convertToArray(answers)
		output = append(output, post)
	}
	if err = UserPost.Err(); err != nil {
		log.Fatal()
	}
	return output
}