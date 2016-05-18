// Copyright 2015 Jenchang Ho All rights reserved.
//
// Use of this source code is governed by a BSD-style license

/*
	main()
		1. cmdline arguments

		=> PreScanPkts()
		=> ScanPkts()
		=> ScanPort()
		=> ShowPortList()

	PreScanPkts()
		1. use (*pcapHandle).ReadPaceket() to walkthrough all pkts
		2. save the pkt # in global variable 'pkt_count'
		3. close pcapHandle

	ScanPkts()
		1. use (*pcapHandle).ReadPaceket() to walkthrough all pkts
		2. create [pkt_count]pkt_fields in global 'pktf'
		3. for each pkt
			=> call ScanLayer(pkt_idx, data, ci.Timestamp)

	ScanLayer(p_idx, data, time.Time)
		1. Decode layers (eth, ip, tcp, udp, arp, payload)
		2. go func() fillup pktf[] 
		3. goroutine sync??

	ScanPort()
		1. loop on pkt_count
		2.   skip if pktf[idx].pkt_idx == 0
		3.   loop on pList ([]portUnit)
		4.     if pkt's ports match pList element
		5.        add the pkt to the end of pkt.next list
		6.     else
		7.        a new portUnit append to pList, and add pkt to it

	ShowPortList()
		1. loop on pList
		2.   print each portUnit's ports, and walkthrough its pkts
*/

/*

*/

package main
