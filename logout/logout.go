package logout

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("pm2zlsz1PdlU8ymTwD4T2UIXpFy6qqzo"))

func LogoutMethodChecker(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		LogoutHandler(w, r)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
