package getbook

import (
	"goroutines/dbs"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("pm2zlsz1PdlU8ymTwD4T2UIXpFy6qqzo"))

func MethodChecker(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/booking", http.StatusSeeOther)
	}
	if r.Method == http.MethodPost {
		GetBookingPage(w, r)
		return
	}
}

func GetBookingPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	uid, ok := session.Values["user-id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	hotel := r.FormValue("hotel")
	parts := strings.Split(hotel, "|")
	hotelname := parts[0]
	location := parts[1]
	startdate := r.FormValue("startdate")
	enddate := r.FormValue("enddate")
	_, err := dbs.DB.Exec("INSERT INTO Hotels (uid, hotelname, location, startdate, enddate) VALUES (?, ?, ?, ?, ?)", uid, hotelname, location, startdate, enddate)
	if err != nil {
		http.Error(w, "Ошибка сохранения: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
