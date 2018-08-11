package web

import (
	"fmt"
	"net/http"
)

func HandleMeow(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	fmt.Fprint(response, "{\"cat\":\"meow\"}")

	return
}