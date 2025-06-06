package layer7

import (
	"net"
	"testing"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func TestPrintDhcpLayer_Request(t *testing.T) {
	eth := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		Version:  4,
		IHL:      5,
		SrcIP:    net.IP{0, 0, 0, 0},
		DstIP:    net.IP{255, 255, 255, 255},
		Protocol: layers.IPProtocolUDP,
	}
	udp := &layers.UDP{
		SrcPort: 68,
		DstPort: 67,
	}
	_ = udp.SetNetworkLayerForChecksum(ip)

	dhcp := &layers.DHCPv4{
		Operation:    layers.DHCPOpRequest,
		YourClientIP: net.IP{0, 0, 0, 0},
		Options: []layers.DHCPOption{
			{Type: 50, Length: 4, Data: []byte{192, 168, 1, 100}},
			{Type: 54, Length: 4, Data: []byte{192, 168, 1, 254}},
		},
	}
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	if err := gopacket.SerializeLayers(buf, opts, eth, ip, udp, dhcp); err != nil {
		t.Fatalf("failed to serialize layers: %v", err)
	}
	packet := gopacket.NewPacket(buf.Bytes(), layers.LayerTypeEthernet, gopacket.Default)

	want := utils.RenderBlock("DHCP Request", []string{
		"Your IP: 0.0.0.0",
		"Subnet Mask: ",
		"Router: ",
		"Lease Time: ",
		"Server Identifier: 192.168.1.254",
	}, color.New(color.FgHiCyan))

	got := PrintDhcpLayer(packet)
	if got != want {
		t.Errorf("PrintDhcpLayer() = %v, want %v", got, want)
	}
}

func TestPrintDhcpLayer_Reply(t *testing.T) {
	eth := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		DstMAC:       net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		Version:  4,
		IHL:      5,
		SrcIP:    net.IP{192, 168, 1, 254},
		DstIP:    net.IP{192, 168, 1, 100},
		Protocol: layers.IPProtocolUDP,
	}
	udp := &layers.UDP{
		SrcPort: 67,
		DstPort: 68,
	}
	_ = udp.SetNetworkLayerForChecksum(ip)

	dhcp := &layers.DHCPv4{
		Operation:    layers.DHCPOpReply,
		YourClientIP: net.IP{192, 168, 1, 100},
		Options: []layers.DHCPOption{
			{Type: 1, Length: 4, Data: []byte{255, 255, 255, 0}},
			{Type: 3, Length: 4, Data: []byte{192, 168, 1, 1}},
			{Type: 6, Length: 4, Data: []byte{8, 8, 8, 8}},
			{Type: 15, Length: uint8(len("example.local")), Data: []byte("example.local")},
			{Type: 51, Length: 4, Data: []byte{0x00, 0x01, 0x51, 0x80}},
			{Type: 54, Length: 4, Data: []byte{192, 168, 1, 254}},
		},
	}
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	if err := gopacket.SerializeLayers(buf, opts, eth, ip, udp, dhcp); err != nil {
		t.Fatalf("failed to serialize layers: %v", err)
	}
	packet := gopacket.NewPacket(buf.Bytes(), layers.LayerTypeEthernet, gopacket.Default)

	want := utils.RenderBlock("DHCP Reply", []string{
		"Your IP: 192.168.1.100",
		"Subnet Mask: 255.255.255.0",
		"Router: 192.168.1.1",
		"Lease Time: 86400 seconds",
		"Server Identifier: 192.168.1.254",
		"Domain Name: example.local",
		"Domain Name Servers: 8.8.8.8",
	}, color.New(color.FgHiCyan))

	got := PrintDhcpLayer(packet)
	if got != want {
		t.Errorf("PrintDhcpLayer() = %v, want %v", got, want)
	}
}
