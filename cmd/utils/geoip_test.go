package utils

import (
	"os"
	"testing"

	"github.com/fatih/color"
)

func TestInitAndCloseGeoIP(t *testing.T) {
	// DBファイルが存在しない場合でもエラー出力のみでpanicしないことを確認
	InitGeoIP()
	CloseGeoIP()
}

func TestLookupCountry(t *testing.T) {
	if os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("Skipping test in CI environment")
	}

	InitGeoIP()
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
