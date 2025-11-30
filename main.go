package main

import (
	"goroutines/dbs"
	"goroutines/login"
	"goroutines/logout"
	"goroutines/middleware"
	"goroutines/parser"
	"goroutines/register"
	"log"
	"net/http"
)

func main() {
	dbs.DBopen()
	defer dbs.DB.Close()
	authmiddleware := middleware.AuthMiddleWare(http.HandlerFunc(MainPage))
	http.Handle("/", authmiddleware)
	http.HandleFunc("/register", register.MethodChecker)
	http.HandleFunc("/login", login.LoginMethodChecker)
	http.HandleFunc("/logout", logout.LogoutMethodChecker)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	parser.HotelsHandler(w, r)
}
