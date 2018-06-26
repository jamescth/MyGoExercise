package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
)

type DVXSoftware struct {
	BaseFields

	// The following is the json fmt for da_alarms_show output
	LastRefreshTime     int64  `json:"lastRefreshTime"`
	LastSuccRefreshTime int64  `json:"lastSuccRefreshTime"`
	RemoteImageMsg      string `json:"remoteImageMsg"`
	RemoteStatus        int    `json:"remoteStatus"`
	SoftwareImages      []struct {
		BuildTime      string `json:"buildTime"`
		ReleaseNoteURL string `json:"releaseNoteUrl"`
		Status         string `json:"status"`
		Version        string `json:"version"`
	} `json:"softwareImages"`
	TaskHistory []struct {
		CompletionTime int64 `json:"completionTime"`

		// "DVX system upgraded"
		// "Agent preinstall check failed. See events for details"
		// "Missing agents n1415-b13n2.datrium.com."
		// "Software image downloaded"
		Description string `json:"description"`

		Error struct {
			Attributes []interface{} `json:"attributes"`
			ErrorCodes []interface{} `json:"errorCodes"`
		} `json:"error"`
		ID        string        `json:"id"`
		KeyValues []interface{} `json:"keyValues"`

		// Upgrade,
		// "Download"
		Kind string `json:"kind"`

		Progress  int   `json:"progress"`
		StartTime int64 `json:"startTime"`

		// "Success"
		// "Error"
		State                string `json:"state"`
		TaskParam            string `json:"taskParam"`
		TaskStateDetailsData struct {
			KeyValues []interface{} `json:"keyValues"`
		} `json:"taskStateDetailsData"`
	} `json:"taskHistory"`
}

func NewDVXSoftware(p string) (DaOut, error) {
	content, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewDVXSoftware"))
	}
	var d DVXSoftware
	if err := json.Unmarshal(content, &d); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewDVXSoftware %s", p))
	}

	d.FileName = p
	idx := len(d.TaskHistory)
	if idx != 0 {
		d.StartTime = d.TaskHistory[idx-1].CompletionTime
		d.EndTime = d.TaskHistory[0].CompletionTime
	}

	var da DaOut = &d
	return da, nil
}

func (d *DVXSoftware) ListIssue(req *Request) (*bytes.Buffer, error) {
	c := req.Component
	start := req.Start
	end := req.End
	var b bytes.Buffer
	var all bool = false

	if d.Start() == 0 {
		b.WriteString(fmt.Sprintf("DVX Software has no taskhistory\n"))
		return &b, nil
	}

	if c == "" {
		all = true
	}
	if start == 0 {
		start = d.Start()
	}
	if end == 0 {
		end = d.End()
	}

	// b.WriteString(fmt.Sprintf("file: %s, time: %s\n", d.Name(), daNanosecond(d.End())))
	b.WriteString(fmt.Sprintf("%10s %40s %3s %10s %30s %30s\n", "Kind", "Des", "%", "State", "Start Time", "Completion time"))
	b.WriteString(fmt.Sprintf("========== ======================================== === ========== ============================== ==============================\n"))

	/*
		fmt.Fprintf(w, "file: %s, time: %s\n", d.Name(), daNanosecond(d.End()))
		fmt.Fprintf(w, "%10s %40s %3s %10s %30s %30s\n", "Kind", "Des", "%", "State", "Start Time", "Completion time")
		fmt.Fprintf(w, "========== ======================================== === ========== ============================== ==============================\n")
	*/
	for _, t := range d.TaskHistory {
		if !all && t.Kind != c {
			continue
		}

		if req.Severity != "" && req.Severity != t.State {
			continue
		}

		if start <= t.CompletionTime && end >= t.CompletionTime {
			sTime := time.Unix(t.StartTime/timedivider, t.StartTime%timedivider)
			eTime := time.Unix(t.CompletionTime/timedivider, t.CompletionTime%timedivider)
			//fmt.Fprintf(w, "%10s %40s %3d %10s %30s %30s\n",
			b.WriteString(fmt.Sprintf("%10s %40s %3d %10s %30s %30s\n",
				t.Kind,
				t.Description,
				t.Progress,
				t.State,
				sTime.In(time.UTC).Format(dalogTimeFmt),
				eTime.In(time.UTC).Format(dalogTimeFmt)))
		}
	}
	return &b, nil
}
