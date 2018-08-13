package web

import (
	"net/http"

	"github.com/juju/loggo"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	accessLog :=  loggo.GetLogger("puphaus.access")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		accessLog.Infof("%s - - \"%s %s %s\" - - \"%s\" \"%s\"", r.RemoteAddr, r.Method, r.RequestURI, r.Proto,
			r.Referer(), r.UserAgent())

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}