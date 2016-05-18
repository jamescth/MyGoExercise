/*
 * This tool is to analyze pcap files.
 *
 * func:
 * 1. show the date of the connection
 * 2. Based on L4, show it's ave. RTT & the outliers (>300ms?)
 * 3. Show Re-transmission or double ack'ed
 * 4. 0 window
 * 5. list sync, fin, and rst
 *
 * To build:
 * go build -a --ldflags '-extldflags "-lm -lstdc++ -static"' -i
 *
 * gdb
 * gdb pkt_anlzer
 *	// do it after setup breakpointers
 * 	run -pcap=./ddmc.fw.pcap
 * b main.main
 * b main.(*PcapFile).parse_pkt
 * p *parser
 * b main.(*L4List).CheckPkt
 * b /auto/home12/hoj9/golang/src/pkt_anlzer/pkts.go:474
 *
 * bug 139439: big pcap
 *
 * Profiling
 *
 * -base /path/to/heap.profile 
 *	allows to compare current profile with some base (what allocations are made over time)
 * -inuse_space
 *	display an amount of memory in use
 * -inuse_objects
 *	display a nuber of objects in use
 *
 * example:
 * 	go tool pprof -inuse_objects -base base.heap /path/binary /tmp/current.heap
 *
 * pprof cmds:
 *	top	show where most allocations happened
 *	list	show annotated source
 *	web	display profile graph in a browser
 *
 *	./pkt_anlzer -pcap=ddmc.fw.pcap -cprof=prof_out.prof -mprof=mprof_out.prof
 *	go tool pprof pkt_anlzer prof_out.prof
 *	(pprof) top10
 *  go tool pprof --text pkt_anlzer prof_out.prof
 *
 * go tool pprof --inuse_objects pkt_anlzer mprof_out.prof
 * go tool pprof --inuse_space pkt_anlzer mprof_out.prof
 * go tool pprof --text pkt_anlzer mprof_out.prof
 *
 * # this requires firefox to be installed
 * go tool pprof ./test /tmp/cprof.prof 
 * (pprof) web
 * 
 * # create a pdf output
 * go tool pprof --pdf ./test /tmp/cprof.prof > callgraph.pdf
 * 
 * # list pprof commands:
 * go tool pprof 2>&1 | less
 * 
 */
package main

