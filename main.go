package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <interface>", os.Args[0])
	}
	iface := os.Args[1]
	snaplen := int32(1600)
	promisc := false
	timeout := pcap.BlockForever

	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	fmt.Printf("Listening on %s...\n", iface)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Printf("[%s] Captured packet: %s\n", time.Now().Format(time.RFC3339), packet)
	}
}
