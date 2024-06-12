package main

import (
	"client"
	"fmt"
	"net/http"
)

/*
main()

This is the main function of the application. It initializes the database, sets up HTTP routes,
and starts the HTTP server to handle incoming requests.

It initializes the database, sets up HTTP routes for various functionalities like user authentication,
profile management, post management, etc., and starts the HTTP server to listen on port 8080.

The function also serves static files for resources like CSS, JavaScript, etc.

Output: none
*/
func main() {
	db := client.OpenDb("./DATA/User_data.db")
	client.InitDb(db)
	defer db.Close()
	fmt.Println("server successfully up, go to http://localhost:8080")

	// Serve static files for resources like CSS, JavaScript, etc.
	fs := http.FileServer(http.Dir("VIEWS/static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	// Set URL handlers for different routes.
	http.HandleFunc("/", client.HomePage)
	http.HandleFunc("/logout", client.Logout)
	http.HandleFunc("/signup", client.Signup)
	http.HandleFunc("/login", client.Login)

	http.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		client.Account(w, r)
		client.UserPost(w, r)
		client.Like(w, r)
		client.Dislike(w, r)
		client.PostSupr(w, r)
		client.PostEdit(w, r)
		client.GetCookie(w, r)
	})

	http.HandleFunc("/profile", client.Profile)
	http.HandleFunc("/pwd", client.ChangePwd)
	http.HandleFunc("/username", client.ChangeUsername)
	http.HandleFunc("/delete", client.Delete)
	http.HandleFunc("/pp", client.ChangePP)
	http.HandleFunc("/404", client.Error404)
	// Start the HTTP server on port 8080.
	http.ListenAndServe(":8080", nil)
}
