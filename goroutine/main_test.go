package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func run(ch chan<- *bytes.Buffer) {
	defer close(ch)

	var b bytes.Buffer
	b.WriteString("Hello\n")
	b.WriteString("World\n")
	ch <- &b
}

func TestHelloWorld(t *testing.T) {
	// t.Fatal("not implemented")
	chans := make([]chan *bytes.Buffer, 4)

	for i, _ := range chans {
		chans[i] = make(chan *bytes.Buffer, 1)
		go func(idx int) {
			run(chans[idx])
		}(i)
	}

	for i, _ := range chans {
		r := <-chans[i]
		fmt.Println(r.WriteTo(os.Stdout))
	}
}
