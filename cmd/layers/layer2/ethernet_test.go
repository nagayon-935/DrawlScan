package layer2

import (
	"net"
	"testing"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func TestPrintEthernetLayer(t *testing.T) {
	// Ethernetパケットを生成
	eth := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		DstMAC:       net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	if err := eth.SerializeTo(buf, opts); err != nil {
		t.Fatalf("failed to serialize ethernet layer: %v", err)
	}
	packet := gopacket.NewPacket(buf.Bytes(), layers.LayerTypeEthernet, gopacket.Default)

	want := utils.RenderBlock("Ethernet Frame", []string{
		"Src MAC: 00:11:22:33:44:55",
		"Dst MAC: aa:bb:cc:dd:ee:ff",
		"Type: IPv4",
	}, color.New(color.FgCyan))

	got := PrintEthernetLayer(packet)
	if got != want {
		t.Errorf("PrintEthernetLayer() = %v, want %v", got, want)
	}
}
