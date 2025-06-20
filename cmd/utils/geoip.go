package utils

import (
	"embed"
	"fmt"
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

// IPアドレスから国情報を取得
func LookupCountry(ipStr string) string {
	ip := net.ParseIP(ipStr)

	CountryRecord, _ := geoipCityDB.Country(ip)
	AsRecord, _ := geoipAsDB.ASN(ip)

	geoipInfo := []string{
		"IP: " + ipStr,
		"Country: " + CountryRecord.Country.Names["en"],
		"Organization: " + AsRecord.AutonomousSystemOrganization,
	}

	return RenderBlock(fmt.Sprintf("GeoIP"), geoipInfo, color.New(color.FgHiRed))
}
