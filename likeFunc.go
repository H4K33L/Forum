package authentification

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func Like(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db := OpenDb("./DATA/User_data.db")
		defer db.Close()
		ID := r.FormValue("like")
		likes := getPostByID(db, ID)
		nbLikes := likes.Like + 1
		_, err := db.Query("UPDATE post SET like =? WHERE ID =? ", nbLikes, ID)
		if err != nil {
			log.Fatal("err rows like :", err)
		}
	}
}

func Dislike(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db := OpenDb("./DATA/User_data.db")
		defer db.Close()
		ID := r.FormValue("dislike")
		dislikes := getPostByID(db, ID)
		nbDislikes := dislikes.Dislike + 1
		_, err := db.Query("UPDATE post SET dislike =? WHERE ID =? ", nbDislikes, ID)
		if err != nil {
			log.Fatal("err rows dislike :", err)
		}
	}
}

func getPostByID(db *sql.DB, ID string) Post {
	output := Post{}
	UserPost, err := db.Query("SELECT * FROM post WHERE ID=?", ID)
	if err != nil {
		log.Fatal("error in hash")
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var UUID string
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&output.ID, &UUID, &output.Username, &output.Message, &output.Image, &output.Date, &chanel, &target, &answers, &output.Like, &output.Dislike); err != nil {
			log.Fatal("error in reading", err)
		}
		output.Chanel = convertToArray(chanel)
		output.Target = convertToArray(target)
		output.Answers = convertToArray(answers)
	}
	if err = UserPost.Err(); err != nil {
		log.Fatal("erreur jsp ou")
	}
	return output
}
