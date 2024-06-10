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

type Post struct {
	ID         			int
	Uuid     			  string
	PostUuid   string
	IsUserMadePost     	bool
	IsUserLikePost 	bool
	IsUserDislikePost	bool
	Username   			string
	Message    			string
	Document   string
	TypeDoc    			string
	Date       			string
	Like       			int
	Dislike    			int
	Chanel   			  []string
	Target   			  []string
	Answers  			  []string
}

func UserPost(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("./DATA/User_data.db")
	defer db.Close()
	var post Post
	if r.Method == "POST" {
		uid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Fatal("cookie not found userpost")
			}
			log.Fatal("Error retrieving cookie uuid :", err)
		}
		var username string
		err1 := db.QueryRow("SELECT username FROM user WHERE uuid=?", uid.Value).Scan(&username)
		if err1 != nil {
			if err1 == sql.ErrNoRows {
				log.Fatal("sql user post :", err1)
			}
			log.Fatal(err1)
		}
		post.Uuid = uid.Value
		post.Username = username
		post.Message = r.FormValue("message")
		u, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("failed to generate UUID: %v", err)
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
						log.Fatal("post image:", err)
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
						}
					}
				}
			} else {
				post.Document = r.FormValue("document")
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
							path := "/VIEWS/static/stylsheet/VIDEO/" + post.PostUuid + ext
							if _, err := os.Stat("." + path); errors.Is(err, os.ErrNotExist) {
								// file does not exist
							} else {
								e := os.Remove("." + path)
								if e != nil {
									log.Fatal(e)
								}
							}

							f, err := os.OpenFile("."+path, os.O_WRONLY|os.O_CREATE, 0666)
							if err != nil {
								fmt.Println(err)
								return
							}
							defer f.Close()
							io.Copy(f, file)
							post.Document = "." + path
						}
					}
				}
			} else {
				post.Document = r.FormValue("document")
			}
		} else {
			post.TypeDoc = ""
			post.Document = ""
		}

		then := time.Now()
		post.Date = strconv.Itoa(then.Year()) + "/" + then.Month().String() + "/" + strconv.Itoa(then.Day()) + " " + strconv.Itoa(then.Hour()) + ":" + strconv.Itoa(then.Minute())
		post.Chanel = strings.Split(r.FormValue("chanel"), "R/")
		post.Target = strings.Split(r.FormValue("target"), "\\\\")
		if r.FormValue("chanel") != "" {
			AddPost(OpenDb("./DATA/User_data.db"), post)
		}
	}
}

func AddPost(db *sql.DB, post Post) {
	statement, err := db.Prepare("INSERT INTO post(uuid, postuuid, username, message, document, typedoc, date, chanel, target, answers, like, dislike) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("sql add post", err)
	}
	chanel := convertToString(post.Chanel)
	target := convertToString(post.Target)
	answers := convertToString(post.Answers)
	statement.Exec(post.Uuid, post.PostUuid, post.Username, post.Message, post.Document, post.TypeDoc, post.Date, chanel, target, answers, post.Like, post.Dislike)
	defer db.Close()
}

func convertToString(array []string) string {
	sort.Strings(array)
	return strings.Join(array, "|\\/|-_-|\\/|+{}")
}
