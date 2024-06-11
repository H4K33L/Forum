package authentification

import (
	"database/sql"
	"fmt"
	"unicode"

	"html/template"
	"log"
	"net/http"

	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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
func Accueil(w http.ResponseWriter, r *http.Request) {
	openpage := template.Must(template.ParseFiles("./VIEWS/html/accueil.html"))
	openpage.Execute(w, users)
}

func Compte(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("./DATA/User_data.db")
	InitDbProfile(db)
	InitDbpost(db)
	InitDbLike(db)
	defer db.Close()
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
	statement, err := db.Prepare("INSERT INTO user(uuid, email, username, pwd) VALUES(?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Adduser", err)
		return "error Prepare new user"
	}
	statement.Exec(user.uid, user.email, user.username, user.pwd)
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
func Connexion(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("./DATA/User_data.db")
	defer db.Close()

	openpage := template.Must(template.ParseFiles("./VIEWS/html/connexion.html"))
	var userconnect user
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
			log.Fatal("Error retrieving cookie 'uuid' :", err)
		}
	}
	iscreate, err := IsUserCreate(uid.Value, db)
	if err != nil {
		log.Fatal("error user created but there is a probleme")
	}
	if uid.Value != "" && iscreate {
		http.Redirect(w, r, "/compte", http.StatusSeeOther)
	} else if r.Method == "POST" {
		userconnect.email = r.FormValue("usermailconn")
		userconnect.username = r.FormValue("usermailconn")
		userconnect.pwd = r.FormValue("pwdconn")
		err1 := db.QueryRow("SELECT uuid FROM user WHERE username=? OR email=?", userconnect.username, userconnect.email).Scan(&userconnect.uid)
		if err1 != nil {
			if err1 == sql.ErrNoRows {
				log.Fatal("sql connexion :", err1)
			}
			log.Fatal(err1)
		}
		booleanUser, err := VerifieNameOrEmail(userconnect.email, db)
		if err != nil {
			log.Fatal("conn email ", err)
		}
		booleanName, err2 := VerifieNameOrEmail(userconnect.username, db)
		if err2 != nil {
			log.Fatal("conn name ", err2)
		}
		booleanPwd, err1 := VerifiePwd(userconnect.email, userconnect.pwd, db)
		if err1 != nil {
			log.Fatal("conn pwd ", err1)
		}
		if !booleanPwd {
			fmt.Println("this password is wrong:", userconnect.pwd)
		} else if booleanUser || booleanName {
			cookieUuid := &http.Cookie{
				Name:    "UUID",
				Value:   userconnect.uid,
				Path:    "/",
				Expires: time.Now().Add(24 * time.Hour),
			}
			http.SetCookie(w, cookieUuid)
			http.Redirect(w, r, "/compte", http.StatusSeeOther)
		} else {
			fmt.Println("this account does not exist")
		}
	}
	openpage.Execute(w, users)
}

/*
Inscription(w, r)

This function handles the user registration process.
It renders the HTML template for the registration page and processes the registration form submission.

Input:

w : http.ResponseWriter, used to write the HTTP response.

r : *http.Request, used to read the HTTP request.

Output: none
*/
func Inscription(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("./DATA/User_data.db")
	defer db.Close()

	openpage := template.Must(template.ParseFiles("./VIEWS/html/inscription.html"))
	var userToAdd user

	if r.Method == "POST" {
		newEmail := r.FormValue("usermail")
		newUserName := r.FormValue("username")
		newPwd := r.FormValue("pwd")
		newPwd2 := r.FormValue("pwd2")
		u, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("failed to generate UUID: %v", err)
		}
		booleanEmail, _ := VerifieNameOrEmail(newEmail, db)
		booleanName, _ := VerifieNameOrEmail(newUserName, db)
		if newPwd != newPwd2 {
			fmt.Println("the passwords are not equal")
		} else if len(newPwd) < 8 {
			fmt.Println("this password is not long enough")
		} else if !isCorrectPassword(newPwd) {
			fmt.Println("this password is not secure enough")
		} else if booleanEmail {
			fmt.Println("this user already exists")
		} else if booleanName {
			fmt.Println("this name is already used")
		} else if newUserName != newPwd && newEmail != newPwd {
			userToAdd.email = newEmail
			userToAdd.username = newUserName
			userToAdd.pwd, err = HashPassword(newPwd)
			if err != nil {
				log.Fatal("erreur hash password inscription ")
			}
			userToAdd.uid = u.String()
			errors := Adduser(db, userToAdd)
			if errors == "" {
				cookieUuid := &http.Cookie{
					Name:    "UUID",
					Value:   userToAdd.uid,
					Path:    "/",
					Expires: time.Now().Add(24 * time.Hour),
				}
				http.SetCookie(w, cookieUuid)
				http.Redirect(w, r, "/compte", http.StatusSeeOther)
				return
			} else {
				fmt.Println("error in adduser")
			}
		} else {
			fmt.Println("you can't take your username as your password")
		}
	}
	openpage.Execute(w, users)
}

/*
VerifieNameOrEmail(Input string, db *sql.DB) (bool, error)

This function checks if the given input (either username or email) already exists in the database.
It queries the database to search for a matching username or email.

Input:

Input : string, the username or email to be checked.

db : *sql.DB, a pointer to the database connection.

Output:

bool : true if the username or email already exists, false otherwise.

error : an error if there is an issue with the database query.
*/
func VerifieNameOrEmail(Input string, db *sql.DB) (bool, error) {
	var NameOrEmail string
	err := db.QueryRow("SELECT username FROM user WHERE email=? OR username=?", Input, Input).Scan(&NameOrEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

/*
IsUserCreate(Input string, db *sql.DB) (bool, error)

This function checks if a user with the given UUID exists in the database.
It queries the database to search for a matching UUID.

Input:

Input : string, the UUID to be checked.

db : *sql.DB, a pointer to the database connection.

Output:

bool : true if the user with the given UUID exists, false otherwise.

error : an error if there is an issue with the database query.
*/
func IsUserCreate(Input string, db *sql.DB) (bool, error) {
	var uuid string
	err := db.QueryRow("SELECT username FROM user WHERE uuid=?", Input).Scan(&uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

/*
VerifiePwd(Input string, Password string, db *sql.DB) (bool, error)

This function verifies if the provided password matches the hashed password stored in the database for the given username or email.
It queries the database to retrieve the hashed password for the provided username or email and then compares it with the provided password.

Input:

Input : string, the username or email for which the password needs to be verified.

Password : string, the password to be verified.

db : *sql.DB, a pointer to the database connection.

Output:

bool : true if the provided password matches the hashed password, false otherwise.

error : an error if there is an issue with the database query.
*/
func VerifiePwd(Input string, Password string, db *sql.DB) (bool, error) {
	var hashedPwd string
	err := db.QueryRow("SELECT pwd FROM user WHERE email=? OR username=?", Input, Input).Scan(&hashedPwd)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return CheckPasswordHash(Password, hashedPwd), nil
}

/*
HashPassword(password string) (string, error)

This function hashes the provided password using bcrypt with a cost factor of 14.

Input:

password : string, the password to be hashed.

Output:

string : the hashed password.

error : an error if there is an issue during the hashing process.
*/
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

/*
CheckPasswordHash(password, hash string) bool

This function compares a password with its hashed version to verify if they match.

Input:

password : string, the password to be checked.

hash : string, the hashed password to compare against.

Output:

bool : true if the password matches the hashed version, false otherwise.
*/
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
func Deconnexion(w http.ResponseWriter, r *http.Request) {
	cookieUuid := &http.Cookie{
		Name:    "UUID",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookieUuid)
	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
}

/*
isCorrectPassword(password string) bool

This function checks if the provided password meets the criteria for being considered secure.
It verifies if the password contains at least one uppercase letter, one lowercase letter, one digit, and one special character.

Input:

password : string, the password to be checked.

Output:

bool : true if the password meets the criteria, false otherwise.
*/
func isCorrectPassword(password string) bool {
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
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
	db := OpenDb("./DATA/User_data.db")
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
			log.Fatal("Error retrieving cookie 'uuid' :", err)
		}
	}
	_, err = db.Exec("DELETE FROM user WHERE uuid=?", uid.Value)
	if err != nil {
		fmt.Println("err delete :", err)
	}
	Deconnexion(w, r)
}
