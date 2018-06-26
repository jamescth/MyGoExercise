package parse

import (
	"bytes"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type DAHead struct {
	BaseFields

	Entries []Entry
}

func NewDaHeadDa(p string) (DaOut, error) {
	dahead, err := NewDaHead(p)
	var d DaOut = dahead
	return d, err

}

func NewDaHead(p string) (*DAHead, error) {
	var dahead DAHead
	dahead.FileName = p
	startTime, endTime, err := subNewObj(&dahead.Entries, p, daheadTimeFmt)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewDAHead"))
	}
	dahead.StartTime = startTime
	dahead.EndTime = endTime
	return &dahead, nil

	/*
		var (
			startTime, endTime, startPos, endPos int64
			dahead                               DAHead
		)

		dahead.FileName = p

		f, err := os.Open(p)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("NewDAHead"))
		}
		defer f.Close()

		scanner, err := getScanner(f)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("NewDAHead"))
		}

		// dahead.FileName = p
		// dahead.Entries = make([]Entry, 0, 100)

		// timestamp in da_head*
		timeFmtLen := len("2018-03-05T19:08:18")
		for scanner.Scan() {
			var unixnano int64

			str := scanner.Text()
			endPos += int64(len(str))

			if len(str) > timeFmtLen {
				parseTime, err := time.Parse(daheadTimeFmt, str[:timeFmtLen])
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

			dahead.Entries = append(dahead.Entries, Entry{
				Time:  unixnano,
				Start: startPos,
				End:   endPos,
			})

			startPos = endPos + 1
		}

		dahead.StartTime = startTime
		dahead.EndTime = endTime
		return &dahead, nil
	*/
}

func (d *DAHead) ListIssue(req *Request) (*bytes.Buffer, error) {
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

		buf.WriteString(str + "\n")
	}

	return &buf, nil
}
