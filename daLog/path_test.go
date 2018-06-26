package main

import (
	"fmt"
	"os"
	"testing"
)

func TestSearchPathLink(t *testing.T) {
	cfg := Conf{
		MountPoint: "/home/jho/colo00",
		LogFiles: []LogFile{
			LogFile{
				Prefix: "core.7.python2.7.6438.gz",
				Type:   "core",
			},
		},
	}

	if err := SearchPath("/cores/bug_29940", &cfg); err != nil {
		t.Fatal(err)
	}

	if err := cfg.FileList(os.Stdout); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func TestSearchPathUnknown(t *testing.T) {
	cfg := Conf{
		MountPoint: "/home/jho/colo00",
		LogFiles: []LogFile{
			LogFile{
				Prefix: "da_setup.log",
				Type:   "plain",
			},
		},
	}

	if err := SearchPath("/cores/bug_29940", &cfg); err != nil {
		t.Fatal(err)
	}

	for _, logf := range cfg.LogFiles {
		t.Logf("%v\n", logf.DaFiles)
	}

	if err := cfg.FileList(os.Stdout); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

}

func TestSearchPath29940(t *testing.T) {
	cfg, err := readConf("config.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := SearchPath("/cores/bug_29940", cfg); err != nil {
		t.Fatal(err)
	}

	for _, logf := range cfg.LogFiles {
		t.Logf("%v\n", logf.DaFiles)
	}
}
