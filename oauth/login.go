package oauth

import (
	"html/template"
	"net/http"
	"regexp"

	"../models"
	"../web"
	"github.com/jinzhu/gorm"
)

const loginPageTemplate string = `<!doctype html>
<html class="no-js" lang="">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="x-ua-compatible" content="ie=edge">
        <title>Login</title>
        <meta name="theme-color" content="#000000">
        <meta name="viewport" content="width=device-width, initial-scale=1">

        <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
        <style>
          #loginbox {
            width: 300px;

            margin-top: 75px;
            margin-right: auto;
            margin-left: auto;
            padding: 10px;

            border-radius: 25px;
            border: 3px solid #000;

            text-align: center;
          }

          #loginbox label {
            position: relative;
            top: 4px;
          }

          #loginbox input {
            width: 200px;
          }

          #loginbox table {
            width: 100%;
          }

          #loginbox table td {
            padding-top: 10px;
            vertical-align: middle;
          }
          #loginbox .logo {
            position: relative;
            width: 150px;
            height: 100px;
            margin-right: auto;
            margin-left: auto;
          }

          #loginbox .logo .pup {
            position: absolute;
            left: 10px;
            top: 10px;

            width: 75px;
          }

          #loginbox .logo .haus {
            position: absolute;
            top: 10px;
            width: 75px;
          }
        </style>
    </head>
    <body>
        <div id='loginbox'>
          <div class="logo">
            <img src="https://o.pup.haus/public/images/noun_134522.svg" class="pup" />
            <img src="https://o.pup.haus/public/images/noun_8503.svg" class="haus" />
          </div>

          <form action="" method="post">
            <table>
{{if .Error}}
              </tr>
                <td colspan="2"><p class="bg-danger">{{.Error}}</p></td>
              </tr>
{{end}}
              <tr>
                <td><label><b>Username</b></label></td>
                <td><input type="text" class="form-control" placeholder="Enter Username" name="username" required></td>
              </tr>
              <tr>
                <td><label><b>Password</b></label></td>
                <td><input type="password" class="form-control" placeholder="Enter Password" name="password" required></td>
              </tr>
              <tr>
                <td colspan="2"><button type="submit" class="btn btn-default">Login</button></td>
              </tr>
            </table>
          </form>
        </div>
    </body>
</html>
`

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
		user, err := models.GetUser(formUsername)
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

		us.Set("LoggedInUserID", "000000")
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