---
title: Demo
---

This page shows examples of DrawlScan in use.

## Drawlscan

### Help

```bash
User$ ./drawlscan --help
Usage: drawlscan [OPTIONS]

OPTIONS:
    -c, --count <NUM>              Capture only a specified number of packets
    -f, --filter <REGX>            Filter packets using a BPF (Berkeley Packet Filter) expression.
                                   You can specify filters such as:
                                     - ip src 192.168.1.1
                                     - ip dst 192.168.1.2
                                     - ip host 192.168.1.1 and ip host 192.168.1.2
                                     - tcp port 80
                                     - udp port 53
                                     - icmp or icmp6
                                     - vlan 100
                                     - ip host 192.168.1.1 and tcp port 80
    -g, --geoip                    Show GeoIP information for source and destination IP addresses
    -h, --help                     Display this help message
    -i, --interface <INTERFACE>    Specify the network interface to capture packets from (e.g., eth0, wlan0)
    -o, --output <FILE>            Save the captured packets to a file in PCAP format
    -r, --read <FILE>              Read packets from a PCAP file instead of capturing live traffic
    -t, --time <TIME>              Stop capturing after a specified number of seconds
    -v, --version                  Show version information
    --no-ascii                     Disable ASCII-art output

```

### Version

```bash
User$ ./drawlscan --version
Version: 0.5.0
```

### Count

```bash
User$ ./drawlscan --count 3
Using interface:  en0
+------------------------------+  +---------------------------------+  
| Ethernet Frame               |  | ARP Packet                      |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Sender MAC: 22:be:b7:af:ac:b7   |  
| Dst MAC: ff:ff:ff:ff:ff:ff   |  | Sender IP: 10.70.70.235         |  
| Type: ARP                    |  | Target MAC: 00:00:00:00:00:00   |  
+------------------------------+  | Target IP: 169.254.169.254      |  
                                  +---------------------------------+  
                                                                       

+------------------------------+  +---------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet               |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235      |  | Src Port: 60601   |  | Encrypted Payload   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 142.250.206.234   |  | Dst Port: 443     |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP             |  | Flags: PSH ACK    |                           
+------------------------------+  +---------------------------+  +-------------------+                           
                                                                                                                 

+------------------------------+  +---------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet               |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 52:54:00:9d:8d:0c   |  | Src IP: 142.250.206.234   |  | Src Port: 443     |  | Encrypted Payload   |  
| Dst MAC: 22:be:b7:af:ac:b7   |  | Dst IP: 10.70.70.235      |  | Dst Port: 60601   |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP             |  | Flags: ACK        |                           
+------------------------------+  +---------------------------+  +-------------------+                           
                                                                                                                 

Captured 3 packets
Capture duration: 0.78 seconds
```

### Filter

```bash
User$ ./drawlscan --filter "icmp"
Using interface:  en0
+------------------------------+  +--------------------------+  +--------------------------------------+  
| Ethernet Frame               |  | IPv4 Packet              |  | ICMP Packet                          |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235     |  | Type: DestinationUnreachable(Port)   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 172.64.148.171   |  +--------------------------------------+  
| Type: IPv4                   |  | Protocol: ICMPv4         |                                            
+------------------------------+  +--------------------------+                                            
                                                                                                          

+------------------------------+  +------------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | ICMP Packet         |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Type: EchoRequest   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 10.70.71.254   |  +---------------------+  
| Type: IPv4                   |  | Protocol: ICMPv4       |                           
+------------------------------+  +------------------------+                           
                                                                                       

+------------------------------+  +------------------------+  +-------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | ICMP Packet       |  
| Src MAC: 52:54:00:9d:8d:0c   |  | Src IP: 10.70.71.254   |  | Type: EchoReply   |  
| Dst MAC: 22:be:b7:af:ac:b7   |  | Dst IP: 10.70.70.235   |  +-------------------+  
| Type: IPv4                   |  | Protocol: ICMPv4       |                         
+------------------------------+  +------------------------+                         
                                                                                     
```

### GeoIP

```bash
User$ ./drawlscan --geoip
Using interface:  en0
+------------------------------+  +------------------------+  +-------------------+  +---------------------+  +---------------------------------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  | GeoIP                                       |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 60974   |  | Encrypted Payload   |  | IP: 20.205.69.80                            |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 20.205.69.80   |  | Dst Port: 443     |  +---------------------+  | Country: Hong Kong                          |  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: ACK        |                           | Organization: MICROSOFT-CORP-MSN-AS-BLOCK   |  
+------------------------------+  +------------------------+  +-------------------+                           +---------------------------------------------+  
                                                                                                                                                               

+------------------------------+  +--------------------------+  +---------------------+  +--------------------------+  
| Ethernet Frame               |  | IPv4 Packet              |  | ICMP Packet         |  | GeoIP                    |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235     |  | Type: EchoRequest   |  | IP: 172.217.25.164       |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 172.217.25.164   |  +---------------------+  | Country: United States   |  
| Type: IPv4                   |  | Protocol: ICMPv4         |                           | Organization: GOOGLE     |  
+------------------------------+  +--------------------------+                           +--------------------------+  
                                                                                                                       

+------------------------------+  +--------------------------+  +-------------------+  +--------------------------+  
| Ethernet Frame               |  | IPv4 Packet              |  | ICMP Packet       |  | GeoIP                    |  
| Src MAC: 52:54:00:9d:8d:0c   |  | Src IP: 172.217.25.164   |  | Type: EchoReply   |  | IP: 172.217.25.164       |  
| Dst MAC: 22:be:b7:af:ac:b7   |  | Dst IP: 10.70.70.235     |  +-------------------+  | Country: United States   |  
| Type: IPv4                   |  | Protocol: ICMPv4         |                         | Organization: GOOGLE     |  
+------------------------------+  +--------------------------+                         +--------------------------+  
                                                                                                                     

+------------------------------+  +------------------------+  +-------------------+  +---------------------+  +---------------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  | GeoIP                     |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 59309   |  | Encrypted Payload   |  | IP: 35.74.215.78          |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 35.74.215.78   |  | Dst Port: 443     |  +---------------------+  | Country: Japan            |  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |                           | Organization: AMAZON-02   |  
+------------------------------+  +------------------------+  +-------------------+                           +---------------------------+  
```

### Interface

```bash
User$ ./drawlscan --interface en18
Using interface:  en18
+------------------------------+  +---------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet               |  | UDP Packet        |  | QUIC                |  
| Src MAC: ae:c3:8d:25:11:e5   |  | Src IP: 10.70.70.247      |  | Src Port: 49309   |  | Encrypted Payload   |  
| Dst MAC: 01:00:5e:7f:ff:fa   |  | Dst IP: 239.255.255.250   |  | Dst Port: 1900    |  +---------------------+  
| Type: IPv4                   |  | Protocol: UDP             |  +-------------------+                           
+------------------------------+  +---------------------------+                                                  
                                                                                                                 

+------------------------------+  +---------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet               |  | UDP Packet        |  | QUIC                |  
| Src MAC: ae:c3:8d:25:11:e5   |  | Src IP: 10.70.70.247      |  | Src Port: 49309   |  | Encrypted Payload   |  
| Dst MAC: 01:00:5e:7f:ff:fa   |  | Dst IP: 239.255.255.250   |  | Dst Port: 1900    |  +---------------------+  
| Type: IPv4                   |  | Protocol: UDP             |  +-------------------+                           
+------------------------------+  +---------------------------+                                                  
```

### Output

```bash
User$ ./drawlscan --output test.pcap --count 5
Using interface:  en0
+------------------------------+  +------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 59309   |  | Encrypted Payload   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 35.74.215.78   |  | Dst Port: 443     |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |                           
+------------------------------+  +------------------------+  +-------------------+                           
                                                                                                              

+------------------------------+  +------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 59313   |  | Encrypted Payload   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 35.74.215.78   |  | Dst Port: 443     |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |                           
+------------------------------+  +------------------------+  +-------------------+                           
                                                                                                              

+------------------------------+  +------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 59331   |  | Encrypted Payload   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 35.74.215.78   |  | Dst Port: 443     |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |                           
+------------------------------+  +------------------------+  +-------------------+                           
                                                                                                              

+------------------------------+  +------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 59332   |  | Encrypted Payload   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 35.74.215.78   |  | Dst Port: 443     |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |                           
+------------------------------+  +------------------------+  +-------------------+                           
                                                                                                              

+------------------------------+  +------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 59310   |  | Encrypted Payload   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 35.74.215.78   |  | Dst Port: 443     |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |                           
+------------------------------+  +------------------------+  +-------------------+                           
                                                                                                              

Captured 5 packets
Capture duration: 0.27 seconds
```

### Read

```bash
User$ ./drawlscan --read test.pcap 
+------------------------------+  +------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 59309   |  | Encrypted Payload   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 35.74.215.78   |  | Dst Port: 443     |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |                           
+------------------------------+  +------------------------+  +-------------------+                           
                                                                                                              

+------------------------------+  +------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 59313   |  | Encrypted Payload   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 35.74.215.78   |  | Dst Port: 443     |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |                           
+------------------------------+  +------------------------+  +-------------------+                           
                                                                                                              

+------------------------------+  +------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 59331   |  | Encrypted Payload   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 35.74.215.78   |  | Dst Port: 443     |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |                           
+------------------------------+  +------------------------+  +-------------------+                           
                                                                                                              

+------------------------------+  +------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 59332   |  | Encrypted Payload   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 35.74.215.78   |  | Dst Port: 443     |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |                           
+------------------------------+  +------------------------+  +-------------------+                           
                                                                                                              

+------------------------------+  +------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 59310   |  | Encrypted Payload   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 35.74.215.78   |  | Dst Port: 443     |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |                           
+------------------------------+  +------------------------+  +-------------------+                           
                                                                                                              

Read 5 packets from test.pcap
```

### Time

```bash
User$ ./drawlscan --time 1
Using interface:  en0
+------------------------------+  +------------------------+  +----------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet           |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 7000       |  
| Dst MAC: 3e:aa:c5:8d:5b:f0   |  | Dst IP: 10.70.70.249   |  | Dst Port: 60406      |  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: SYN ACK ECE   |  
+------------------------------+  +------------------------+  +----------------------+  
                                                                                        

+------------------------------+  +------------------------+  +-------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  
| Src MAC: 3e:aa:c5:8d:5b:f0   |  | Src IP: 10.70.70.249   |  | Src Port: 60406   |  
| Dst MAC: 22:be:b7:af:ac:b7   |  | Dst IP: 10.70.70.235   |  | Dst Port: 7000    |  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: ACK        |  
+------------------------------+  +------------------------+  +-------------------+  
                                                                                     

+------------------------------+  +------------------------+  +-------------------+  +----------------------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTP                             |  
| Src MAC: 3e:aa:c5:8d:5b:f0   |  | Src IP: 10.70.70.249   |  | Src Port: 60406   |  | Method: GET                      |  
| Dst MAC: 22:be:b7:af:ac:b7   |  | Dst IP: 10.70.70.235   |  | Dst Port: 7000    |  | Path: /info?txtAirPlay&txtRAOP   |  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |  +----------------------------------+  
+------------------------------+  +------------------------+  +-------------------+                                        
                                                                                                                           

+------------------------------+  +------------------------+  +-------------------+  +---------------------+  
| Ethernet Frame               |  | IPv4 Packet            |  | TCP Packet        |  | HTTPS               |  
| Src MAC: 22:be:b7:af:ac:b7   |  | Src IP: 10.70.70.235   |  | Src Port: 59323   |  | Encrypted Payload   |  
| Dst MAC: 52:54:00:9d:8d:0c   |  | Dst IP: 35.74.215.78   |  | Dst Port: 443     |  +---------------------+  
| Type: IPv4                   |  | Protocol: TCP          |  | Flags: PSH ACK    |                           
+------------------------------+  +------------------------+  +-------------------+                           
                                                                                                              

Captured 4 packets
Capture duration: 1.00 seconds
```

### No-ASCII

```bash
User$ ./drawlscan --no-ascii  --count 3
Using interface:  en0
PACKET: 106 bytes, wire length 106 cap length 106 @ 2025-07-01 14:39:07.464201 +0900 JST
- Layer 1 (14 bytes) = Ethernet {Contents=[..14..] Payload=[..92..] SrcMAC=52:54:00:9d:8d:0c DstMAC=22:be:b7:af:ac:b7 EthernetType=IPv4 Length=0}
- Layer 2 (20 bytes) = IPv4     {Contents=[..20..] Payload=[..72..] Version=4 IHL=5 TOS=0 Length=92 Id=6850 Flags= FragOffset=0 TTL=120 Protocol=TCP Checksum=35551 SrcIP=34.117.41.85 DstIP=10.70.70.235 Options=[] Padding=[]}
- Layer 3 (32 bytes) = TCP      {Contents=[..32..] Payload=[..40..] SrcPort=443(https) DstPort=61188 Seq=2686491988 Ack=3952425489 DataOffset=8 FIN=false SYN=false RST=false PSH=true ACK=true URG=false ECE=false CWR=false NS=false Window=1014 Checksum=34901 Urgent=0 Options=[TCPOption(NOP:), TCPOption(NOP:), TCPOption(Timestamps:3926470033/114974378 0xea092d9106da5eaa)] Padding=[]}
- Layer 4 (40 bytes) = Payload  40 byte(s)

PACKET: 106 bytes, wire length 106 cap length 106 @ 2025-07-01 14:39:07.464204 +0900 JST
- Layer 1 (14 bytes) = Ethernet {Contents=[..14..] Payload=[..92..] SrcMAC=52:54:00:9d:8d:0c DstMAC=22:be:b7:af:ac:b7 EthernetType=IPv4 Length=0}
- Layer 2 (20 bytes) = IPv4     {Contents=[..20..] Payload=[..72..] Version=4 IHL=5 TOS=0 Length=92 Id=22995 Flags= FragOffset=0 TTL=120 Protocol=TCP Checksum=19406 SrcIP=34.117.41.85 DstIP=10.70.70.235 Options=[] Padding=[]}
- Layer 3 (32 bytes) = TCP      {Contents=[..32..] Payload=[..40..] SrcPort=443(https) DstPort=61190 Seq=2064359072 Ack=1280250472 DataOffset=8 FIN=false SYN=false RST=false PSH=true ACK=true URG=false ECE=false CWR=false NS=false Window=1014 Checksum=1387 Urgent=0 Options=[TCPOption(NOP:), TCPOption(NOP:), TCPOption(Timestamps:1891989485/234819839 0x70c577ed0dff10ff)] Padding=[]}
- Layer 4 (40 bytes) = Payload  40 byte(s)

PACKET: 66 bytes, wire length 66 cap length 66 @ 2025-07-01 14:39:07.464608 +0900 JST
- Layer 1 (14 bytes) = Ethernet {Contents=[..14..] Payload=[..52..] SrcMAC=22:be:b7:af:ac:b7 DstMAC=52:54:00:9d:8d:0c EthernetType=IPv4 Length=0}
- Layer 2 (20 bytes) = IPv4     {Contents=[..20..] Payload=[..32..] Version=4 IHL=5 TOS=0 Length=52 Id=0 Flags= FragOffset=0 TTL=64 Protocol=TCP Checksum=56777 SrcIP=10.70.70.235 DstIP=34.117.41.85 Options=[] Padding=[]}
- Layer 3 (32 bytes) = TCP      {Contents=[..32..] Payload=[] SrcPort=61188 DstPort=443(https) Seq=3952425489 Ack=2686492028 DataOffset=8 FIN=false SYN=false RST=false PSH=false ACK=true URG=false ECE=false CWR=false NS=false Window=2048 Checksum=62725 Urgent=0 Options=[TCPOption(NOP:), TCPOption(NOP:), TCPOption(Timestamps:114977347/3926470033 0x06da6a43ea092d91)] Padding=[]}

Captured 3 packets
Capture duration: 0.87 seconds
```
