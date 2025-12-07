package booking

import (
	"goroutines/dbs"
	"html/template"
	"net/http"
)

type Hotel struct {
	Name     string
	Location string
}

var bookingTmpl = template.Must(template.ParseFiles("templates/booking.html"))

func BookingPage(w http.ResponseWriter, r *http.Request) {
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
