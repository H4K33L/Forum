package main

import (
	"authentification"
	"fmt"

	//"html/template"
	"net/http"
)

// main is the main function of the program.
func main() {
	db := authentification.OpenDb("../DATA/User_data.db")
	authentification.InitDb(db)
	authentification.InitDbProfile(db)
	defer db.Close()
	fmt.Println("server successfully up, go to http://localhost:8080")

	// Serve static files for resources like CSS, JavaScript, etc.
	fs := http.FileServer(http.Dir("VIEWS/static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	// Set URL handlers for different routes.
	http.HandleFunc("/", authentification.Accueil)
	http.HandleFunc("/inscription", func(w http.ResponseWriter, r *http.Request) {
		authentification.Inscription(w, r)
	})
	http.HandleFunc("/connexion", func(w http.ResponseWriter, r *http.Request) {
		authentification.Connexion(w, r)
	})
	http.HandleFunc("/compte", func(w http.ResponseWriter, r *http.Request) {
		authentification.Compte(w, r)
		authentification.UserPost(w, r)
		/*posts := authentification.GetPost(w,r)
		var tmpl = template.Must(template.ParseFiles("template/homePage.html"))
		tmpl.Execute(w, posts)*/
	})
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		authentification.Profile(w, r)
	})
	// Start the HTTP server on port 8080.
	http.ListenAndServe(":8080", nil)
}
