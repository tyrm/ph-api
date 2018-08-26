package geoip

import (
	"fmt"
	"net/http"

	"github.com/juju/loggo"
)

var logger *loggo.Logger

const URLBase = "http://geolite.maxmind.com/download/geoip/database"
const URLCityFile = "GeoLite2-City.tar.gz"
const URLCountryFile = "GeoLite2-Country.tar.gz"
const URLASNFile = "GeoLite2-ASN.tar.gz"

func checkFiles() {
	cityFileHead, err := http.Head(fmt.Sprintf("%s/%s", URLBase, URLCityFile))
	if err != nil {
		logger.Errorf("Error HEADing City file.")
	}
	logger.Debugf("%v", cityFileHead.Header.Get("Content-Disposition"))
}

func Init() {
	newLogger :=  loggo.GetLogger("puphaus.geoip")
	logger = &newLogger

	checkFiles()
}