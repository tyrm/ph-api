package oauth

import (
	"net/http"
	"os"

	"github.com/juju/loggo"
	"gopkg.in/go-oauth2/redis.v1"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	omodels "gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
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

	clientStore := store.NewClientStore()
	clientStore.Set("222222", &omodels.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:9094",
	})
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

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
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