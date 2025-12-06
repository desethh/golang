package login

import (
	"goroutines/dbs"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
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

var store = sessions.NewCookieStore([]byte("pm2zlsz1PdlU8ymTwD4T2UIXpFy6qqzo"))

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
	userid := dbs.DB.QueryRow("SELECT uid from Users WHERE username = ?", username)
	var uid int
	userid.Scan(&uid)
	session, _ := store.Get(r, "user-session")
	session.Values["authenticated"] = true
	session.Values["user-id"] = uid
	session.Values["username"] = username

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Cannot save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
