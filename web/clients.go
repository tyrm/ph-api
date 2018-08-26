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

func HandleGetClient(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", jsonapi.MediaType)
	vars := mux.Vars(request)
	client, err := models.GetClient(vars["client"])

	if err == gorm.ErrRecordNotFound {
		logger.Errorf("%s", err)
		MakeErrorResponse(response, 404, vars["messageId"], 0)
		return
	} else if err != nil {
		logger.Errorf("%s", err)
		MakeErrorResponse(response, 400, err.Error(), 0)
		return
	}

	if err := jsonapi.MarshalPayload(response, &client); err != nil {
		logger.Errorf("%s", err)
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}
}

func HandleGetClientUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", jsonapi.MediaType)
	vars := mux.Vars(request)
	client, err := models.GetClient(vars["client"])

	if err == gorm.ErrRecordNotFound {
		logger.Errorf("%s", err)
		MakeErrorResponse(response, 404, vars["messageId"], 0)
		return
	} else if err != nil {
		logger.Errorf("%s", err)
		MakeErrorResponse(response, 400, err.Error(), 0)
		return
	}

	user := *client.User
	if err := jsonapi.MarshalPayload(response, &user); err != nil {
		logger.Errorf("%s", err)
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}
}

func HandleGetClientList(response http.ResponseWriter, request *http.Request) {
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
	clients, err := models.GetClientPage(count, page)
	if err != nil {
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}

	if len(clients) == 0 {
		MakeErrorResponse(response, http.StatusNotFound, fmt.Sprintf("page %s", queryPage[0]), 0)
		return
	}

	// Format
	var clientList []interface{}
	for _, aClient := range(clients) {
		newObj := aClient
		clientList = append(clientList, &newObj)
	}

	// Build Response
	if err := jsonapi.MarshalPayload(response, clientList); err != nil {
		logger.Errorf("%s", err)
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}

	return
}