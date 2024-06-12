package client

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type profile struct {
	Username string
	Uid      string
	Pp       string
	Ext      string
}

var profiles profile

/*
Profile(w, r)

This function retrieves and displays the user's profile.
It retrieves the user's UUID from a cookie and queries the database to fetch the profile information.

Input: w : http.ResponseWriter, used to write the HTTP response. /// r : *http.Request, used to read the HTTP request.

Output: none
*/

func Profile(w http.ResponseWriter, r *http.Request) {
	// Open a connection to the user database
	db := OpenDb("./DATA/User_data.db")
	defer db.Close()

	// Retrieve the UUID cookie from the request
	uid, err := r.Cookie("UUID")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Fatal("profile profile cookie not found :", err)
		}
		log.Fatal("profile profile Error retrieving cookie UUID:", err)
	}

	// Query the profile table to retrieve the profile data based on the UUID
	err1 := db.QueryRow("SELECT * FROM profile WHERE uuid=?", uid.Value).Scan(&profiles.Uid, &profiles.Username, &profiles.Pp)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			log.Fatal("profile profile sql :", err1)
		}
		log.Fatal("profile profile error scan :", err1)
	}

	// Parse the profile page template and execute it with the retrieved profile data
	openpage := template.Must(template.ParseFiles("./VIEWS/html/profilePage.html"))
	openpage.Execute(w, profiles)
}

// createProfile(w, r)
//
// This function creates a user profile when an account is created.
// It uses a cookie to retrieve the user's UUID and interacts with the database to store the profile information.
//
// Input :
//
// w : http.ResponseWriter, used to write the HTTP response.
//
// r : *http.Request, used to read the HTTP request.
//
// Output : none/
func createProfile(w http.ResponseWriter, r *http.Request) {
	// Open a connection to the user database
	db := OpenDb("./DATA/User_data.db")
	defer db.Close()

	// Retrieve the UUID cookie from the request
	uid, err := r.Cookie("UUID")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Fatal("profile createProfile cookie not found ")
		}
		log.Fatal("profile createProfile Error retrieving cookie UUID:", err)
	}

	// Retrieve the username associated with the UUID
	var username string
	err1 := db.QueryRow("SELECT username FROM user WHERE uuid=?", uid.Value).Scan(&username)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			log.Fatal("profile createProfile sql:", err1)
		}
		log.Fatal("profile createprofile error scan :", err1)
	}

	// Check if the username exists in the database
	booleanEmail, _ := VerifieNameOrEmail(username, db)
	booleanName, _ := VerifieNameOrEmail(username, db)
	if booleanEmail || booleanName {
		// Create a new profile instance
		var userProfile profile
		userProfile.Username = username
		userProfile.Uid = uid.Value
		userProfile.Pp = "../static/stylsheet/IMAGES/PP/Avatar.jpg"

		// Prepare and execute the SQL statement to insert the new profile
		statement, err := db.Prepare("INSERT INTO profile(uuid, username, profilepicture) VALUES(?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			log.Fatal("profile createProfile error Prepare new profile")
		}
		statement.Exec(userProfile.Uid, userProfile.Username, userProfile.Pp)
	}
}
