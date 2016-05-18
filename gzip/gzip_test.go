// go test -bench=.
package main

import (
	"testing"
)

func BenchmarkGzip(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test_cgo()
	}
}
