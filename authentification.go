package authentification

import (
	"database/sql"
	"fmt"

	"html/template"
	"log"
	"net/http"

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

func OpenDb(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("OpenDb 1:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("OpenDb 2:", err)
	}
	return db
}

func InitDb(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS user (
				uuid VARCHAR(80) NOT NULL UNIQUE,
				email VARCHAR(80) NOT NULL UNIQUE,
				username VARCHAR(10) NOT NULL UNIQUE,
				pwd VARCHAR(255) NOT NULL,
			);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal("InitDb :", (dberr.Error()))
	}
}

func Accueil(w http.ResponseWriter, r *http.Request) {
	openpage := template.Must(template.ParseFiles("../VIEWS/html/accueil.html"))
	openpage.Execute(w, users)
}

func Compte(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("../DATA/User_data.db")
	InitDbpost(db)
	defer db.Close()
	createProfile(w, r)
	openpage := template.Must(template.ParseFiles("../VIEWS/html/homePage.html"))
	openpage.Execute(w, users)
}

func Adduser(db *sql.DB, user user) string {
	statement, err := db.Prepare("INSERT INTO user(uuid, email, username, pwd) VALUES(?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Adduser", err)
		return "error Prepare new user"
	}
	statement.Exec(user.uid, user.email, user.username, user.pwd)
	return ""
}

func Connexion(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("../DATA/User_data.db")
	defer db.Close()

	openpage := template.Must(template.ParseFiles("../VIEWS/html/connexion.html"))
	var userconnect user

	if r.Method == "POST" {
		uid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				// Si le cookie n'existe pas
				log.Fatal("Cookie username not found")
				// Vous pouvez gérer ce cas en définissant une valeur par défaut, en renvoyant une erreur HTTP, etc.
				log.Fatal("Cookie 'uuid' not found", http.StatusUnauthorized)
				return
			}
			// Si une autre erreur s'est produite
			log.Fatal("Error retrieving cookie 'uuid' :", err)
		}
		userconnect.uid = uid.Value
		userconnect.email = r.FormValue("usermailconn")
		userconnect.username = r.FormValue("usermailconn")
		userconnect.pwd = r.FormValue("pwdconn")
		if userconnect.uid != "" {
			http.Redirect(w, r, "/compte", http.StatusSeeOther)
		} else {
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
				cookieUuid := &http.Cookie{
					Name:  "UUID",
					Value: userconnect.uid,
					Path:  "/",
				}
				http.SetCookie(w, cookieUuid)
				http.Redirect(w, r, "/compte", http.StatusSeeOther)
				return
			} else {
				fmt.Println("this account does not exist")
			}
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
		u, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("failed to generate UUID: %v", err)
		}
		//log.Printf("generated Version 4 UUID %v", u)
		booleanEmail, _ := VerifieEmail(newEmail, db)
		booleanName, _ := VerifieName(newUserName, db)

		if newPwd != newPwd2 {
			fmt.Println("the passwords are not equal")
		} else if booleanEmail {
			fmt.Println("this user already exists")
		} else if booleanName {
			fmt.Println("this name is already used")
		} else if newUserName != newPwd && newEmail != newPwd {
			userToAdd.email = newEmail
			userToAdd.username = newUserName
			userToAdd.pwd, err = HashPassword(newPwd)
			if err != nil {
				log.Fatal("erreur hash password inscription :", err)
			}
			userToAdd.uid = u.String()
			errors := Adduser(db, userToAdd)
			if errors == "" {
				cookieUuid := &http.Cookie{
					Name:  "UUID",
					Value: userToAdd.uid,
					Path:  "/",
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
