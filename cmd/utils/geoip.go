package utils

import (
	"fmt"
	"log"
	"net"

	"github.com/fatih/color"
	"github.com/oschwald/geoip2-golang"
)

var geoipCityDB *geoip2.Reader
var geoipAsDB *geoip2.Reader

func InitGeoIP() {
	var err error
	if geoipCityDB, err = geoip2.Open("./cmd/utils/geodata/GeoLite2-City.mmdb"); err != nil {
		log.Fatalf("GeoIP DB open error: %v", err)
	}
	if geoipAsDB, err = geoip2.Open("./cmd/utils/geodata/GeoLite2-ASN.mmdb"); err != nil {
		log.Fatalf("GeoIP DB open error: %v", err)
	}
}

// IPアドレスから国情報を取得
func LookupCountry(ipStr string) string {
	if geoipCityDB == nil {
		return ""
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}

	CountryRecord, err := geoipCityDB.Country(ip)
	if err != nil {
		return ""
	}

	AsRecord, err := geoipAsDB.ASN(ip)
	if err != nil {
		return ""
	}

	if CountryRecord.Country.Names["en"] == "" || AsRecord.AutonomousSystemOrganization == "" {
		return RenderBlock(fmt.Sprintf("GeoIP"), []string{"None"}, color.New(color.FgHiRed))
	}

	geoipInfo := []string{
		"IP: " + ipStr,
		"Country: " + CountryRecord.Country.Names["en"],
		"Organization: " + AsRecord.AutonomousSystemOrganization,
	}

	return RenderBlock(fmt.Sprintf("GeoIP"), geoipInfo, color.New(color.FgHiRed))
}
