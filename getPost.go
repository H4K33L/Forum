package authentification

import (
	"database/sql"
	"log"
	"net/http"
	"slices"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func GetPost(w http.ResponseWriter, r *http.Request) []Post {
	if r.Method == "GET" {
		uid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Fatal("cookie not found userpost")
			}
			log.Fatal("Error retrieving cookie uuid :", err)
		}

		usename := r.FormValue("username")
		chanels := r.FormValue("chanels")
		return GetPostByBoth(OpenDb("./DATA/User_data.db"), usename, chanels, uid)
	}
	return nil
}

func getPostByUser(db *sql.DB, username string, uid *http.Cookie) []Post {
	output := []Post{}
	UserPost, err := db.Query("SELECT * FROM post WHERE username=?", username)
	if err != nil {
		log.Fatal("error in hash")
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var post Post
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&post.ID, &post.Uuid, &post.Username, &post.Message, &post.Document, &post.Date, &chanel, &target, &answers, &post.Like, &post.Dislike); err != nil {
			log.Fatal("error in reading", err)
		}
		post.Chanel = convertToArray(chanel)
		post.Target = convertToArray(target)
		post.Answers = convertToArray(answers)

		post.IsUserMadePost = (uid.Value == post.Uuid)
		post.IsUserLikePost, post.IsUserDislikePost = getLikedPost(db, post.ID, uid.Value)

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

func getPostByChanel(db *sql.DB, chanel string, uid *http.Cookie) []Post {
	array := strings.Split(chanel, "R/")
	slices.Sort(array)
	chanel = convertToString(array)
	output := []Post{}
	UserPost, err := db.Query("SELECT * FROM post WHERE chanel LIKE %?%", chanel)
	if err != nil {
		log.Fatal("error in hash")
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var post Post
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&post.ID, &post.Uuid, &post.PostUuid, &post.Username, &post.Message, &post.Document, &post.Ext, &post.TypeDoc, &post.Date, &chanel, &target, &answers, &post.Like, &post.Dislike); err != nil {
			log.Fatal("error in reading", err)
		}
		post.Chanel = convertToArray(chanel)
		post.Target = convertToArray(target)
		post.Answers = convertToArray(answers)

		post.IsUserMadePost = (uid.Value == post.Uuid)

		output = append(output, post)
	}
	if err = UserPost.Err(); err != nil {
		log.Fatal("error jsp ou")
	}
	return output
}

func GetPostByBoth(db *sql.DB, username string, chanel string, uid *http.Cookie) []Post {
	postList1 := getPostByUser(db, username, uid)
	if chanel == "" {
		return postList1
	}
	postList2 := getPostByChanel(db, chanel, uid)
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
