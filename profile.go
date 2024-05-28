package authentification

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type profile struct {
	username string
	uid      string
	pp       string
	follow   []string
	follower []string
}

var profiles profile

func Profile(w http.ResponseWriter, r *http.Request) {
	// open the first web page openPage.html
	openpage := template.Must(template.ParseFiles("../VIEWS/html/profilePage.html"))
	// execute the modification of the page
	openpage.Execute(w, profiles)
}

func InitDbProfile(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS profile (
				uuid VARCHAR(80) NOT NULL UNIQUE,
				username VARCHAR(80) NOT NULL UNIQUE,
				profilepicture VARCHAR(80),
				follow VARCHAR(80),
				follower VARCHAR(80),
				PRIMARY KEY (uuid, username),
				FOREIGN KEY (uuid, username) REFERENCES user(uuid, username)
			);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal("InitDbProfile :", (dberr.Error()))
	}
}

func createProfile(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("../DATA/User_data.db")
	uid, err := r.Cookie("UUID")
	if err != nil {
		if err == http.ErrNoCookie {
			// Si le cookie n'existe pas
			log.Fatal("Cookie UUID connexion not found")
		}
		log.Fatal("Error retrieving cookie UUID:", err)
	}
	var username string
	err1 := db.QueryRow("SELECT username FROM user WHERE uuid=?", uid.Value).Scan(&username)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			log.Fatal("sql create profile:", err1)
		}
		log.Fatal(err1)
	}
	booleanEmail, _ := VerifieEmail(username, db)
	booleanName, _ := VerifieName(username, db)
	if booleanEmail || booleanName {
		var userProfile profile
		userProfile.username = username
		userProfile.uid = uid.Value
		follow := convertToString(userProfile.follow)
		follower := convertToString(userProfile.follower)
		statement, err := db.Prepare("INSERT INTO profile(uuid, username, profilepicture, follow, follower) VALUES(?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			log.Fatal("error Prepare new profile")
		}
		statement.Exec(userProfile.uid, userProfile.username, userProfile.pp, follow, follower)
		defer db.Close()
	}
}
