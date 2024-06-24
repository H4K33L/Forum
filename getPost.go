package client

import (
	"database/sql"
	"fmt"
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
	// Check if the request method is GET
	if r.Method == "GET" {
		// Retrieve the UUID cookie from the request
		uid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("getpostt GetPost, cookie not found userpost :", err)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return nil
			}
			fmt.Println("getpost GetPost, Error retrieving cookie uuid :", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return nil
		}

		// Retrieve the username and channels from the request form
		usename := r.FormValue("username")
		chanels := r.FormValue("chanels")

		// Return the posts obtained by the given username, channels, and UUID
		return GetPostByBoth(OpenDb("./DATA/User_data.db", w, r), usename, chanels, uid, w, r)
	}
	return nil
}

/*
The function take a DB, a string reprsenting the username and the user uuid
the function get all post whith coresponding username and return it in array of Post.

input : db *sql.DB, username string, uid *http.Cookie

output : array of Post struct
*/
func getPostByUser(db *sql.DB, username string, uid *http.Cookie, w http.ResponseWriter, r *http.Request) []Post {
	output := []Post{}
	// Query posts by the given username
	UserPost, err := db.Query("SELECT * FROM post WHERE username=?", username)
	if err != nil {
		fmt.Println("getpost getPostByUser, error in hash :", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return nil
	}
	defer UserPost.Close()
	for UserPost.Next() {
		var post Post
		var chanel string
		var target string
		if err := UserPost.Scan(&post.ID, &post.Uuid, &post.Username, &post.Message, &post.Document, &post.Date, &chanel, &target, &post.Like, &post.Dislike); err != nil {
			fmt.Println("getpostt getPostByUser, error in reading :", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return nil
		}
		post.Chanel = convertToArray(chanel)
		post.Target = convertToArray(target)

		// Check if the current user made the post
		post.IsUserMadePost = (uid.Value == post.Uuid)
		// Check if the current user liked or disliked the post
		post.IsUserLikePost, post.IsUserDislikePost = getLikedPost(db, post.ID, uid.Value, w, r)

		output = append(output, post)
	}
	if err = UserPost.Err(); err != nil {
		fmt.Println("getpost getPostByUser, error in reading :", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return nil
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
func getPostByChanel(db *sql.DB, chanel string, uid *http.Cookie, w http.ResponseWriter, r *http.Request) []Post {
	// Split the channel string into an array of channel names.
	array := strings.Split(chanel, "R/")
	// Sort the array of channel names.
	sort.Strings(array)
	// Convert the sorted array of channel names back to a string.
	chanel = convertToString(array)
	// Initialize an empty slice to store the retrieved posts.
	output := []Post{}
	// Query the database to select posts where the channel column contains the specified channel string.
	UserPost, err := db.Query("SELECT * FROM post WHERE chanel LIKE ?", "%"+chanel+"%")
	if err != nil {
		// Log a fatal error if there is an issue with the database query.
		fmt.Println("getpost getPostByChanel, error in Get post request :", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return nil
	}
	// Defer closing the UserPost rows after the function returns.
	defer UserPost.Close()
	// Iterate over each row in the UserPost result set.
	for UserPost.Next() {
		// Initialize a Post struct to store the data of each post.
		var post Post
		// Initialize variables to store the channel and target strings.
		var chanel string
		var target string
		// Scan the columns of the current row into the fields of the Post struct.
		if err := UserPost.Scan(&post.ID, &post.Uuid, &post.PostUuid, &post.Username, &post.Message, &post.Document, &post.Ext, &post.TypeDoc, &post.Date, &chanel, &target, &post.Like, &post.Dislike); err != nil {
			// Log a fatal error if there is an issue scanning the row.
			fmt.Println("getpost getPostByChanel, error in reading :", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return nil
		}
		// Convert the channel and target strings into arrays and assign them to the corresponding fields of the Post struct.
		post.Chanel = convertToArray(chanel)
		post.Target = convertToArray(target)
		// Check if the current post was made by the logged-in user and assign the result to the IsUserMadePost field of the Post struct.
		post.IsUserMadePost = (uid.Value == post.Uuid)
		// Append the current post to the output slice.
		output = append(output, post)
	}
	// Check for any errors that occurred during iteration over the UserPost result set.
	if err = UserPost.Err(); err != nil {
		// Log a fatal error if there was an error.
		fmt.Println("getpost getPostByChanel, error in reading :", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return nil
	}
	// Return the slice containing the retrieved posts.
	return output
}

/*
The function take a db, string representing name and chanel, the uuid user
and use getPostByChanel and getPostByUser to return an array of Post.

input : db *sql.DB, username string, chanel string, uid *http.Cookie

output : array of Post struct
*/
func GetPostByBoth(db *sql.DB, username string, chanel string, uid *http.Cookie, w http.ResponseWriter, r *http.Request) []Post {
	// Retrieve posts by username.
	postList1 := getPostByUser(db, username, uid, w, r)
	// If no channel filter is provided, return posts filtered by username only.
	if chanel == "" {
		return postList1
	}
	// Retrieve posts by channel.
	postList2 := getPostByChanel(db, chanel, uid, w, r)
	// If no username filter is provided, return posts filtered by channel only.
	if username == "" {
		return postList2
	}
	// Initialize an empty slice to store the posts that meet both username and channel criteria.
	output := []Post{}
	// Iterate over each post from the first list.
	for _, post1 := range postList1 {
		// Iterate over each post from the second list.
		for _, post2 := range postList2 {
			// Check if the current post from the first list matches the current post from the second list.
			if comparePost(post1, post2) {
				// If the posts match, add the current post from the first list to the output slice.
				output = append(output, post1)
			}
		}
	}
	// Return the slice containing the posts that meet both username and channel criteria.
	return output
}

/*
The function take a DB, a string reprsenting the ID and return the post coresponding whith the ID.

input : db *sql.DB, ID string

output : a Post
*/
func getPostByID(db *sql.DB, ID string, w http.ResponseWriter, r *http.Request) Post {
	output := Post{}

	// Query the post table to retrieve the post details based on its ID
	UserPost, err := db.Query("SELECT * FROM post WHERE ID=?", ID)
	if err != nil {
		fmt.Println("getPost getPostByID, error in hash :", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return output
	}
	defer UserPost.Close()

	// Iterate over the query results
	for UserPost.Next() {
		var chanel string
		var target string
		var answers string
		// Scan the query results into the output struct
		if err := UserPost.Scan(&output.ID, &output.Uuid, &output.Username, &output.Message, &output.Document, &output.Date, &chanel, &target, &answers, &output.Like, &output.Dislike); err != nil {
			fmt.Println("getPost getPostByID, error in reading :", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return output
		}
		// Convert the channel and target strings to arrays
		output.Chanel = convertToArray(chanel)
		output.Target = convertToArray(target)
	}

	// Check for any errors during iteration
	if err = UserPost.Err(); err != nil {
		fmt.Println("getPost getPostByID, error in reading :", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return output
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
func getLikedPost(db *sql.DB, ID int, uuid string, w http.ResponseWriter, r *http.Request) (bool, bool) {
	// Query the like table to check if the post with the given ID and UUID has been liked or disliked
	liked, err := db.Query("SELECT * FROM like WHERE id=? AND uuid=?", ID, uuid)
	if err != nil {
		fmt.Println("getPost getPostByID, error in hash like :", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return false, false
	}
	defer liked.Close()

	var Liked bool
	var Disliked bool

	// Iterate over the query results
	for liked.Next() {
		var Id int
		var Uuid string
		// Scan the query results into variables
		if err := liked.Scan(&Id, &Uuid, &Liked, &Disliked); err != nil {
			fmt.Println("getPost getPostByID, error in reading :", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return false, false
		}
	}

	// Check for any errors during iteration
	if err = liked.Err(); err != nil {
		fmt.Println("getPost getPostByID, error in reading :", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return false, false
	}

	return Liked, Disliked
}
