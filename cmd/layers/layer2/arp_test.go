package layer2

import (
	"testing"

	"net"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func TestPrintARPLayer(t *testing.T) {
	// ARPパケットを生成
	eth := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		SourceProtAddress: []byte{192, 168, 1, 10},
		DstHwAddress:      []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstProtAddress:    []byte{192, 168, 1, 1},
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	if err := gopacket.SerializeLayers(buf, opts, eth, arp); err != nil {
		t.Fatalf("failed to serialize layers: %v", err)
	}
	packet := gopacket.NewPacket(buf.Bytes(), layers.LayerTypeEthernet, gopacket.Default)

	want := utils.RenderBlock("ARP Packet", []string{
		"Sender MAC: 00:11:22:33:44:55",
		"Sender IP: 192.168.1.10",
		"Target MAC: 00:00:00:00:00:00",
		"Target IP: 192.168.1.1",
	}, color.New(color.FgHiYellow))

	got := PrintARPLayer(packet)
	if got != want {
		t.Errorf("PrintARPLayer() = %v, want %v", got, want)
	}
}
