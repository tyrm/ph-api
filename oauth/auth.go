package oauth

import (
	"html/template"
	"net/http"
	"net/url"

	"../models"
	"../web"
	"github.com/jinzhu/gorm"
)

type AuthPageData struct {
	ApplicationName string
}

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
	if us.Get("ReturnUri") == nil {
		web.MakeErrorResponse(w, http.StatusBadRequest, "Missing client data", 0)
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

	clientID := us.Get("ReturnUri").(url.Values).Get("client_id")
	client, err := models.GetClient(clientID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			web.MakeErrorResponse(w, http.StatusForbidden, "Invalid client id", 0)
			return
		} else if err != nil {
			web.MakeErrorResponse(w, 400, err.Error(), 0)
			return
		}
	}

	outputAuthPage(w, r,  client.Name)
}

func outputAuthPage(w http.ResponseWriter, req *http.Request, appName string) {
	// Load Templates
	t := template.New("auth page")
	t, err := t.Parse(authPageTemplate)
	if err != nil {
		logger.Errorf("Error parsing auth page template: %s", err)
	}

	t.Execute(w, &AuthPageData{
		ApplicationName: appName,
	})
	return
}