# DrawlScan

Observe your network traffic in real-time, rendered as ASCII art.

![Go Version1.23](https://img.shields.io/badge/go-v1.22-blue "Go Version1.23")
![MIT License](https://img.shields.io/badge/license-MIT-blue "MIT License")

## Overview

DrawlScan is a CUI-based network watcher tool written in Go.
It captures packets from a specified network interface and visualizes their structure and origin using ASCII art, GeoIP, and reverse DNS lookup.

It’s like tcpdump, but with an artistic flair.

Key Features:  
    •   🎨 Visualize packet structures (Ethernet/IP/TCP/UDP/etc.) as ASCII diagrams  
    •   🌍 GeoIP-based source/destination display  
    •   🔎 Reverse DNS lookup of IPs  
    •   🧭 Lightweight, TUI-style interface — no GUI required  
    •   🐧 Perfect for learning, demos, or simply keeping an eye on your machine  

## Usage

```bash
drawlscan [OPTION]
OPTION
    -c, --count <NUM>              Capture only a specified number of packets
    -d, --detail                   Show detailed packet information, including header fields and metadata
    -g, --geoip                    Show GeoIP information for source and destination IP addresses
    -h, --help                     Help message
    -i, --interface <INTERFACE>    Specify the network interface to capture packets from (e.g., eth0, wlan0).
    -o, --output <FILE>            Save the captured packets to a file in PCAP format
    -r, --rdns                     Perform reverse DNS lookups on source and destination IP addresses
    -s, --summary                  Display a summary of captured packets by protocol, source, etc
    -t, --timeout <TIME>           Stop capturing after a specified number of seconds
    -v, --version                  Version information
    --ascii,                       Enable ASCII-art visualization of packets and traffic (Default is enable)
    --no-ascii,                    Disable ASCII-art output
```

## Installation

### Homebrew

```bash
🚧 under construction 🚧
```

### Compile yourself

```bash
🚧 under construction 🚧
```

## About

### Author

* [nagayon-935](https://github.com/nagayon-935)

### icon

<img width="200" src="docs/image/image.png", alt="DrawlScan Icon">

### The project name(**DrawlScan**) comes from?

DrawlScan is a coined word that combines Draw (to represent packet visualization), Owl (a symbol of wisdom and observation), and Scan (for network traffic analysis)
