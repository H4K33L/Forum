package client

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

/*
This structure is used acros the project to represent a post.
*/
type Post struct {
	ID                int
	Uuid              string
	PostUuid          string
	IsUserMadePost    bool
	IsUserLikePost    bool
	IsUserDislikePost bool
	Username          string
	Message           string
	Document          string
	TypeDoc           string
	Ext               string
	Date              string
	Like              int
	Dislike           int
	Chanel            []string
	Target            []string
}

/*
The function http request and http respose, get the information in the
form Post to create a post struct and user the function AddPost to add
the post to db.

input : w http.ResponseWriter, r *http.Request

output : nothing
*/
func UserPost(w http.ResponseWriter, r *http.Request) {
	// Open a connection to the user database
	db := OpenDb("./DATA/User_data.db")
	defer db.Close() // Ensure the database connection is closed when the function exits
	var post Post    // Define a variable to hold the post data
	// Check if the request method is POST
	if r.Method == "POST" {
		// Retrieve the UUID cookie from the request
		uid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Fatal(" post UserPost, cookie not found userpost :", err)
			}
			log.Fatal("post UserPost, Error retrieving cookie uuid :", err)
		}
		var username string
		// Query the user table to retrieve the username based on the UUID
		err1 := db.QueryRow("SELECT username FROM user WHERE uuid=?", uid.Value).Scan(&username)
		if err1 != nil {
			if err1 == sql.ErrNoRows {
				log.Fatal("post UserPost, sql user post :", err1)
			}
			log.Fatal(err1)
		}
		// Fill the post struct with data from the request
		post.Uuid = uid.Value
		post.Username = username
		post.Message = r.FormValue("message")
		u, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("post UserPost, failed to generate UUID: %v", err)
		}
		// Generate a new UUID for the post
		post.PostUuid = u.String()
		// Check if the post type is image
		if r.FormValue("type") == "image" {
			post.TypeDoc = "image"
			// Check if the image is uploaded as a file
			if r.FormValue("typedoc") == "file" {
				// Retrieve the file from the form data
				file, handler, err := r.FormFile("documentFile")
				if err != nil {
					// Handle the case where no file is uploaded
					if err == http.ErrMissingFile {
						post.Document = ""
					} else {
						log.Fatal("post UserPost, post image:", err)
					}
				} else {
					// Extract the file extension and validate it
					extension := strings.LastIndex(handler.Filename, ".")
					if extension == -1 {
						log.Fatal("post UserPost image : there is no extension to the file")
					} else {
						ext := handler.Filename[extension:]
						e := strings.ToLower(ext)
						// Check if the file extension is valid for an image
						if e == ".png" || e == ".jpeg" || e == ".jpg" || e == ".gif" || e == ".svg" || e == ".avif" || e == ".apng" || e == ".webp" {
							// Define the path for storing the image file
							path := "/static/stylsheet/IMAGES/POST/" + post.PostUuid + ext
							// Check if the file already exists and remove it
							if _, err := os.Stat("./VIEWS" + path); errors.Is(err, os.ErrNotExist) {
								log.Fatal("post userpost no extension os.error :", err)
							} else {
								err = os.Remove("./VIEWS" + path)
								if err != nil {
									log.Fatal("post userpost can't remove the path :", err)
								}
							}

							// Create and open the file for writing
							f, err := os.OpenFile("./VIEWS"+path, os.O_WRONLY|os.O_CREATE, 0666)
							if err != nil {
								log.Fatal("post userpost can't open the file :", err)
							}
							defer f.Close()
							// Copy the uploaded file data to the destination file
							_, err = io.Copy(f, file)
							if err != nil {
								log.Fatal("post userpost can't copy the file :", err)
							}
							// Set the post's document path and extension
							post.Document = path
							post.Ext = "file"
						}
					}
				}
			} else {
				// If the image is provided as a URL, set the document URL and extension accordingly
				post.Document = r.FormValue("document")
				post.Ext = "url"
			}
		} else // Check if the post type is video
		if r.FormValue("type") == "video" {
			post.TypeDoc = "video"
			// Check if the video is uploaded as a file
			if r.FormValue("typedoc") == "file" {
				// Retrieve the file from the form data
				file, handler, err := r.FormFile("documentFile")
				if err != nil {
					// Handle the case where no file is uploaded
					if err == http.ErrMissingFile {
						post.Document = ""
					} else {
						log.Fatal("post userpost video :", err)
					}
				} else {
					// Extract the file extension and validate it
					extension := strings.LastIndex(handler.Filename, ".")
					if extension == -1 {
						fmt.Println("post user post video : there is no extension to the file")
					} else {
						ext := handler.Filename[extension:]
						e := strings.ToLower(ext)
						// Check if the file extension is valid for a video
						if e == ".mp4" || e == ".webm" || e == ".ogg" {
							// Define the path for storing the video file
							path := "/static/stylsheet/VIDEO/" + post.PostUuid + ext
							// Check if the file already exists and remove it
							if _, err := os.Stat("./VIEWS" + path); errors.Is(err, os.ErrNotExist) {
								log.Fatal("post userpost no extention video :", err)
							} else {
								err = os.Remove("./VIEWS" + path)
								if err != nil {
									log.Fatal("post userpost can't remove the path :", err)
								}
							}

							// Create and open the file for writing
							f, err := os.OpenFile("./VIEWS"+path, os.O_WRONLY|os.O_CREATE, 0666)
							if err != nil {
								log.Fatal("post userpost ", err)
								return
							}
							defer f.Close()
							// Copy the uploaded file data to the destination file
							io.Copy(f, file)
							// Set the post's document path and extension
							post.Document = path
							post.Ext = "file"
						}
					}
				}
			} else {
				// If the video is provided as a URL, set the document URL and extension accordingly
				post.Document = r.FormValue("document")
				post.Ext = "url"
			}
		} else {
			// If the post type is neither image nor video, set document-related fields to empty
			post.TypeDoc = ""
			post.Document = ""
			post.Ext = ""
		}

		// Record the current date and time of the post
		then := time.Now()
		post.Date = strconv.Itoa(then.Year()) + "/" + then.Month().String() + "/" + strconv.Itoa(then.Day()) + " " + strconv.Itoa(then.Hour()) + ":" + strconv.Itoa(then.Minute())

		// Split channel and target values from the form
		post.Chanel = strings.Split(r.FormValue("chanel"), "R/")
		post.Target = strings.Split(r.FormValue("target"), "\\\\")

		// Check if a channel is specified before adding the post to the database
		if r.FormValue("chanel") != "" {
			// Open a connection to the user database
			db := OpenDb("./DATA/User_data.db")
			defer db.Close()
			// Add the post to the database
			AddPost(db, post)
		}
	}
}

/*
The function take db an post as argument and insert post in db.

input : db *sql.DB, post Post

output : none
*/
func AddPost(db *sql.DB, post Post) {
	// Prepare the SQL statement for inserting a new post into the database
	statement, err := db.Prepare("INSERT INTO post(uuid, postuuid, username, message, document, ext, typedoc, date, chanel, target, like, dislike) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("post AddPost, sql add post :", err)
	}
	// Convert channel and target slices to string format
	chanel := convertToString(post.Chanel)
	target := convertToString(post.Target)
	// Execute the SQL statement with post data
	statement.Exec(post.Uuid, post.PostUuid, post.Username, post.Message, post.Document, post.Ext, post.TypeDoc, post.Date, chanel, target, post.Like, post.Dislike)
}

/*
The function take an array of sting and convert it to sting whith "|\\/|-_-|\\/|+{}" joiner,
this function is only suposed to be used to convert array of string in string to be stocked in db.

input : array []string

output : string
*/
func convertToString(array []string) string {
	// Sort the elements of the array alphabetically
	sort.Strings(array)
	// Join the sorted array elements into a single string using a custom separator
	return strings.Join(array, "|\\/|-_-|\\/|+{}")
}
