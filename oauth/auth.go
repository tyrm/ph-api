package oauth

import (
	"html/template"
	"net/http"
	"net/url"

	"../web"
)
const authPageTemplate string = `<!doctype html>
<html class="no-js" lang="">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="x-ua-compatible" content="ie=edge">
        <title>Authorize</title>
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

          <form action="" method="POST">
            <table>
              <tr>
                <td><h1>Authorize</h1></td>
              </tr>
              <tr>
                <td><p>The client would like to perform actions on your behalf.</p></td>
              </tr>
              <tr>
                <td><button type="submit" class="btn btn-default">Allow</button></td>
              </tr>
            </table>
          </form>
        </div>
    </body>
</html>
`

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
		web.MakeErrorResponse(w, http.StatusBadRequest, "Missing ReturnURI", 0)
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

	outputAuthPage(w, r)
}

func outputAuthPage(w http.ResponseWriter, req *http.Request) {
	// Load Templates
	t := template.New("auth page")
	t, err := t.Parse(authPageTemplate)
	if err != nil {
		logger.Errorf("Error parsing auth page template: %s", err)
	}

	t.Execute(w, &LoginPageData{})
	return
}