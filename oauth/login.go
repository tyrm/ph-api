package oauth

import (
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		us, err := globalSessions.SessionStart(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		us.Set("LoggedInUserID", "000000")
		w.Header().Set("Location", "/oauth/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/login.html")
}
