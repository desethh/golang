package register

import (
	"goroutines/dbs"
	"html/template"
	"log"
	"net/http"
)

func MethodChecker(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/reg.html"))

		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if r.Method == http.MethodPost {
		RegPage(w, r)
	}
}

func RegPage(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	var IsEx int
	err := dbs.DB.QueryRow("SELECT COUNT(*) FROM Users WHERE username = ?", username).Scan(&IsEx)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if IsEx > 0 {
		http.Error(w, "This username already exists", http.StatusBadRequest)
		return
	}

	_, err = dbs.DB.Exec("INSERT INTO Users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		http.Error(w, "Ошибка сохранения: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "http://localhost:8080/login", http.StatusSeeOther)
}
