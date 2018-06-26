package config

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	tests := []struct {
		num  int
		time string
	}{
		{1, "2018-01-30T06:20:20.641+0000"},
		{2, "2018-01-30T06:20:21.760+0000"},
	}

	for _, tc := range tests {
		//ti, err := time.Parse(time.RFC3339, tc.time)
		ti, err := time.Parse("2006-01-02T15:04:05+0000", tc.time)
		if err != nil {
			t.Fatalf("time parse failed %d %s %v", tc.num, tc.time, err)
		}

		t.Logf("time is %s", ti)
		t.Log(ti.Unix())
		t.Log(ti.UnixNano())
	}
}
