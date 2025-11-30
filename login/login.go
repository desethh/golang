package login

import (
	"goroutines/dbs"
	"log"
	"net/http"
	"text/template"
)

func LoginMethodChecker(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))

		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if r.Method == http.MethodPost {
		LoginPage(w, r)
		return
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var isEx int
	row := dbs.DB.QueryRow("SELECT COUNT(*) FROM Users WHERE username = ? AND password = ?", username, password)
	row.Scan(&isEx)

	if isEx == 0 {
		http.Error(w, "Wrong Username or Password", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: "logged_in",
		Path:  "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: username,
		Path:  "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
