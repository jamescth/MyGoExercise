package main

import(
	"fmt"
	"flag"
	"os"
	"runtime"
	"log"

	"io/ioutil"

	// typeof
	"reflect"

	// profiling
	"runtime/pprof"
	"net/http"
	_ "net/http/pprof"

)

// variables for logging
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// global variables
var(
	tcp_count int64 = 0
	//ktf []pkt_fields
	//ptList []portUnit
)

func main() {
	// input arguments
	var(
		pcap1 = flag.String("pcap1", "", "The tcpdump pcap file")
		pcap2 = flag.String("pcap2", "", "The tcpdump pcap file")
		cprof = flag.String("cprof", "", "profile output")
		mprof = flag.String("mprof", "", "write memory profile")
		info = flag.String("info", "", "the detail output")
		trace = flag.String("trace", "", "the trace output")
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

	var(
		err	error
		finfo	*os.File
	)

	// if no pcap file, exit
	if *pcap1 == "" || *pcap2 =="" {
		flag.Usage()
		os.Exit(1)
	}

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	// set up the log env
	if *info != "" {
		finfo, err = os.Create(*info)
		check(err)
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

	if *trace != "" {
		finfo, err = os.Create(*trace)
		check(err)
		defer finfo.Close()

		Trace = log.New(finfo,
			"TRACE: ",
			log.Lshortfile)

	} else {
		Trace = log.New(ioutil.Discard,
			"TRACE: ",
			log.Lshortfile)
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

	runtime.GOMAXPROCS(runtime.NumCPU())
	pkt1_count := PreScanPkts(pcap1)
	pkt2_count := PreScanPkts(pcap2)

	g_pkt1_array := make([]pkt_fields, pkt1_count)
	g_pkt2_array := make([]pkt_fields, pkt2_count)

	ScanPkts(pcap1, g_pkt1_array)
	ScanPkts(pcap2, g_pkt2_array)

	port_array1 := ScanPort(g_pkt1_array, pkt1_count)
	port_array2 := ScanPort(g_pkt2_array, pkt2_count)

	comp_arrays(port_array1, port_array2)

	ShowPortList(port_array1)
	ShowPortList(port_array2)

	// fmt.Printf("pkt1_count is %d, pkt2_count is %d\n", pkt1_count, pkt2_count)
	fmt.Println("")
}

/*
 * This func is for debug purpose.  It shows the type of the given var.
 * This helps getting more useful info how the var can be utilized.
 */
func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
