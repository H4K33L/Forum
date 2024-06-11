package authentification

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/*
The function http request and http respose, get the information in the 
form postSupr and delete the post concerned.

input : w http.ResponseWriter, r *http.Request

output : nothing
*/
func PostSupr(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db := OpenDb("./DATA/User_data.db")
		defer db.Close()
		ID := r.FormValue("ToDelID")
		post := getPostByID(db, ID)
		Uuid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Fatal("PostSupr, cookie not found Uuid :", err)
			}
			log.Fatal("PostSupr, Error retrieving cookie Uuid :", err)
		}
		if ID != "" && post.Uuid == Uuid.Value {
			i, err := strconv.Atoi(ID)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("DELETE FROM `post` WHERE ID =? ", i)
			if err != nil {
				log.Fatal("PostSupr, err deleting post :", err)
			}
			_, err = db.Exec("DELETE FROM like WHERE ID =?", i)
			if err != nil {
				log.Fatal("PostSupr, err deleting post :", err)
			}
		}
	}
}

/*
The function http request and http respose, get the information in the 
form Edit and edit the post concerned with the modification chose by the user.

input : w http.ResponseWriter, r *http.Request

output : nothing
*/
func PostEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db := OpenDb("./DATA/User_data.db")
		defer db.Close()
		ID := r.FormValue("ToEditID")
		post := getPostByID(db, ID)
		Uuid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Fatal("PostEdit, cookie not found Uuid :", err)
			}
			log.Fatal("PostEdit, Error retrieving cookie Uuid :", err)
		}
		if ID != "" && post.Uuid == Uuid.Value {
			message := r.FormValue("messageEdit")
			if message == "" {
				message = post.Message
			}

			image := r.FormValue("imageEdit")
			if image == "" {
				image = post.Document
			}

			chanel := convertToString(strings.Split(r.FormValue("chanelEdit"), "R/"))
			if chanel == "" {
				chanel = convertToString(post.Chanel)
			}

			target := convertToString(strings.Split(r.FormValue("targetEdit"), "\\\\"))
			if target == "" {
				target = convertToString(post.Target)
			}

			then := time.Now()
			date := strconv.Itoa(then.Year()) + "/" + then.Month().String() + "/" + strconv.Itoa(then.Day()) + "/" + strconv.Itoa(then.Hour()) + "/" + strconv.Itoa(then.Minute()) + "/" + strconv.Itoa(then.Second())

			i, err := strconv.Atoi(ID)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("UPDATE post SET message =?, image =?, date =?, chanel =?, target =? WHERE ID =? ", message, image, date, chanel, target, i)
			if err != nil {
				log.Fatal("PostEdit, err Editing post :", err)
			}
		}
	}
}
