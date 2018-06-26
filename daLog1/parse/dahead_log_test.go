package parse

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

func TestDaHeadLog(t *testing.T) {
	tests := []struct {
		num int
		p   string
	}{
		{1, "../tests/vesx00_da_00_09_39_13_01_debug_support_id_5e7d8f05059011e88da0b3c4f72695bf_2018-01-30T07_37_06/var/run/log/da_head.log"},
	}

	for _, tc := range tests {
		//da, err := NewDalogEntries(tc.p)
		da, err := NewDaHead(tc.p)
		if err != nil {
			t.Errorf("%+v", err)
		}
		fmt.Println(da.Name())

		var wg sync.WaitGroup
		ch := make(chan *bytes.Buffer, 1)
		wg.Add(1)
		go func() {
			defer wg.Done()
			// fmt.Println(da.Name())
			buf, err := da.ListIssue(&Request{Start: 0, End: 0})
			if err != nil {
				t.Errorf("%+v", err)
			}
			ch <- buf
			close(ch)
		}()

		wg.Wait()
		if da.Start() == 0 || da.End() == 0 {
			t.Fatalf("start %d end %d\n", da.Start(), da.End())
		}

		buf := <-ch
		t.Log(string(buf.Bytes()))
		t.Logf("%s\n%s %s\n", da.Name(), daNanosecond(da.Start()), daNanosecond(da.End()))

	} // for

}
