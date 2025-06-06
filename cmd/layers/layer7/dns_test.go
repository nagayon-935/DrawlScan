package layer7

import (
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func TestPrintDnsLayer(t *testing.T) {
	eth := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		DstMAC:       net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		Version:  4,
		IHL:      5,
		SrcIP:    net.IP{192, 0, 2, 1},
		DstIP:    net.IP{8, 8, 8, 8},
		Protocol: layers.IPProtocolUDP,
	}
	udp := &layers.UDP{
		SrcPort: 12345,
		DstPort: 53,
	}
	_ = udp.SetNetworkLayerForChecksum(ip)

	dns := &layers.DNS{
		ID:     0x1234,
		QR:     true,
		OpCode: layers.DNSOpCodeQuery,
		Questions: []layers.DNSQuestion{
			{Name: []byte("example.com"), Type: layers.DNSTypeA},
		},
		Answers: []layers.DNSResourceRecord{
			{
				Name: []byte("example.com"),
				Type: layers.DNSTypeA,
				IP:   net.IP{93, 184, 216, 34},
			},
		},
		ResponseCode: layers.DNSResponseCodeNoErr,
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	err := gopacket.SerializeLayers(buf, opts, eth, ip, udp, dns)
	if err != nil {
		t.Fatalf("failed to serialize layers: %v", err)
	}

	packet := gopacket.NewPacket(buf.Bytes(), layers.LayerTypeEthernet, gopacket.Default)

	got := PrintDnsLayer(packet)
	if got == "" {
		t.Error("PrintDnsLayer() returned empty string, expected DNS block output")
	} else {
		t.Log("DNS block:\n" + got)
	}
}
