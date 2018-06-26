package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	logpath = flag.String("logpath", "", "log path to analyze")
	conf    = flag.String("conf", "config.json", "config file in JSON")
	errOut  = flag.String("errout", "run.err", "runtime errors")
	retOut  = flag.String("retout", "ret.out", "result output")
	jsonCFG Config

	gDalog []Dalog
)

func main() {
	logErr, err := os.OpenFile(*errOut, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logErr.Close()

	logOut, err := os.OpenFile(*retOut, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logOut.Close()

	multi := io.MultiWriter(logErr, os.Stdout)

	//MyFile := log.New(multi, "PREFIX: ", log.Lshortfile)
	log.SetFlags(log.Lshortfile)
	log.SetOutput(multi)

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("   %s -logpath /cores/bug_30073 -conf config.json\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *logpath == "" || *conf == "" {
		flag.Usage()
		os.Exit(1)
	}

	fcfg, err := ioutil.ReadFile(*conf)
	if err != nil {
		log.Fatalf("Read %s Error %v\n", *conf, err)
	}

	if err := json.Unmarshal(fcfg, &jsonCFG); err != nil {
		log.Fatalf("Unmarshal %s failed %v\n", *conf, err)
	}

	if err := searchPath(*logpath); err != nil {
		log.Fatalf("searchPath failed %v\n", err)
	}

	log.Printf("gDalog len is %d\n", len(gDalog))

	fnTime := func(d1, d2 *Dalog) bool {
		return d1.TimeStamp < d2.TimeStamp
	}

	By(fnTime).Sort(gDalog)
	for _, dlog := range gDalog {
		if out := dlog.OutputErr(); out != "" {
			log.Println(out)
		}
	}
	//fmt.Println(flag.Args()[0])
	//if err := ReadGzipFile(flag.Args()[0]); err != nil {
	//	log.Fatal(err)
	//}
	//log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(os.Args[1]))))
}

func searchPath(root string) error {
	eachPath := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Error %v %s", err, path)
		}

		if info.IsDir() {
			return nil
		}

		// we are dealing w/ files
		idxLastSlash := strings.LastIndex(path, "/")

		if idxLastSlash == -1 {
			// any file under the root will not be out target anyway.
			return nil
		}

		// file path is path[:idxLastSlash+1]
		// file name is path[idxLastSlash+1:]
		p := path[:idxLastSlash+1]
		base := path[idxLastSlash+1:]
		for _, l := range jsonCFG.LogFiles {
			if strings.HasPrefix(base, l.Prefix) {
				if len(base) < 4 {
					log.Println(root+"/"+p+base, "type unknown")
				} else if base[len(base)-3:] == ".gz" {
					log.Printf("type: .gz %s\n", path)
					if err := ReadGzipFile(path); err != nil {
						log.Fatalf("ReadGzipFile failed %v", err)
					}
				} else if base[len(base)-4:] == ".log" {
					log.Println(root+"/"+p+base, "type plain")
				} else {
					log.Println(root+"/"+p+base, "type unknown")
				}
			}
		}
		//if strings.HasSuffix(path, "/da/data/var/log") {
		//	fmt.Println(path)
		//}

		return nil
	}
	return filepath.Walk(root, eachPath)
}

func scanFile(fname string, scanner *bufio.Scanner) error {
	const (
		LogEntryPrefix = "@cee"
	)

	fnChopStr := func(str string, pattern byte, n int) int {
		idx := 0
		l := len(str)
		for idx < l {
			if str[idx] == pattern {
				if n == 1 {
					return idx
				}
				n--
			}
			idx++
		}

		return -1
	}

	i := 0
	for scanner.Scan() {
		i++
		str := scanner.Text()
		if len(str) < len(LogEntryPrefix) {
			log.Printf("Invalid Prefix. %s line %d\n  %s\n", fname, i, str)
			continue
		}

		if str[:len(LogEntryPrefix)] != LogEntryPrefix {
			// need to handle
			//	ZOO_ERROR, ZOO_WARN, ZOO_INFO
			//  2017-07-26T20:32:30.308886+0000:4098(0x7ffff4437700):ZOO_WARN@zookeeper_interest@1773: Exceeded deadline by 3335ms
			if !strings.Contains(str, "ZOO_") {
				// handle incomplete logging
				log.Printf("Unknown Prefix. %s line %d\n  %s\n", fname, i, str)
				//log.Fatalf("Unknown Prefix. %s line %d %s\n", fname, i, str)
				//return fmt.Errorf("Unknown Prefix. %s line %d %s\n", fname, i, str)
				continue
			}

			timeIdx := fnChopStr(str, ':', 3)
			if timeIdx < 0 {
				log.Fatalf("fnChopStr timeIdx error %s line %d : 3 string %s", fname, i, str)
			}
			thID := fnChopStr(str, ':', 4)
			if thID < 0 {
				log.Fatalf("fnChopStr thID error %s line %d : 4 string %s", fname, i, str)
			}
			locIdx := fnChopStr(str, ':', 5)
			if locIdx < 0 {
				log.Fatalf("fnChopStr locID error %s line %d : 5 string %s", fname, i, str)
			}
			prioLoc := str[thID+1 : locIdx]
			prioIdx := strings.Index(prioLoc, "@")
			if prioIdx < 0 {
				log.Println("time:", str[:timeIdx], "TID", str[timeIdx+1:thID], "MSG", str[locIdx+1:])
				log.Fatalf("priority finding error %s line %d \n  prioLoc %v thID %d locIdx %d\n   string %s",
					fname, i, prioLoc, thID, locIdx, str)
			}

			prio := prioLoc[:prioIdx]
			loc := prioLoc[prioIdx+1:]
			newDalog := Dalog{
				TimeStamp:    str[:timeIdx],
				Msg:          str[locIdx+1:],
				ThreadID:     str[timeIdx+1 : thID],
				Priority:     prio,
				CodeLocation: loc,
			}
			gDalog = append(gDalog, newDalog)
			continue
		}
		// trim the prefix
		str = str[len(LogEntryPrefix):]

		res := Dalog{}
		if err := json.Unmarshal([]byte(strings.Replace(str, "\t", "    ", -1)), &res); err != nil {
			// need to handle imcompleted JSON fmt
			log.Printf("Unmarshal %v line:%d file:%s string: %s\n", err, i, fname, str)
			continue
		}
		//fmt.Println(res.String())
		gDalog = append(gDalog, res)
	}
	return nil
}

func ReadGzipFile(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Open %v %s", err, fname)
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		log.Fatalf("NewReader %v %s", err, fname)
	}
	defer gzf.Close()

	scanner := bufio.NewScanner(gzf)

	return scanFile(fname, scanner)
}
