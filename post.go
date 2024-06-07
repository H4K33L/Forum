package authentification

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	ID       int
	Uuid     string
	IsUserMadePost     bool
	Username string
	Message  string
	Image    string
	Date     string
	Like     int
	Dislike  int
	Chanel   []string
	Target   []string
	Answers  []string
}

func InitDbpost(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS post
	(
	id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
	uuid VARCHAR(80) NOT NULL,
	username VARCHAR(80) NOT NULL,
	message LONG VARCHAR,
	image VARCHAR(80),
	date VARCHAR(80),
	chanel VARCHAR(80),
	target VARCHAR(80),
	answers LONG VARCHAR,
	like INTEGER,
	dislike INTEGER,
	FOREIGN KEY(uuid) 
		REFERENCES user(uuid)
		ON DELETE CASCADE
		ON UPDATE CASCADE
	);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal("InitDbpost :", dberr.Error())
	}
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
		post.Image = r.FormValue("image")
		then := time.Now()
		post.Date = strconv.Itoa(then.Year()) + "/" + then.Month().String() + "/" + strconv.Itoa(then.Day()) + "/" + strconv.Itoa(then.Hour()) + "/" + strconv.Itoa(then.Minute()) + "/" + strconv.Itoa(then.Second())
		post.Chanel = strings.Split(r.FormValue("chanel"), "R/")
		post.Target = strings.Split(r.FormValue("target"), "\\\\")
		if r.FormValue("chanel") != "" {
			AddPost(OpenDb("./DATA/User_data.db"), post)
		}
	}
}

func AddPost(db *sql.DB, post Post) {
	statement, err := db.Prepare("INSERT INTO post(uuid, username, message, image, date, chanel, target, answers, like, dislike) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("sql add post", err)
	}
	chanel := convertToString(post.Chanel)
	target := convertToString(post.Target)
	answers := convertToString(post.Answers)
	statement.Exec(post.Uuid, post.Username, post.Message, post.Image, post.Date, chanel, target, answers, post.Like, post.Dislike)
	defer db.Close()
}

func convertToString(array []string) string {
	return strings.Join(array, "|\\/|-_-|\\/|+{}")
}