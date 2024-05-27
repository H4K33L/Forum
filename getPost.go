package authentification

import (
	"database/sql"
	"strings"

	"log"

	_ "github.com/mattn/go-sqlite3"

)

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