package main

import (
	"goroutines/booking"
	"goroutines/dbs"

	"goroutines/getbook"
	"goroutines/login"
	"goroutines/logout"
	"goroutines/middleware"
	"goroutines/parser"
	"goroutines/register"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("pm2zlsz1PdlU8ymTwD4T2UIXpFy6qqzo"))

func main() {
	dbs.DBopen()
	defer dbs.DB.Close()
	authmiddleware := middleware.AuthMiddleWare(http.HandlerFunc(MainPage))
	http.Handle("/", authmiddleware)
	http.HandleFunc("/register", register.MethodChecker)
	http.HandleFunc("/login", login.LoginMethodChecker)
	http.HandleFunc("/logout", logout.LogoutMethodChecker)
	http.HandleFunc("/booking", booking.BookingPage)
	http.HandleFunc("/book", getbook.GetBookingPage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	parser.HotelsHandler(w, r)
}
