package layer4

import (
	"net"
	"testing"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func TestPrintTCPLayer(t *testing.T) {
	// TCPパケットを生成
	eth := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		DstMAC:       net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		SrcIP:    net.IP{192, 168, 1, 10},
		DstIP:    net.IP{192, 168, 1, 20},
		Protocol: layers.IPProtocolTCP,
		Version:  4,
		IHL:      5,
	}
	tcp := &layers.TCP{
		SrcPort: 12345,
		DstPort: 80,
		SYN:     true,
		ACK:     true,
	}
	_ = tcp.SetNetworkLayerForChecksum(ip)

	// Application payload を追加
	payload := gopacket.Payload([]byte("GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"))

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	if err := gopacket.SerializeLayers(buf, opts, eth, ip, tcp, payload); err != nil {
		t.Fatalf("failed to serialize layers: %v", err)
	}
	packet := gopacket.NewPacket(buf.Bytes(), layers.LayerTypeEthernet, gopacket.Default)

	want := utils.RenderBlock("TCP Packet", []string{
		"Src Port: 12345",
		"Dst Port: 80",
		"Flags: SYN ACK",
	}, color.New(color.FgMagenta))

	got := PrintTcpLayer(packet)
	if got != want {
		t.Errorf("PrintTCPLayer() = %v, want %v", got, want)
	}
}

func TestTcpFlagsString(t *testing.T) {
	tcp := &layers.TCP{
		SYN: true,
		ACK: true,
		FIN: false,
	}
	flags := tcpFlagsString(tcp)
	want := "SYN ACK"
	if flags != want {
		t.Errorf("tcpFlagsString() = %v, want %v", flags, want)
	}
}
