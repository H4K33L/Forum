package client

import (
	"html/template"
	"log"
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
	// Ouvre une connexion à la base de données et assure sa fermeture à la fin de la fonction.
	db := OpenDb("./DATA/User_data.db")
	defer db.Close()

	var posts []Post

	// Si des publications sont disponibles pour la requête actuelle
	if GetPost(w, r) != nil {
		// Crée et configure un cookie pour enregistrer le dernier nom d'utilisateur de la requête.
		requestPostName := &http.Cookie{
			Name:    "LastUsername",
			Value:   r.FormValue("username"),
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		}
		// Envoie le cookie au client.
		http.SetCookie(w, requestPostName)

		// Crée et configure un cookie pour enregistrer le dernier canal de la requête.
		requestPostChanel := &http.Cookie{
			Name:    "LastChanel",
			Value:   r.FormValue("chanels"),
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		}
		// Envoie le cookie au client.
		http.SetCookie(w, requestPostChanel)

		// Obtient les publications pour la requête actuelle.
		posts = GetPost(w, r)
	} else {
		// Récupère le cookie "LastUsername".
		lastUsername, err := r.Cookie("LastUsername")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Fatal("cookie not found LastUsername") // Le cookie "LastUsername" est introuvable.
			}
			log.Fatal("Error retrieving cookie LastUsername:", err) // Erreur lors de la récupération du cookie "LastUsername".
		}

		// Récupère le cookie "LastChanel".
		lastChanel, err := r.Cookie("LastChanel")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Fatal("cookie not found LastChanel") // Le cookie "LastChanel" est introuvable.
			}
			log.Fatal("Error retrieving cookie LastChanel:", err) // Erreur lors de la récupération du cookie "LastChanel".
		}

		// Récupère le cookie "UUID".
		uid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Fatal("cookie not found userpost") // Le cookie "UUID" est introuvable.
			}
			log.Fatal("Error retrieving cookie uuid:", err) // Erreur lors de la récupération du cookie "UUID".
		}

		// Obtient les publications en fonction du dernier nom d'utilisateur et du dernier canal.
		posts = GetPostByBoth(db, lastUsername.Value, lastChanel.Value, uid)
	}

	// Charge et exécute le modèle de la page d'accueil avec les publications obtenues.
	openpage := template.Must(template.ParseFiles("./VIEWS/html/homePage.html"))
	if err := openpage.Execute(w, posts); err != nil {
		log.Fatal("erreur lors de l'envoi:", err) // Erreur lors de l'exécution du modèle.
	}
}
