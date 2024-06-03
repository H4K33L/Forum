package authentification

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func Like(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ID := r.FormValue("like")
		if ID != "" {
			db := OpenDb("./DATA/User_data.db")
			defer db.Close()
			likes := getPostByID(db, ID)
			nbLikes := likes.Like + 1
			i, err := strconv.Atoi(ID)
    		if err != nil {
        		log.Fatal(err)
    		}
			_, err = db.Exec("UPDATE post SET like =? WHERE ID =? ", nbLikes, i)
			if err != nil {
				log.Fatal("err rows like :", err)
			}
		}
	}
}

func Dislike(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ID := r.FormValue("dislike")
		if ID != "" {
			db := OpenDb("./DATA/User_data.db")
			defer db.Close()
			dislikes := getPostByID(db, ID)
			nbDislikes := dislikes.Dislike + 1
			i, err := strconv.Atoi(ID)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("UPDATE post SET dislike =? WHERE ID =? ", nbDislikes, i)
			if err != nil {
				log.Fatal("err rows dislike :", err)
			}
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
