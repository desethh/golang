package parser

import (
	"goroutines/dbs"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("templates/mainpage.html"))

type Hotel struct {
	Name     string
	Location string
	Price    int
}

type PageData struct {
	Username string
	Hotels   []Hotel
}

func HotelsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := dbs.DB.Query("SELECT name, location, price FROM Hotels")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var hotels []Hotel
	for rows.Next() {
		var h Hotel
		rows.Scan(&h.Name, &h.Location, &h.Price)
		hotels = append(hotels, h)
	}
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "http://localhost:8080/login", http.StatusSeeOther)
		return
	}

	data := PageData{
		Username: cookie.Value,
		Hotels:   hotels,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
