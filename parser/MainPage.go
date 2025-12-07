package parser

import (
	"goroutines/dbs"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

var tmpl = template.Must(template.ParseFiles("templates/mainpage.html"))

type Hotel struct {
	Name      string
	Location  string
	Startdate time.Time
	Enddate   time.Time
}

type PageData struct {
	Username string
	Hotels   []Hotel
}

var store = sessions.NewCookieStore([]byte("pm2zlsz1PdlU8ymTwD4T2UIXpFy6qqzo"))

func HotelsHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "user-session")

	uid, ok := session.Values["user-id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := dbs.DB.Query("SELECT h.hotelname, h.location, h.startdate, h.enddate FROM Users u JOIN Hotels h ON u.uid = h.uid WHERE u.uid = ?", uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var hotels []Hotel

	for rows.Next() {
		var h Hotel
		rows.Scan(&h.Name, &h.Location, &h.Startdate, &h.Enddate)

		hotels = append(hotels, h)
	}

	cookie, ok := session.Values["username"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := PageData{
		Username: cookie,
		Hotels:   hotels,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
