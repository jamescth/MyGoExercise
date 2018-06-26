package main

import "io"

type FS interface {
	Open(p string) (io.ReadCloser, error)
	Create(p string) (io.WriteCloser, error)
}
