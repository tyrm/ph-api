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

	// Create Top Router
	r := mux.NewRouter()
	r.Use(web.LoggingMiddleware)

	if config.HTTPCorsOrigin != "" {
		logger.Infof("CORS origin found, adding 'Access-Control-Allow-Origin' header.")
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", config.HTTPCorsOrigin)

				next.ServeHTTP(w, r)
			})
		})
	}

	// Oauth Router
	oauth.InitOath(config.RedisAddr)
	rOauth := r.PathPrefix("/oauth").Subrouter()
	rOauth.HandleFunc("/auth", oauth.HandleAuth)
	rOauth.HandleFunc("/authorize", oauth.HandleAuthorize)
	rOauth.HandleFunc("/haus.svg", oauth.HandleSVGHaus).Methods("GET")
	rOauth.HandleFunc("/login", oauth.HandleLogin)
	rOauth.HandleFunc("/pup.svg", oauth.HandleSVGPup).Methods("GET")
	rOauth.HandleFunc("/token", oauth.HandleToken)

	// API Router
	rApi := r.PathPrefix("/api").Subrouter()
	rApi.Use(oauth.ProtectMiddleware) // Require Valid Bearer
	rApi.HandleFunc("/meow", web.HandleMeow)
	rApi.HandleFunc("/users", web.HandleGetUserList).Methods("GET")
	rApi.HandleFunc("/users", web.HandlePostUser).Methods("POST")
	rApi.HandleFunc("/users/{username}", web.HandleGetUser).Methods("GET")

	// Catchall for API Router so we throw 403 for all requests to api without valid token to prevent scans
	rApi.PathPrefix("/").HandlerFunc(web.HandleNotImplemented)
	r.PathPrefix("/").HandlerFunc(web.HandleNotImplemented) // Top Router Catch All

	go http.ListenAndServe(":8080", r)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	logger.Infof("%s", <-nch)

	logger.Infof("Done!")
}

