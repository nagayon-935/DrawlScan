package layer7

import (
	"net"
	"testing"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

// HTTP: TCP/80, HTTPメソッドで始まるペイロード
func createHttpTestPacket() gopacket.Packet {
	eth := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		DstMAC:       net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		Version:  4,
		IHL:      5,
		TOS:      0,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
		SrcIP:    net.IP{192, 168, 0, 1},
		DstIP:    net.IP{192, 168, 0, 2},
	}
	tcp := &layers.TCP{
		SrcPort: 12345,
		DstPort: 80,
		Seq:     1105024978,
		ACK:     true,
		SYN:     true,
		Window:  14600,
	}
	_ = tcp.SetNetworkLayerForChecksum(ip)
	httpPayload := []byte("GET /index.html HTTP/1.1\r\nHost: example.com\r\n\r\n")
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	_ = gopacket.SerializeLayers(buf, opts, eth, ip, tcp, gopacket.Payload(httpPayload))
	return gopacket.NewPacket(buf.Bytes(), layers.LayerTypeEthernet, gopacket.Default)
}

func createHTTPSPacket(srcPort, dstPort layers.TCPPort, payload []byte) gopacket.Packet {
	ip := &layers.IPv4{
		Version:  4,
		IHL:      5,
		TOS:      0,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
		SrcIP:    net.IP{192, 168, 0, 1},
		DstIP:    net.IP{192, 168, 0, 2},
	}
	tcp := &layers.TCP{
		SrcPort: srcPort,
		DstPort: dstPort,
		Seq:     1000,
		ACK:     true,
		Window:  14600,
	}
	_ = tcp.SetNetworkLayerForChecksum(ip)

	// TLS handshake先頭バイト（0x16）を付けたペイロードを作成
	tlsPayload := append([]byte{0x16}, payload...)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	err := gopacket.SerializeLayers(buf, opts, ip, tcp, gopacket.Payload(tlsPayload))
	if err != nil {
		panic(err)
	}

	return gopacket.NewPacket(buf.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
}

func createQUICPacket(srcPort, dstPort layers.UDPPort, payload []byte) gopacket.Packet {
	ip := &layers.IPv4{
		Version:  4,
		IHL:      5,
		TOS:      0,
		TTL:      64,
		Protocol: layers.IPProtocolUDP,
		SrcIP:    net.IP{10, 0, 0, 1},
		DstIP:    net.IP{10, 0, 0, 2},
	}
	udp := &layers.UDP{
		SrcPort: srcPort,
		DstPort: dstPort,
	}
	_ = udp.SetNetworkLayerForChecksum(ip)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	err := gopacket.SerializeLayers(buf, opts, ip, udp, gopacket.Payload(payload))
	if err != nil {
		panic(err)
	}

	return gopacket.NewPacket(buf.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
}

// Unknown: TCP/1234, HTTPメソッド等で始まらないペイロード
func createUnknownTestPacket() gopacket.Packet {
	ip := &layers.IPv4{
		SrcIP:    []byte{192, 168, 1, 10},
		DstIP:    []byte{192, 168, 1, 20},
		Protocol: layers.IPProtocolTCP,
	}
	tcp := &layers.TCP{
		SrcPort: 12345,
		DstPort: 1234,
		Seq:     1,
	}
	tcp.Payload = []byte("SOMETHING ELSE")
	_ = tcp.SetNetworkLayerForChecksum(ip)
	buf := gopacket.NewSerializeBuffer()
	_ = gopacket.SerializeLayers(buf, gopacket.SerializeOptions{}, ip, tcp)
	return gopacket.NewPacket(buf.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
}

func TestDetectAppProtocol(t *testing.T) {
	tests := []struct {
		name     string
		packet   gopacket.Packet
		expected string
	}{
		{
			name:     "HTTP Request",
			packet:   createHttpTestPacket(),
			expected: "HTTP",
		},
		{
			name:     "HTTPS",
			packet:   createHTTPSPacket(1234, 443, []byte{0x16, 0x03, 0x01, 0x02}),
			expected: "HTTPS",
		},
		{
			name:     "QUIC",
			packet:   createQUICPacket(1234, 443, []byte{0x00, 0x01, 0x02, 0x03}),
			expected: "QUIC",
		},
		{
			name:     "Unknown",
			packet:   createUnknownTestPacket(),
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectAppProtocol(tt.packet)
			if got != tt.expected {
				t.Errorf("DetectAppProtocol() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPrintHttpInfo(t *testing.T) {
	packet := createHttpTestPacket()

	got := printHttpInfo(packet)
	want := utils.RenderBlock("HTTP", []string{
		"Method: GET",
		"Path: /index.html",
		"Host: example.com",
	}, color.New(color.FgHiYellow))

	if got != want {
		t.Errorf("printHttpInfo() = \n%v\nwant:\n%v", got, want)
	}
}
