package layer3

import (
	"testing"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func TestPrintICMPLayer(t *testing.T) {
	icmp := &layers.ICMPv4{
		TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
		Id:       1234,
		Seq:      1,
	}
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	if err := icmp.SerializeTo(buf, opts); err != nil {
		t.Fatalf("failed to serialize icmpv4 layer: %v", err)
	}
	packet := gopacket.NewPacket(buf.Bytes(), layers.LayerTypeICMPv4, gopacket.Default)

	want := utils.RenderBlock("ICMP Packet", []string{
		"Type: EchoRequest",
	}, color.New(color.FgYellow))

	got := PrintIcmpLayer(packet)
	if got != want {
		t.Errorf("PrintICMPLayer() = %v, want %v", got, want)
	}
}
