package client

import (
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
	// Check if the request method is POST
	if r.Method == "POST" {
		// Retrieve the post ID from the form data
		ID := r.FormValue("like")
		// Proceed if the post ID is not empty
		if ID != "" {
			// Open a connection to the user database
			db := OpenDb("./DATA/User_data.db")
			defer db.Close()

			// Get the post details based on its ID
			likes := getPostByID(db, ID)

			// Retrieve the UUID cookie from the request
			Uuid, err := r.Cookie("UUID")
			if err != nil {
				if err == http.ErrNoCookie {
					log.Fatal("likefunc Like, cookie not found Uuid", err)
				}
				log.Fatal("likefunc Like, Error retrieving cookie Uuid :", err)
			}

			// Check if the user has already liked or disliked the post
			liked, disliked := getLikedPost(db, likes.ID, Uuid.Value)

			// If the user has already liked the post, undo the like
			if liked {
				nbLikes := likes.Like - 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal("likefunc Like, error during Atoi conversion :", err)
				}
				// Update the post's like count
				_, err = db.Exec("UPDATE post SET like =? WHERE ID =? ", nbLikes, i)
				if err != nil {
					log.Fatal("likefunc Like, err rows like :", err)
				}

				// Remove the like entry from the like table
				_, err = db.Exec("DELETE FROM like WHERE ID =? AND uuid=? ", i, Uuid.Value)
				if err != nil {
					log.Fatal("likefunc Like, err deleting post :", err)
				}
				// If the user has disliked the post, change the dislike to like
			} else if disliked {
				nbLikes := likes.Like + 1
				nbDislikes := likes.Dislike - 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal("likefunc Like, error during Atoi conversion :", err)
				}
				// Update the post's like and dislike counts
				_, err = db.Exec("UPDATE post SET like =?, dislike =? WHERE ID =? ", nbLikes, nbDislikes, i)
				if err != nil {
					log.Fatal("likefunc Like, err rows like :", err)
				}

				// Update the like entry in the like table
				_, err = db.Exec("UPDATE like SET liked=?, disliked=? WHERE ID=? AND uuid=?", true, false, ID, Uuid.Value)
				if err != nil {
					log.Fatal("likefunc Like, err rows like :", err)
				}
				// If the user has not liked or disliked the post, like it
			} else {
				nbLikes := likes.Like + 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal("likefunc Like, error during Atoi conversion :", err)
				}
				// Update the post's like count
				_, err = db.Exec("UPDATE post SET like =? WHERE ID =? ", nbLikes, i)
				if err != nil {
					log.Fatal("likefunc Like, err rows like :", err)
				}

				// Insert a new like entry in the like table
				statement, err := db.Prepare("INSERT INTO like(id, uuid, liked, disliked) VALUES(?, ?, ?, ?)")
				if err != nil {
					log.Fatal("likefunc Like, sql add post", err)
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
	// Check if the request method is POST
	if r.Method == "POST" {
		// Retrieve the post ID from the form data
		ID := r.FormValue("dislike")
		// Proceed if the post ID is not empty
		if ID != "" {
			// Open a connection to the user database
			db := OpenDb("./DATA/User_data.db")
			defer db.Close()

			// Get the post details based on its ID
			dislikes := getPostByID(db, ID)

			// Retrieve the UUID cookie from the request
			Uuid, err := r.Cookie("UUID")
			if err != nil {
				if err == http.ErrNoCookie {
					log.Fatal("likefunc Dislike, cookie not found Uuid", err)
				}
				log.Fatal("likefunc Dislike, Error retrieving cookie Uuid :", err)
			}

			// Check if the user has already liked or disliked the post
			liked, disliked := getLikedPost(db, dislikes.ID, Uuid.Value)

			// If the user has already liked the post, change the like to dislike
			if liked {
				nbLikes := dislikes.Like - 1
				nbDislikes := dislikes.Dislike + 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal("likefunc Dislike, error during Atoi conversion :", err)
				}
				// Update the post's like and dislike counts
				_, err = db.Exec("UPDATE post SET like =?, dislike =? WHERE ID =? ", nbLikes, nbDislikes, i)
				if err != nil {
					log.Fatal("likefunc Dislike, err rows like :", err)
				}

				// Update the like entry in the like table
				_, err = db.Exec("UPDATE like SET liked=?, disliked=? WHERE ID=? AND uuid=?", false, true, ID, Uuid.Value)
				if err != nil {
					log.Fatal("likefunc Dislike, err rows like :", err)
				}
				// If the user has already disliked the post, undo the dislike
			} else if disliked {
				nbDislikes := dislikes.Dislike - 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal("likefunc Dislike, error during Atoi conversion :", err)
				}
				// Update the post's dislike count
				_, err = db.Exec("UPDATE post SET dislike =? WHERE ID =? ", nbDislikes, i)
				if err != nil {
					log.Fatal("likefunc Dislike, err rows like :", err)
				}

				// Remove the like entry from the like table
				_, err = db.Exec("DELETE FROM like WHERE ID =? AND uuid=? ", i, Uuid.Value)
				if err != nil {
					log.Fatal("likefunc Dislike, err deleting post :", err)
				}
				// If the user has not liked or disliked the post, dislike it
			} else {
				nbDislikes := dislikes.Dislike + 1
				i, err := strconv.Atoi(ID)
				if err != nil {
					log.Fatal("likefunc Dislike, error during Atoi conversion :", err)
				}
				// Update the post's dislike count
				_, err = db.Exec("UPDATE post SET dislike =? WHERE ID =? ", nbDislikes, i)
				if err != nil {
					log.Fatal("likefunc Dislike, err rows like :", err)
				}

				// Insert a new like entry in the like table
				statement, err := db.Prepare("INSERT INTO like(id, uuid, liked, disliked) VALUES(?, ?, ?, ?)")
				if err != nil {
					log.Fatal("likefunc Dislike, sql add like", err)
				}
				statement.Exec(i, Uuid.Value, false, true)
			}
		}
	}
}
