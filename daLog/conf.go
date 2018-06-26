package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	// dalog "./dalog_entry"
	dalog "github.com/jamescth/MyGoExercise/daLog/dalog_entry"
	"github.com/pkg/errors"
)

/*
	esx:
		{path}/var/run/log/
			auth.log	=> ssh
			dacli.log	=> dlog
			da_head.log	=> plain/no timeStamp
			datrium_bootstrap_stage*.log
			datrium_config_agent.log
			esxupdate.log
			hostd.log
			shell.log
			syslog.log
			vmkernel.log
			vmkernel.*.gz
			vmkeventd.log
			vobd.log => esx.problem
			vpxa.log

	da:
		host
			{path}/var/log/
				the following logs link to /asupdata-internal/..../var/log/{file}
				da_setup.log
				esx_platmgr.log
				procmgr_cleanup.log
				syslog.log
				vmkernel.log
			{path}/var/log/datrium
				files here may point to /asupdata-internal/......

			{path}/var/log/datrium/traces/default

		ctrler:
			{path}/da/data/var/log
			{path}/da/data/var/log/traces/default



	{
		"mountPoint":"",
		"logFiles":[
			{
				"prefix":"upgrade_mgr.log",
				"type":"dalog",
				"patterns":["checkpoint

		]

	}
*/

type DaFile struct {
	File   string
	LinkTo string
}

func NewDaFile(f string) DaFile {
	return DaFile{File: f}
}

func NewDaFileWithLink(f, l string) DaFile {
	return DaFile{File: f, LinkTo: l}
}

func (da *DaFile) scanFile(w io.Writer, patterns []string) error {
	const (
		LogEntryPrefix = "@cee"
	)

	fileOpened := ""
	if da.LinkTo != "" {
		fileOpened = da.LinkTo
	} else {
		fileOpened = da.File
	}

	f, err := os.Open(fileOpened)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("scanFile open %s", fileOpened))
	}
	defer f.Close()

	// REVISIT: filesize < 2
	buf := make([]byte, 2)
	_, err = f.Read(buf)
	if err != nil {
		if err == io.EOF {
			return err
		}
		return errors.Wrap(err, fmt.Sprintf("scanFile read %s", fileOpened))
	}
	if _, err := f.Seek(0, 0); err != nil {
		return errors.Wrap(err, fmt.Sprintf("scanFile seek %s", fileOpened))
	}

	// check if it's gz
	var (
		scanner *bufio.Scanner
	)

	if buf[0] == 31 && buf[1] == 139 {
		// it's a gz file
		gzf, err := gzip.NewReader(f)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("gz NewReader %s", fileOpened))
		}
		defer gzf.Close()

		scanner = bufio.NewScanner(gzf)
	} else {
		// plain file
		scanner = bufio.NewScanner(f)
	}

	base := filepath.Base(fileOpened)
	if len(base) > 10 {
		base = base[:10]
	}
	i := 0
	for scanner.Scan() {
		i++
		str := scanner.Text()
		//if len(str) < len(LogEntryPrefix) {
		//	w.Write([]byte(fmt.Sprintf("Invalid Prefix. line %d\n  %s\n", i, str)))
		//	continue
		//}

		// @cee fmt
		if len(str) > len(LogEntryPrefix) && str[:len(LogEntryPrefix)] == LogEntryPrefix {
			// trim the prefix
			str = str[len(LogEntryPrefix):]

			res := dalog.DalogEntry{}
			if err := json.Unmarshal([]byte(strings.Replace(str, "\t", "    ", -1)), &res); err != nil {
				// need to handle imcompleted JSON fmt
				w.Write([]byte(fmt.Sprintf("Unmarshal %v line:%d string: %s\n", err, i, str)))
				continue
			}

			wFunc := func(w io.Writer, base string, res dalog.DalogEntry) {
				w.Write([]byte(fmt.Sprintf("        %10s/cee: %s [%s]", base, res.TimeStamp, res.Priority)))
				if res.ThreadName != "" {
					w.Write([]byte(fmt.Sprintf(" [%s]", res.ThreadName)))
				}
				if res.ProcessName != "" {
					w.Write([]byte(fmt.Sprintf(" {%s}", res.ProcessName)))
				}
				if res.EventType != "" {
					w.Write([]byte(fmt.Sprintf(" {%s}", res.EventType)))
				}
				if res.ComponentID != "" {
					w.Write([]byte(fmt.Sprintf(" {%s}", res.ComponentID)))
				}
				if res.ExitCode != "" {
					w.Write([]byte(fmt.Sprintf(" {%s}", res.ExitCode)))
				}
				if res.Msg != "" {
					w.Write([]byte(fmt.Sprintf(" %s", res.Msg)))
				}
				if res.ExceptionInfo != "" {
					w.Write([]byte(fmt.Sprintf(" %s", res.ExceptionInfo)))
				}
				if res.SrcVer != "" {
					w.Write([]byte(fmt.Sprintf(" %s -> %s", res.SrcVer, res.DstVer)))
				}
				if res.OP != "" {
					w.Write([]byte(fmt.Sprintf(" %s", res.OP)))
				}

				w.Write([]byte(fmt.Sprintf("\n")))
			}

			// Print Exceptions
			if res.ExceptionInfo != "" {
				w.Write([]byte(fmt.Sprintf("        %10s/cee: %s [%s] %s %s\n", base, res.TimeStamp, res.Priority, res.Msg, res.ExceptionInfo)))
				continue
			}

			if res.Priority == "DA_LOG_ERR" ||
				res.Priority == "ERROR" || // procmgr.log
				res.Priority == "DA_LOG_WARNING" ||
				res.Priority == "WARNING" || // procmgr.log
				res.Priority == "DA_LOG_EMERG" ||
				res.Severity == "EMERGENCY" ||
				res.Severity == "WARNING" ||
				res.Severity == "ERROR" {
				// DA_LOG_ERR DA_LOG_WARNING WARNING

				//if res.EventType != "" {
				//w.Write([]byte(fmt.Sprintf("        %10s/cee: %s [%s] %s [%s] %s\n", base, res.TimeStamp, res.Pid, res.Priority, res.EventType, res.Msg)))
				//continue
				//}
				// else
				//w.Write([]byte(fmt.Sprintf("        %10s/cee: %s [%s] %s %s\n", base, res.TimeStamp, res.Pid, res.Priority, res.Msg)))
				wFunc(w, base, res)
				continue
			}

			// let's print all content of non-cee fmt first
			if patterns == nil {
				wFunc(w, base, res)
				/*
					w.Write([]byte(fmt.Sprintf("        %10s/cee: %s [%s]", base, res.TimeStamp, res.Priority)))
					if res.ThreadName != "" {
						w.Write([]byte(fmt.Sprintf(" [%s]", res.ThreadName)))
					}
					if res.ProcessName != "" {
						w.Write([]byte(fmt.Sprintf(" {%s}", res.ProcessName)))
					}
					if res.EventType != "" {
						w.Write([]byte(fmt.Sprintf(" {%s}", res.EventType)))
					}
					if res.ComponentID != "" {
						w.Write([]byte(fmt.Sprintf(" {%s}", res.ComponentID)))
					}
					if res.ExitCode != "" {
						w.Write([]byte(fmt.Sprintf(" {%s}", res.ExitCode)))
					}
					if res.Msg != "" {
						w.Write([]byte(fmt.Sprintf(" %s", res.Msg)))
					}
					if res.SrcVer != "" {
						w.Write([]byte(fmt.Sprintf(" %s -> %s", res.SrcVer, res.DstVer)))
					}
					if res.OP != "" {
						w.Write([]byte(fmt.Sprintf(" %s", res.OP)))
					}

					w.Write([]byte(fmt.Sprintf("\n")))
				*/
				continue
			}
			for _, p := range patterns {
				if strings.Contains(str, p) {
					w.Write([]byte(fmt.Sprintf("        %10s/cee: %s [%s] %s %s\n", base, res.TimeStamp, res.Pid, res.Priority, res.Msg)))
				}
			}
			// REVISIT:
			// what to do w/ logs
			//da.Dalogs = append(da.Dalogs, res)
			continue
		}

		// let's print all content of non-cee fmt first
		if patterns == nil {
			w.Write([]byte(fmt.Sprintf("    %10s/non-cee: %s\n", base, str)))
			continue
		}

		for _, p := range patterns {
			if strings.Contains(str, p) {
				w.Write([]byte(fmt.Sprintf("    %10s/non-cee: %s\n", base, str)))
			}
		}

		/*
			newDalog := dalog.DalogEntry{
				TimeStamp:    str[:timeIdx],
				Msg:          str[locIdx+1:],
				ThreadID:     str[timeIdx+1 : thID],
				Priority:     prio,
				CodeLocation: loc,
			}
		*/

		// REVISIT:
		// what to do w/ logs
		//da.Dalogs = append(da.Dalogs, res)
		//gDalog = append(gDalog, newDalog)
		// continue
	}

	return nil

}

type LogFile struct {
	// Prefix of the scanned log files.
	// for Example, upgrade_mgr.log includes upgrade_mgr.log.<date>.<num>.gz
	Prefix string `json:"prefix"`

	// Log files format. gz, core, zcore, plain
	// For core and zcore, the current version will not scan the content.
	//   The tool will trigger gdb/dagdb to get the stacktrace in the future version.
	// For gz and plain, it has no impact for the current version.
	//   The tool reads the first 2 bytes to determine if the file is a gz format
	Type string `json:"type"`

	// Patterns define the keywords to be searched in the log files.
	// if Patters is not defined in the JSON file, all content will be included.
	Patterns []string `json:"patterns"`

	DaFiles []DaFile
	DaLogs  []dalog.DalogEntry
}

type Conf struct {
	MountPoint string    `json:"mountPoint"`
	LogFiles   []LogFile `json:"logFiles"`
}

func (c *Conf) FileList(w io.Writer) error {
	for _, conf := range c.LogFiles {
		if _, err := w.Write([]byte(conf.Prefix + "\n")); err != nil {
			return errors.Wrap(err, fmt.Sprintf("File: "))
		}

		for _, f := range conf.DaFiles {
			if _, err := w.Write([]byte("  " + f.File + "\n")); err != nil {
				return errors.Wrap(err, fmt.Sprintf("File: "))
			}

			if f.LinkTo != "" {
				if _, err := w.Write([]byte("    " + f.LinkTo + "\n")); err != nil {
					return errors.Wrap(err, fmt.Sprintf("LinkTo: "))
				}

			}
		}
	}

	return nil
}

func (c *Conf) Output(w io.Writer) error {
	for _, lgf := range c.LogFiles {
		if _, err := w.Write([]byte("Prefix: " + lgf.Prefix + "\n")); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Write Prefix: "))
		}

		for _, f := range lgf.DaFiles {
			if _, err := w.Write([]byte("Sub:" + f.File + "\n")); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Write Sub: "))
			}

			if f.LinkTo != "" {
				if _, err := w.Write([]byte("    " + f.LinkTo + "\n")); err != nil {
					return errors.Wrap(err, fmt.Sprintf("Write LinkTo: "))
				}

			}

			// No file scan for core files
			if lgf.Type == "core" || lgf.Type == "zcore" {
				continue
			}

			if err := f.scanFile(w, lgf.Patterns); err != nil {
				if err == io.EOF {
					if _, err := w.Write([]byte("File size 0: " + "\n")); err != nil {
						return errors.Wrap(err, fmt.Sprintf("Write file size 0"))
					}
					continue
				}
				return errors.Wrap(err, fmt.Sprintf("scanFile: "))
			}
		}
		w.Write([]byte(fmt.Sprintf("***************************************************************************************\n")))
	}

	return nil

}
