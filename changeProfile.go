package client

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
ChangePwd(w, r)

This function handles the password change process for a user.
It serves an HTML form for changing the password and processes the form submission.

Input:

w : http.ResponseWriter, used to write the HTTP response.

r : *http.Request, used to read the HTTP request.

Output: none
*/
func ChangePwd(w http.ResponseWriter, r *http.Request) {
	// Parse the HTML template file for the password change page
	openpage := template.Must(template.ParseFiles("./VIEWS/html/pwd.html"))

	// Define a variable to hold user information
	var userChangePwd user

	// Open a connection to the user database
	db := OpenDb("./DATA/User_data.db")
	defer db.Close()

	// Retrieve the UUID cookie from the request
	uid, err := r.Cookie("UUID")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Fatal("profile changepwd cookie not found")
		}
		log.Fatal(" profile changepwd  Error retrieving cookie UUID:", err)
	}

	// Handle the POST method for changing the password
	if r.Method == "POST" {
		// Retrieve the actual password, new password, and confirmation of the new password from the form
		pwd := r.FormValue("actualpwd")
		newPwd := r.FormValue("newPwd")
		newPwd2 := r.FormValue("newPwd2")

		// Retrieve user information from the database based on the UUID
		err1 := db.QueryRow("SELECT username, email, pwd FROM user WHERE uuid=?", uid.Value).Scan(&userChangePwd.username, &userChangePwd.email, &userChangePwd.pwd)
		if err1 != nil {
			if err1 == sql.ErrNoRows {
				log.Fatal("profile changepwd  sql create:", err1)
			}
			log.Fatal(err1)
		}

		// Check the conditions for changing the password
		if newPwd != newPwd2 {
			fmt.Println("the new passwords are not equal")
		} else if newPwd == pwd {
			fmt.Println("you can't replace the actual password by the actual password")
		} else if !CheckPasswordHash(pwd, userChangePwd.pwd) {
			fmt.Println("the actual password is wrong")
		} else if !isCorrectPassword(pwd) {
			fmt.Println("the password is wrongfully written, you need at least one uppercase letter, one lowercase letter, one number, and one special character")
		} else {
			// Hash the new password and update it in the database
			hashed, err := HashPassword(pwd)
			if err != nil {
				log.Fatal("profile changepwd  err hash :", err)
			}
			_, err = db.Exec("UPDATE user SET pwd =? WHERE UUID =? ", hashed, uid.Value)
			if err != nil {
				log.Fatal("profile changepwd  err rows :", err)
			}
		}
	}

	// Execute the HTML template with the user information
	openpage.Execute(w, userChangePwd)
}

/*
ChangeUsername(w, r)

This function handles the username change process for a user.
It serves an HTML form for changing the username and processes the form submission.

Input:

w : http.ResponseWriter, used to write the HTTP response.

r : *http.Request, used to read the HTTP request.

Output: none
*/
func ChangeUsername(w http.ResponseWriter, r *http.Request) {
	// Parse the HTML template file for the username change page
	openpage := template.Must(template.ParseFiles("./VIEWS/html/username.html"))

	// Define a variable to hold user information
	var userChangeUsername user

	// Open a connection to the user database
	db := OpenDb("./DATA/User_data.db")
	defer db.Close()

	// Retrieve the UUID cookie from the request
	uid, err := r.Cookie("UUID")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Fatal("profile changeusername cookie not found")
		}
		log.Fatal("profile changeusername Error retrieving cookie UUID:", err)
	}

	// Handle the POST method for changing the username
	if r.Method == "POST" {
		// Retrieve the current username and new username from the form
		username := r.FormValue("username")
		newUsername := r.FormValue("newUsername")
		newUsername2 := r.FormValue("newUsername2")

		// Retrieve user information from the database based on the UUID
		err1 := db.QueryRow("SELECT username, email, pwd FROM user WHERE uuid=?", uid.Value).Scan(&userChangeUsername.username, &userChangeUsername.email, &userChangeUsername.pwd)
		if err1 != nil {
			if err1 == sql.ErrNoRows {
				log.Fatal("profile changeusername sql create:", err1)
			}
			log.Fatal(err1)
		}

		// Check the conditions for changing the username
		if newUsername != newUsername2 {
			fmt.Println("the new usernames are not equal")
		} else if newUsername == username {
			fmt.Println("you can't replace the actual username by the actual username")
		} else {
			// Update the username in the database
			if err != nil {
				log.Fatal("profile changeusername err hash :", err)
			}
			_, err = db.Exec("UPDATE user SET username =? WHERE UUID =? ", username, uid.Value)
			if err != nil {
				log.Fatal("profile changeusername err rows :", err)
			}
		}
	}

	// Execute the HTML template with the user information
	openpage.Execute(w, userChangeUsername)
}

/*
ChangePP(w, r)

This function handles the process of changing the user's profile picture.
It allows the user to upload a new profile picture or provide a URL to an existing one,
and updates the profile picture information in the database accordingly.

Input: w :

http.ResponseWriter, used to write the HTTP response.

r : *http.Request, used to read the HTTP request.

Output: none
*/
func ChangePP(w http.ResponseWriter, r *http.Request) {
	// Open a connection to the user database
	db := OpenDb("./DATA/User_data.db")
	defer db.Close()

	// Retrieve the UUID cookie from the request
	uid, err := r.Cookie("UUID")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Fatal("profile changepp cookie not found")
		}
		log.Fatal("profile changepp Error retrieving cookie UUID:", err)
	}

	// Parse the HTML template file for the profile picture change page
	openpage := template.Must(template.ParseFiles("./VIEWS/html/pp.html"))

	// Define a variable to hold profile picture information
	var ppProfile profile

	// Handle the file upload for changing the profile picture
	if r.FormValue("typedoc") == "file" {
		file, handler, err := r.FormFile("documentFile")
		if err != nil {
			if err == http.ErrMissingFile {
				ppProfile.Pp = "../static/stylsheet/IMAGES/PP/Avatar.jpg"
			} else {
				log.Fatal("profile changepp ppProfile image:", err)
			}
		} else {
			extension := strings.LastIndex(handler.Filename, ".")
			if extension == -1 {
				fmt.Println("profile changepp : there is no extension to the file")
			} else {
				ext := handler.Filename[extension:]
				e := strings.ToLower(ext)
				if e == ".png" || e == ".jpeg" || e == ".jpg" || e == ".gif" || e == ".svg" || e == ".avif" || e == ".apng" || e == ".webp" {
					path := "/static/stylsheet/IMAGES/PP/" + uid.Value + ext
					if _, err := os.Stat("./VIEWS" + path); errors.Is(err, os.ErrNotExist) {
						log.Fatal("profile changepp no extension :", err)
					} else {
						err = os.Remove("./VIEWS" + path)
						if err != nil {
							log.Fatal("profile changepp can't remove path", err)
						}
					}

					f, err := os.OpenFile("./VIEWS"+path, os.O_WRONLY|os.O_CREATE, 0666)
					if err != nil {
						log.Fatal("profile changepp can't open file", err)
					}
					defer f.Close()
					_, err = io.Copy(f, file)
					if err != nil {
						log.Fatal("profile changepp can't copy file", err)
					}
					ppProfile.Pp = path
					ppProfile.Ext = "file"
				}
			}
		}
	} else {
		ppProfile.Pp = r.FormValue("document")
		ppProfile.Ext = "url"
	}

	// Update the profile picture in the database
	_, err = db.Exec("UPDATE profile SET profilepicture =? WHERE UUID =? ", ppProfile.Pp, uid.Value)
	if err != nil {
		log.Fatal("profile changepp err rows :", err)
	}

	// Execute the HTML template with the profile picture information
	openpage.Execute(w, ppProfile)
}
