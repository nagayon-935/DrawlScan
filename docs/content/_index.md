---
title: DrawlScan
toc: false
---

This is the DrawlScan page.

## Explore

{{< cards >}}
  {{< card link="demo" title="Demo" icon="terminal" >}}
  {{< card link="docs" title="Docs" icon="book-open">}}
{{< /cards >}}

## Documentation

DrawlScan is a CUI-based network watcher tool written in Go.
It captures packets from a specified network interface and visualizes their structure and origin using ASCII art, GeoIP, and reverse DNS lookup.

It’s like tcpdump, but with an artistic flair.

Key Features:  
    •   🎨 Visualize packet structures (Ethernet/IP/TCP/UDP/etc.) as ASCII diagrams  
    •   🌍 GeoIP-based source/destination display  
    •   🧭 Lightweight, TUI-style interface — no GUI required  
    •   🐧 Perfect for learning, demos, or simply keeping an eye on your machine  

## Usage

```bash
drawlscan [OPTION]
OPTION
    -c, --count <NUM>              Capture only a specified number of packets
    -f, --filter <REGX>            Filter packets using a BPF (Berkeley Packet Filter) expression.
                                   You can specify filters such as:
                                     - "ip src 192.168.1.1"
                                     - "ip dst 192.168.1.2"
                                     - "ip host 192.168.1.1 and ip host 192.168.1.2"
                                     - "tcp port 80"
                                     - "udp port 53"
                                     - "icmp or icmp6"
                                     - "vlan 100"
                                     - "ip host 192.168.1.1 and tcp port 80"
    -g, --geoip                    Show GeoIP information for source and destination IP addresses
    -h, --help                     Display this help message
    -i, --interface <INTERFACE>    Specify the network interface to capture packets from (e.g., eth0, wlan0)
    -o, --output <FILE>            Save the captured packets to a file in PCAP format
    -r, --read <FILE>              Read packets from a PCAP file instead of capturing live traffic
    -t, --time <TIME>              Stop capturing after a specified number of seconds
    -v, --version                  Show version information
    --no-ascii                     Disable ASCII-art output
```

## Installation

> [!WARNING]
>Recommend installing via Homebrew.
>When building from source, note that the software depends on C language libraries. Depending on your OS, additional system modules may be > > required.

### Homebrew

```bash
brew install nagayon-935/tap/drawlscan
sudo ln -s "$(which drawlscan)" /usr/local/bin/drawlscan
sudo drawlscan 
```

### Compile yourself

```bash
git clone https://github.com/nagayon-935/DrawlScan.git
cd DrawkScan
go mod tidy
go build -o drawlscan cmd/main/drawlscan.go cmd/main/version.go
sudo ./drawlscan
```

> [!WARNING]
> On Linux, you may need to install the gcc compiler and libpcap-dev.
> Also, you may need to set CGO_ENABLED=1 during the build process.


### Docker

```bash
docker pull ghcr.io/nagayon-935/drawlscan:latest
docker run -it --rm  drawlscan 
```

## About

### Icon

![DrawlScan Icon](./logo.png "DrawlScan Icon")

### The project name(**DrawlScan**) comes from?

DrawlScan is a coined word that combines Draw (to represent packet visualization), Owl (a symbol of wisdom and observation), and Scan (for network traffic analysis)
