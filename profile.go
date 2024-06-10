package authentification

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

	_ "github.com/mattn/go-sqlite3"
)

type profile struct {
	Username string
	Uid      string
	Pp       string
	Ext      string
}

var profiles profile

func Profile(w http.ResponseWriter, r *http.Request) {
	// open the first web page openPage.html
	openpage := template.Must(template.ParseFiles("./VIEWS/html/profilePage.html"))
	// execute the modification of the page
	openpage.Execute(w, profiles)
}

func createProfile(w http.ResponseWriter, r *http.Request) {
	db := OpenDb("./DATA/User_data.db")
	uid, err := r.Cookie("UUID")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Fatal("cookie not found createProfile")
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
	booleanEmail, _ := VerifieNameOrEmail(username, db)
	booleanName, _ := VerifieNameOrEmail(username, db)
	if booleanEmail || booleanName {
		var userProfile profile
		userProfile.Username = username
		userProfile.Uid = uid.Value
		userProfile.Pp = "../static/stylsheet/IMAGES/PP/Avatar.jpg"
		statement, err := db.Prepare("INSERT INTO profile(uuid, username, profilepicture) VALUES(?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			log.Fatal("error Prepare new profile")
		}
		statement.Exec(userProfile.Uid, userProfile.Username, userProfile.Pp)
		defer db.Close()
	}
}

func ChangePwd(w http.ResponseWriter, r *http.Request) {
	openpage := template.Must(template.ParseFiles("./VIEWS/html/pwd.html"))
	var userChangePwd user
	db := OpenDb("./DATA/User_data.db")
	uid, err := r.Cookie("UUID")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Fatal("cookie not found createProfile")
		}
		log.Fatal("Error retrieving cookie UUID:", err)
	}
	if r.Method == "POST" {
		pwd := r.FormValue("actualpwd")
		newPwd := r.FormValue("newPwd")
		newPwd2 := r.FormValue("newPwd2")
		err1 := db.QueryRow("SELECT username, email, pwd FROM user WHERE uuid=?", uid.Value).Scan(&userChangePwd.username, &userChangePwd.email, &userChangePwd.pwd)
		if err1 != nil {
			if err1 == sql.ErrNoRows {
				log.Fatal("sql create profile:", err1)
			}
			log.Fatal(err1)
		}
		if newPwd != newPwd2 {
			fmt.Println("the new passwords are not equal")
		} else if newPwd == pwd {
			fmt.Println("you can't replace the actual password by the actual password")
		} else if !CheckPasswordHash(pwd, userChangePwd.pwd) {
			fmt.Println("the actual password is wrong")
		} else if !isCorrectPassword(pwd) {
			fmt.Println("the password is wrongfully writen you need at least one maj one min one number and one special character")
		} else {
			hashed, err := HashPassword(pwd)
			if err != nil {
				log.Fatal("err hash profile :", err)
			}
			_, err = db.Exec("UPDATE user SET pwd =? WHERE UUID =? ", hashed, uid.Value)
			if err != nil {
				log.Fatal("err rows profile :", err)
			}
		}
	}
	openpage.Execute(w, userChangePwd)
}

func ChangeUsername(w http.ResponseWriter, r *http.Request) {
	openpage := template.Must(template.ParseFiles("./VIEWS/html/username.html"))
	var userChangeUsername user
	db := OpenDb("./DATA/User_data.db")
	uid, err := r.Cookie("UUID")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Fatal("cookie not found createProfile")
		}
		log.Fatal("Error retrieving cookie UUID:", err)
	}
	if r.Method == "POST" {
		username := r.FormValue("username")
		newUsername := r.FormValue("newUsername")
		newUsername2 := r.FormValue("newUsername2")
		err1 := db.QueryRow("SELECT username, email, pwd FROM user WHERE uuid=?", uid.Value).Scan(&userChangeUsername.username, &userChangeUsername.email, &userChangeUsername.pwd)
		if err1 != nil {
			if err1 == sql.ErrNoRows {
				log.Fatal("sql create profile:", err1)
			}
			log.Fatal(err1)
		}
		if newUsername != newUsername2 {
			fmt.Println("the new passwords are not equal")
		} else if newUsername == username {
			fmt.Println("you can't replace the actual password by the actual password")
		} else {
			if err != nil {
				log.Fatal("err hash profile :", err)
			}
			_, err = db.Exec("UPDATE user SET username =? WHERE UUID =? ", username, uid.Value)
			if err != nil {
				log.Fatal("err rows profile :", err)
			}
		}
	}
	openpage.Execute(w, userChangeUsername)
}

func ChangePP(w http.ResponseWriter, r *http.Request) {
	openpage := template.Must(template.ParseFiles("./VIEWS/html/pwd.html"))
	var ppProfile profile
	if r.FormValue("typedoc") == "file" {
		file, handler, err := r.FormFile("documentFile")
		if err != nil {
			if err == http.ErrMissingFile {
				fmt.Println("no file uploaded")
				ppProfile.Pp = "../static/stylsheet/IMAGES/PP/Avatar.jpg"
			} else {
				log.Fatal("ppProfile image:", err)
			}
		} else {
			extension := strings.LastIndex(handler.Filename, ".")
			if extension == -1 {
				fmt.Println("there is no extension to the file")
			} else {
				ext := handler.Filename[extension:]
				e := strings.ToLower(ext)
				if e == ".png" || e == ".jpeg" || e == ".jpg" || e == ".gif" || e == ".svg" || e == ".avif" || e == ".apng" || e == ".webp" {
					path := "/static/stylsheet/IMAGES/PP/" + ppProfile.Uid + ext
					if _, err := os.Stat("./VIEWS" + path); errors.Is(err, os.ErrNotExist) {
						// file does not exist
					} else {
						e := os.Remove("./VIEWS" + path)
						if e != nil {
							log.Fatal(e)
						}
					}

					f, err := os.OpenFile("./VIEWS"+path, os.O_WRONLY|os.O_CREATE, 0666)
					if err != nil {
						fmt.Println(err)
						return
					}
					defer f.Close()
					_, err = io.Copy(f, file)
					if err != nil {
						fmt.Println(err)
						return
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
	openpage.Execute(w, ppProfile)
}
