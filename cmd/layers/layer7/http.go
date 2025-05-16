package layer7

import (
	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintHttpLayer(packet gopacket.Packet) string {
	tcp := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)

	if tcp.SrcPort == 80 || tcp.DstPort == 80 {
		if app := packet.ApplicationLayer(); app != nil {
			payload := string(app.Payload())
			if len(payload) > 64 {
				payload = payload[:64] + "..."
			}
			return utils.RenderBlock("HTTP", []string{payload}, color.New(color.FgHiWhite))
		}
	}

	return ""
}
