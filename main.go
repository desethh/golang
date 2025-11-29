package main

import (
	"goroutines/dbs"
	"goroutines/login"
	"text/template"

	"goroutines/middleware"
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
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/mainpage.html"))

	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
