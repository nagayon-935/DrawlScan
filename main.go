package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const (
	snapshotLen = 1024
	promiscuous = true
	timeout     = pcap.BlockForever
)

func main() {
	var iface string
	flag.StringVar(&iface, "i", "", "network interface to capture packets from")
	flag.Parse()

	if iface == "" {
		iface = autoSelectInterface()
		if iface == "" {
			log.Fatal("No suitable interface found")
		}
		fmt.Printf("Using interface: %s\n", iface)
	}

	handle, err := pcap.OpenLive(iface, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		printPacketSummary(packet)
	}
}

func autoSelectInterface() string {
	ifs, err := net.Interfaces()
	if err != nil {
		log.Fatal("Failed to get interfaces:", err)
	}

	// インターフェイスをインデックスの昇順でソート
	sort.Slice(ifs, func(i, j int) bool {
		return ifs[i].Index < ifs[j].Index
	})

	for _, iface := range ifs {

		// インターフェイスが有効で、ループバックやトンネルでないものを選択
		if (iface.Flags&net.FlagUp != 0) && (iface.Flags&net.FlagLoopback == 0) && !strings.HasPrefix(iface.Name, "utun") {
			// 接続状態を確認
			if isInterfaceConnected(iface.Name) {
				return iface.Name
			}
		}
	}

	// 適切なインターフェイスが見つからない場合
	fmt.Println("No suitable interface found.")
	return ""
}

func isInterfaceConnected(ifaceName string) bool {
	// pcap.FindAllDevs を使用して、利用可能なインターフェイスを取得
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return false
	}

	// 指定されたインターフェイスがデバイスリストに含まれているか確認
	for _, dev := range devices {
		if dev.Name == ifaceName {
			// アドレスが存在する場合、接続されているとみなす
			if len(dev.Addresses) > 0 {
				return true
			}
			return false
		}
	}

	return false
}

func printPacketSummary(packet gopacket.Packet) {
	var blocks []string

	if ethLayer := packet.Layer(layers.LayerTypeEthernet); ethLayer != nil {
		eth := ethLayer.(*layers.Ethernet)
		blocks = append(blocks, renderBlock("Ethernet Frame", []string{
			"Src MAC: " + eth.SrcMAC.String(),
			"Dst MAC: " + eth.DstMAC.String(),
			"Type: " + eth.EthernetType.String(),
		}, color.New(color.FgCyan)))
	}

	if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
		arp := arpLayer.(*layers.ARP)
		blocks = append(blocks, renderBlock("ARP Packet", []string{
			"Sender IP: " + net.IP(arp.SourceProtAddress).String(),
			"Target IP: " + net.IP(arp.DstProtAddress).String(),
		}, color.New(color.FgHiYellow)))
	}

	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip := ipLayer.(*layers.IPv4)
		blocks = append(blocks, renderBlock("IPv4 Header", []string{
			"Src IP: " + ip.SrcIP.String(),
			"Dst IP: " + ip.DstIP.String(),
			"Protocol: " + ip.Protocol.String(),
		}, color.New(color.FgGreen)))
	}

	if icmpLayer := packet.Layer(layers.LayerTypeICMPv4); icmpLayer != nil {
		icmp := icmpLayer.(*layers.ICMPv4)
		blocks = append(blocks, renderBlock("ICMP", []string{
			"Type: " + icmp.TypeCode.String(),
		}, color.New(color.FgHiMagenta)))
	}

	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp := udpLayer.(*layers.UDP)
		blocks = append(blocks, renderBlock("UDP Datagram", []string{
			fmt.Sprintf("Src Port: %d", udp.SrcPort),
			fmt.Sprintf("Dst Port: %d", udp.DstPort),
		}, color.New(color.FgBlue)))

		// DHCP (UDP 67,68)
		if udp.SrcPort == 67 || udp.DstPort == 67 || udp.SrcPort == 68 || udp.DstPort == 68 {
			if dhcpLayer := packet.Layer(layers.LayerTypeDHCPv4); dhcpLayer != nil {
				dhcp := dhcpLayer.(*layers.DHCPv4)
				blocks = append(blocks, renderBlock("DHCP", []string{
					"Op: " + dhcp.Operation.String(),
					"Client IP: " + dhcp.ClientIP.String(),
					"Your IP: " + dhcp.YourClientIP.String(),
				}, color.New(color.FgHiCyan)))
			}
		}

		// DNS (UDP 53)
		if udp.SrcPort == 53 || udp.DstPort == 53 {
			if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
				dns := dnsLayer.(*layers.DNS)
				var queries []string
				for _, q := range dns.Questions {
					queries = append(queries, string(q.Name))
				}
				if len(queries) == 0 {
					queries = append(queries, "(no query)")
				}
				blocks = append(blocks, renderBlock("DNS Query", queries, color.New(color.FgHiYellow)))
			}
		}
	}

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp := tcpLayer.(*layers.TCP)

		// フラグとその状態をマップで管理
		flagMap := map[string]bool{
			"SYN": tcp.SYN,
			"ACK": tcp.ACK,
			"FIN": tcp.FIN,
			"RST": tcp.RST,
			"PSH": tcp.PSH,
			"URG": tcp.URG,
		}

		// True のフラグだけを収集
		var activeFlags []string
		for flag, isActive := range flagMap {
			if isActive {
				activeFlags = append(activeFlags, flag)
			}
		}

		blocks = append(blocks, renderBlock("TCP Segment", []string{
			fmt.Sprintf("Src Port: %d", tcp.SrcPort),
			fmt.Sprintf("Dst Port: %d", tcp.DstPort),
			"Flags: " + strings.Join(activeFlags, ", "),
		}, color.New(color.FgMagenta)))
	}

	if app := packet.ApplicationLayer(); app != nil {
		payload := string(app.Payload())
		if len(payload) > 64 {
			payload = payload[:64] + "..."
		}
		payload = strings.ReplaceAll(payload, "\n", "\\n")
		blocks = append(blocks, renderBlock("Payload", []string{payload}, color.New(color.FgWhite)))
	}

	printHorizontalBlocks(blocks)
}

func renderBlock(title string, lines []string, c *color.Color) string {
	var b strings.Builder
	c.Fprintf(&b, "+-------------------------------+\n")
	c.Fprintf(&b, "| %-29s |\n", title)
	for _, line := range lines {
		c.Fprintf(&b, "| %-29s |\n", line)
	}
	c.Fprintf(&b, "+-------------------------------+\n")
	return b.String()
}

func printHorizontalBlocks(blocks []string) {
	if len(blocks) == 0 {
		return
	}

	lines := make([][]string, len(blocks))
	maxLines := 0

	for i, b := range blocks {
		blockLines := strings.Split(strings.TrimSuffix(b, "\n"), "\n")
		lines[i] = blockLines
		if len(blockLines) > maxLines {
			maxLines = len(blockLines)
		}
	}

	for l := 0; l < maxLines; l++ {
		for i := 0; i < len(lines); i++ {
			if l < len(lines[i]) {
				fmt.Print(lines[i][l])
			} else {
				fmt.Print(strings.Repeat(" ", len(lines[i][0])))
			}
			fmt.Print("  ")
		}
		fmt.Println()
	}
	fmt.Println()
}
