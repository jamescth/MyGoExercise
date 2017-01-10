package main

import (
	"os"
	"testing"
)

//***************************************************************************

func TestProcessPkts(t *testing.T) {
	processPkts("./ddmc.fw.pcap", os.Stdout)
}
