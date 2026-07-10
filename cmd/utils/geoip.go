package utils

import (
	"embed"
	"net"

	"github.com/fatih/color"
	"github.com/oschwald/geoip2-golang"
)

var geoipCityDB *geoip2.Reader
var geoipAsDB *geoip2.Reader

//go:embed geodata/*
var files embed.FS

func InitGeoIP() {
	cityDB, _ := files.ReadFile("geodata/GeoLite2-City.mmdb")
	asDB, _ := files.ReadFile("geodata/GeoLite2-ASN.mmdb")
	geoipCityDB, _ = geoip2.FromBytes(cityDB)
	geoipAsDB, _ = geoip2.FromBytes(asDB)
}

func CloseGeoIP() {
	if geoipCityDB != nil {
		_ = geoipCityDB.Close()
	}
	if geoipAsDB != nil {
		_ = geoipAsDB.Close()
	}
}

// LookupCountry returns a rendered GeoIP block for the given IP address and
// whether GeoIP data was found for it.
func LookupCountry(ipStr string) (string, bool) {
	ip := net.ParseIP(ipStr)

	CountryRecord, _ := geoipCityDB.Country(ip)
	AsRecord, _ := geoipAsDB.ASN(ip)

	country := CountryRecord.Country.Names["en"]
	org := AsRecord.AutonomousSystemOrganization

	if country == "" && org == "" {
		return "", false
	}

	geoipInfo := []string{
		"IP: " + ipStr,
		"Country: " + country,
		"Organization: " + org,
	}

	return RenderBlock("GeoIP", geoipInfo, color.New(color.FgHiRed)), true
}
