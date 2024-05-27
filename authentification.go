package authentification

import (
	"database/sql"
	"fmt"

	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	username string
	email    string
	pwd      string
}

var users user

func OpenDb(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func InitDb(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS user (
				email VARCHAR(80) NOT NULL UNIQUE,
				username VARCHAR(80) NOT NULL UNIQUE,
				pwd VARCHAR(80) NOT NULL,
				PRIMARY KEY (username, pwd)
			);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal((dberr.Error()))
	}
}

func Accueil(w http.ResponseWriter, r *http.Request) {
	openpage := template.Must(template.ParseFiles("../VIEWS/html/accueil.html"))
	openpage.Execute(w, users)
}

func Compte(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("../DATA/User_data.db")
	InitDb(db)
	defer db.Close()
	openpage := template.Must(template.ParseFiles("../VIEWS/html/homePage.html"))
	openpage.Execute(w, users)
}

func Adduser(db *sql.DB, user user) string {
	statement, err := db.Prepare("INSERT INTO user(email, username, pwd) VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return "error Prepare new user"
	}
	statement.Exec(user.email, user.username, user.pwd)
	return ""
}

func Connexion(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("../DATA/User_data.db")
	defer db.Close()

	openpage := template.Must(template.ParseFiles("../VIEWS/html/connexion.html"))
	var userconnect user

	if r.Method == "POST" {
		userconnect.email = r.FormValue("usermailconn")
		userconnect.username = r.FormValue("usermailconn")
		userconnect.pwd = r.FormValue("pwdconn")

		booleanUser, err := VerifieEmail(userconnect.email, db)
		if err != nil {
			log.Fatal("conn email ", err)
		}
		booleanName, err2 := VerifieName(userconnect.username, db)
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
			cookieName := &http.Cookie{
				Name:  "Username",
				Value: userconnect.username,
				Path:  "/",
			}
			http.SetCookie(w, cookieName)
			cookiePwd := &http.Cookie{
				Name:  "Pwd",
				Value: userconnect.pwd,
				Path:  "/",
			}
			http.SetCookie(w, cookiePwd)
			http.Redirect(w, r, "/compte", http.StatusSeeOther)
			return
		} else {
			fmt.Println("this account does not exist")
		}
	}
	openpage.Execute(w, users)
}

func Inscription(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("../DATA/User_data.db")
	defer db.Close()

	openpage := template.Must(template.ParseFiles("../VIEWS/html/inscription.html"))
	var userToAdd user

	if r.Method == "POST" {
		newEmail := r.FormValue("usermail")
		newUserName := r.FormValue("username")
		newPwd := r.FormValue("pwd")
		newPwd2 := r.FormValue("pwd2")

		booleanEmail, _ := VerifieEmail(newEmail, db)
		booleanName, _ := VerifieName(newUserName, db)

		if newPwd != newPwd2 {
			fmt.Println("the passwords are not equal")
		} else if booleanEmail {
			fmt.Println("this user already exists")
		} else if booleanName {
			fmt.Println("this name is already used")
		} else {
			userToAdd.email = newEmail
			userToAdd.username = newUserName
			userToAdd.pwd, _ = HashPassword(newPwd)
			errors := Adduser(db, userToAdd)
			if errors == "" {
				cookieName := &http.Cookie{
					Name:  "Username",
					Value: userToAdd.username,
					Path:  "/",
				}
				http.SetCookie(w, cookieName)
				cookiePwd := &http.Cookie{
					Name:  "Pwd",
					Value: newPwd,
					Path:  "/",
				}
				http.SetCookie(w, cookiePwd)
				http.Redirect(w, r, "/compte", http.StatusSeeOther)
				return
			} else {
				fmt.Println("error in adduser")
			}
		}
	}
	openpage.Execute(w, users)
}

func VerifieEmail(Email string, db *sql.DB) (bool, error) {
	var email string
	err := db.QueryRow("SELECT email FROM user WHERE email=?", Email).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func VerifieName(Name string, db *sql.DB) (bool, error) {
	var username string
	err := db.QueryRow("SELECT username FROM user WHERE username=?", Name).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
