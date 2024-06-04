package authentification

import (
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func PostSupr(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db := OpenDb("./DATA/User_data.db")
		defer db.Close()
		ID := r.FormValue("ToDelID")
		post := getPostByID(db,ID)
		Uuid, err := r.Cookie("UUID")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Fatal("cookie not found Uuid")
			}
			// Autre erreur
			log.Fatal("Error retrieving cookie Uuid :", err)
		}
		if ID != "" && post.Uuid == Uuid.Value {
			i, err := strconv.Atoi(ID)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("DELETE FROM `post` WHERE ID =? ",i)
			if err != nil {
				log.Fatal("err deleting post :", err)
			}
		}
	}
}