package oauth

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"regexp"

	"../models"
	"../web"
	"github.com/jinzhu/gorm"
)

type LoginPageData struct {
	Error string
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		us, err := globalSessions.SessionStart(w, r)
		if err != nil {
			web.MakeErrorResponse(w, http.StatusInternalServerError, err.Error(), 0)
			return
		}

		// Check Client ID
		returnUri := us.Get("ReturnUri")
		if returnUri == nil {
			outputLoginPage(w, r, "invalid configuration")
			logger.Errorf("Login request contains no client data")
			return
		}

		clientID := returnUri.(url.Values).Get("client_id")
		_, err = models.GetClient(clientID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				outputLoginPage(w, r, "invalid configuration")
				logger.Errorf("Login request contains invalid client id: %s", clientID)
				return
			} else if err != nil {
				web.MakeErrorResponse(w, http.StatusInternalServerError, err.Error(), 0)
				return
			}
		}

		r.ParseForm()
		// Check Username
		formUsername := r.Form["username"][0]
		reUsername, err := regexp.Compile(`^[a-zA-Z0-9_]+$`)
		if err != nil {
			web.MakeErrorResponse(w, http.StatusInternalServerError, err.Error(), 0)
			return
		}
		if !reUsername.MatchString(formUsername) {
			outputLoginPage(w, r, "invalid username")
			return
		}

		// Find User
		user, err := models.GetUserByUsername(formUsername)
		if err == gorm.ErrRecordNotFound {
			outputLoginPage(w, r, "username/password not recognized")
			return
		} else if err != nil {
			logger.Errorf("Error retrieving user %s: %s", r.Form["username"][0], err)
			web.MakeErrorResponse(w, http.StatusInternalServerError, err.Error(), 0)
			return
		}

		// Verify Password
		if !user.CheckPassword(r.Form["password"][0]) {
			outputLoginPage(w, r, "username/password not recognized")
			return
		}

		us.Set("LoggedInUserID", fmt.Sprint(user.ID))
		w.Header().Set("Location", "/oauth/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputLoginPage(w, r, "")
}

func outputLoginPage(w http.ResponseWriter, req *http.Request, errorString string) {
	// Load Templates
	t := template.New("login page")
	t, err := t.Parse(loginPageTemplate)
	if err != nil {
		logger.Errorf("Error parsing login page template: %s", err)
	}

	t.Execute(w, &LoginPageData{
		Error: errorString,
	})
	return
}