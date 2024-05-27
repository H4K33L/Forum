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
	pwd      string
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
				username VARCHAR(80) NOT NULL UNIQUE,
				pwd VARCHAR(80) NOT NULL,
				profilepicture VARCHAR(80),
				follow VARCHAR(80),
				follower VARCHAR(80),
				PRIMARY KEY (username, pwd),
				FOREIGN KEY (username, pwd) REFERENCES user(username, pwd)
			);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal((dberr.Error()))
	}
}

func createProfile(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("../DATA/User_data.db")
	cookieName, err := r.Cookie("Username")
	if err != nil {
		if err == http.ErrNoCookie {
			// Si le cookie n'existe pas
			log.Fatal("Cookie username connexion not found")
		}
		log.Fatal("Error retrieving cookie:", err)
	}
	cookiePwd, err := r.Cookie("Pwd")
	if err != nil {
		if err == http.ErrNoCookie {
			// Si le cookie n'existe pas
			log.Fatal("Cookie pwd connexion not found")
		}
		log.Fatal("Error retrieving cookie:", err)
	}
	booleanEmail, _ := VerifieEmail(cookieName.Value, db)
	booleanName, _ := VerifieName(cookiePwd.Value, db)
	if booleanEmail || booleanName {
		var userProfile profile
		userProfile.username = cookieName.Value
		userProfile.pwd, err = HashPassword(cookiePwd.Value)
		if err != nil {
			log.Fatal(err)
		}
		follow := convertToString(userProfile.follow)
		follower := convertToString(userProfile.follower)
		statement, err := db.Prepare("INSERT INTO profile(username, pwd, profilepicture, follow, follower) VALUES(?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			log.Fatal("error Prepare new user")
		}
		statement.Exec(userProfile.username, userProfile.pwd, userProfile.pp, follow, follower)
		defer db.Close()
	}
}
