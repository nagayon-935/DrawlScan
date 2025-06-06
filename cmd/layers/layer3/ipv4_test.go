package layer3

import (
	"net"
	"testing"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func TestPrintIPv4Layer(t *testing.T) {
	ip := &layers.IPv4{
		SrcIP:    net.IP{192, 168, 1, 10},
		DstIP:    net.IP{192, 168, 1, 20},
		Protocol: layers.IPProtocolTCP,
	}
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	if err := ip.SerializeTo(buf, opts); err != nil {
		t.Fatalf("failed to serialize ipv4 layer: %v", err)
	}
	packet := gopacket.NewPacket(buf.Bytes(), layers.LayerTypeIPv4, gopacket.Default)

	want := utils.RenderBlock("IPv4 Packet", []string{
		"Src IP: 192.168.1.10",
		"Dst IP: 192.168.1.20",
		"Protocol: TCP",
	}, color.New(color.FgGreen))

	got := PrintIPv4Layer(packet)
	if got != want {
		t.Errorf("PrintIPv4Layer() = %v, want %v", got, want)
	}
}
