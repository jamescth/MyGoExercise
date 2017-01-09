package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type myPacket struct {
	SrcMac  net.HardwareAddr
	DstMac  net.HardwareAddr
	SrcIP   net.IP
	DstIP   net.IP
	SrcPort layers.TCPPort
	DstPort layers.TCPPort

	// Length is the IP length. use it to subtract ip header and tcp header for payload size
	Length uint16
	Time   time.Time
	Seq    uint32
	Ack    uint32
	Window uint16
	Size   int
}

type PktStream struct {
	SrcMac  net.HardwareAddr
	DstMac  net.HardwareAddr
	SrcIP   net.IP
	DstIP   net.IP
	SrcPort layers.TCPPort
	DstPort layers.TCPPort

	WinScale int
	pkts     []myPacket
}

func processPkts() {
	var (
		ethLayer layers.Ethernet
		ipLayer  layers.IPv4
		tcpLayer layers.TCP
	)

	pktnum := 0
	tcpnum := 0

	h, err := pcap.OpenOffline("./ddmc.fw.pcap")
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	// Loop through packets in file
	pktSrc := gopacket.NewPacketSource(h, h.LinkType())

	var pp []PktStream

	for packet := range pktSrc.Packets() {
		pktnum++

		// parsing packet
		parser := gopacket.NewDecodingLayerParser(
			layers.LayerTypeEthernet,
			&ethLayer,
			&ipLayer,
			&tcpLayer,
		)
		foundLayerTypes := []gopacket.LayerType{}

		if err := parser.DecodeLayers(packet.Data(), &foundLayerTypes); err != nil {
			// if we cant decode, then, it's not part of the list above
			continue
		}

		tcpnum++
		// if we are here, this is a TCP pkt

		var p myPacket

		p.SrcIP = ipLayer.SrcIP
		p.DstIP = ipLayer.DstIP
		p.SrcPort = tcpLayer.SrcPort
		p.DstPort = tcpLayer.DstPort
		p.Ack = tcpLayer.Ack
		p.Seq = tcpLayer.Seq
		p.Window = tcpLayer.Window
		p.Size = len(tcpLayer.Payload)
		p.Time = packet.Metadata().Timestamp
		p.Length = ipLayer.Length

		// if there is a tcp stream for this packet, insert the pkt to the stream and return
		f := func() bool {
			ppLen := len(pp)
			for i := 0; i < ppLen; i++ {
				switch {
				case pp[i].SrcPort == p.SrcPort && pp[i].DstPort == p.DstPort:
					pp[i].pkts = append(pp[i].pkts, p)
					return true
				case pp[i].SrcPort == p.DstPort && pp[i].DstPort == p.SrcPort:
					pp[i].pkts = append(pp[i].pkts, p)
					return true
				}
				continue
			}
			return false
		}

		if f() {
			continue
		}

		var newPStream PktStream

		newPStream.SrcPort = p.SrcPort
		newPStream.DstPort = p.DstPort
		newPStream.SrcIP = p.SrcIP
		newPStream.DstIP = p.DstIP
		newPStream.pkts = append(newPStream.pkts, p)
		pp = append(pp, newPStream)
		// pp.pkts = append(pp.pkts, p)
	}

	for _, ps := range pp {
		fmt.Printf("%s - %s, port %s - %s pkts %d\n", ps.SrcIP, ps.DstIP, ps.SrcPort, ps.DstPort, len(ps.pkts))

		var t time.Time
		for i, pkt := range ps.pkts {

			var du time.Duration = 0
			if i != 0 {
				du = pkt.Time.Sub(t)
			}

			fmt.Printf("   %d - %d, seq %10d ack %10d window %4d size %d time %s\n", pkt.SrcPort, pkt.DstPort, pkt.Seq, pkt.Ack, pkt.Window, pkt.Length /*pkt.Size*/, du)
			t = pkt.Time
		}
	}
	fmt.Println(pktnum, tcpnum, len(pp))
}

func main() {
	countPcapPkt()
}

func countPcapPkt() error {
	pktCount := 0

	h, err := pcap.OpenOffline("./ddmc.fw.pcap")
	defer h.Close()
	if err != nil {
		return err
	}

	for {
		_, _, err := h.ReadPacketData()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("total pkts %d\n", pktCount)
				return nil
			}
			return err
		}
		pktCount++
	}
}

func scanPkts() error {
	// pktCount := 0
	h, err := pcap.OpenOffline("./ddmc.fw.pcap")
	if err != nil {
		return err
	}
	defer h.Close()

	// Loop through packets in file
	pktSrc := gopacket.NewPacketSource(h, h.LinkType())
	for packet := range pktSrc.Packets() {
		packetInfo(packet)
	}

	return nil
}

func packetInfo(pkt gopacket.Packet) {
	var (
		ethLayer layers.Ethernet
		ipLayer  layers.IPv4
		tcpLayer layers.TCP
	)
	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeEthernet,
		&ethLayer,
		&ipLayer,
		&tcpLayer,
	)

	foundLayerTypes := []gopacket.LayerType{}

	if err := parser.DecodeLayers(pkt.Data(), &foundLayerTypes); err != nil {
		// fmt.Printf("err %s\n", err)
		// if we cant decode, then, it's not part of the list above
		return
	}

	if tcpLayer.SYN {
		for _, opt := range tcpLayer.Options {
			if opt.OptionType == layers.TCPOptionKindWindowScale {
				fmt.Printf("%s to %s, %s to %s len %d 0x%s\n", ipLayer.SrcIP, ipLayer.DstIP, tcpLayer.SrcPort, tcpLayer.DstPort, opt.OptionLength, hex.EncodeToString(opt.OptionData))
			}
		}
	}
}
