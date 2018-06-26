package parse

import (
	"bytes"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type TestLog struct {
	BaseFields

	Entries []Entry
}

func NewTestLogDa(p string) (DaOut, error) {
	testlog, err := NewTestLog(p)
	var d DaOut = testlog
	return d, err

}

func NewTestLog(p string) (*TestLog, error) {
	var testlog TestLog
	testlog.FileName = p
	startTime, endTime, err := subNewObj(&testlog.Entries, p, testlogTimeFmt)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewESXSysLog"))
	}
	testlog.StartTime = startTime
	testlog.EndTime = endTime
	return &testlog, nil
}

func (d *TestLog) ListIssue(req *Request) (*bytes.Buffer, error) {
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
			fmt.Sprintf("TestLog.ListIssue"))
	}

	// open the file and read buf
	f, err := os.Open(d.Name())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("TestLog.ListIssue"))
	}
	defer f.Close()

	var buf bytes.Buffer
	scanner, err := getScanner(f)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("TestLog.ListIssue"))
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

		buf.WriteString(str + "\n")
	}

	return &buf, nil
}
