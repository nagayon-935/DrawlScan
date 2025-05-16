package layer7

import (
	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintHttpsLayer(packet gopacket.Packet) []string {
	// HTTPS (TCP 443)
	var blocks []string
	blocks = append(blocks, utils.RenderBlock("HTTPS", []string{
		"Encrypted Traffic",
	}, color.New(color.FgHiBlue)))
	return blocks
}
