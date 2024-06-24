package client

import (
	"fmt"
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
	// Check if the request method is POST
	if r.Method == "POST" {
		// Open a connection to the user database
		db := OpenDb("./DATA/User_data.db", w, r)
		defer db.Close()

		// Retrieve the ID of the post to be deleted from the form
		ID := r.FormValue("ToDelID")

		// Get the post information from the database based on the ID
		post := getPostByID(db, ID, w, r)

		// Retrieve the UUID cookie from the request
		Uuid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("postmodification PostSupr, cookie not found Uuid :", err)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
			fmt.Println("postmodification PostSupr, Error retrieving cookie Uuid :", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}

		// Check if the ID is not empty and the post UUID matches the user UUID
		if ID != "" && post.Uuid == Uuid.Value {
			// Convert the ID to integer
			i, err := strconv.Atoi(ID)
			if err != nil {
				fmt.Println("postmodification postsupr, err atoi :", err)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}

			// Delete the post from the 'post' table
			_, err = db.Exec("DELETE FROM `post` WHERE ID =? ", i)
			if err != nil {
				fmt.Println("postmodification PostSupr, err deleting post :", err)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}

			// Delete associated likes from the 'like' table
			_, err = db.Exec("DELETE FROM like WHERE ID =?", i)
			if err != nil {
				fmt.Println("postmodification PostSupr, err deleting post :", err)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
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
	// Check if the request method is POST
	if r.Method == "POST" {
		// Open a connection to the user database
		db := OpenDb("./DATA/User_data.db", w, r)
		defer db.Close()

		// Retrieve the ID of the post to be edited from the form
		ID := r.FormValue("ToEditID")

		// Get the post information from the database based on the ID
		post := getPostByID(db, ID, w, r)

		// Retrieve the UUID cookie from the request
		Uuid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("postmodification PostEdit, cookie not found Uuid :", err)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
			fmt.Println("postmodification PostEdit, Error retrieving cookie Uuid :", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}

		// Check if the ID is not empty and the post UUID matches the user UUID
		if ID != "" && post.Uuid == Uuid.Value {
			// Retrieve the edited message from the form, or use the existing message if empty
			message := r.FormValue("messageEdit")
			if message == "" {
				message = post.Message
			}

			// Retrieve the edited image from the form, or use the existing image if empty
			image := r.FormValue("imageEdit")
			if image == "" {
				image = post.Document
			}

			// Retrieve the edited channels from the form, or use the existing channels if empty
			chanel := convertToString(strings.Split(r.FormValue("chanelEdit"), "R/"))
			if chanel == "" {
				chanel = convertToString(post.Chanel)
			}

			// Retrieve the edited targets from the form, or use the existing targets if empty
			target := convertToString(strings.Split(r.FormValue("targetEdit"), "\\\\"))
			if target == "" {
				target = convertToString(post.Target)
			}

			// Get the current timestamp
			then := time.Now()
			date := strconv.Itoa(then.Year()) + "/" + then.Month().String() + "/" + strconv.Itoa(then.Day()) + "/" + strconv.Itoa(then.Hour()) + "/" + strconv.Itoa(then.Minute()) + "/" + strconv.Itoa(then.Second())

			// Convert the ID to integer
			i, err := strconv.Atoi(ID)
			if err != nil {
				fmt.Println("postmodification postedit, error atoi:", err)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}

			// Update the post information in the 'post' table
			_, err = db.Exec("UPDATE post SET message =?, image =?, date =?, chanel =?, target =? WHERE ID =? ", message, image, date, chanel, target, i)
			if err != nil {
				fmt.Println("postmodification PostEdit, err Editing post :", err)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
		}
	}
}
