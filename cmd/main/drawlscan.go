package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/nagayon-935/DrawlScan/cmd/handler"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

const (
	snapshotLen = 1024
	promiscuous = true
	timeout     = pcap.BlockForever
)

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

		if (iface.Flags&net.FlagUp != 0) && (iface.Flags&net.FlagLoopback == 0) && !strings.HasPrefix(iface.Name, "utun") {
			if isInterfaceConnected(iface.Name) {
				return iface.Name
			}
		}
	}

	fmt.Println("No suitable interface found.")
	return ""
}

func isInterfaceConnected(ifaceName string) bool {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return false
	}

	for _, dev := range devices {
		if dev.Name == ifaceName {
			if len(dev.Addresses) > 0 {
				return true
			}
			return false
		}
	}

	return false
}

func goMain(args []string) int {
	optionMap := handler.OptionHandler(args)

	iface := ""
	if v, ok := optionMap["InterfaceName"].(string); ok && v != "" {
		iface = v
	} else {
		iface = autoSelectInterface()
		if iface == "" {
			log.Fatal("No suitable interface found")
		}
		fmt.Printf("Using interface: %s\n", iface)
	}

	count := 0
	if v, ok := optionMap["Count"].(int); ok {
		count = v
	}
	timeoutSec := 0
	if v, ok := optionMap["Timeout"].(int); ok {
		timeoutSec = v
	}

	handle, err := pcap.OpenLive(iface, 65535, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetChan := packetSource.Packets()

	var timeoutCh <-chan time.Time
	if timeoutSec > 0 {
		timeoutCh = time.After(time.Duration(timeoutSec) * time.Second)
	}

	received := 0
loop:
	for {
		if os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
			fmt.Println("CI環境のためパケットキャプチャをスキップします")
		}
		select {
		case packet, ok := <-packetChan:
			if !ok {
				break loop
			}
			var blocks []string
			for _, h := range handler.Handlers {
				if packet.Layer(h.LayerType) != nil {
					blocks = append(blocks, h.Handler(packet))
				}
			}
			utils.PrintHorizontalBlocks(blocks)
			received++
			if count > 0 && received >= count {
				break loop
			}
		case <-timeoutCh:
			fmt.Println("Timeout reached")
			break loop
		}
		// ループを抜ける条件
		if (count > 0 && received >= count) || (timeoutSec > 0 && timeoutCh != nil) {
			break loop
		}
	}

	fmt.Printf("Captured %d packets\n", received)
	return 0
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
