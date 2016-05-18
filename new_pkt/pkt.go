package main

import (
	"fmt"
	"io"
	"net"
	"time"

	// pcap
	"github.com/google/gopacket/pcap"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type list struct {
	next      *pkt_fields
	next_0win *pkt_fields
	//next_ack  *pkt_fields
	next_rst *pkt_fields
	next_syn *pkt_fields
	next_fin *pkt_fields
	next_ooo *pkt_fields
	// next_dup *pkt_fields
}

type s_eth struct {
	sMac net.HardwareAddr
	dMac net.HardwareAddr
}

type s_ip struct {
	sIP net.IP
	dIP net.IP
}

type s_port struct {
	sPort layers.TCPPort
	dPort layers.TCPPort
}

type s_pkt struct {
	Seq     uint32
	Ack     uint32
	ack_pkt int64 // this pkt has been ack'ed by #
	pl_len  int
	f_syn   bool
	f_ack   bool
	f_fin   bool
	f_psh   bool
	f_rst   bool
	f_0win  uint16
}

type s_pkt_time struct {
	ts  time.Time
	rtt time.Duration
}

type pkt_fields struct {
	list
	pkt_idx int64
	s_eth
	s_ip
	s_port
	s_pkt
	s_pkt_time
}

func PreScanPkts(pfile *string) int64 {
	var pkt_count int64

	h, err := pcap.OpenOffline(*pfile)
	check(err)

	// pre Scan for # of pkts
	for {
		_, _, err := h.ReadPacketData()
		if err == io.EOF {
			// REVISIT: can we not to close the file.
			//          but use seek() instead?
			h.Close()
			return pkt_count
		}
		check(err)
		pkt_count++
	}

	// we should never reach here
	// we either return if io.EOF or panic for everything else
}

func ScanPkts(pfile *string, pkt_array []pkt_fields) {
	h, err := pcap.OpenOffline(*pfile)
	check(err)
	defer h.Close()

	var pkt_idx int64 = 0

	for {
		data, ci, err := h.ReadPacketData()
		if err == io.EOF {
			fmt.Printf("pkts %d tcp %d\n", pkt_idx, tcp_count)
			Info.Printf("pkts %d tcp %d\n", pkt_idx, tcp_count)
			return
		}
		check(err)
		pkt_idx++
		if pkt_idx%50000 == 0 {
			fmt.Printf("%d pkts processed\n", pkt_idx)
		}

		// REVISIT: we should use go scanLayer() instead
		//          need to create a channel error for this
		if err = scanLayer(pkt_array, pkt_idx, data, ci.Timestamp); err != nil {
			// do something??
			continue
		}
	}
}

func scanLayer(
	pkt_array []pkt_fields,
	pkt_idx int64,
	data []uint8,
	ts time.Time) error {

	var (
		eth layers.Ethernet
		ip4 layers.IPv4
		ip6 layers.IPv6
		tcp layers.TCP
		udp layers.UDP
		arp layers.ARP
		//llc layers.LLC
		payload gopacket.Payload
	)

	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp, &arp,
		&payload)

	decodedLayers := make([]gopacket.LayerType, 0, 10)

	if err := parser.DecodeLayers(data, &decodedLayers); err != nil {
		// Info.Printf("Parseing pkt %d error %s\n", pkt_idx, err)
		return err
	}

	for _, typ := range decodedLayers {
		switch typ {
		case layers.LayerTypeTCP:
			tcp_count++
			// REVISIT: cant use goroutine here due to race condition
			func() {
				idx := pkt_idx - 1
				pkt_array[idx].pkt_idx = pkt_idx
				pkt_array[idx].sMac = eth.SrcMAC
				pkt_array[idx].dMac = eth.DstMAC
				pkt_array[idx].sIP = ip4.SrcIP
				pkt_array[idx].dIP = ip4.DstIP
				pkt_array[idx].sPort = tcp.SrcPort
				pkt_array[idx].dPort = tcp.DstPort
				pkt_array[idx].Seq = tcp.Seq
				pkt_array[idx].Ack = tcp.Ack
				pkt_array[idx].pl_len = len(tcp.Payload)
				if tcp.SYN {
					pkt_array[idx].f_syn = true
				}
				if tcp.ACK {
					pkt_array[idx].f_ack = true
				}
				if tcp.FIN {
					pkt_array[idx].f_fin = true
				}
				if tcp.PSH {
					pkt_array[idx].f_psh = true
				}
				if tcp.RST {
					pkt_array[idx].f_rst = true
				}
				pkt_array[idx].f_0win = tcp.Window
				pkt_array[idx].ts = ts

				Trace.Printf("%d [%s %s] [%s %s] [%s %s] %d",
					pkt_array[idx].pkt_idx,
					pkt_array[idx].sMac, pkt_array[idx].dMac,
					pkt_array[idx].sIP, pkt_array[idx].dIP,
					pkt_array[idx].sPort, pkt_array[idx].dPort,
					pkt_array[idx].pl_len)
				Trace.Println("  ", pkt_array[idx].s_pkt, pkt_array[idx].s_pkt_time)
			}()
		} // switch typ
	} // for decodedLayers

	return nil
}
