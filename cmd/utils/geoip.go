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

func InitGeoIP() error {
	cityDB, err := files.ReadFile("geodata/GeoLite2-City.mmdb")
	if err != nil {
		return fmt.Errorf("failed to read GeoLite2-City.mmdb: %w", err)
	}
	asDB, err := files.ReadFile("geodata/GeoLite2-ASN.mmdb")
	if err != nil {
		return fmt.Errorf("failed to read GeoLite2-ASN.mmdb: %w", err)
	}
	geoipCityDB, err = geoip2.FromBytes(cityDB)
	if err != nil {
		return fmt.Errorf("failed to load GeoLite2-City database: %w", err)
	}
	geoipAsDB, err = geoip2.FromBytes(asDB)
	if err != nil {
		return fmt.Errorf("failed to load GeoLite2-ASN database: %w", err)
	}
	return nil
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
	if geoipCityDB == nil || geoipAsDB == nil {
		return "invisible"
	}

	ip := net.ParseIP(ipStr)

	CountryRecord, _ := geoipCityDB.Country(ip)
	AsRecord, _ := geoipAsDB.ASN(ip)

	country := CountryRecord.Country.Names["en"]
	org := AsRecord.AutonomousSystemOrganization

	if country == "" && org == "" {
		return "invisible"
	}

	geoipInfo := []string{
		"IP: " + ipStr,
		"Country: " + country,
		"Organization: " + org,
	}

	return RenderBlock("GeoIP", geoipInfo, color.New(color.FgHiRed))
}
