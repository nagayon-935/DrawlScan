package layer4

import (
	"net"
	"testing"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func TestPrintUDPLayer(t *testing.T) {
	eth := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		DstMAC:       net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		Version:  4,
		IHL:      5,
		SrcIP:    net.IP{192, 168, 1, 10},
		DstIP:    net.IP{192, 168, 1, 20},
		Protocol: layers.IPProtocolUDP,
	}
	udp := &layers.UDP{
		SrcPort: 12345,
		DstPort: 53,
	}
	_ = udp.SetNetworkLayerForChecksum(ip)

	payload := gopacket.Payload([]byte("dummy"))

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	if err := gopacket.SerializeLayers(buf, opts, eth, ip, udp, payload); err != nil {
		t.Fatalf("failed to serialize layers: %v", err)
	}

	packet := gopacket.NewPacket(buf.Bytes(), layers.LayerTypeEthernet, gopacket.Default)

	want := utils.RenderBlock("UDP Packet", []string{
		"Src Port: 12345",
		"Dst Port: 53",
	}, color.New(color.FgBlue))

	got := PrintUdpLayer(packet)
	if got != want {
		t.Errorf("PrintUDPLayer() = %v, want %v", got, want)
	}
}
