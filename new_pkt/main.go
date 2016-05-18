package main
/*
 *
 * to build it statically:
 * 	go build -a --ldflags '-extldflags "-lm -lstdc++ -static"' -i
*/

import(
	"fmt"
	"flag"
	"os"
	"runtime"

	"io"
	"io/ioutil"
	"log"

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
		pcapfile = flag.String("pcap", "", "The tcpdump pcap file")
		cprof = flag.String("cprof", "", "profile output")
		mprof = flag.String("mprof", "", "write memory profile")
		info = flag.String("info", "", "the detail output")
		trace = flag.String("trace", "", "the trace output")
		preview = flag.Bool("preview", false, "scan pkts only")
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
	if *pcapfile == "" {
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

	// set up the log env
	if *trace != "" {
		finfo, err = os.Create(*trace)
		check(err)
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

	runtime.GOMAXPROCS(runtime.NumCPU())
	pkt_count := PreScanPkts(pcapfile)

	fmt.Println("pkt_count is ", pkt_count)

	if *preview {
		return
	}

	g_pkt_array := make([]pkt_fields, pkt_count)
	ScanPkts(pcapfile, g_pkt_array)

	port_array := ScanPort(g_pkt_array, pkt_count)

	ProcPortLists(port_array)

	ShowPortList(port_array)
}

/*
 * This func is called the 1st line in main().
 * Set up the logger
 *
 * for additional logger info, check
 * http://www.goinggo.net/2013/11/using-log-package-in-go.html
 */
func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		//log.Ldate | log.Ltime | log.Lshortfile)
		log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		//log.Ldate|log.Ltime|log.Lshortfile)
		log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
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
