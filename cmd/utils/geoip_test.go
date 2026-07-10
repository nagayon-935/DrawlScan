package utils

import (
	"os"
	"testing"

	"github.com/fatih/color"
)

func TestInitAndCloseGeoIP(t *testing.T) {
	if err := InitGeoIP(); err != nil {
		t.Fatalf("InitGeoIP() error = %v, want nil", err)
	}
	CloseGeoIP()
}

func TestLookupCountry_NilDB(t *testing.T) {
	// Arrange: simulate InitGeoIP not having been called (or having failed)
	prevCity, prevAs := geoipCityDB, geoipAsDB
	geoipCityDB, geoipAsDB = nil, nil
	defer func() { geoipCityDB, geoipAsDB = prevCity, prevAs }()

	// Act
	got := LookupCountry("133.220.131.100")

	// Assert
	if got != "invisible" {
		t.Errorf("LookupCountry() = %v, want %q when GeoIP DB is not initialized", got, "invisible")
	}
}

func TestLookupCountry(t *testing.T) {
	if os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("Skipping test in CI environment")
	}

	if err := InitGeoIP(); err != nil {
		t.Fatalf("InitGeoIP() error = %v, want nil", err)
	}
	defer CloseGeoIP()

	want := RenderBlock("GeoIP", []string{
		"IP: 133.220.131.100",
		"Country: Japan",
		"Organization: Research Organization of Information and Systems, National Institute of Informatics",
	}, color.New(color.FgHiRed))

	ipStr := "133.220.131.100"
	got := LookupCountry(ipStr)
	if got != want {
		t.Errorf("LookupCountry() = %v, want %v", got, want)
	}
}
