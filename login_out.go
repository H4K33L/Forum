package client

import (
	"database/sql"
	"fmt"

	"html/template"
	"net/http"

	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type user struct {
	uid      string
	username string
	email    string
	pwd      string
}

var users user

/*
Accueil(w, r)

This function serves the homepage of the application.
It renders the HTML template for the homepage and passes user data to it for rendering.

Input:

w : http.ResponseWriter, used to write the HTTP response.

r : *http.Request, used to read the HTTP request.

Output:
none
*/
func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}
	// Parse the HTML template named "accueil.html".
	openpage := template.Must(template.ParseFiles("./VIEWS/html/accueil.html"))

	// Execute the parsed template, passing any necessary data (users) to it.
	openpage.Execute(w, users) // The 'users' variable is referenced but not defined.
}

/*
Accueil(w, r)

This function serves the homepage of the application.
It renders the HTML template for the homepage and passes user data to it for rendering.
It's run all table to be created if not existe.
Input:

w : http.ResponseWriter, used to write the HTTP response.

r : *http.Request, used to read the HTTP request.

Output:
none
*/
func Account(w http.ResponseWriter, r *http.Request) {
	// Open the database connection.
	db := OpenDb("./DATA/User_data.db", w, r)
	// Initialize database tables if they don't exist.
	InitDbProfile(db, w, r)
	InitDbpost(db, w, r)
	InitDbLike(db, w, r)
	// Defer closing the database connection until the function returns.
	defer db.Close()
	// Create a user profile if it doesn't exist.
	createProfile(w, r)
}

/*
Adduser(db *sql.DB, user user) string

This function adds a new user to the database.
It prepares and executes an SQL statement to insert the user's information into the database.

Input:

db : *sql.DB, a pointer to the database connection.

user : user, the user object containing user information.

Output:
string, an empty string if the operation is successful, otherwise an error message.
*/
func Adduser(db *sql.DB, user user) string {
	// Prepare the SQL statement for inserting a new user.
	statement, err := db.Prepare("INSERT INTO user(uuid, email, username, pwd) VALUES(?, ?, ?, ?)")
	if err != nil {
		// If there's an error preparing the statement, return an error message.
		fmt.Println("login_out Adduser :", err)
		return "authentification adduser error Prepare new user"
	}
	// Execute the prepared statement to insert the new user into the database.
	statement.Exec(user.uid, user.email, user.username, user.pwd)
	// Return an empty string to indicate success (no error).
	return ""
}

/*
Connexion(w, r)

This function handles the user login process.
It renders the HTML template for the login page and processes the login form submission.

Input:

w : http.ResponseWriter, used to write the HTTP response.

r : *http.Request, used to read the HTTP request.

Output: none
*/
func Login(w http.ResponseWriter, r *http.Request) {
	// Open the database connection.
	db := OpenDb("./DATA/User_data.db", w, r)
	defer db.Close()

	// Parse the HTML template for the login page.
	openpage := template.Must(template.ParseFiles("./VIEWS/html/connexion.html"))

	// Initialize a user struct to store login information.
	var userconnect user

	// Retrieve the UUID cookie from the request.
	uid, err := r.Cookie("UUID")
	if err != nil {
		// If the UUID cookie is not found, create a new one.
		if err == http.ErrNoCookie {
			cookieUuid := &http.Cookie{
				Name:    "UUID",
				Value:   "",
				Path:    "/",
				Expires: time.Now().Add(24 * time.Hour),
			}
			http.SetCookie(w, cookieUuid)
			uid = cookieUuid
		} else {
			fmt.Println("login_out login Error retrieving cookie 'UUID' :", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}
	}

	// Check if the user is already logged in and has a profile.
	iscreate, err := IsUserCreate(uid.Value, db)
	if err != nil {
		fmt.Println("login_out login error user created but there is a problem :", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}

	// If the user is already logged in and has a profile, redirect to the account page.
	if uid.Value != "" && iscreate {
		http.Redirect(w, r, "/account", http.StatusSeeOther)
	} else if r.Method == "POST" {
		// If the request method is POST, it means the user is attempting to log in.

		// Retrieve the login credentials from the form.
		userconnect.email = r.FormValue("usermailconn")
		userconnect.username = r.FormValue("usermailconn")
		userconnect.pwd = r.FormValue("pwdconn")

		// Query the database to check if the user exists and retrieve their UUID.
		err1 := db.QueryRow("SELECT uuid FROM user WHERE username=? OR email=?", userconnect.username, userconnect.email).Scan(&userconnect.uid)
		if err1 != nil {
			if err1 == sql.ErrNoRows {
				fmt.Println("login_out login sql :", err1)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
			fmt.Println(err1)
		}

		// Verify the login credentials.
		booleanUser, err := VerifieNameOrEmail(userconnect.email, db)
		if err != nil {
			fmt.Println("login_out login :", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}
		booleanName, err2 := VerifieNameOrEmail(userconnect.username, db)
		if err2 != nil {
			fmt.Println("login_out login :", err2)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}
		booleanPwd, err1 := VerifiePwd(userconnect.email, userconnect.pwd, db)
		if err1 != nil {
			fmt.Println("login_out login :", err1)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}

		// If the credentials are valid, set the UUID cookie and redirect to the account page.
		if !booleanPwd {
			fmt.Println("login_out login, this password is wrong :", userconnect.pwd)
		} else if booleanUser || booleanName {
			cookieUuid := &http.Cookie{
				Name:    "UUID",
				Value:   userconnect.uid,
				Path:    "/",
				Expires: time.Now().Add(24 * time.Hour),
			}
			http.SetCookie(w, cookieUuid)
			http.Redirect(w, r, "/account", http.StatusSeeOther)
		} else {
			fmt.Println("login_out login, this account does not exist")
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}
	}
	// Execute the login page template.
	openpage.Execute(w, users)
}

/*
signup(w, r)

This function handles the user registration process.
It renders the HTML template for the registration page and processes the registration form submission.

Input:

w : http.ResponseWriter, used to write the HTTP response.

r : *http.Request, used to read the HTTP request.

Output: none
*/
func Signup(w http.ResponseWriter, r *http.Request) {
	// Open the database connection.
	db := OpenDb("./DATA/User_data.db", w, r)
	defer db.Close()

	// Parse the HTML template for the registration page.
	openpage := template.Must(template.ParseFiles("./VIEWS/html/inscription.html"))

	// Initialize a user struct to store registration information.
	var userToAdd user

	// If the request method is POST, it means the user is attempting to register.
	if r.Method == "POST" {
		// Retrieve the registration information from the form.
		newEmail := r.FormValue("usermail")
		newUserName := r.FormValue("username")
		newPwd := r.FormValue("pwd")
		newPwd2 := r.FormValue("pwd2")

		// Generate a UUID for the new user.
		u, err := uuid.NewV4()
		if err != nil {
			fmt.Println("login_out signup failed to generate UUID: %v", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}

		// Check if the email and username are already in use.
		booleanEmail, _ := VerifieNameOrEmail(newEmail, db)
		booleanName, _ := VerifieNameOrEmail(newUserName, db)

		// Validate the password.
		if newPwd != newPwd2 {
			fmt.Println("login_out signup, the passwords are not equal")
		} else if len(newPwd) < 8 {
			fmt.Println("login_out signup, this password is not long enough")
		} else if !isCorrectPassword(newPwd) {
			fmt.Println("login_out signup, this password is not secure enough")
		} else if booleanEmail {
			fmt.Println("login_out signup, this user already exists")
		} else if booleanName {
			fmt.Println("login_out signup, this name is already used")
		} else if newUserName != newPwd && newEmail != newPwd {
			// If all checks pass, create a new user entry in the database.
			userToAdd.email = newEmail
			userToAdd.username = newUserName
			userToAdd.pwd, err = HashPassword(newPwd)
			if err != nil {
				fmt.Println("login_out signup, error hashing password during registration :", err)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
			userToAdd.uid = u.String()

			// Add the new user to the database.
			errors := Adduser(db, userToAdd)
			if errors == "" {
				// If user creation is successful, set the UUID cookie and redirect to the account page.
				cookieUuid := &http.Cookie{
					Name:    "UUID",
					Value:   userToAdd.uid,
					Path:    "/",
					Expires: time.Now().Add(24 * time.Hour),
				}
				http.SetCookie(w, cookieUuid)
				http.Redirect(w, r, "/account", http.StatusSeeOther)
				return
			} else {
				fmt.Println("login_out signup, error in adduser", err)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
		} else {
			fmt.Println("login_out signup, you can't take your username as your password")
		}
	}
	// Execute the registration page template.
	openpage.Execute(w, users)
}

/*
Deconnexion(w http.ResponseWriter, r *http.Request)

This function handles the user logout process.
It clears the UUID cookie, effectively logging out the user, and redirects them to the home page.

Input:

w : http.ResponseWriter, used to write the HTTP response.

r : *http.Request, used to read the HTTP request.

Output: none
*/
func Logout(w http.ResponseWriter, r *http.Request) {
	cookieUuid := &http.Cookie{
		Name:    "UUID",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookieUuid)
	http.Redirect(w, r, "/", http.StatusSeeOther) // Redirection to home page
}

/*
Delete(w http.ResponseWriter, r *http.Request)

This function handles the deletion of a user account.
It first retrieves the UUID cookie from the request, then deletes the user from the database using the UUID.
After deletion, it clears the UUID cookie and redirects the user to the home page.

Input:

w : http.ResponseWriter, used to write the HTTP response.

r : *http.Request, used to read the HTTP request.

Output: none
*/
func Delete(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("./DATA/User_data.db", w, r)
	defer db.Close()
	uid, err := r.Cookie("UUID")
	if err != nil {
		if err == http.ErrNoCookie {
			cookieUuid := &http.Cookie{
				Name:    "UUID",
				Value:   "",
				Path:    "/",
				Expires: time.Now().Add(24 * time.Hour),
			}
			http.SetCookie(w, cookieUuid)
			uid = cookieUuid
		} else {
			fmt.Println("login_out delete, Error retrieving cookie 'UUID' :", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}
	}
	_, err = db.Exec("DELETE FROM user WHERE uuid=?", uid.Value)
	if err != nil {
		fmt.Println("login_out delete, Error deleting user :", err)
	}
	Logout(w, r) // logout of the user
}
