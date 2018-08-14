package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type UserResponse struct {
	User  *models.User  `json:"user"`
}


func HandleGetUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(request)
	user, error := models.GetUserByUsername(vars["username"])

	if error == gorm.ErrRecordNotFound {
		MakeErrorResponse(response, 404, vars["messageId"], 0)
		return
	} else if error != nil {
		MakeErrorResponse(response, 400, error.Error(), 0)
		return
	}

	b, _ := json.Marshal(user)
	fmt.Fprintf(response, "%s", b)

	return
}