package booking

import (
	"goroutines/dbs"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

type Hotel struct {
	Name     string
	Location string
}

var store = sessions.NewCookieStore([]byte("pm2zlsz1PdlU8ymTwD4T2UIXpFy6qqzo"))

var bookingTmpl = template.Must(template.ParseFiles("templates/booking.html"))

func BookingPage(w http.ResponseWriter, r *http.Request) {
	sessions, _ := store.Get(r, "user-session")
	_, ok := sessions.Values["authenticated"].(bool)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	rows, err := dbs.DB.Query(`SELECT name, location FROM avhotels`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var hotels []Hotel
	for rows.Next() {
		var h Hotel
		rows.Scan(&h.Name, &h.Location)
		hotels = append(hotels, h)
	}

	data := map[string]interface{}{
		"Hotels": hotels,
	}

	bookingTmpl.Execute(w, data)
}
