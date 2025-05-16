package handler

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/layers/layer2"
	"github.com/nagayon-935/DrawlScan/cmd/layers/layer3"
	"github.com/nagayon-935/DrawlScan/cmd/layers/layer4"
	"github.com/nagayon-935/DrawlScan/cmd/layers/layer7"
)

type LayerHandler func(gopacket.Packet) string

var Handlers = []struct {
	LayerType gopacket.LayerType
	Handler   LayerHandler
}{
	{layers.LayerTypeEthernet, layer2.PrintEthernetLayer},
	{layers.LayerTypeIPv4, layer3.PrintIPv4Layer},
	{layers.LayerTypeICMPv4, layer3.PrintIcmpLayer},
	{layers.LayerTypeARP, layer2.PrintARPLayer},
	{layers.LayerTypeTCP, layer4.PrintTcpLayer},
	{layers.LayerTypeUDP, layer4.PrintUdpLayer},
	{layers.LayerTypeDNS, layer7.PrintDnsLayer},
	{layers.LayerTypeDHCPv4, layer7.PrintDhcpLayer},
	{layers.LayerTypeTCP, layer7.PrintHttpLayer},
	{layers.LayerTypeTCP, layer7.PrintHttpsLayer},
}
