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

			Uuid, err := r.Cookie("UUID")
			if err != nil {
				if err == http.ErrNoCookie {
					log.Fatal("cookie not found Uuid")
				}
				log.Fatal("Error retrieving cookie Uuid :", err)
			}

			liked, disliked := getLikedPost(db, likes.ID, Uuid.Value)
			if liked {
				nbLikes := likes.Like - 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("UPDATE post SET like =? WHERE ID =? ", nbLikes, i)
				if err != nil {
					log.Fatal("err rows like :", err)
				}

				_, err = db.Exec("DELETE FROM like WHERE ID =? AND uuid=? ", i, Uuid.Value)
				if err != nil {
					log.Fatal("err deleting post :", err)
				}
			} else if disliked {
				nbLikes := likes.Like + 1
				nbDislikes := likes.Dislike - 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("UPDATE post SET like =?, dislike =? WHERE ID =? ", nbLikes, nbDislikes, i)
				if err != nil {
					log.Fatal("err rows like :", err)
				}

				_, err = db.Exec("UPDATE like SET liked=?, disliked=? WHERE ID=? AND uuid=?", true, false, ID, Uuid.Value)
				if err != nil {
					log.Fatal("err rows like :", err)
				}
			} else {
				nbLikes := likes.Like + 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("UPDATE post SET like =? WHERE ID =? ", nbLikes, i)
				if err != nil {
					log.Fatal("err rows like :", err)
				}

				statement, err := db.Prepare("INSERT INTO like(id, uuid, liked, disliked) VALUES(?, ?, ?, ?)")
				if err != nil {
					log.Fatal("sql add post", err)
				}
				statement.Exec(i, Uuid.Value, true, false)
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

			Uuid, err := r.Cookie("UUID")
			if err != nil {
				if err == http.ErrNoCookie {
					log.Fatal("cookie not found Uuid")
				}
				log.Fatal("Error retrieving cookie Uuid :", err)
			}

			liked, disliked := getLikedPost(db, dislikes.ID, Uuid.Value)
			if liked {
				nbLikes := dislikes.Like - 1
				nbDislikes := dislikes.Dislike + 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("UPDATE post SET like =?, dislike =? WHERE ID =? ", nbLikes, nbDislikes, i)
				if err != nil {
					log.Fatal("err rows like :", err)
				}

				_, err = db.Exec("UPDATE like SET liked=?, disliked=? WHERE ID=? AND uuid=?", false, true, ID, Uuid.Value)
				if err != nil {
					log.Fatal("err rows like :", err)
				}
			} else if disliked {
				nbDislikes := dislikes.Dislike - 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("UPDATE post SET dislike =? WHERE ID =? ", nbDislikes, i)
				if err != nil {
					log.Fatal("err rows like :", err)
				}

				_, err = db.Exec("DELETE FROM like WHERE ID =? AND uuid=? ", i, Uuid.Value)
				if err != nil {
					log.Fatal("err deleting post :", err)
				}
			} else {
				nbDislikes := dislikes.Dislike + 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Exec("UPDATE post SET dislike =? WHERE ID =? ", nbDislikes, i)
				if err != nil {
					log.Fatal("err rows like :", err)
				}

				statement, err := db.Prepare("INSERT INTO like(id, uuid, liked, disliked) VALUES(?, ?, ?, ?)")
				if err != nil {
					log.Fatal("sql add like", err)
				}
				statement.Exec(i, Uuid.Value, false, true)
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
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&output.ID, &output.Uuid, &output.Username, &output.Message, &output.Document, &output.Date, &chanel, &target, &answers, &output.Like, &output.Dislike); err != nil {
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

func getLikedPost(db *sql.DB, ID int, uuid string) (bool, bool){
	liked, err := db.Query("SELECT * FROM like WHERE id=? AND uuid=?", ID, uuid)
	if err != nil {
		log.Fatal("error in hash like",err)
	}
	defer liked.Close()
	var Id int
	var Uuid string
	var Liked bool
	var Disliked bool
	for liked.Next() {
		if err := liked.Scan(&Id, &Uuid, &Liked, &Disliked); err != nil {
			log.Fatal("error in reading", err)
		}
	}
	if err = liked.Err(); err != nil {
		log.Fatal("erreur jsp ou")
	}
	return Liked, Disliked
}