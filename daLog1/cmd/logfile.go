package main

import (
	"bytes"
	"io"

	"github.com/jamescth/MyGoExercise/daLog1/parse"
)

type logFile struct {
	f  parse.DaOut
	ch chan *bytes.Buffer
}

type dvxlog struct {
	prefix string
	fn     func(string) (parse.DaOut, error)
	st     int64  // start time
	et     int64  // end time
	c      string // component
	s      string // serverity
	w      io.Writer
}
