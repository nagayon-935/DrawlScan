package layer7

import (
	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintHttpLayer(packet gopacket.Packet) []string {
	var blocks []string
	if app := packet.ApplicationLayer(); app != nil {
		payload := string(app.Payload())
		if len(payload) > 64 {
			payload = payload[:64] + "..."
		}
		blocks = append(blocks, utils.RenderBlock("HTTP", []string{payload}, color.New(color.FgHiWhite)))
	}
	return blocks
}
