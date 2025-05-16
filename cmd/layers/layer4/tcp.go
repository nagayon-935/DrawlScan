package layer4

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintTcpLayer(packet gopacket.Packet) string {
	tcp := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)

	http := PrintTcpLayerWithHttpAndHttps(packet, tcp)

	return utils.RenderBlock("TCP Segment", []string{
		fmt.Sprintf("Src Port: %d", tcp.SrcPort),
		fmt.Sprintf("Dst Port: %d", tcp.DstPort),
		fmt.Sprintf("%s", http),
	}, color.New(color.FgMagenta))
}

func PrintTcpLayerWithHttpAndHttps(packet gopacket.Packet, tcp *layers.TCP) []string {
	var blocks []string
	// HTTP (TCP 80)
	if tcp.SrcPort == 80 || tcp.DstPort == 80 {
		if app := packet.ApplicationLayer(); app != nil {
			payload := string(app.Payload())
			if len(payload) > 64 {
				payload = payload[:64] + "..."
			}
			blocks = append(blocks, utils.RenderBlock("HTTP", []string{payload}, color.New(color.FgHiWhite)))
		}
	}

	// HTTPS (TCP 443)
	if tcp.SrcPort == 443 || tcp.DstPort == 443 {
		blocks = append(blocks, utils.RenderBlock("HTTPS", []string{
			"Encrypted Traffic",
		}, color.New(color.FgHiBlue)))
	}

	return blocks
}
