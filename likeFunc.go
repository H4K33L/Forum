package authentification

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

/*
The function http request and http respose, get the information in the 
form like and add add a like to a post, or remove it if the post is already
like by the user.

input : w http.ResponseWriter, r *http.Request

output : nothing
*/
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
					log.Fatal("Like, cookie not found Uuid", err)
				}
				log.Fatal("Like, Error retrieving cookie Uuid :", err)
			}

			liked, disliked := getLikedPost(db, likes.ID, Uuid.Value)
			if liked {
				nbLikes := likes.Like - 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal("Like, error during Atoi conversion :", err)
				}
				_, err = db.Exec("UPDATE post SET like =? WHERE ID =? ", nbLikes, i)
				if err != nil {
					log.Fatal("Like, err rows like :", err)
				}

				_, err = db.Exec("DELETE FROM like WHERE ID =? AND uuid=? ", i, Uuid.Value)
				if err != nil {
					log.Fatal("Like, err deleting post :", err)
				}
			} else if disliked {
				nbLikes := likes.Like + 1
				nbDislikes := likes.Dislike - 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal("Like, error during Atoi conversion :", err)
				}
				_, err = db.Exec("UPDATE post SET like =?, dislike =? WHERE ID =? ", nbLikes, nbDislikes, i)
				if err != nil {
					log.Fatal("Like, err rows like :", err)
				}

				_, err = db.Exec("UPDATE like SET liked=?, disliked=? WHERE ID=? AND uuid=?", true, false, ID, Uuid.Value)
				if err != nil {
					log.Fatal("Like, err rows like :", err)
				}
			} else {
				nbLikes := likes.Like + 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal("Like, error during Atoi conversion :", err)
				}
				_, err = db.Exec("UPDATE post SET like =? WHERE ID =? ", nbLikes, i)
				if err != nil {
					log.Fatal("Like, err rows like :", err)
				}

				statement, err := db.Prepare("INSERT INTO like(id, uuid, liked, disliked) VALUES(?, ?, ?, ?)")
				if err != nil {
					log.Fatal("Like, sql add post", err)
				}
				statement.Exec(i, Uuid.Value, true, false)
			}
		}
	}
}

/*
The function http request and http respose, get the information in the 
form dislike and add add a dislike to a post, or remove it if the post is already
dislike by the user.

input : w http.ResponseWriter, r *http.Request

output : nothing
*/
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
					log.Fatal("Dislike, cookie not found Uuid", err)
				}
				log.Fatal("Dislike, Error retrieving cookie Uuid :", err)
			}

			liked, disliked := getLikedPost(db, dislikes.ID, Uuid.Value)
			if liked {
				nbLikes := dislikes.Like - 1
				nbDislikes := dislikes.Dislike + 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal("Dislike, error during Atoi conversion :", err)
				}
				_, err = db.Exec("UPDATE post SET like =?, dislike =? WHERE ID =? ", nbLikes, nbDislikes, i)
				if err != nil {
					log.Fatal("Dislike, err rows like :", err)
				}

				_, err = db.Exec("UPDATE like SET liked=?, disliked=? WHERE ID=? AND uuid=?", false, true, ID, Uuid.Value)
				if err != nil {
					log.Fatal("Dislike, err rows like :", err)
				}
			} else if disliked {
				nbDislikes := dislikes.Dislike - 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal("Dislike, error during Atoi conversion :", err)
				}
				_, err = db.Exec("UPDATE post SET dislike =? WHERE ID =? ", nbDislikes, i)
				if err != nil {
					log.Fatal("Dislike, err rows like :", err)
				}

				_, err = db.Exec("DELETE FROM like WHERE ID =? AND uuid=? ", i, Uuid.Value)
				if err != nil {
					log.Fatal("Dislike, err deleting post :", err)
				}
			} else {
				nbDislikes := dislikes.Dislike + 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal("Dislike, error during Atoi conversion :", err)
				}
				_, err = db.Exec("UPDATE post SET dislike =? WHERE ID =? ", nbDislikes, i)
				if err != nil {
					log.Fatal("Dislike, err rows like :", err)
				}

				statement, err := db.Prepare("INSERT INTO like(id, uuid, liked, disliked) VALUES(?, ?, ?, ?)")
				if err != nil {
					log.Fatal("Dislike, sql add like", err)
				}
				statement.Exec(i, Uuid.Value, false, true)
			}
		}
	}
}

/*
The function take a DB, a string reprsenting the ID and return the post coresponding whith the ID.

input : db *sql.DB, ID string

output : a Post
*/
func getPostByID(db *sql.DB, ID string) Post {
	output := Post{}
	UserPost, err := db.Query("SELECT * FROM post WHERE ID=?", ID)
	if err != nil {
		log.Fatal("getPostByID, error in hash :", err)
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var chanel string
		var target string
		var answers string
		if err := UserPost.Scan(&output.ID, &output.Uuid, &output.Username, &output.Message, &output.Document, &output.Date, &chanel, &target, &answers, &output.Like, &output.Dislike); err != nil {
			log.Fatal("getPostByID, error in reading :", err)
		}
		output.Chanel = convertToArray(chanel)
		output.Target = convertToArray(target)
		output.Answers = convertToArray(answers)
	}
	if err = UserPost.Err(); err != nil {
		log.Fatal("getPostByID, error in reading :", err)
	}
	return output
}

/*
The function take a DB, a int representing the ID and a string representing th uuid, the function
get the in like db table coresponding whith ID and uuid and return the two boolean representing 
if the user like or dislike the post.

input : db *sql.DB, ID int, uuid string

output : boolean , boolean
*/
func getLikedPost(db *sql.DB, ID int, uuid string) (bool, bool){
	liked, err := db.Query("SELECT * FROM like WHERE id=? AND uuid=?", ID, uuid)
	if err != nil {
		log.Fatal("getPostByID, error in hash like :",err)
	}
	defer liked.Close()
	var Id int
	var Uuid string
	var Liked bool
	var Disliked bool
	for liked.Next() {
		if err := liked.Scan(&Id, &Uuid, &Liked, &Disliked); err != nil {
			log.Fatal("getPostByID, error in reading :", err)
		}
	}
	if err = liked.Err(); err != nil {
		log.Fatal("getPostByID, error in reading :", err)
	}
	return Liked, Disliked
}