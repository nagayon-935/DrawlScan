package utils

import (
	"fmt"
	"net"

	"github.com/fatih/color"
	"github.com/oschwald/geoip2-golang"
)

var geoipCityDB *geoip2.Reader
var geoipAsDB *geoip2.Reader

func InitGeoIP() {
	var err error
	if geoipCityDB, err = geoip2.Open("./geodata/GeoLite2-City.mmdb"); err != nil {
		fmt.Println("GeoIP DB open error: ", err)
	}
	if geoipAsDB, err = geoip2.Open("./geodata/GeoLite2-ASN.mmdb"); err != nil {
		fmt.Println("GeoIP DB open error :", err)
	}
}

func CloseGeoIP() {
	if geoipCityDB != nil {
		_ = geoipCityDB.Close()
	}
	if geoipAsDB != nil {
		_ = geoipAsDB.Close()
	}
}

// IPアドレスから国情報を取得
func LookupCountry(ipStr string) string {
	ip := net.ParseIP(ipStr)

	CountryRecord, _ := geoipCityDB.Country(ip)
	AsRecord, _ := geoipAsDB.ASN(ip)

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
