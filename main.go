package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"./models"
	"./oauth"
	"./web"

	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

var logger *loggo.Logger

func main() {
	loggo.ConfigureLoggers("<root>=TRACE")

	newLogger :=  loggo.GetLogger("puphaus")
	logger = &newLogger

	config := CollectConfig()

	// Connect DB
	models.InitDB(config.DBEngine)
	defer models.CloseDB()

	// Init Oauth
	oauth.InitOath(config.RedisAddr)

	r := mux.NewRouter()
	rApi := r.PathPrefix("/api").Subrouter()

	rApi.HandleFunc("/envelope/{messageId}", web.HandleNotImplemented)
	rApi.HandleFunc("/envelope", web.HandleNotImplemented)

	// Meow
	rApi.HandleFunc("/meow", web.HandleMeow)

	// Oauth
	rOauth := r.PathPrefix("/oauth").Subrouter()
	rOauth.HandleFunc("/auth", oauth.HandleAuth)
	rOauth.HandleFunc("/authorize", oauth.HandleAuthorize)
	rOauth.HandleFunc("/login", oauth.HandleLogin)
	rOauth.HandleFunc("/token", oauth.HandleToken)

	r.HandleFunc("/oauth/token", oauth.HandleToken)

	// 404 handler
	r.PathPrefix("/").HandlerFunc(web.HandleNotFound)

	rApi.Use(web.LoggingMiddleware)

	go http.ListenAndServe(":8080", r)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	logger.Infof("%s", <-nch)

	logger.Infof("Done!")
}

