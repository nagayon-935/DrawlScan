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

	// タイムアウト用
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

// packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
// for packet := range packetSource.Packets() {
// 	var blocks []string
// 	for _, h := range handler.Handlers {
// 		if packet.Layer(h.LayerType) != nil {
// 			blocks = append(blocks, h.Handler(packet))
// 		}
// 	}
// 	utils.PrintHorizontalBlocks(blocks)
// }

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
