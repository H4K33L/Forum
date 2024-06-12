package client

import (
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
	// Open a connection to the database and ensure it is closed at the end of the function.
	db := OpenDb("./DATA/User_data.db", w, r)
	defer db.Close()

	var posts []Post

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
		posts = GetPost(w, r)
	} else {
		// Retrieve the "LastUsername" cookie.
		lastUsername, err := r.Cookie("LastUsername")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("cookie not found LastUsername") // The "LastUsername" cookie is not found.
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
			fmt.Println("Error retrieving cookie LastUsername:", err) // Error retrieving the "LastUsername" cookie.
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}

		// Retrieve the "LastChanel" cookie.
		lastChanel, err := r.Cookie("LastChanel")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("cookie not found LastChanel") // The "LastChanel" cookie is not found.
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
			fmt.Println("Error retrieving cookie LastChanel:", err) // Error retrieving the "LastChanel" cookie.
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}

		// Retrieve the "UUID" cookie.
		uid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("cookie not found userpost") // The "UUID" cookie is not found.
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
			fmt.Println("Error retrieving cookie uuid:", err) // Error retrieving the "UUID" cookie.
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}

		// Get the posts based on the last username and last channel.
		posts = GetPostByBoth(db, lastUsername.Value, lastChanel.Value, uid,w,r)
	}

	// Load and execute the homepage template with the obtained posts.
	openpage := template.Must(template.ParseFiles("./VIEWS/html/homePage.html"))
	if err := openpage.Execute(w, posts); err != nil {
		fmt.Println("error executing template:", err) // Error executing the template.
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
}
