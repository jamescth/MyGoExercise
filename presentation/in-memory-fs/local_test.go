package main

import (
	"log"
	"testing"
)

func TestLocal(t *testing.T) {
	rc, err := Local("/tmp").Create("helloworld")
	if err != nil {
		t.Errorf("Local Create failed %v", err)
	}
	defer rc.Close()

	ro, err := Local("/tmp").Open("helloworld")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer ro.Close()

}
