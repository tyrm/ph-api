package web

import (
	"fmt"
	"net/http"
	"strconv"

	"../models"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func HandleGetUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", jsonapi.MediaType)
	vars := mux.Vars(request)
	user, err := models.GetUserByUsername(vars["username"])

	if err == gorm.ErrRecordNotFound {
		logger.Errorf("%s", err)
		MakeErrorResponse(response, 404, vars["messageId"], 0)
		return
	} else if err != nil {
		logger.Errorf("%s", err)
		MakeErrorResponse(response, 400, err.Error(), 0)
		return
	}

	// Build Response
	user.Password = "" // Don't send password
	if err := jsonapi.MarshalPayload(response, &user); err != nil {
		logger.Errorf("%s", err)
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}

	return
}

func HandleGetUserList(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", jsonapi.MediaType)
	queries := request.URL.Query()

	// Get Page Number
	var page = 0
	queryPage, hasPage := queries["page[number]"]
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
	queryCount, hasCount := queries["page[size]"]
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

	if len(users) == 0 {
		MakeErrorResponse(response, http.StatusNotFound, fmt.Sprintf("page %s", queryPage[0]), 0)
		return
	}

	// Format
	var userList []interface{}
	for _, aUser := range(users) {
		newObj := aUser
		userList = append(userList, &newObj)
	}

	// Build Response
	if err := jsonapi.MarshalPayload(response, userList); err != nil {
		logger.Errorf("%s", err)
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}

	return
}

func HandlePostUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", jsonapi.MediaType)

	// Decode JSON
	user := new(models.User)
	if err := jsonapi.UnmarshalPayload(request.Body, user); err != nil {
		MakeErrorResponse(response, http.StatusBadRequest, err.Error(), 1)
		return
	}

	// Validate
	if user.Username == "" {
		MakeErrorResponse(response, http.StatusUnprocessableEntity, "username", 2201)
		return
	}
	if user.Password == "" {
		MakeErrorResponse(response, http.StatusUnprocessableEntity, "password", 2201)
		return
	}

	exists, err := models.UserExists(user.Username)
	if err != nil {
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}
	if exists {
		MakeErrorResponse(response, http.StatusConflict, fmt.Sprintf("User '%s' already exists", user.Username), 0)
		return
	}

	// Add user to DB
	newUser, err := models.SetUser(*user)
	if err != nil {
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}

	// Build Response
	if err := jsonapi.MarshalPayload(response, &newUser); err != nil {
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}
	return
}
