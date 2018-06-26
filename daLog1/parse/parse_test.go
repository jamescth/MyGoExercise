package parse

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

var (
	// /cores/bug_31759/a/dva-controller-0050cc1d0ea4/commands
	// b9-mgmt /net/fs-02/tank1/autosupport/internal/2017/09/2017_09_30/30304d35-1612-2200-95a5-95a980fce604/dva-controller-0050cc1d0ea4/daily/2017-09-30T04_00_07/commands
	//
	testpath string = "../tests/commands"
)

func TestMain(t *testing.T) {
	if err := filepath.Walk(testpath, fnWorkLogs); err != nil {
		t.Errorf("%+v\n", err)
	}
}

func fnWorkLogs(p string, info os.FileInfo, err error) error {
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("fnWorkLogs"))
	}

	if info.IsDir() {
		return nil
	}

	das := []struct {
		prefix string
		fn     func(string) (DaOut, error)
		st     int64  // start time
		et     int64  // end time
		c      string // component
		s      string // serverity
		w      io.Writer
	}{
		{"da_alarms_show.", NewAlarms, 0, 0, "", "", os.Stdout},
		{"da_dvx_software_show.", NewDVXSoftware, 0, 0, "", "", os.Stdout},
		{"da_hosts_show.", NewHosts, 0, 0, "", "", os.Stdout},
		{"da_nodes_show.", NewNodes, 0, 0, "", "", os.Stdout},
	}

	var da DaOut
	lenArray := len(das)

	checkChans := false
	chans := make([]chan *bytes.Buffer, lenArray)

	// var wg sync.WaitGroup
	for i, d := range das {
		if strings.HasPrefix(path.Base(p), d.prefix) {

			// init the struct
			da, err = d.fn(p)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("fnWorkLogs"))
			}

			//ch := make(chan *bytes.Buffer, 1)
			if !checkChans {
				checkChans = true
			}

			chans[i] = make(chan *bytes.Buffer, 1)

			go func(da DaOut, idx int) {
				defer close(chans[idx])

				buf, err := da.ListIssue(&Request{das[idx].st, das[idx].et, das[idx].c, das[idx].s})
				if err != nil {
					//return errors.Wrap(err, fmt.Sprintf("fnWorkLogs"))
					log.Fatalf("%+v", err)
				}
				chans[idx] <- buf
			}(da, i)
		}
	}

	if checkChans {
		for i, _ := range chans {
			if chans[i] == nil {
				continue
			}
			buf := <-chans[i]
			log.Print(string(buf.Bytes()))
		}
	}

	/*
		content, err := ioutil.ReadFile(p)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("fnWorkLogs"))
		}

		var first interface{}
		if err := json.Unmarshal(content, &first); err != nil {
			// log.Printf("Unmarshal. It's not JSON fmt %s", path)
			return nil
			//return errors.Wrap(err, fmt.Sprintf("Unmarshal. It's not JSON fmt %s", path))
		}

		dec := json.NewDecoder(bytes.NewReader(content))
		switch first.(type) {
		// decode array of stats in JSON fmt
		case []interface{}:
			var v []map[string]interface{}
			if err := dec.Decode(&v); err != nil {
				// t.Errorf("Decode %s: %v\n", file.Name(), err)
				return errors.Wrap(err, fmt.Sprintf("Unmarshal. Decode %s\n", p))
			}
			// log.Printf("%v\n", v)
			// log.Printf("[]map[striing]interface{}: %s\n", p)
		case map[string]interface{}:
			var v map[string]interface{}
			if err := dec.Decode(&v); err != nil {
				// t.Errorf("Decode %s: %v\n", file.Name(), err)
				return errors.Wrap(err, fmt.Sprintf("Unmarshal. Decode %s\n", p))
			}
			// log.Printf("map[striing]interface{}: %s\n", p)

		default:
			log.Printf("%s %v\n", p, reflect.TypeOf(first))
		}
	*/

	return nil
}
