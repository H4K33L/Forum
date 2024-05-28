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
		posts = getPostByBoth(OpenDb("../DATA/User_data.db"), usename, chanels)
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
		var ID int
		var post post
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&ID,&post.username, &post.message, &post.image, &post.date, &chanel, &target, &answers, &post.like, &post.dislike); err != nil {
			log.Fatal("error in reading",err)
		}
		post.chanel = convertToArray(chanel)
		post.target = convertToArray(target)
		post.answers = convertToArray(answers)
		output = append(output, post)
	}
	if err = UserPost.Err(); err != nil {
		log.Fatal("erreur jsp ou")
	}
	return output
}

func convertToArray(str string) []string {
	return strings.Split(str, "|\\/|-_-|\\/|+{}")
}

func getPostByChanel(db *sql.DB, chanel string) []post {
	chanel = convertToString(strings.Split(chanel, "R/"))
	output := []post{}
	UserPost, err := db.Query("SELECT * FROM post WHERE chanel=?", chanel)
	if err != nil {
		log.Fatal("error in hash")
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var ID int
		var post post
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&ID, &post.username, &post.message, &post.image, &post.date, &chanel, &target, &answers, &post.like, &post.dislike); err != nil {
			log.Fatal("error in reading",err)
		}
		post.chanel = convertToArray(chanel)
		post.target = convertToArray(target)
		post.answers = convertToArray(answers)
		output = append(output, post)
	}
	if err = UserPost.Err(); err != nil {
		log.Fatal("error jsp ou")
	}
	return output
}

func getPostByBoth(db *sql.DB, username string, chanel string) []post {
	postList1 := getPostByUser(db, username)
	if chanel == "" {
		return postList1
	}
	postList2 := getPostByChanel(db, chanel)
	if username == "" {
		return postList2
	}
	output := []post{}
	for _, post1 := range postList1 {
		for _, post2 := range postList2 {
			if comparePost(post1,post2) {
				output = append(output, post1)
			}
		}
	}
	return output
}

func comparePost(post1 post, post2 post) bool {
	return post1.username == post2.username && post1.message == post2.message && post1.image == post2.image && post1.date == post2.date && post1.like == post2.like && post1.dislike == post2.dislike &&convertToString(post1.chanel) == convertToString(post2.chanel) && convertToString(post1.target) == convertToString(post2.target) && convertToString(post1.answers) == convertToString(post2.answers)
}