package oauth

import (
	"net/http"

	"../web"
)

func ProtectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := srv.ValidationBearerToken(r)
		if err != nil {
			web.MakeErrorResponse(w, http.StatusForbidden, err.Error(), 0)
			return
		}

		next.ServeHTTP(w, r)
	})
}