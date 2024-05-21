package client

import (
	"database/sql"
	"fmt"

	//"html/template"
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

//var users user

func OpenDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func InitDb(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS user
	(
	username VARCHAR(80) NOT NULL UNIQUE
	email VARCHAR(80) NOT NULL UNIQUE,
	pwd VARCHAR(80) NOT NULL,
	pp VARCHAR(80),
	follow VARCHAR(80),
	);`
	_, dberr := db.Exec(table)
	if dberr != nil {
		log.Fatal(dberr.Error())
	}
}

func Adduser(db *sql.DB, user user) string {
	statement, err := db.Prepare("INSERT INTO user(email, pwd) VALUES(?, ?)")
	if err != nil {
		fmt.Println(err)
		return "error Prepare new user"
	}
	fmt.Println(user.email)
	statement.Exec(user.email, user.pwd)
	defer db.Close()
	return ""
}

func Connexion(w http.ResponseWriter, r *http.Request) {
	db := OpenDb()
	// open the first web page openPage.html
	//openpage := template.Must(template.ParseFiles("template/connexion.html"))
	var userconnect user
	if r.Method == "POST" {
		Email := r.FormValue("usermailconn")
		Password := r.FormValue("pwdconn")
		userconnect.email = Email
		userconnect.pwd = Password
		booleanUser, _ := VerifieUser(userconnect.email, db)
		booleanPwd, _ := VerifiePwd(userconnect.email, userconnect.pwd, db)
		if !booleanPwd {
			fmt.Println("this password is  wrong")
		} else if booleanUser {
			//http.Redirect(w, r, "/compte", http.StatusSeeOther)
		} else {
			fmt.Println("this compte does not exist")
		}
	}
	defer db.Close()
	//openpage.Execute(w, users)
}

func Inscription(w http.ResponseWriter, r *http.Request) {
	db := OpenDb()
	// open the first web page openPage.html
	//openpage := template.Must(template.ParseFiles("template/inscription.html"))
	var userToAdd user
	if r.Method == "POST" {
		newEmail := r.FormValue("usermail")
		newPwd := r.FormValue("pwd")
		newPwd2 := r.FormValue("pwd2")
		booleanUser, _ := VerifieUser(newEmail, db)
		if newPwd != newPwd2 {
			fmt.Println("the password are not equal")
		} else if booleanUser {
			fmt.Println("this user already exist")
		} else {
			userToAdd.email = newEmail
			userToAdd.pwd, _ = HashPassword(newPwd)
			errors := Adduser(db, userToAdd)
			if errors == "" {
				//http.Redirect(w, r, "/compte", http.StatusSeeOther)
			} else {
				fmt.Println("error in adduser")
			}
		}
	}
	defer db.Close()
	//openpage.Execute(w, users)
}

func VerifieUser(Email string, db *sql.DB) (bool, error) {
	var tableUser []string
	Globaluser, err := db.Query("SELECT * FROM user WHERE email=?", Email)
	if err != nil {
		fmt.Println("error in hash")
	}
	defer Globaluser.Close()
	for Globaluser.Next() {
		var trueUser user
		if err := Globaluser.Scan(&trueUser.username, &trueUser.email, &trueUser.pwd); err != nil {
			return false, err
		}
		tableUser = append(tableUser, trueUser.email)
	}
	if err = Globaluser.Err(); err != nil {
		return false, err
	}
	for i := range tableUser {
		if tableUser[i] == Email {
			return true, nil
		}
	}
	return false, nil
}

func VerifiePwd(Email string, Password string, db *sql.DB) (bool, error) {
	var tableUser []string
	var tablePwd []string
	Globaluser, err := db.Query("SELECT * FROM user WHERE email=?", Email)
	if err != nil {
		fmt.Println("error in hash")
	}
	defer Globaluser.Close()
	for Globaluser.Next() {
		var trueUser user
		if err := Globaluser.Scan(&trueUser.username, &trueUser.email, &trueUser.pwd); err != nil {
			return false, err
		}
		tableUser = append(tableUser, trueUser.email)
		tablePwd = append(tablePwd, trueUser.pwd)
	}
	if err = Globaluser.Err(); err != nil {
		return false, err
	}
	for i := range tableUser {
		if tableUser[i] == Email {
			hash := tablePwd[i]
			return CheckPasswordHash(Password, hash), nil
		}
	}
	return false, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
