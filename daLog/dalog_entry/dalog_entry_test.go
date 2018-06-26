package dalog_entry

import "testing"

func TestDalogEntrySort(t *testing.T) {
	dalogs := []DalogEntry{
		{TimeStamp: "20170406:08:00:10", Msg: "HelloA", Priority: "DA_LOG_ERR", Host: "A", Path: "/fake/test", LineNum: 1},
		{TimeStamp: "20170406:08:00:00", Msg: "HelloB", Priority: "DA_LOG_INFO", Host: "B", Path: "/fake/test", LineNum: 10},
		{TimeStamp: "20170406:08:00:20", Msg: "Hello", Priority: "DA_LOG_EMERG", Host: "A", Path: "/fake/test", LineNum: 11},
		{TimeStamp: "20170406:08:00:30", Msg: "HelloC", Priority: "DA_LOG_CRIT", Host: "B", Path: "/fake/test", LineNum: 12},
		{TimeStamp: "20170406:08:00:10", Msg: "HelloE", Priority: "DA_LOG_DEBUG", Host: "C", Path: "/fake/test", LineNum: 31},
		{TimeStamp: "20170406:08:00:150", Msg: "Hello", Priority: "DA_LOG_DEFAULT", Host: "A", Path: "/fake/test", LineNum: 11},
		{TimeStamp: "20170406:08:00:120", Msg: "HelloG", Priority: "DA_LOG_ERR", Host: "D", Path: "/fake/test", LineNum: 91},
		{TimeStamp: "20170406:08:00:130", Msg: "HelloZ", Priority: "DA_LOG_ALERT", Host: "A", Path: "/fake/test", LineNum: 51},
		{TimeStamp: "20170406:08:00:0", Msg: "HelloE", Priority: "DA_LOG_ERR", Host: "D", Path: "/fake/test", LineNum: 1},
		{TimeStamp: "20170406:08:00:01", Msg: "HelloF", Priority: "DA_LOG_WARNING", Host: "A", Path: "/fake/test", LineNum: 4},
	}

	fnSortByTime := func(d1, d2 *DalogEntry) bool {
		return d1.TimeStamp < d2.TimeStamp
	}

	By(fnSortByTime).Sort(dalogs)

	for idx, _ := range dalogs {
		if idx >= len(dalogs)-1 {
			break
		}

		if dalogs[idx].TimeStamp > dalogs[idx+1].TimeStamp {
			t.Fatalf("Test fail: \n  dalogs[%d].TimeStamp %s\n  dalogs[%d].TimeStamp %s\n",
				idx, dalogs[idx].TimeStamp, idx+1, dalogs[idx+1].TimeStamp)
		}
	}
}
