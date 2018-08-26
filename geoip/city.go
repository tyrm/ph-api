package geoip

import (
	"errors"
	"net"

	"github.com/oschwald/geoip2-golang"
)

const geoIPDBFile = "data/GeoLite2-City.mmdb"

var ErrParsingIP = errors.New("could not parse IP")

func GetGeoIPCityRecord(addr string) (record *geoip2.City, err error) {
	db, err := geoip2.Open(geoIPDBFile)
	if err != nil {
		logger.Errorf("Error opening geoip2 db: %s", err)
		return
	}
	defer db.Close()

	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(addr)
	if ip == nil {
		logger.Errorf("Error parsing ip: %s", addr)
		err = ErrParsingIP
		return
	}
	record, err = db.City(ip)
	if err != nil {
		logger.Errorf("Error getting geoip2 record for %s: %s", addr, err)
	}

	return
}