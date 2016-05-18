package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket/layers"
)

const (
	PKT_PROGRESS = 50000
)

type pkt_list struct {
	head *pkt_fields
	tail *pkt_fields
}

/*
 * the func insert an pkt to the linklist's tail.
 * if the linklist is empty, both head & tail will be set
 * to the given pkt.
 */
func (pl *pkt_list) insert(pkt *pkt_fields) {
	if pl.head == nil {
		pl.head = pkt
		pl.tail = pkt
		return
	}

	// REVISIT:
	// can't use next.  that's for pkt chain.
	pl.tail.next = pkt
	return
}

type portUnit struct {
	numPkt1 uint
	numPkt2 uint
	sPort   layers.TCPPort
	dPort   layers.TCPPort

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
	// oooList  pkt_list
	oooList *pkt_fields
	// dupList  *pkt_fields
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

		if i != 0 && i%PKT_PROGRESS == 0 {
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
			sPort:    pkt_array[i].sPort,
			dPort:    pkt_array[i].dPort,
			plist1_h: &pkt_array[i],
			plist1_t: &pkt_array[i],
			numPkt1:  1}

		// REVISIT: replace append() w/ something else to reduce GC
		port_array = append(port_array, pu)
	} // for pkt_count

	return port_array
}

func ProcPortLists(port_array []portUnit) {
	for pidx, _ := range port_array {
		ProcessList(&port_array[pidx], port_array[pidx].plist1_h)
		ProcessList(&port_array[pidx], port_array[pidx].plist2_h)

		// do the RTT calculation seperately.  Do a cleanup later
		check_pkts(&port_array[pidx])
	}
}

// we gonna walkthrough both lists for a connection
func check_pkts(pu *portUnit) {
	// for sender, check:
	// seq# = last pkt(seq# + len(payload))
	//   a. if seq# > last pkt(seq# + len(payload))
	//      missing pkt
	//   b. if seq# < last pkt(seq# + len(payload))
	//      a dup pkt or out-of-order
	//
	// for receiver (another list):
	// 1. make ack'er
	// 2. calculate rtt

	/*
		assuming we have the following pkts:

		    idx  dir seq  ack  pkt-len  current previous last rtt
		                                  index   index   rec
		      1   l1   1    1       10        1     nil   nil  n
		      2   l1  11    1       20        2       1   nil  n
		      3   l1  31    1       10        3       2   nil  n
		      4   l2   1   41        0        4     nil     3  y
		      5   l2   1   41       20        5       4     3  n
		      6   l1  41   21        0        6       3     5  y
	*/

	// the way we implement the list is to put the lowest idx pkt
	// in the list1.  So, we should be able to apply this when
	// walking through seq/ack in both lists

	// also, use pkt_idx to advance to the next pkt.
	// since this is a tcpdump, pkt_idx is first come first serve.
	if pu.plist2_h == nil {
		// if plist2 is nil, we only have 1 way conmunication. it's bad.
		fmt.Printf("  *** this is 1 way conmunication. ***\n")
		return
	}

	// we will have at least 1 pkt in plist1.
	// so, l1_pre_pkt will not be nil
	var (
		l1_cur_pkt *pkt_fields = pu.plist1_h
		l2_cur_pkt *pkt_fields = pu.plist2_h
		l1_pre_pkt *pkt_fields
		l2_pre_pkt *pkt_fields
	)

	// this more like an assertion
	// l1's pkt_idx must be smailler than l2
	if l1_cur_pkt.pkt_idx >= l2_cur_pkt.pkt_idx {
		log.Fatal("l1 pkt_idx ", l1_cur_pkt.pkt_idx,
			" >= l2 pkt_idx ", l2_cur_pkt.pkt_idx)
	}

	// compare the idx from both l4 connections, call
	// check_ooo() for the smaller idx pkt (FIFO)
	for l1_cur_pkt != nil && l2_cur_pkt != nil {
		if l1_pre_pkt == nil {
			// if previous pkt is nil, we are the 1st pkt
			l1_pre_pkt = l1_cur_pkt
			l1_cur_pkt = l1_cur_pkt.next
			continue
		}

		// assertion
		if l1_cur_pkt.pkt_idx == l2_cur_pkt.pkt_idx {
			log.Fatal("l1 pkt_idx ", l1_cur_pkt.pkt_idx,
				" == l2 pkt_idx ", l2_cur_pkt.pkt_idx)
		}

		if l1_cur_pkt.pkt_idx < l2_cur_pkt.pkt_idx {
			check_ooo(pu, l1_cur_pkt, l1_pre_pkt, l2_pre_pkt)
			l1_pre_pkt = l1_cur_pkt
			l1_cur_pkt = l1_cur_pkt.next
			continue
		}

		// l1_cur_pkt.pkt_idx > l2_cur_pkt.pkt_idx
		check_ooo(pu, l2_cur_pkt, l2_pre_pkt, l1_pre_pkt)
		l2_pre_pkt = l2_cur_pkt
		l2_cur_pkt = l2_cur_pkt.next
	}

}

func append_ooo(pu *portUnit, p *pkt_fields) {
	if pu.oooList == nil {
		pu.oooList = p
		return
	}

	// work through till the last one and insert
	var p_next *pkt_fields = pu.oooList
	for p_next.next_ooo != nil {
		p_next = p_next.next_ooo
	}
	p_next.next_ooo = p
}

// check out-of-order pkt
func check_ooo(pu *portUnit, p, pre, p_ack *pkt_fields) {

	// if pre == nil ???
	if pre != nil {
		ret := func() bool {
			seq := pre.Seq + uint32(pre.pl_len)

			if pre.f_syn || pre.f_fin {
				seq++
			}

			// if seq # is correct, return.
			if p.Seq != seq {
				append_ooo(pu, p)
				Trace.Printf("%d p.Seq %d %d pre.Seq %d pre.pl_len %d\n",
					p.pkt_idx, p.Seq, pre.pkt_idx, pre.Seq, pre.pl_len)
				return true
			}

			return false
		}()

		if ret {
			return
		}
	}

	if p_ack == nil {
		return
	}

	// REVISIT: check pkt 1100, 1102, 1104, 1106 ...
	// this pkt is ack.
	if p.Seq == p_ack.Ack {
		if p_ack.ack_pkt == 0 {
			p_ack.ack_pkt = p.pkt_idx
			p_ack.rtt = p.ts.Sub(p_ack.ts)
			if p_ack.rtt < 0 {
				p_ack.rtt = 1
			}
			return
		}
		return
	}

	ack := p_ack.Seq + uint32(p_ack.pl_len)
	if p_ack.f_syn || p_ack.f_fin {
		ack++
	}

	// the pkt's ack # should be the same as
	// the other direction's last pkt seq # from
	if p.Ack != ack {
		// append into OOO if it's not
		append_ooo(pu, p)
		Trace.Printf("%d p Seq %d Ack %d : %d p_ack Seq %d Ack %d\n",
			p.pkt_idx, p.Seq, p.Ack, p_ack.pkt_idx, p_ack.Seq, p_ack.Ack)
		return
	}
}

/*
func check_sender(pu *portUnit, p, pre *pkt_fields) {
	// sender
	// p = list1
	// verify_seq = last_p.seq + last_p.pl_len
	// if last_p.f_syn verify_seq++
	// if p.seq != verify_seq
	//    append OOO

	seq := pre.Seq + uint32(pre.pl_len)

	if pre.f_syn {
		seq++
	}

	// if seq # is correct, return.
	if p.Seq == seq {
		return
	}

	// if the head is nil, inseart to the head.
	// pu.oooList.insert(p)

	if pu.oooList == nil {
		pu.oooList = p
		return
	}

	// work through till the last one and insert
	var p_next *pkt_fields = pu.oooList
	for p_next.next_ooo != nil {
		p_next = p_next.next_ooo
	}
	p_next.next_ooo = p
}

func check_receiver(pu *portUnit, p, pre *pkt_fields) {
	// have another func to handle receiver check
	// receiver
	// verify_ack = last pkt(seq# + len(payload))
	// if last_p.f_syn verify_seq++
	// if p.ack == verify_ack
	ack := pre.Seq + uint32(pre.pl_len)
	if pre.f_syn {
		ack++
	}

	if p.Ack == ack && pre.ack_pkt == 0 {
		pre.ack_pkt = p.pkt_idx
		pre.rtt = p.ts.Sub(pre.ts)

		return
	} else if p.Ack != ack {
		if !p.f_ack && p.f_0win != 0 {
			// pu.oooList.insert(p)
			if pu.oooList == nil {
				pu.oooList = p
				return
			}

			// work through till the last one and insert
			var p_next *pkt_fields = pu.oooList
			for p_next.next_ooo != nil {
				p_next = p_next.next_ooo
			}
			p_next.next_ooo = p
		}
	}
	// if p.Ack == ack && pre.ack_pkt != 0
	// It's normal packet.  let go for now
}
*/

// func (pu *portUnit) ProcessList() {
func ProcessList(pu *portUnit, lst *pkt_fields) {
	var p *pkt_fields = lst
	// we need to walktrhough all pkts in lst (pktList.next)
	for {
		if p == nil {
			// if nil, we done for the list
			return
		}

		// fmt.Printf(" %d %v %x\n", p.pkt_idx, p.f_0win, pu.zeroList)
		// Trace.Printf("0win %d %v\n", p.pkt_idx, p.f_0win)
		if p.f_0win == 0 {
			if pu.zeroList == nil {
				pu.zeroList = p
			} else {
				// loop on next_0win
				// var p_next *pkt_fields = p.next_0win
				var p_next *pkt_fields = pu.zeroList

				for p_next.next_0win != nil {
					p_next = p_next.next_0win
				}

				p_next.next_0win = p
			}
		} // f_0wind
		/*
			Trace.Printf("ACK %d %v\n", p.pkt_idx, p.f_ack)
			if p.f_ack {
				if pu.ackList == nil {
					pu.ackList = p
				} else {
					// loop on next_0win
					var p_next *pkt_fields = pu.ackList

					for p_next.next_ack != nil {
						p_next = p_next.next_ack
					}

					p_next.next_ack = p
				}
			} // f_ack
		*/

		// Trace.Printf("syn %d %v\n", p.pkt_idx, p.f_syn)
		if p.f_syn {
			if pu.synList == nil {
				pu.synList = p
			} else {
				// loop on next_syn
				var p_next *pkt_fields = pu.synList

				for p_next.next_syn != nil {
					p_next = p_next.next_syn
				}

				p_next.next_syn = p
				//Trace.Printf("  p last %d next %d", p_next.pkt_idx,
				//	p_next.next_syn.pkt_idx)
			}
			//Trace.Printf("  synList %d", pu.synList.pkt_idx)
		} // f_syn

		//Trace.Printf("fin %d %v\n", p.pkt_idx, p.f_fin)
		if p.f_fin {
			if pu.finList == nil {
				pu.finList = p
			} else {
				// loop on next_fin
				var p_next *pkt_fields = pu.finList

				for p_next.next_fin != nil {
					p_next = p_next.next_fin
				}

				p_next.next_fin = p
				//Trace.Printf("  p last %d next %d",
				//	p_next.pkt_idx, p_next.next_fin.pkt_idx)
			}
			//Trace.Printf("  finList %d", pu.finList.pkt_idx)
		} // f_fin

		//Trace.Printf("rst %d %v\n", p.pkt_idx, p.f_rst)
		if p.f_rst {
			if pu.rstList == nil {
				pu.rstList = p
			} else {
				// loop on next_rst
				var p_next *pkt_fields = pu.rstList

				for p_next.next_rst != nil {
					p_next = p_next.next_rst
				}

				p_next.next_rst = p
				//Trace.Printf("  p last %d next %d",
				//	p_next.pkt_idx, p_next.next_rst.pkt_idx)
			}
			//Trace.Printf("  rstList %d", pu.rstList.pkt_idx)
		} // f_rst

		p = p.next
	} // for
}

func ShowPortList(port_array []portUnit) {
	Trace.Println("ShowPortList()")

	for pidx, pu := range port_array {
		p := port_array[pidx].plist1_h

		fmt.Printf("Mac %s %s IP %s %s\n", p.sMac, p.dMac, p.sIP, p.dIP)
		fmt.Printf("port %s %s list1 pkts %d\n",
			pu.sPort, pu.dPort, pu.numPkt1)

		Info.Printf("\nMac %s %s IP %s %s port %s %s pkts %d %d\n",
			p.sMac, p.dMac, p.sIP, p.dIP, pu.sPort,
			pu.dPort, pu.numPkt1, pu.numPkt2)
		Info.Printf("   Seq       Ack    ack_pkt pl_len f_syn f_ack f_fin f_rst f_0win ts rtt\n")
		Info.Println("  ", p.pkt_idx, p.s_pkt, p.s_pkt_time)
		// REVISIT: how can we improve the repeatition?
		fmt.Printf(" %d ", p.pkt_idx)
		for p.next != nil {
			p = p.next
			Info.Println("  ", p.pkt_idx, p.s_pkt, p.s_pkt_time)
			fmt.Printf("%d ", p.pkt_idx)
		}
		fmt.Println("")
		Info.Println("the other direction")

		fmt.Println("")
		fmt.Printf("port %s %s list2 pkts %d\n",
			pu.dPort, pu.sPort, pu.numPkt2)
		p = port_array[pidx].plist2_h
		Info.Println("  ", p.pkt_idx, p.s_pkt, p.s_pkt_time)
		fmt.Printf(" %d ", p.pkt_idx)
		for p.next != nil {
			p = p.next
			Info.Println("  ", p.pkt_idx, p.s_pkt, p.s_pkt_time)
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

		fmt.Println("")
		if port_array[pidx].oooList != nil {
			fmt.Printf(" OOO %d ", port_array[pidx].oooList.pkt_idx)

			pz := port_array[pidx].oooList.next_ooo
			for pz != nil {
				fmt.Printf("%d ", pz.pkt_idx)
				pz = pz.next_ooo
			}
			fmt.Println("")
		}

		fmt.Println("\n**************************\n")

	}
}
