package layer7

import (
	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintHttpsLayer(packet gopacket.Packet) string {
	// HTTPS (TCP 443)
	return utils.RenderBlock("HTTPS", []string{
		"Encrypted Traffic",
	}, color.New(color.FgHiBlue))
}
