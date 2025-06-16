package main

import (
	"fmt"
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

func processAndPrintPacket(packet gopacket.Packet, geoip bool, isAscii bool) {
	var blocks []string
	for _, h := range handler.Handlers {
		if packet.Layer(h.LayerType) != nil {
			blocks = append(blocks, h.Handler(packet))
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
}

func goMain(args []string) int {
	optionMap := handler.Options(args)
	var (
		count         = optionMap["Count"].(int)
		filter        = optionMap["Filter"].(string)
		geoip         = optionMap["Geoip"].(bool)
		help          = optionMap["Help"].(bool)
		iface         = optionMap["InterfaceName"].(string)
		writeFilePath = optionMap["OutputFile"].(string)
		readFilePath  = optionMap["ReadFile"].(string)
		timeSec       = optionMap["Time"].(int)
		version       = optionMap["Version"].(bool)
		isAscii       = !optionMap["NoAscii"].(bool)
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

	if readFilePath != "" {
		handle, err := pcap.OpenOffline(readFilePath)
		if filter != "" {
			if err := handle.SetBPFFilter(filter); err != nil {
				fmt.Println("Failed to set BPF filter: ", err)
				return 1
			}
		}
		if err != nil {
			fmt.Println("Failed to open pcap file:", err)
			return 1
		}
		defer handle.Close()
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		received := 0
		for packet := range packetSource.Packets() {
			processAndPrintPacket(packet, geoip, isAscii)
			received++
			if count > 0 && received >= count {
				break
			}
		}
		fmt.Printf("Read %d packets from %s\n", received, readFilePath)
		return 0
	}

	if iface == "" {
		iface = utils.AutoSelectInterface()
		if iface == "" {
			fmt.Println("No suitable interface found")
			return 1
		}
	}

	handle, err := pcap.OpenLive(iface, 65535, true, pcap.BlockForever)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer handle.Close()

	if filter != "" {
		if err := handle.SetBPFFilter(filter); err != nil {
			fmt.Println("Failed to set BPF filter: ", err)
			return 1
		}
	}

	var distFile *os.File
	var pcapw *pcapgo.Writer
	if writeFilePath != "" {
		ext := strings.ToLower(filepath.Ext(writeFilePath))
		if ext != ".pcap" && ext != ".pcapng" {
			fmt.Println("Output file must have .pcap or .pcapng extension: ", writeFilePath)
			return 1
		}
		var err error
		distFile, err = os.Create(writeFilePath)
		if err != nil {
			fmt.Println("Failed to create output file: ", err)
			return 1
		}
		defer distFile.Close()
		pcapw = pcapgo.NewWriter(distFile)
		if err := pcapw.WriteFileHeader(1600, handle.LinkType()); err != nil {
			fmt.Println("WriteFileHeader: ", err)
			return 1
		}
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetChan := packetSource.Packets()

	var timeCh <-chan time.Time
	if timeSec > 0 {
		timeCh = time.After(time.Duration(timeSec) * time.Second)
	}
	received := 0

	fmt.Println("Using interface: ", iface)

	done := false
	start := time.Now()

	for !done {
		select {
		case packet, ok := <-packetChan:
			if !ok {
				done = true
			}
			if distFile != nil {
				if err := pcapw.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
					fmt.Println("pcap.WritePacket(): ", err)
					return 1
				}
			}
			processAndPrintPacket(packet, geoip, isAscii)
			received++
		case <-timeCh:
			done = true
		}
		if count > 0 && received >= count {
			done = true
		}
	}

	utils.CloseGeoIP()
	elapsed := time.Since(start)
	fmt.Printf("Captured %d packets\n", received)
	fmt.Printf("Capture duration: %.2f seconds\n", elapsed.Seconds())
	return 0
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
