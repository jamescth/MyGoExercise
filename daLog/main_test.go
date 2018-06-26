package main

import (
	"fmt"
	"os"
	"testing"
)

/*
func TestMain(t *testing.T) {
	cfg0, err := readConf("./config.json")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		num  int
		root string
		cfg  *Conf
		err  error
	}{
		{1, "/cores/bug_30049/", cfg0, nil},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d", tc.num), func(t *testing.T) {
			if err := SearchPath(tc.root, tc.cfg); err != tc.err {
				t.Fatal(err)
			}
		})
	}

	for _, f := range flist {
		t.Log(f)
	}
}
*/

func TestGZ(t *testing.T) {
	tests := []struct {
		num      int
		ret      bool
		mnt      string
		filePath string
	}{
		{1, false, "/home/jho/colo00/",
			"net/fs-01/tank0/autosupport/internal/asupdata/2017_07_22/31653337-1707-2220-a8e0-1c72fd6a0c00/vesx00_00_50_56_a4_28_87/debug-cae029c06f1a11e79af202dafd6a0c00/2017-07-22T20_17_34/var/log/da_setup.log"},
		{2, true, "/home/jho/colo00/",
			"net/fs-01/tank0/autosupport/internal/asupdata/2017_07_22/31653337-1707-2220-a8e0-1c72fd6a0c00/vesx00_00_50_56_a4_28_87/debug-cae029c06f1a11e79af202dafd6a0c00/2017-07-22T20_17_34/var/log/datrium/upgrade_mgr.log.2017-07-22T20-14-34.53728.0.gz"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d", tc.num), func(t *testing.T) {
			buf := make([]byte, 2)
			f, err := os.Open(tc.mnt + tc.filePath)
			if err != nil {
				t.Fatal(err)
			}

			_, err = f.Read(buf)
			if err != nil {
				t.Fatal(err)
			}
			if buf[0] == 31 && buf[1] == 139 {
				if tc.ret == false {
					t.Fatal("%d wanted %v got false\n", tc.num, tc.ret)
				}
			}
		})
	}
}

/*
func TestBusy(t *testing.T) {
	fPath := "/home/jho/esx/vmfs/volumes/596fd509-22533401-54bc-0050569f639e/log/esx_platmgr_dasys.log"

	f, err := os.OpenFile(fPath, os.O_RDONLY|os.O_DIRECT, 0644)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
}
*/
