package main

import(
	"fmt"
	_ "log"
	"strings"

	"github.com/google/gopacket/layers"
)

const(
	PKT_PROGRESS = 50000
)

type portUnit struct {
	numPkt1  uint
	numPkt2  uint
	sPort    layers.TCPPort
	dPort    layers.TCPPort

	// link list for all pkts
	plist1_h *pkt_fields
	plist1_t *pkt_fields
	plist2_h *pkt_fields
	plist2_t *pkt_fields

	//ackList  *pkt_fields
	rstList  *pkt_fields
	finList  *pkt_fields
	synList  *pkt_fields
	zeroList *pkt_fields
}

/*
 This func is to generate Port Array based on pkt's IP & port number
*/
func ScanPort(
	pkt_array []pkt_fields,
	pkt_count int64) []portUnit {

	var i int64
	var port_array []portUnit

	// walkthrough all the pkts
	for i = 0; i < pkt_count; i++ {
		
		if i!=0 && i%PKT_PROGRESS == 0 {
			fmt.Printf("%d pkts scaned\n", i)
		}

		// if pkt_idx is 0, the pkt fmt (such as IPv6) is not support yet.
		if pkt_array[i].pkt_idx == 0 {
			continue
		}

		// At this point, all the pkts are where we have to handle.
		// We either create a new portUnit, or chain it to an existing linklist

		// If True, We need to create a new one
		// If False, we already chain it.  So, back to the top to get the next one
		ret := func() bool {
			for pt_idx, c := range port_array {
				if c.sPort == pkt_array[i].sPort && c.dPort == pkt_array[i].dPort {
					// we only work on the tail
					port_array[pt_idx].plist1_t.next = &pkt_array[i]
					port_array[pt_idx].plist1_t = &pkt_array[i]
					port_array[pt_idx].numPkt1++

					return false

				} else if c.dPort == pkt_array[i].sPort &&
					c.sPort == pkt_array[i].dPort {
					// list 2??

					if port_array[pt_idx].plist2_t == nil {
						// tail is nil, this only occurs list2
						// bcoz list1 will at least have 1 element from append()
						port_array[pt_idx].plist2_h = &pkt_array[i]
						port_array[pt_idx].plist2_t = &pkt_array[i]
						port_array[pt_idx].numPkt2++
						return false
					}

					port_array[pt_idx].plist2_t.next = &pkt_array[i]
					port_array[pt_idx].plist2_t = &pkt_array[i]
					port_array[pt_idx].numPkt2++

					return false

				} // else if
			} // for port_array
			return true
		}()

		if !ret {
			continue
		}

		// When we reach here, we need to create a new portUnit entry
		pu := portUnit{
			sPort:pkt_array[i].sPort,
			dPort:pkt_array[i].dPort,
			plist1_h:&pkt_array[i],
			plist1_t:&pkt_array[i],
			numPkt1:1}

		// REVISIT: replace append() w/ something else to reduce GC
		port_array = append(port_array, pu)
	} // for pkt_count

	return port_array
}

func comp_arrays(p_array1, p_array2 []portUnit) {
	var(
		p1_idx, p2_idx int
		pu1, pu2 portUnit
	)

	// loop through port arraies on both pcap files
	// and compare if IPs match.
	for p1_idx, pu1 = range p_array1 {
		p1_l1 := p_array1[p1_idx].plist1_h
		// p1_l2 := p_array1[p1_idx].plist2_h

		for p2_idx, pu2 = range p_array2 {
			p2_l1 := p_array2[p2_idx].plist1_h
			p2_l2 := p_array2[p2_idx].plist2_h

			//if p1_l1.sIP == p2_l1.sIP && p1_l1.dIP == p2_l1.dIP {
			if strings.Contains(string(p1_l1.sIP), string(p2_l1.sIP)) &&
			   strings.Contains(string(p1_l1.dIP), string(p2_l1.dIP)) {
				Info.Println("pcap1 sIP", p1_l1.sIP, "sPort", p1_l1.sPort, 
							"dIP", p1_l1.dIP, "dPort", p1_l1.dPort, pu1.numPkt1, pu1.numPkt2)
				Info.Println("pcap2 sIP", p2_l1.sIP, "sPort", p2_l1.sPort, 
							"dIP", p2_l1.dIP, "dPort", p2_l1.dPort, pu2.numPkt1, pu2.numPkt2)
			}
			//if p1_l1.sIP == p2_l2.sIP && p1_l1.dIP == p2_l2.dIP {
			if strings.Contains(string(p1_l1.sIP), string(p2_l2.sIP)) &&
			   strings.Contains(string(p1_l1.dIP), string(p2_l2.dIP)) {
				Info.Println("pcap1 sIP", p1_l1.sIP, "sPort", p1_l1.sPort, 
							"dIP", p1_l1.dIP, "dPort", p1_l1.dPort, pu1.numPkt1, pu1.numPkt2)
				Info.Println("pcap2 sIP", p2_l2.sIP, "sPort", p2_l2.sPort, 
							"dIP", p2_l2.dIP, "dPort", p2_l2.dPort, pu2.numPkt1, pu2.numPkt2)
			}
		} // range p_array2 

	} // range p_array1
}

func ShowPortList(port_array []portUnit) {
	Info.Println("ShowPortList()")

	for pidx, pu := range port_array {
		p := port_array[pidx].plist1_h

		fmt.Printf("Mac %s %s IP %s %s\n", p.sMac, p.dMac, p.sIP, p.dIP)
		fmt.Printf("port %s %s list1 pkts %d\n",
			pu.sPort, pu.dPort, pu.numPkt1)

		Info.Printf("\nMac %s %s IP %s %s port %s %s pkts %d %d\n",
			p.sMac, p.dMac, p.sIP, p.dIP, pu.sPort,
			pu.dPort, pu.numPkt1, pu.numPkt2)
		Info.Printf("  Seq       Ack    ack_pkt pl_len f_syn f_ack f_fin f_rst f_0win ts rtt\n")
		Info.Println("  ", p.s_pkt)
		// REVISIT: how can we improve the repeatition?
		fmt.Printf(" %d ", p.pkt_idx)
		for p.next != nil {
			p = p.next
			Info.Println("  ", p.s_pkt)
			fmt.Printf("%d ", p.pkt_idx)
		}
		fmt.Println("")
		Info.Println("")

		fmt.Println("")
		fmt.Printf("port %s %s list2 pkts %d\n",
			pu.dPort, pu.sPort, pu.numPkt2)
		p = port_array[pidx].plist2_h
		Info.Println("  ", p.s_pkt)
		fmt.Printf(" %d ", p.pkt_idx)
		for p.next != nil {
			p = p.next
			Info.Println("  ", p.s_pkt)
			fmt.Printf("%d ", p.pkt_idx)
		}
		fmt.Println("")

		fmt.Println("")
		if port_array[pidx].zeroList != nil {
			fmt.Printf(" 0 Win %d ", port_array[pidx].zeroList.pkt_idx)

			pz := port_array[pidx].zeroList.next_0win
			for pz != nil {
				fmt.Printf("%d ", pz.pkt_idx)
				pz = pz.next_0win
			}
			fmt.Println("")
		}

		/*
		fmt.Println("")
		if port_array[pidx].ackList != nil {
			fmt.Printf(" ACK %d ", port_array[pidx].ackList.pkt_idx)

			pz := port_array[pidx].ackList.next_ack
			for pz != nil {
				fmt.Printf("%d ", pz.pkt_idx)
				pz = pz.next_ack
			}
			fmt.Println("")
		}
		*/

		fmt.Println("")
		if port_array[pidx].synList != nil {
			fmt.Printf(" SYN %d ", port_array[pidx].synList.pkt_idx)

			pz := port_array[pidx].synList.next_syn
			for pz != nil {
				fmt.Printf("%d ", pz.pkt_idx)
				pz = pz.next_syn
			}
			fmt.Println("")
		}

		fmt.Println("")
		if port_array[pidx].finList != nil {
			fmt.Printf(" FIN %d ", port_array[pidx].finList.pkt_idx)

			pz := port_array[pidx].finList.next_fin
			for pz != nil {
				fmt.Printf("%d ", pz.pkt_idx)
				pz = pz.next_fin
			}
			fmt.Println("")
		}

		fmt.Println("")
		if port_array[pidx].rstList != nil {
			fmt.Printf(" RST %d ", port_array[pidx].rstList.pkt_idx)

			pz := port_array[pidx].rstList.next_rst
			for pz != nil {
				fmt.Printf("%d ", pz.pkt_idx)
				pz = pz.next_rst
			}
			fmt.Println("")
		}

		fmt.Println("\n**************************\n")

	}
}


