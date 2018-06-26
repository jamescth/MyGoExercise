package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

type Mem map[string][]byte

func (m Mem) Open(path string) (io.ReadCloser, error) {
	b, ok := m[path]
	if !ok {
		return nil, fmt.Errorf("%s not found", path)
	}
	return ioutil.NopCloser(bytes.NewReader(b)), nil
}

func (m Mem) Read(path string, n int) ([]byte, error) {
	c, ok := m[path]
	if !ok {
		return nil, fmt.Errorf("%s not found", path)
	}

	if n > len(c) {
		return nil, fmt.Errorf("%s size %d < request size %d", path, len(c), n)
	}

	return c[:n], nil
}
