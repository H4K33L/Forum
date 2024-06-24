package client

import (
	"html/template"
	"net/http"
)

func Confidential(w http.ResponseWriter, r *http.Request) {
	openpage := template.Must(template.ParseFiles("./VIEWS/html/AnnexePageConfidentialité.html"))
	openpage.Execute(w, users)
}

func About(w http.ResponseWriter, r *http.Request) {
	openpage := template.Must(template.ParseFiles("./VIEWS/html/AnnexePageContenu.html"))
	openpage.Execute(w, users)
}
