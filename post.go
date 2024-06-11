package authentification

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
    "sort"
	"log"
	"net/http"
	"os"
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
	Answers           []string
}

/*
The function http request and http respose, get the information in the 
form Post to create a post struct and user the function AddPost to add 
the post to db.

input : w http.ResponseWriter, r *http.Request

output : nothing
*/
func UserPost(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("./DATA/User_data.db")
	defer db.Close()
	var post Post
	if r.Method == "POST" {
		uid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Fatal("UserPost, cookie not found userpost :", err)
			}
			log.Fatal("UserPost, Error retrieving cookie uuid :", err)
		}
		var username string
		err1 := db.QueryRow("SELECT username FROM user WHERE uuid=?", uid.Value).Scan(&username)
		if err1 != nil {
			if err1 == sql.ErrNoRows {
				log.Fatal("UserPost, sql user post :", err1)
			}
			log.Fatal(err1)
		}
		post.Uuid = uid.Value
		post.Username = username
		post.Message = r.FormValue("message")
		u, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("UserPost, failed to generate UUID: %v", err)
		}
		post.PostUuid = u.String()
		if r.FormValue("type") == "image" {
			post.TypeDoc = "image"
			if r.FormValue("typedoc") == "file" {
				file, handler, err := r.FormFile("documentFile")
				if err != nil {
					if err == http.ErrMissingFile {
						fmt.Println("no file uploaded")
						post.Document = ""
					} else {
						log.Fatal("UserPost, post image:", err)
					}

				} else {
					extension := strings.LastIndex(handler.Filename, ".")
					if extension == -1 {
						fmt.Println("there is no extension to the file")
					} else {
						ext := handler.Filename[extension:]
						e := strings.ToLower(ext)
						if e == ".png" || e == ".jpeg" || e == ".jpg" || e == ".gif" || e == ".svg" || e == ".avif" || e == ".apng" || e == ".webp" {
							path := "/static/stylsheet/IMAGES/POST/" + post.PostUuid + ext
							if _, err := os.Stat("./VIEWS" + path); errors.Is(err, os.ErrNotExist) {
								// file does not exist
							} else {
								e := os.Remove("./VIEWS" + path)
								if e != nil {
									log.Fatal(e)
								}
							}

							f, err := os.OpenFile("./VIEWS"+path, os.O_WRONLY|os.O_CREATE, 0666)
							if err != nil {
								fmt.Println(err)
								return
							}
							defer f.Close()
							_, err = io.Copy(f, file)
							if err != nil {
								fmt.Println(err)
								return
							}
							post.Document = path
							post.Ext = "file"
						}
					}
				}
			} else {
				post.Document = r.FormValue("document")
				post.Ext = "url"
			}
		} else if r.FormValue("type") == "video" {
			post.TypeDoc = "video"
			if r.FormValue("typedoc") == "file" {
				file, handler, err := r.FormFile("documentFile")
				if err != nil {
					if err == http.ErrMissingFile {
						fmt.Println("no file uploaded")
						post.Document = ""
					} else {
						log.Fatal("post video:", err)
					}
				} else {
					extension := strings.LastIndex(handler.Filename, ".")
					if extension == -1 {
						fmt.Println("there is no extension to the file")
					} else {
						ext := handler.Filename[extension:]
						e := strings.ToLower(ext)
						if e == ".mp4" || e == ".webm" || e == ".ogg" {
							path := "/static/stylsheet/VIDEO/" + post.PostUuid + ext
							if _, err := os.Stat("./VIEWS" + path); errors.Is(err, os.ErrNotExist) {
								// file does not exist
							} else {
								e := os.Remove("./VIEWS" + path)
								if e != nil {
									log.Fatal(e)
								}
							}

							f, err := os.OpenFile("./VIEWS"+path, os.O_WRONLY|os.O_CREATE, 0666)
							if err != nil {
								fmt.Println(err)
								return
							}
							defer f.Close()
							io.Copy(f, file)
							post.Document = path
							post.Ext = "file"
						}
					}
				}
			} else {
				post.Document = r.FormValue("document")
				post.Ext = "url"
			}
		} else {
			post.TypeDoc = ""
			post.Document = ""
			post.Ext = ""
		}

		then := time.Now()
		post.Date = strconv.Itoa(then.Year()) + "/" + then.Month().String() + "/" + strconv.Itoa(then.Day()) + " " + strconv.Itoa(then.Hour()) + ":" + strconv.Itoa(then.Minute())
		post.Chanel = strings.Split(r.FormValue("chanel"), "R/")
		post.Target = strings.Split(r.FormValue("target"), "\\\\")
		if r.FormValue("chanel") != "" {
			db := OpenDb("./DATA/User_data.db")
			defer db.Close()
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
	statement, err := db.Prepare("INSERT INTO post(uuid, postuuid, username, message, document, ext, typedoc, date, chanel, target, answers, like, dislike) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("AddPost, sql add post", err)
	}
	chanel := convertToString(post.Chanel)
	target := convertToString(post.Target)
	answers := convertToString(post.Answers)
	statement.Exec(post.Uuid, post.PostUuid, post.Username, post.Message, post.Document, post.Ext, post.TypeDoc, post.Date, chanel, target, answers, post.Like, post.Dislike)
}

/*
The function take an array of sting and convert it to sting whith "|\\/|-_-|\\/|+{}" joiner,
this function is only suposed to be used to convert array of string in string to be stocked in db.

input : array []string

output : string
*/
func convertToString(array []string) string {
	sort.Strings(array)
	return strings.Join(array, "|\\/|-_-|\\/|+{}")
}
