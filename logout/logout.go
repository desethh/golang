package logout

import (
	"net/http"
)

func LogoutMethodChecker(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		LogoutHandler(w, r)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
