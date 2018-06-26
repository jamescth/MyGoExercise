package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Local string

func (l Local) Open(p string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(string(l), p))
}

func (l Local) Create(p string) (io.WriteCloser, error) {
	return os.Create(filepath.Join(string(l), p))
}

func main() {
	fmt.Println("vim-go")
}
