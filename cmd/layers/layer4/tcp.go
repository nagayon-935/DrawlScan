package layer4

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintTcpLayer(packet gopacket.Packet) string {
	tcp := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)

	return utils.RenderBlock("TCP Packet", []string{
		fmt.Sprintf("Src Port: %d", tcp.SrcPort),
		fmt.Sprintf("Dst Port: %d", tcp.DstPort),
		fmt.Sprintf("Flags: %s", tcpFlagsString(tcp)),
	}, color.New(color.FgMagenta))
}

func tcpFlagsString(tcp *layers.TCP) string {
	// 順序付きでフラグ名と値を並べる
	names := []string{"NS", "CWR", "ECE", "URG", "ACK", "PSH", "RST", "SYN", "FIN"}
	bools := []bool{tcp.NS, tcp.CWR, tcp.ECE, tcp.URG, tcp.ACK, tcp.PSH, tcp.RST, tcp.SYN, tcp.FIN}
	// strings.FieldsFuncで一気に抽出
	return strings.Join(
		func() []string {
			out := make([]string, 0, len(names))
			for i, v := range bools {
				if v {
					out = append(out, names[i])
				}
			}
			return out
		}(),
		" ",
	)
}
