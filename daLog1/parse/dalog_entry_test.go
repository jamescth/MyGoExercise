package parse

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestTimeConvert(t *testing.T) {
	str := "../tests/commands/da_hosts_show.2017-09-30T04_02_49.out"
	base := path.Base(str)
	t.Logf("base %s\n", base)
	base = strings.TrimPrefix(base, "da_hosts_show.")
	t.Logf("base %s\n", base)
	base = strings.TrimRight(base, ".out")
	t.Logf("base %s\n", base)
	sTime, err := time.Parse("2006-01-02T15_04_05", base)
	if err != nil {
		t.Errorf("Error %s %v\n", base, err)
	}
	t.Logf("sTime %s %d\n", sTime, sTime.Unix())
}

func TestReadLine(t *testing.T) {
	tests := []struct {
		num int
		p   string
	}{
		{1, "../tests/vesx00_da_00_09_39_13_01_debug_support_id_5e7d8f05059011e88da0b3c4f72695bf_2018-01-30T07_37_06/var/log/esx_platmgr.log"},
		{2, "../tests/vesx00_da_00_09_39_13_01_debug_support_id_5e7d8f05059011e88da0b3c4f72695bf_2018-01-30T07_37_06/var/log/upgrade_mgr.log.2018-03-05T19-04-02.3212.0.gz"},
		{3, "../tests/FrontEnd0/vesx00_da_00_20_72_8a_de/var/log/datrium/upgrade_mgr.log.2018-03-03T04-12-16.150898.0.gz"},
	}

	// fmt.Println(len(tests))
	for _, tc := range tests {
		//da, err := NewDalogEntries(tc.p)
		da, err := NewDalogEntriesDaOut(tc.p)
		if err != nil {
			t.Errorf("%+v", err)
		}

		//for i, e := range da.Entries {
		// t.Logf("%d %s %d %d\n", i+1, daNanosecond(e.Time), e.Start, e.End)
		//}

		var wg sync.WaitGroup
		ch := make(chan *bytes.Buffer, 1)
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(da.Name())
			buf, err := da.ListIssue(&Request{Start: 0, End: 0})
			if err != nil {
				t.Errorf("%+v", err)
			}
			ch <- buf
			close(ch)
		}()

		wg.Wait()
		buf := <-ch
		t.Log(string(buf.Bytes()))

	} // for
}
