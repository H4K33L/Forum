package authentification

import (
	"database/sql"
	"log"
	"net/http"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

/*
The function take http request and http respose, get the information in the 
form getpost and use GetPostByBoth to return an array of Post struct.

input : w http.ResponseWriter, r *http.Request

output : array of Post struct
*/
func GetPost(w http.ResponseWriter, r *http.Request) []Post {
	if r.Method == "GET" {
		uid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Fatal("GetPost, cookie not found userpost :", err)
			}
			log.Fatal("GetPost, Error retrieving cookie uuid :", err)
		}

		usename := r.FormValue("username")
		chanels := r.FormValue("chanels")
		return GetPostByBoth(OpenDb("./DATA/User_data.db"), usename, chanels, uid)
	}
	return nil
}

/*
The function take a DB, a string reprsenting the username and the user uuid
the function get all post whith coresponding username and return it in array of Post.

input : db *sql.DB, username string, uid *http.Cookie

output : array of Post struct
*/
func getPostByUser(db *sql.DB, username string, uid *http.Cookie) []Post {
	output := []Post{}
	UserPost, err := db.Query("SELECT * FROM post WHERE username=?", username)
	if err != nil {
		log.Fatal("getPostByUser, error in hash :", err)
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var post Post
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&post.ID, &post.Uuid, &post.Username, &post.Message, &post.Document, &post.Date, &chanel, &target, &answers, &post.Like, &post.Dislike); err != nil {
			log.Fatal("getPostByUser, error in reading :", err)
		}
		post.Chanel = convertToArray(chanel)
		post.Target = convertToArray(target)
		post.Answers = convertToArray(answers)

		post.IsUserMadePost = (uid.Value == post.Uuid)
		post.IsUserLikePost, post.IsUserDislikePost = getLikedPost(db, post.ID, uid.Value)

		output = append(output, post)
	}
	if err = UserPost.Err(); err != nil {
		log.Fatal("getPostByUser, error in reading :", err)
	}
	return output
}

/*
The function take a string and convert it to array, the separator is "|\\/|-_-|\\/|+{}", 
this function is only used for convert string coming out of the DB.

input : str string

output : array of string
*/
func convertToArray(str string) []string {
	return strings.Split(str, "|\\/|-_-|\\/|+{}")
}

/*
The function take a DB, a string reprsenting the chanels and the user uuid
the function get all post whith coresponding chanels and return it in array of Post.

input : db *sql.DB, chanel string, uid *http.Cookie

output : array of Post struct
*/
func getPostByChanel(db *sql.DB, chanel string, uid *http.Cookie) []Post {
	array := strings.Split(chanel, "R/")
	sort.Strings(array)
	chanel = convertToString(array)
	output := []Post{}
	UserPost, err := db.Query("SELECT * FROM post WHERE chanel LIKE ?", "%"+chanel+"%")
	if err != nil {
		log.Fatal("getPostByChanel, error in Get post request :", err)
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var post Post
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&post.ID, &post.Uuid, &post.PostUuid, &post.Username, &post.Message, &post.Document, &post.Ext, &post.TypeDoc, &post.Date, &chanel, &target, &answers, &post.Like, &post.Dislike); err != nil {
			log.Fatal("getPostByChanel, error in reading", err)
		}
		post.Chanel = convertToArray(chanel)
		post.Target = convertToArray(target)
		post.Answers = convertToArray(answers)

		post.IsUserMadePost = (uid.Value == post.Uuid)

		output = append(output, post)
	}
	if err = UserPost.Err(); err != nil {
		log.Fatal("getPostByChanel, error in reading", err)
	}
	return output
}

/*
The function take a db, string representing name and chanel, the uuid user
and use getPostByChanel and getPostByUser to return an array of Post.

input : db *sql.DB, username string, chanel string, uid *http.Cookie

output : array of Post struct
*/
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

/*
The function take two Post struct and compare there ID and return a boolean.

input : post1 Post, post2 Post

output : boolean
*/
func comparePost(post1 Post, post2 Post) bool {
	return post1.ID == post2.ID
}
