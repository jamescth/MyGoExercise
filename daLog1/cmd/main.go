package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jamescth/MyGoExercise/daLog1/parse"
	"github.com/pkg/errors"
)

var first []dvxlog = []dvxlog{
	{"da_alarms_show.", parse.NewAlarms, 0, 0, "", "", os.Stdout},
	{"da_dvx_software_show.", parse.NewDVXSoftware, 0, 0, "", "", os.Stdout},
	{"da_hosts_show.", parse.NewHosts, 0, 0, "", "", os.Stdout},
	{"da_nodes_show.", parse.NewNodes, 0, 0, "", "", os.Stdout},
}

var second []dvxlog = []dvxlog{
	{"test.log", parse.NewTestLogDa, 0, 0, "", "", os.Stdout},

	{"hostd.log", parse.NewESXLogDa, 0, 0, "", "", os.Stdout},
	{"vmkernel.", parse.NewESXLogDa, 0, 0, "", "", os.Stdout},
	{"vmkwarning.log", parse.NewESXLogDa, 0, 0, "", "", os.Stdout},
	{"vmksummary.log", parse.NewESXSysLogDa, 0, 0, "", "", os.Stdout},
	{"syslog.log", parse.NewESXSysLogDa, 0, 0, "", "", os.Stdout},
	{"dvaEvents.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"da_head_rp.", parse.NewDaHeadDa, 0, 0, "", "", os.Stdout},
	{"da_head_init.", parse.NewDaHeadDa, 0, 0, "", "", os.Stdout},
	{"da_head.", parse.NewDaHeadDa, 0, 0, "", "", os.Stdout},
	{"da_head0.", parse.NewDaHeadDa, 0, 0, "", "", os.Stdout},
	{"da_setup.", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"da_setup_dasys.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"da_setup_directio_dasys.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"handle_mount_pybridge.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"da_pre_install.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"da_post_install.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"host_preinstall_hook.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"dvx_preinstall_hook.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"esx_platmgr.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"esx_platmgr_dasys.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"esx_platmgr_copy.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"hamgr.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"hamgr_cli.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"upgrade_mgr.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"upgrade_util.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"procmgr.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"procmgr_cleanup.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	// controller
	{"platmgr.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"shutdown.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	{"fe.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
	// KVM
	{"SetupLinuxHost.log", parse.NewDalogEntriesDaOut, 0, 0, "", "", os.Stdout},
}
var (
	version    string
	minversion string

	lpath  = flag.String("path", "", "log path to analyze")
	iBug   = flag.String("bug", "", "bugxxxxx to analyze")
	output = flag.String("out", "result.out", "output file")

	//chans []chan *bytes.Buffer
	firstlogs  []logFile
	secondlogs []logFile
	omitLogs   []string
)

func init() {
	log.SetFlags(log.Lshortfile)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Version %s\nMinversion %s\n\n", version, minversion)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "   %s -bug 30073 -conf config.json\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func checkPrefix(p string, das []dvxlog, logs *[]logFile, omit *[]string) (err error) {
	var da parse.DaOut

	for i, d := range das {
		if !strings.HasPrefix(path.Base(p), d.prefix) {
			continue
		}

		// init the struct
		da, err = d.fn(p)
		if err != nil {
			if err == io.EOF {
				fmt.Fprintf(os.Stderr, "empty file: %s\n", p)
				continue
			}
			fmt.Fprintf(os.Stderr, "error file: %s err:%+v\n", p, err)
			continue
			// return errors.Wrap(err, fmt.Sprintf("checkPrefix %s", p))
		}

		ch := make(chan *bytes.Buffer, 1)
		*logs = append(*logs, logFile{da, ch})

		go func(da parse.DaOut, idx int) {
			defer close(ch)

			buf, err := da.ListIssue(&parse.Request{das[idx].st, das[idx].et, das[idx].c, das[idx].s})
			if err != nil {
				//return errors.Wrap(err, fmt.Sprintf("fnWorkLogs"))
				//log.Fatalf("%+v", err)
				fmt.Fprintf(os.Stderr, "%+v", err)

			}
			ch <- buf
		}(da, i)
	}

	if da != nil {
		return nil
	}
	if omit != nil {
		*omit = append(*omit, p)
	}
	return nil
}

// fnWorkLogs is a function for filepath.Walk() to run each time a path is discovered.
func fnWorkLogs(p string, info os.FileInfo, err error) error {
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("fnWorkLogs"))
	}

	if info.IsDir() {
		return nil
	}

	// TODO: can we checkPrefix once?
	if err := checkPrefix(p, first, &firstlogs, nil); err != nil {
		return errors.Wrap(err, fmt.Sprintf("fnWorkLogs"))
	}
	if err := checkPrefix(p, second, &secondlogs, &omitLogs); err != nil {
		return errors.Wrap(err, fmt.Sprintf("fnWorkLogs"))
	}

	return nil
}

func main() {
	searchPath := ""
	if *iBug != "" {
		searchPath = "/cores/bug_" + *iBug
	} else if *lpath != "" {
		searchPath = *lpath
	}

	if err := filepath.Walk(searchPath, fnWorkLogs); err != nil {
		log.Printf("%+v\n", err)
	}

	f, err := os.OpenFile(*output, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outputToFile: %s %v", *output, err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := f.WriteString("search path: " + searchPath + "\n"); err != nil {
		fmt.Fprintf(os.Stderr, "WriteString searchPath %v %s", *output, err)
		os.Exit(1)
	}

	// TODO: we should consolidate 1st and 2nd loops
	for _, l := range firstlogs {
		buf := <-l.ch
		if _, err := f.WriteString("dafile: " + l.f.Name() + "\n"); err != nil {
			fmt.Fprintf(os.Stderr, "WriteString finename %s %v", *output, err)
			os.Exit(1)
		}
		if _, err := f.Write(buf.Bytes()); err != nil {
			fmt.Fprintf(os.Stderr, "WriteString content %s %v", *output, err)
			os.Exit(1)
		}
		f.Write([]byte("\n"))
	}

	for _, l := range secondlogs {
		buf := <-l.ch
		if _, err := f.WriteString("dafile: " + l.f.Name() + "\n"); err != nil {
			fmt.Fprintf(os.Stderr, "WriteString finename %s %v", *output, err)
			os.Exit(1)
		}
		if _, err := f.Write(buf.Bytes()); err != nil {
			fmt.Fprintf(os.Stderr, "WriteString content %s %v", *output, err)
			os.Exit(1)
		}
		f.Write([]byte("\n"))
	}

	for _, s := range omitLogs {
		if _, err := f.WriteString("not-yet: " + s + "\n"); err != nil {
			fmt.Fprintf(os.Stderr, "WriteString omitlogs %s %v", *output, err)
			os.Exit(1)
		}

	}
	return
}
