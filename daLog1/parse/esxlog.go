package parse

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
)

type ESXLog struct {
	BaseFields

	Entries []Entry
}

func NewESXLogDa(p string) (DaOut, error) {
	esxlog, err := NewESXLog(p)
	var d DaOut = esxlog
	return d, err

}

func NewESXLog(p string) (*ESXLog, error) {
	var esxlog ESXLog
	esxlog.FileName = p
	startTime, endTime, err := subNewObj(&esxlog.Entries, p, esxlogTimeFmt)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewESXLog"))
	}
	esxlog.StartTime = startTime
	esxlog.EndTime = endTime
	return &esxlog, nil
}

/*
 subNewObj() receives a pointer to an Entry array for appending entry, and a file name
 string to work on. Also, timestamp format for parsing at the begin of each line.

 after scan the file, it retures the earliest and latest timestamp back to the caller
 while updating its Entry array.
*/
func subNewObj(entries *[]Entry, p, timeFmt string) (startTime, endTime int64, err error) {
	var (
		startPos, endPos int64
	)

	f, err := os.Open(p)
	if err != nil {
		return 0, 0, errors.Wrap(err, fmt.Sprintf("subNewObj"))
	}
	defer f.Close()

	scanner, err := getScanner(f)
	if err != nil {
		return 0, 0, errors.Wrap(err, fmt.Sprintf("subNewObj"))
	}

	for scanner.Scan() {
		var unixnano int64

		str := scanner.Text()
		endPos += int64(len(str))

		timeFmtLen := len(timeFmt)
		if len(str) > timeFmtLen {
			parseTime, err := time.Parse(timeFmt, str[:timeFmtLen])
			if err == nil {
				unixnano = parseTime.UnixNano()
			}

			// No error return if da_head* line w/out timestamp
			// move on
		}

		if startTime == 0 && unixnano != 0 {
			startTime = unixnano
		}

		if unixnano != 0 {
			endTime = unixnano
		}

		*entries = append(*entries, Entry{
			Time:  unixnano,
			Start: startPos,
			End:   endPos,
		})

		startPos = endPos + 1
	}

	return startTime, endTime, nil

}
func (d *ESXLog) ListIssue(req *Request) (*bytes.Buffer, error) {
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
			fmt.Sprintf("ESXLog.ListIssue"))
	}

	// open the file and read buf
	f, err := os.Open(d.Name())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("ESXLog.ListIssue"))
	}
	defer f.Close()

	var buf bytes.Buffer
	scanner, err := getScanner(f)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("ESXLog.ListIssue"))
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
