package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var codeTitle = map[int]string{
	1:    "Malformed JSON Body",
	2201: "Missing Required Attribute",
	2202: "Requested Relationship Not Found",
}

type ErrorMessage struct {
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Status int    `json:"status,omitempty"`
	Code   int    `json:"code,omitempty"`
}

type ErrorResponse struct {
	Error  *ErrorMessage `json:"error"`
}

func HandleNotFound(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	MakeErrorResponse(response, http.StatusNotFound, request.URL.Path, 0)
	return
}

func MakeErrorResponse(response http.ResponseWriter, status int, detail string, code int) {

	// Get Title
	var title string
	if code == 0 { // code 0 means no code
		title = http.StatusText(status)
	} else {
		title = codeTitle[code]
	}

	// Send Response
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)

	m := ErrorResponse{
		Error: &ErrorMessage{
			Title:  title,
			Detail: detail,
			Status: status,
			Code:   code,
		},
	}
	b, _ := json.Marshal(m)
	fmt.Fprintf(response, "%s", b)

	return
}

func HandleNotImplemented(response http.ResponseWriter, request *http.Request) {

	MakeErrorResponse(response, http.StatusMethodNotAllowed, request.Method, 0)
	return
}