package client

import (
	"html/template"
	"net/http"
)

func Error404(w http.ResponseWriter, r *http.Request) {
	// Parse the HTML template named "accueil.html".
	openpage := template.Must(template.ParseFiles("./VIEWS/html/404.html"))

	// Execute the parsed template, passing any necessary data (users) to it.
	openpage.Execute(w, users) // The 'users' variable is referenced but not defined.
}
