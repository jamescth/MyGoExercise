package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/pprof"

	"pkta/util"

	"github.com/google/gopacket/pcap"
)

// variables for logging
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func main() {
	// input arguments
	var (
		pcapfile = flag.String("pcap", "", "The tcpdump pcap file")
		cprof    = flag.String("cprof", "", "profile output")
		mprof    = flag.String("mprof", "", "write memory profile")
		info     = flag.String("info", "", "the detail output")
		trace    = flag.String("trace", "", "the trace output")
		// preview  = flag.Bool("preview", false, "scan pkts only")
	)

	flag.Usage = func() {
		fmt.Printf("usage of %s:\n", os.Args[0])
		fmt.Println("  pkt_anlyzer -pcap=<pcap file> [-out=<output file>]")
		flag.PrintDefaults()
	}

	//*****************************************************************
	// parsing the cmdline inputs
	//*****************************************************************
	flag.Parse()

	// if no pcap file, exit
	if *pcapfile == "" {
		flag.Usage()
		os.Exit(1)
	}

	var (
		err   error
		finfo *os.File
	)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	// set up the log env
	if *info != "" {
		finfo, err = os.Create(*info)
		util.Check(err)
		defer finfo.Close()
		// Init(finfo, finfo, os.Stdout, os.Stderr)

		Info = log.New(finfo,
			"INFO: ",
			//log.Ldate|log.Ltime|log.Lshortfile)
			log.Lshortfile)

	} else {
		Info = log.New(ioutil.Discard,
			"INFO: ",
			//log.Ldate|log.Ltime|log.Lshortfile)
			log.Lshortfile)
		//Init(ioutil.Discard, ioutil.Discard, os.Stdout, os.Stderr)
	}

	// set up the log env
	if *trace != "" {
		finfo, err = os.Create(*trace)
		util.Check(err)
		defer finfo.Close()
		// Init(finfo, finfo, os.Stdout, os.Stderr)

		Trace = log.New(finfo,
			"TRACE: ",
			//log.Ldate|log.Ltime|log.Lshortfile)
			log.Lshortfile)

	} else {
		Trace = log.New(ioutil.Discard,
			"TRACE: ",
			//log.Ldate|log.Ltime|log.Lshortfile)
			log.Lshortfile)
		//Init(ioutil.Discard, ioutil.Discard, os.Stdout, os.Stderr)
	}

	if *cprof != "" {
		f, err := os.Create(*cprof)
		if err != nil {
			log.Fatal("Profile create failed, ", err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *mprof != "" {
		f, err := os.Create(*mprof)
		if err != nil {
			log.Fatal("memory profile creation failed, ", err)
		}
		pprof.WriteHeapProfile(f)
		defer f.Close()
	}

	// go tool pprof http://localhost:6060/debug/pprof/heap
	// curl -s http://localhost:6060/debug/pprof/heap > base.heap
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

}

func init() {
	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

}

func PreScanPkts(pfile string) (int64, error) {
	h, err := pcap.OpenOffline(pfile)
	if err != nil {
		return 0, err
	}

	var pkt_cnt int64
	// quick scan the file to get pkt count
	for {
		_, _, err := h.ReadPacketData()
		if err == io.EOF {
			h.Close()
			return pkt_cnt, nil
		} else if err != nil {
			return 0, err
		}
		pkt_cnt++
	}
}
