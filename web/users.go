package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRequestBody struct {
	User UserRequest `json:"user"`
}

type UserPageResponse struct {
	Users *[]models.User `json:"users"`
}

type UserResponse struct {
	User *models.User `json:"user"`
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

	b, _ := json.Marshal(&UserResponse{User: &user})
	fmt.Fprintf(response, "%s", b)

	return
}

func HandleGetUserList(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	queries := request.URL.Query()

	// Get Page Number
	var page = 0
	queryPage, hasPage := queries["page"]
	if hasPage {
		queryPageInt, err := strconv.Atoi(queryPage[0])
		if err != nil || queryPageInt < 1 {
			MakeErrorResponse(response, http.StatusBadRequest, queryPage[0], 0)
			return
		}

		page = queryPageInt - 1
	}

	// Get Count Number
	var count = 100
	queryCount, hasCount := queries["count"]
	if hasCount {
		queryCountInt, err := strconv.Atoi(queryCount[0])
		if err != nil || queryCountInt < 1 {
			MakeErrorResponse(response, http.StatusBadRequest, queryCount[0], 0)
			return
		}

		count = queryCountInt
	}

	// Get User List
	users, err := models.GetUsersPage(count, page)
	if err != nil {
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}

	b, _ := json.Marshal(&UserPageResponse{Users: &users})
	fmt.Fprintf(response, "%s", b)

	return
}

func HandlePostUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	// Decode JSON
	var userRequest UserRequestBody
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&userRequest)
	if err != nil {
		MakeErrorResponse(response, http.StatusBadRequest, err.Error(), 1)
		return
	}

	// Validate
	if userRequest.User.Username == "" {
		MakeErrorResponse(response, http.StatusUnprocessableEntity, "username", 2201)
		return
	}
	if userRequest.User.Password == "" {
		MakeErrorResponse(response, http.StatusUnprocessableEntity, "password", 2201)
		return
	}

	exists, err := models.UserExists(userRequest.User.Username)
	if err != nil {
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}
	if exists {
		MakeErrorResponse(response, http.StatusConflict, fmt.Sprintf("User '%s' already exists", userRequest.User.Username), 0)
		return
	}

	// Add user to DB
	newUser, err := models.SetUser(models.User{
		Username: userRequest.User.Username,
		Password: userRequest.User.Password,
	})
	if err != nil {
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}

	b, _ := json.Marshal(&UserResponse{User: &newUser})
	fmt.Fprintf(response, "%s", b)

	return
}
