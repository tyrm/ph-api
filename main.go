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
	r.Use(web.LoggingMiddleware)

	rApi := r.PathPrefix("/api").Subrouter()
	rApi.Use(oauth.ProtectMiddleware) // Require Valid Bearer

	rApi.HandleFunc("/envelope/{messageId}", web.HandleNotImplemented)

	// Meow
	rApi.HandleFunc("/meow", web.HandleMeow)

	// 404 handler
	rApi.PathPrefix("/").HandlerFunc(web.HandleNotFound)

	// Oauth
	rOauth := r.PathPrefix("/oauth").Subrouter()
	rOauth.HandleFunc("/auth", oauth.HandleAuth)
	rOauth.HandleFunc("/authorize", oauth.HandleAuthorize)
	rOauth.HandleFunc("/haus.svg", oauth.HandleSVGHaus)
	rOauth.HandleFunc("/login", oauth.HandleLogin)
	rOauth.HandleFunc("/pup.svg", oauth.HandleSVGPup)
	rOauth.HandleFunc("/token", oauth.HandleToken)

	r.HandleFunc("/oauth/token", oauth.HandleToken)

	// 404 handler
	r.PathPrefix("/").HandlerFunc(web.HandleNotImplemented)


	go http.ListenAndServe(":8080", r)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	logger.Infof("%s", <-nch)

	logger.Infof("Done!")
}

