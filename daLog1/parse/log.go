package parse

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
)

var (
	timedivider int64 = 1000000000

	dalogTimeFmt     string = "2006-01-02T15:04:05.000+0000"
	daheadTimeFmt    string = "2006-01-02T15:04:05"
	esxlogTimeFmt    string = "2006-01-02T15:04:05.000Z"
	esxSyslogTimeFmt string = "2006-01-02T15:04:05Z"
	testlogTimeFmt   string = "2006-01-02T15:04:05.000000+0000"
)

type DaOut interface {
	Name() string
	Start() int64
	End() int64
	ListIssue(*Request) (*bytes.Buffer, error)
}

type BaseFields struct {
	FileName  string
	StartTime int64
	EndTime   int64
}

type Out struct {
	DaOut
}

func (b *BaseFields) Name() string {
	return b.FileName
}

func (b *BaseFields) Start() int64 {
	return b.StartTime
}

func (b *BaseFields) End() int64 {
	return b.EndTime
}

func daNanosecond(t int64) string {
	return time.Unix(t/timedivider, t%timedivider).In(time.UTC).Format(dalogTimeFmt)
}

func daSecond(t int64) string {
	return time.Unix(t, 0).In(time.UTC).Format(dalogTimeFmt)
}

func getScanner(f *os.File) (*bufio.Scanner, error) {
	// REVISIT: filesize < 2
	buf := make([]byte, 2)
	_, err := f.Read(buf)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, errors.Wrap(err, fmt.Sprintf("getScanner"))
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("getScanner"))
	}

	var scanner *bufio.Scanner = bufio.NewScanner(f)

	if buf[0] == 31 && buf[1] == 139 {
		// it's a gz file
		gzf, err := gzip.NewReader(f)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("getScanner"))
		}
		defer gzf.Close()
		scanner = bufio.NewScanner(gzf)
	}

	return scanner, nil
}
