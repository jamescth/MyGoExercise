package dalog_entry

import (
	"fmt"
	"sort"
)

/*
bug 31200
http://tracker.datrium.com/show_bug.cgi?id=31200
time skew.  check net_heartbeat_svr.log
*/
type DalogEntry struct {
	TimeStamp     string `json:"timestamp"`
	Msg           string `json:"msg"`
	CodeLocation  string `json:"codelocation"`
	ThreadName    string `json:"threadname:`
	ThreadID      string `json:"threadID"`
	Priority      string `json:"priority"`
	Host          string `json:"host"`
	Proc          string `json:"proc"`
	Pid           string `json:"pid"`
	ExceptionInfo string `json:"exceptionInfo"`
	Severity      string `json:"severity"`      // dvaEvents
	EventID       string `json:"eventId"`       // dvaEvents
	EventType     string `json:"eventType"`     // dvaEvents
	ProcessName   string `json:"processName"`   // dvaEvents
	ComponentKind string `json:"componentKind"` // dvaEvents
	ComponentID   string `json:"componentId"`   // dvaEvents
	SrcVer        string `json:"srcVer"`        // dvaEvents
	DstVer        string `json:"dstVer"`        // dvaEvents
	OP            string `json:"op"`            // dvaEvents
	ExitCode      string `json:"exitCode"`      // dvaEvents
	Path          string `json:"path"`
	LineNum       int    `json:"lineNum"`
}

func (d *DalogEntry) String() string {
	return fmt.Sprintf("%s %s %s %s %d %s",
		d.TimeStamp, d.Host, d.Priority, d.Msg, d.LineNum, d.Path)
}

type By func(d1, d2 *DalogEntry) bool

func (by By) Sort(lgs []DalogEntry) {
	da := &dalogSorter{
		dalogs: lgs,
		by:     by, // The sort method's receiver is the func (closure) that defines the sort order
	}
	sort.Sort(da)
}

type dalogSorter struct {
	dalogs []DalogEntry
	by     func(d1, d2 *DalogEntry) bool
}

func (d *dalogSorter) Len() int {
	return len(d.dalogs)
}

func (d *dalogSorter) Swap(i, j int) {
	d.dalogs[i], d.dalogs[j] = d.dalogs[j], d.dalogs[i]
}

func (d *dalogSorter) Less(i, j int) bool {
	return d.by(&d.dalogs[i], &d.dalogs[j])
}
