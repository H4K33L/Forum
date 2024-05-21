package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// tmpl is an HTML template for displaying data.
var tmpl = template.Must(template.ParseFiles("VIEWS/html/connection.html"))

// main is the main function of the program.
func main() {

	fmt.Println("server successfully up, go to http://localhost:8080")

	// Serve static files for resources like CSS, JavaScript, etc.
	fs := http.FileServer(http.Dir("VIEWS/static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	// Set URL handlers for different routes.
	http.HandleFunc("/", load)

	// Start the HTTP server on port 8080.
	http.ListenAndServe(":8080", nil)
}

func load(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w,"")
}