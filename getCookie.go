package client

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

/*
GetCookie(w,r)

# This function get the cookie for the function Getpost

input :

w http.ReponseWriter

r *http.Request

output : none
*/
func GetCookie(w http.ResponseWriter, r *http.Request) {

	var AllStruct allStruct
	// Open a connection to the database and ensure it is closed at the end of the function.
	db := OpenDb("./DATA/User_data.db", w, r)
	defer db.Close()

	// If posts are available for the current request
	if GetPost(w, r) != nil {
		// Create and configure a cookie to store the last username from the request.
		requestPostName := &http.Cookie{
			Name:    "LastUsername",
			Value:   r.FormValue("username"),
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		}
		// Send the cookie to the client.
		http.SetCookie(w, requestPostName)

		// Create and configure a cookie to store the last channel from the request.
		requestPostChanel := &http.Cookie{
			Name:    "LastChanel",
			Value:   r.FormValue("chanels"),
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		}
		// Send the cookie to the client.
		http.SetCookie(w, requestPostChanel)

		// Get the posts for the current request.
		AllStruct.Posts = GetPost(w, r)
	} else {
		// Retrieve the "LastUsername" cookie.
		lastUsername, err := r.Cookie("LastUsername")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("GetCookie GetCookie cookie not found LastUsername") // The "LastUsername" cookie is not found.
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
			fmt.Println("GetCookie GetCookie Error retrieving cookie LastUsername :", err) // Error retrieving the "LastUsername" cookie.
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}

		// Retrieve the "LastChanel" cookie.
		lastChanel, err := r.Cookie("LastChanel")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("GetCookie GetCookie cookie not found LastChanel") // The "LastChanel" cookie is not found.
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
			fmt.Println("GetCookie GetCookie Error retrieving cookie LastChanel :", err) // Error retrieving the "LastChanel" cookie.
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}

		// Retrieve the "UUID" cookie.
		uid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("GetCookie GetCookie cookie not found userpost") // The "UUID" cookie is not found.
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
			fmt.Println("GetCookie GetCookie Error retrieving cookie uuid :", err) // Error retrieving the "UUID" cookie.
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}

		// Get the posts based on the last username and last channel.
		AllStruct.Posts = GetPostByBoth(db, lastUsername.Value, lastChanel.Value, uid, w, r)

	}
	uid, err := r.Cookie("UUID")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("GetCookie GetCookie cookie not found userpost") // The "UUID" cookie is not found.
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}
		fmt.Println("GetCookie GetCookie Error retrieving cookie uuid :", err) // Error retrieving the "UUID" cookie.
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	err1 := db.QueryRow("SELECT uuid,username,profilepicture FROM profile WHERE uuid=?", uid.Value).Scan(&AllStruct.Profile.Uid, &AllStruct.Profile.Username, &AllStruct.Profile.Pp)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			fmt.Println("profile profile sql :", err1)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}
		fmt.Println("profile profile error scan :", err1)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	// Load and execute the homepage template with the obtained posts.
	openpage := template.Must(template.ParseFiles("./VIEWS/html/homePage.html"))
	if err := openpage.Execute(w, AllStruct); err != nil {
		fmt.Println("GetCookie GetCookie error executing template :", err) // Error executing the template.
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
}
