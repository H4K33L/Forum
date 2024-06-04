package authentification

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func GetPost(w http.ResponseWriter, r *http.Request) []Post {
	if r.Method == "GET" {
		usename := r.FormValue("username")
		chanels := r.FormValue("chanels")
		return GetPostByBoth(OpenDb("./DATA/User_data.db"), usename, chanels)
	}
	return nil
}

func getPostByUser(db *sql.DB, username string) []Post {
	output := []Post{}
	UserPost, err := db.Query("SELECT * FROM post WHERE username=?", username)
	if err != nil {
		log.Fatal("error in hash")
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var UUID string
		var post Post
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&post.ID, &UUID, &post.Username, &post.Message, &post.Image, &post.Date, &chanel, &target, &answers, &post.Like, &post.Dislike); err != nil {
			log.Fatal("error in reading", err)
		}
		post.Chanel = convertToArray(chanel)
		post.Target = convertToArray(target)
		post.Answers = convertToArray(answers)
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

func getPostByChanel(db *sql.DB, chanel string) []Post {
	chanel = convertToString(strings.Split(chanel, "R/"))
	output := []Post{}
	UserPost, err := db.Query("SELECT * FROM post WHERE chanel=?", chanel)
	if err != nil {
		log.Fatal("error in hash")
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var UUID string
		var post Post
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&post.ID, &UUID, &post.Username, &post.Message, &post.Image, &post.Date, &chanel, &target, &answers, &post.Like, &post.Dislike); err != nil {
			log.Fatal("error in reading", err)
		}
		post.Chanel = convertToArray(chanel)
		post.Target = convertToArray(target)
		post.Answers = convertToArray(answers)
		output = append(output, post)
	}
	if err = UserPost.Err(); err != nil {
		log.Fatal("error jsp ou")
	}
	return output
}

func GetPostByBoth(db *sql.DB, username string, chanel string) []Post {
	postList1 := getPostByUser(db, username)
	if chanel == "" {
		return postList1
	}
	postList2 := getPostByChanel(db, chanel)
	if username == "" {
		return postList2
	}
	output := []Post{}
	for _, post1 := range postList1 {
		for _, post2 := range postList2 {
			if comparePost(post1, post2) {
				output = append(output, post1)
			}
		}
	}
	return output
}

func comparePost(post1 Post, post2 Post) bool {
	return post1.ID == post2.ID
}
