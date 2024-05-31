package authentification

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func Like(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("./DATA/User_data.db")
	if r.Method == "Post" {
		ID := r.FormValue("like")
		likes := getPostByID(db, ID)

	}
}

func Dislike(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("./DATA/User_data.db")
	if r.Method == "post" {
		ID := r.FormValue("dislike")
		likes := getPostByID(db, ID)

	}
}

func getPostByID(db *sql.DB, ID string) []Post {
	output := Post
	UserPost, err := db.Query("SELECT * FROM post WHERE ID=?", ID)
	if err != nil {
		log.Fatal("error in hash")
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var ID int
		var UUID string
		var post Post
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&ID, &UUID, &post.Username, &post.Message, &post.Image, &post.Date, &chanel, &target, &answers, &post.Like, &post.Dislike); err != nil {
			log.Fatal("error in reading", err)
		}
		post.Chanel = convertToArray(chanel)
		post.Target = convertToArray(target)
		post.Answers = convertToArray(answers)
	}
	if err = UserPost.Err(); err != nil {
		log.Fatal("erreur jsp ou")
	}
	return output
}
