package main

import (
	"authentification"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// main is the main function of the program.
func main() {
	db := authentification.OpenDb("./DATA/User_data.db")
	authentification.InitDb(db)
	defer db.Close()
	fmt.Println("server successfully up, go to http://localhost:8080")

	// Serve static files for resources like CSS, JavaScript, etc.
	fs := http.FileServer(http.Dir("VIEWS/static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	// Set URL handlers for different routes.
	http.HandleFunc("/", authentification.Accueil)
	http.HandleFunc("/deconnexion", authentification.Deconnexion)
	http.HandleFunc("/inscription", func(w http.ResponseWriter, r *http.Request) {
		authentification.Inscription(w, r)
	})
	http.HandleFunc("/connexion", func(w http.ResponseWriter, r *http.Request) {
		authentification.Connexion(w, r)
	})
	http.HandleFunc("/compte", func(w http.ResponseWriter, r *http.Request) {
		authentification.Compte(w, r)
		authentification.UserPost(w, r)
		authentification.Like(w, r)
		authentification.Dislike(w, r)
		authentification.PostSupr(w, r)
		authentification.PostEdit(w, r)

		var posts []authentification.Post
		if authentification.GetPost(w, r) != nil {
			requestPostName := &http.Cookie{
				Name:    "LastUsername",
				Value:   r.FormValue("username"),
				Path:    "/",
				Expires: time.Now().Add(24 * time.Hour),
			}
			http.SetCookie(w, requestPostName)

			requestPostChanel := &http.Cookie{
				Name:    "LastChanel",
				Value:   r.FormValue("chanels"),
				Path:    "/",
				Expires: time.Now().Add(24 * time.Hour),
			}
			http.SetCookie(w, requestPostChanel)

			posts = authentification.GetPost(w, r)
		} else {
			lastUsername, err := r.Cookie("LastUsername")
			if err != nil {
				if err == http.ErrNoCookie {
					log.Fatal("cookie not found LastUsername")
				}
				log.Fatal("Error retrieving cookie LastUsername :", err)
			}

			lastChanel, err := r.Cookie("LastChanel")
			if err != nil {
				if err == http.ErrNoCookie {
					log.Fatal("cookie not found LastChanel")
				}
				log.Fatal("Error retrieving cookie LastChanel :", err)
			}

			uid, err := r.Cookie("UUID")
			if err != nil {
				if err == http.ErrNoCookie {
					log.Fatal("cookie not found userpost")
				}
				log.Fatal("Error retrieving cookie uuid :", err)
			}
			posts = authentification.GetPostByBoth(db,lastUsername.Value,lastChanel.Value, uid)
		}

		openpage := template.Must(template.ParseFiles("./VIEWS/html/homePage.html"))
		if err := openpage.Execute(w, posts); err != nil {
			log.Fatal("erreur lors de l'envois ", err)
		}
	})

	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		authentification.Profile(w, r)
	})
	http.HandleFunc("/pwd", func(w http.ResponseWriter, r *http.Request) {
		authentification.ChangePwd(w, r)
	})
	http.HandleFunc("/username", func(w http.ResponseWriter, r *http.Request) {
		authentification.ChangeUsername(w, r)
	})
	// Start the HTTP server on port 8080.
	http.ListenAndServe(":8080", nil)
}
