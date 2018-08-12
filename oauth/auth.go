package oauth

import (
	"net/http"
	"net/url"

	"../web"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	us, err := globalSessions.SessionStart(w, r)
	if err != nil {
		web.MakeErrorResponse(w, http.StatusInternalServerError, err.Error(), 0)
		return
	}
	if us.Get("LoggedInUserID") == nil {
		w.Header().Set("Location", "/oauth/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	if r.Method == "POST" {
		form := us.Get("ReturnUri").(url.Values)
		u := new(url.URL)
		u.Path = "/oauth/authorize"
		u.RawQuery = form.Encode()
		w.Header().Set("Location", u.String())
		w.WriteHeader(http.StatusFound)
		us.Delete("Form")
		us.Set("UserID", us.Get("LoggedInUserID"))
		return
	}

	outputHTML(w, r, "static/auth.html")
}

func HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	err := srv.HandleAuthorizeRequest(w, r)
	if err != nil {
		web.MakeErrorResponse(w, http.StatusBadRequest, err.Error(), 0)
	}
}

func HandleToken(w http.ResponseWriter, r *http.Request) {
	err := srv.HandleTokenRequest(w, r)
	if err != nil {
		web.MakeErrorResponse(w, http.StatusInternalServerError, err.Error(), 0)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
