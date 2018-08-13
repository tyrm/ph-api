package oauth

import (
	"net/http"

	"../models"
	"../web"
	"github.com/juju/loggo"
	"gopkg.in/go-oauth2/redis.v1"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/session.v1"
)

var (
	globalSessions *session.Manager
	logger         *loggo.Logger
	srv            *server.Server
)

func InitOath(reddisAddr string) {
	newLogger :=  loggo.GetLogger("puphaus.oauth")
	logger = &newLogger

	globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid","gclifetime":3600}`)
	go globalSessions.GC()

	manager := manage.NewDefaultManager()
	// token store
	manager.MustTokenStorage(redis.NewTokenStore(&redis.Config{
		Addr: reddisAddr,
	}))

	models.SetClient(&models.Client{
		ID:     "222222",
		Name:   "Test Client",
		Secret: "22222222",
		Domain: "http://localhost:9094",
		UserID: "1",
	})

	models.SetClient(&models.Client{
		ID:     "postman",
		Name:   "Postman",
		Secret: "postman",
		Domain: "https://www.getpostman.com/oauth2/callback",
		UserID: "1",
	})

	clientStore := NewClientStore()
	manager.MapClientStorage(clientStore)

	srv = server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		logger.Errorf("Internal Error: %s", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		logger.Errorf("Response Error: %s", re.Error.Error())
	})
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
	}
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	us, err := globalSessions.SessionStart(w, r)
	uid := us.Get("UserID")
	if uid == nil {
		if r.Form == nil {
			r.ParseForm()
		}
		us.Set("ReturnUri", r.Form)
		w.Header().Set("Location", "/oauth/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	us.Delete("UserID")
	return
}