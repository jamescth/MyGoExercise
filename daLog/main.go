package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

var (
	lpath  = flag.String("lpath", "", "log path to analyze")
	iBug   = flag.String("bug", "", "bugxxxxx to analyze")
	iPath  = flag.String("path", "", "{src} to analyze")
	sTime  = flag.String("startTime", "", "the starting timestamp")
	conf   = flag.String("conf", "config.json", "config file in JSON")
	flist  = flag.String("flist", "filelist.out", "file list output")
	httpON = flag.Bool("is_http", false, `use it as HTTP server (default "false")`)
	addr   = flag.String("addr", ":8080", "HTTP serving address for dalog")
	rout   = flag.String("rout", "result.out", "result output")
	// flist []string
)

// init is called before main() to parse all the command line arguments
func init() {
	log.SetFlags(log.Lshortfile)
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("   %s -bug 30073 -conf config.json\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func readConf(cf string) (*Conf, error) {
	fcfg, err := ioutil.ReadFile(cf)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("ReadFile: %s", cf))
	}

	cfg := Conf{}
	if err := json.Unmarshal(fcfg, &cfg); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unmarshal: %s", cf))
	}

	return &cfg, nil
}

func outputToFile(fName string, fn func(w io.Writer) error) error {
	f, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("outputToFile: %s", fName))
	}
	defer f.Close()

	return fn(f)
}

func main() {
	var searchPath string

	if (*httpON == false && *iBug == "" && *iPath == "") || *conf == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := readConf(*conf)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	if *iBug != "" {
		searchPath = "/cores/bug_" + *iBug
	} else if *iPath != "" {
		searchPath = *iPath
	}

	if *httpON == false {
		if err := SearchPath(searchPath, cfg); err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}

		f, err := os.OpenFile(*flist, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf("outputToFile: %s", *flist)
			os.Exit(1)
		}
		defer f.Close()
		if err := cfg.FileList(f); err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}

		fOut, err := os.OpenFile(*rout, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf("outputToFile: %s", *rout)
			os.Exit(1)
		}
		defer fOut.Close()
		if err := cfg.Output(fOut); err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/bug/{id}", GetBugOutput).Methods("GET")
	log.Fatal(http.ListenAndServe(*addr, router))
	// http service
}
