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

	r.HandleFunc("/api/envelope/{messageId}", web.HandleNotImplemented)
	r.HandleFunc("/api/envelope", web.HandleNotImplemented)

	// Meow
	r.HandleFunc("/api/meow", web.HandleMeow)

	// Oauth
	r.HandleFunc("/oauth/auth", oauth.HandleAuth)
	r.HandleFunc("/oauth/authorize", oauth.HandleAuthorize)
	r.HandleFunc("/oauth/login", oauth.HandleLogin)
	r.HandleFunc("/oauth/token", oauth.HandleToken)

	// 404 handler
	r.PathPrefix("/").HandlerFunc(web.HandleNotFound)

	r.Use(web.LoggingMiddleware)

	go http.ListenAndServe(":8080", r)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	logger.Infof("%s", <-nch)

	logger.Infof("Done!")
}

