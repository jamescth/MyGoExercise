package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

/*
bug 31200
http://tracker.datrium.com/show_bug.cgi?id=31200
time skew.  check net_heartbeat_svr.log

bug 35068/d
*/
var (
	ErrFileTooSmall error = fmt.Errorf("File is too small")
)

// The following is the json fmt for da_alarms_show output
type DalogEntry struct {
	TimeStamp      string `json:"timestamp"`
	Msg            string `json:"msg"`
	CodeLocation   string `json:"codelocation"`
	ThreadName     string `json:"threadname:`
	ThreadID       string `json:"threadID"`
	Priority       string `json:"priority"`
	Host           string `json:"host"`
	Proc           string `json:"proc"`
	Pid            string `json:"pid"`
	ExceptionInfo  string `json:"exceptionInfo"`
	Severity       string `json:"severity"`       // dvaEvents
	EventID        string `json:"eventId"`        // dvaEvents
	EventType      string `json:"eventType"`      // dvaEvents
	ProcessName    string `json:"processName"`    // dvaEvents
	ComponentKind  string `json:"componentKind"`  // dvaEvents
	ComponentID    string `json:"componentId"`    // dvaEvents
	SrcVer         string `json:"srcVer"`         // dvaEvents
	DstVer         string `json:"dstVer"`         // dvaEvents
	OP             string `json:"op"`             // dvaEvents
	ExitCode       string `json:"exitCode"`       // dvaEvents
	FailedWrites   string `json:"failedWrites"`   // dvaEvents
	FailedReads    string `json:"failedReads"`    // dvaEvents
	FailedReserves string `json:"failedReserves"` // dvaEvents
	DevFsPath      string `json:"devFsPath"`      // dvaEvents
	HBAName        string `json:"hbaName"`        // dvaEvents
	Path           string `json:"path"`
	LineNum        int    `json:"lineNum"`
}

type Entry struct {
	Time  int64
	Start int64
	End   int64
}

type DalogEntries struct {
	BaseFields

	Entries []Entry
}

const (
	cee = "@cee"
)

func NewDalogEntriesDaOut(p string) (DaOut, error) {
	dalog, err := NewDalogEntries(p)
	var d DaOut = dalog
	return d, err

}

//func NewDalogEntries(p string) (DaOut, error) {
func NewDalogEntries(p string) (*DalogEntries, error) {
	var (
		startPos, endPos, startTime, endTime int64
		dalog                                DalogEntries
	)

	f, err := os.Open(p)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewDalogEntries"))
	}
	defer f.Close()

	scanner, err := getScanner(f)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, errors.Wrap(err, fmt.Sprintf("NewDalogEntries"))
	}

	dalog.FileName = p
	dalog.Entries = make([]Entry, 0, 100)

	for scanner.Scan() {
		var unixnano int64

		str := scanner.Text()
		endPos += int64(len(str))
		// log.Print(str)
		search := "@cee{\"timestamp\": \""
		if idx := strings.Index(str, search); idx != -1 {
			if len(str) > len(search)+len(dalogTimeFmt) {
				// t.Log(str[len(search) : len(search)+len(timeformat)])
				timestamp := str[len(search) : len(search)+len(dalogTimeFmt)]
				parseTime, err := time.Parse(dalogTimeFmt, timestamp)
				if err == nil {
					unixnano = parseTime.UnixNano()
					// return nil, errors.Wrap(err, fmt.Sprintf("NewDalogEntries %s", p))
				}

				// No error return if we cant parse timestamp
				// move on
			}
		}

		if startTime == 0 {
			startTime = unixnano
		}
		endTime = unixnano

		dalog.Entries = append(dalog.Entries, Entry{
			Time:  unixnano,
			Start: startPos,
			End:   endPos,
		})

		startPos = endPos + 1

	}
	dalog.StartTime = startTime
	dalog.EndTime = endTime
	return &dalog, nil
	//	var da DaOut
	//	return da, nil
}

//func (d *DalogEntries) ListIssue(startTime, endTime int64, c, s string, w io.Writer) error {
func (d *DalogEntries) ListIssue(req *Request) (*bytes.Buffer, error) {
	// how to list issue?
	// start, end?
	var (
		sidx, eidx int
	)

	startTime := req.Start
	endTime := req.End

	if startTime == 0 {
		startTime = d.StartTime
	}
	if endTime == 0 {
		endTime = d.EndTime
	}

	// calculate where and how many bytes to read
	var startIdx, endIdx int = -1, -1
	for i, e := range d.Entries {
		if startIdx < 0 {
			if startTime >= e.Time {
				startIdx = i
				endIdx = i
			}
			continue
		}

		if endTime < e.Time {
			break
		}
		endIdx = i
	}

	if startIdx < 0 {
		return nil, errors.Wrap(fmt.Errorf("Cant find timestamp"),
			fmt.Sprintf("DalogEntries.ListIssue"))
	}

	// open the file and read buf
	f, err := os.Open(d.Name())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("DalogEntries.ListIssue"))
	}
	defer f.Close()

	var buf bytes.Buffer
	scanner, err := getScanner(f)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("DalogEntries.ListIssue"))
	}

	for scanner.Scan() {
		str := scanner.Text()
		if sidx < startIdx {
			sidx++
			eidx++
			continue
		}

		if eidx > endIdx {
			break
		}
		// fmt.Println("James", i, str)
		// Unmarshal @cee to json fmt
		var da DalogEntry
		if len(str) <= len(cee) {
			buf.WriteString(fmt.Sprintf("ERROR Unmarshal: %s\n", str))
			continue
		}

		// json unmarshal doesn't accept tabs, trim it ahead of unmarshal
		str = strings.Replace(str, "\t", "", -1)

		if err := json.Unmarshal([]byte(str[len(cee):]), &da); err != nil {
			buf.WriteString(fmt.Sprintf("ERROR Unmarshal: %v\n%s\n", err, str))
			// fmt.Fprintf(w, "ERROR Unmarshal: %s\n", str)
			// don't report unmarshal error, just print out the content
			continue
		}
		buf.WriteString(fmt.Sprintf(" %s %8s", da.TimeStamp[:23], da.Priority))
		if da.ThreadName != "" {
			buf.WriteString(fmt.Sprintf(" [%s]", da.ThreadName))
		}
		if da.ProcessName != "" {
			buf.WriteString(fmt.Sprintf(" {%s}", da.ProcessName))
		}
		if da.EventType != "" {
			buf.WriteString(fmt.Sprintf(" {%s}", da.EventType))

			// print out detail info for ESXSSDStatsFailureEvent
			if da.EventType == "ESXSSDStatsFailureEvent" || da.EventType == "ESXHBAStatsFailureEvent" {
				buf.WriteString(fmt.Sprintf(" hbaName %s", da.HBAName))
				buf.WriteString(fmt.Sprintf(" devFsPath %s", da.DevFsPath))
				buf.WriteString(fmt.Sprintf(" failedReads %s", da.FailedReads))
				buf.WriteString(fmt.Sprintf(" failedWrites %s", da.FailedWrites))
				buf.WriteString(fmt.Sprintf(" failedReserves %s", da.FailedReserves))
			}
		}
		if da.ComponentID != "" {
			buf.WriteString(fmt.Sprintf(" {%s}", da.ComponentID))
		}
		if da.ExitCode != "" {
			buf.WriteString(fmt.Sprintf(" {%s}", da.ExitCode))
		}
		if da.Msg != "" {
			buf.WriteString(fmt.Sprintf(" %s", da.Msg))
		}
		if da.ExceptionInfo != "" {
			buf.WriteString(fmt.Sprintf(" %s", da.ExceptionInfo))
		}
		if da.SrcVer != "" {
			buf.WriteString(fmt.Sprintf(" %s -> %s", da.SrcVer, da.DstVer))
		}
		if da.OP != "" {
			buf.WriteString(fmt.Sprintf(" %s", da.OP))
		}

		// the following info print in a newline, and align it
		buf.WriteString(fmt.Sprintf("\n%33s", ""))
		if da.Host != "" {
			buf.WriteString(fmt.Sprintf(" {%s}", da.Host))
		}
		if da.CodeLocation != "" {
			buf.WriteString(fmt.Sprintf(" {%s}", da.CodeLocation))
		}
		if da.Path != "" {
			buf.WriteString(fmt.Sprintf(" %s", da.Path))
		}
		if da.LineNum != 0 {
			buf.WriteString(fmt.Sprintf(" %d", da.LineNum))
		}
		buf.WriteString(fmt.Sprintf("\n"))
		eidx++

		// if we directly write to Writer, perf will be much faster because
		// no GC involvement.
		// But how can we handle concurrency?
		/*
			fmt.Fprintf(w, " %s %8s", da.TimeStamp[:23], da.Priority)
			if da.ThreadName != "" {
				fmt.Fprintf(w, " [%s]", da.ThreadName)
			}
		*/
	}
	return &buf, nil
}
