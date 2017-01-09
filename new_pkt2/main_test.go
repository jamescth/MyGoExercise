package main

import "testing"

func TestReadPcap(t *testing.T) {
	if err := countPcapPkt(); err != nil {
		t.Fatal()
	}
}

func TestScanPkts(t *testing.T) {
	if err := scanPkts(); err != nil {
		t.Fatal()
	}
}

//***************************************************************************

func TestABC(t *testing.T) {
	processPkts()
}
