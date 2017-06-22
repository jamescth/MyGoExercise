package main

import (
	"bytes"
	"testing"
)

func TestExercise(t *testing.T) {
	b := bytes.Split(exercise(), []byte("\n"))
	mapOut := make(map[int]string)

	for idx, str := range b {
		mapOut[idx+1] = string(str)
	}

	tests := []struct {
		num int
		out string
	}{
		{1, "1"},
		{2, "2"},
		{3, "Fizz"},
		{5, "Buzz"},
		{15, "FizzBuzz"},
		{100, "Buzz"},
	}

	for _, tc := range tests {
		if mapOut[tc.num] != tc.out {
			t.Errorf("Test failed %d expect %s got %s\n",
				tc.num, tc.out, mapOut[tc.num])
		}
	}
}

func BenchmarkExercise(b *testing.B) {
	for n := 0; n < b.N; n++ {
		exercise()
	}
}
