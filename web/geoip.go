package web

import (
	"fmt"
	"net/http"

	"../geoip"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
)

func HandleGetGeoIPCity(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", jsonapi.MediaType)
	vars := mux.Vars(request)
	record, err := geoip.GetGeoIPCityRecord(vars["addr"])
	if err != nil {
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}

	fmt.Printf("%v\n", request.Header.Get("Accept-Language"))

	fmt.Printf("%v\n", record)
	//fmt.Fprintf(response, "%v", record)
	if err := jsonapi.MarshalPayload(response, &record); err != nil {
		logger.Errorf("%s", err)
		MakeErrorResponse(response, http.StatusInternalServerError, err.Error(), 0)
		return
	}

	return
}
