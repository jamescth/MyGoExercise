package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// use the following to auto gen the struct from da_alarms_show output
// https://mholt.github.io/json-to-go/

type Alarms struct {
	BaseFields

	// The following is the json fmt for da_alarms_show output
	Alarms []struct {
		// UpgradeSuccess, UpgradeStart
		AlarmType string `json:"alarmType"`
		Event     struct {
			ComponentID string `json:"componentId"`
			// CONTROLLER, DATASTORE, UpgradeMgr
			ComponentKind string `json:"componentKind"`
			ComponentName string `json:"componentName"`
			EventID       string `json:"eventId"`
			EventType     string `json:"eventType"`
			KeyValues     []struct {
				Key string `json:"key"`
				Val struct {
					EnumVal struct {
					} `json:"enumVal"`
					MsgVal struct {
						KeyVal []interface{} `json:"keyVal"`
					} `json:"msgVal"`
					RepeatedStringVal []interface{} `json:"repeatedStringVal"`
					StringVal         string        `json:"stringVal"`
					Type              int           `json:"type"`
				} `json:"val"`
			} `json:"keyValues"`
			Severity    string `json:"severity"`
			TriggerTime int64  `json:"triggerTime"`
		} `json:"event"`
		TriggerTime int64 `json:"triggerTime"`
	} `json:"alarms"`
}

func NewAlarms(p string) (DaOut, error) {
	content, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewAlarms"))
	}
	var a Alarms
	if err := json.Unmarshal(content, &a); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewAlarms %s", p))
	}

	a.FileName = p
	idx := len(a.Alarms)
	if idx == 0 {
		return nil, io.EOF
	}

	a.StartTime = a.Alarms[idx-1].TriggerTime
	a.EndTime = a.Alarms[0].TriggerTime

	var da DaOut = &a
	return da, nil
}

func (a *Alarms) ListIssue(req *Request) (*bytes.Buffer, error) {
	c := req.Component
	start := req.Start
	end := req.End
	var buf bytes.Buffer

	var (
		all bool = false
	)

	if c == "" {
		all = true
	}
	if start == 0 {
		start = a.Start()
	}
	if end == 0 {
		end = a.End()
	}

	//buf.WriteString(fmt.Sprintf("file: %s, time: %s\n", a.Name(), daNanosecond(a.End())))
	buf.WriteString(fmt.Sprintf("%10s %10s %40s %30s %30s\n", "Component", "Severity", "Event Type", "Component Name", "TriggerTime"))
	buf.WriteString(fmt.Sprintf("========== ========== ======================================== ============================== ==============================\n"))
	/*
		fmt.Fprintf(w, "file: %s, time: %s\n", a.Name(), daNanosecond(a.End()))
		fmt.Fprintf(w, "%10s %10s %40s %30s\n", "Component", "Severity", "Event Type", "TriggerTime")
		fmt.Fprintf(w, "========== ========== ======================================== ==============================\n")
	*/
	for _, i := range a.Alarms {
		e := i.Event
		if !all && e.ComponentKind != c {
			continue
		}

		if req.Severity != "" && e.Severity != req.Severity {
			continue
		}

		if start <= e.TriggerTime && end >= e.TriggerTime {
			//fmt.Fprintf(w, "%10s %10s %40s %30s\n",
			buf.WriteString(fmt.Sprintf("%10s %10s %40s %30s %30s\n",
				e.ComponentKind,
				e.Severity,
				e.EventType,
				e.ComponentName,
				daNanosecond(e.TriggerTime)))
		}
	}
	return &buf, nil
}
