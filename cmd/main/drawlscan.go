package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"github.com/nagayon-935/DrawlScan/cmd/handler"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func goMain(args []string) int {
	optionMap := handler.Options(args)
	var (
		help         = optionMap["Help"].(bool)
		version      = optionMap["Version"].(bool)
		geoip        = optionMap["Geoip"].(bool)
		distFilePath = optionMap["OutputFile"].(string)
		filter       = optionMap["Filter"].(string)
		isAscii      = !optionMap["NoAscii"].(bool)
		iface        = optionMap["InterfaceName"].(string)
		count        = optionMap["Count"].(int)
		timeSec      = optionMap["Time"].(int)
	)

	if help {
		fmt.Println(handler.HelpMessage())
		return 0
	}

	if version {
		fmt.Println("Version: " + VERSION)
		return 0
	}

	if geoip {
		utils.InitGeoIP()
	}

	if iface == "" {
		iface = utils.AutoSelectInterface()
		if iface == "" {
			log.Fatal("No suitable interface found")
			return 1
		}
	}

	handle, err := pcap.OpenLive(iface, 65535, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	var distFile *os.File
	var pcapw *pcapgo.Writer
	if distFilePath != "" {
		ext := strings.ToLower(filepath.Ext(distFilePath))
		if ext != ".pcap" && ext != ".pcapng" {
			log.Fatalf("Output file must have .pcap or .pcapng extension: %s", distFilePath)
			return 1
		}
		var err error
		distFile, err = os.Create(distFilePath)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer distFile.Close()
		pcapw = pcapgo.NewWriter(distFile)
		if err := pcapw.WriteFileHeader(1600, handle.LinkType()); err != nil {
			log.Fatalf("WriteFileHeader: %v", err)
		}
	}

	if filter != "" {
		if err := handle.SetBPFFilter(filter); err != nil {
			log.Fatalf("Failed to set BPF filter: %v", err)
		}
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetChan := packetSource.Packets()

	var timeCh <-chan time.Time
	if timeSec > 0 {
		timeCh = time.After(time.Duration(timeSec) * time.Second)
	}
	received := 0

	fmt.Printf("Using interface: %s\n", iface)

	done := false
	start := time.Now()

	for !done {
		if os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
			fmt.Println("CI環境のためパケットキャプチャをスキップします")
			//testdataからpcapファイルをインポートしてテストを行う
			return 0
		}
		select {
		case packet, ok := <-packetChan:
			if !ok {
				done = true
			}
			var blocks []string
			for _, h := range handler.Handlers {
				if packet.Layer(h.LayerType) != nil {
					blocks = append(blocks, h.Handler(packet))
				}
			}
			if distFile != nil {
				if err := pcapw.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
					log.Fatalf("pcap.WritePacket(): %v", err)
				}
			}
			if geoip {
				if netLayer := packet.NetworkLayer(); netLayer != nil {
					src, dst := netLayer.NetworkFlow().Endpoints()
					blocks = append(blocks, utils.LookupCountry(src.String()))
					blocks = append(blocks, utils.LookupCountry(dst.String()))
				}
			}
			if isAscii {
				utils.PrintHorizontalBlocks(blocks)
			} else {
				fmt.Println(packet)
			}
			received++
			if count > 0 && received >= count {
				done = true
			}
		case <-timeCh:
			done = true
		}
		if count > 0 && received >= count {
			done = true
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Captured %d packets\n", received)
	fmt.Printf("Capture duration: %.2f seconds\n", elapsed.Seconds())
	return 0
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
