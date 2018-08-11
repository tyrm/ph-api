package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"./models"
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

	r := mux.NewRouter()
	r.HandleFunc("/envelope/{messageId}", web.HandleNotImplemented)
	r.HandleFunc("/envelope", web.HandleNotImplemented)

	// Meow
	r.HandleFunc("/meow", web.HandleMeow)

	// Users
	r.HandleFunc("/meow", web.HandleMeow)

	// 404 handler
	r.PathPrefix("/").HandlerFunc(web.HandleNotFound)

	go http.ListenAndServe(":8080", r)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-nch)

	logger.Infof("Done!")
}

