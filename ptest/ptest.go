// Print the numbers from 1 to 100. But for multiples
// of three print “Fizz” instead of the number and for
// the multiples of five print “Buzz”. For numbers
// which are multiples of both three and five print
// “FizzBuzz”
package main

import "fmt"

const (
	// ascii code
	newline = 10
	zero    = 48

	// FIZZ multiple of 3
	FIZZ = 3
	// BUZZ multiple of 5
	BUZZ = 5
	// FIZZBUZZ multiple of 15
	FIZZBUZZ = FIZZ * BUZZ
)

// exercise() is the main func to return the result in bytes
// back to caller
//
// To check inline func
//    go build -gcflags=-m main.go
func exercise() []byte {
	start, end := 1, 100

	// since CPU is way faster than memory allocation,
	// it's ok to loop it for buffer size calculation
	// so we only alloc once.

	var size int
	for i := start; i <= end; i++ {
		if i%FIZZBUZZ == 0 {
			size += 9
			continue
		}
		if i%BUZZ == 0 || i%FIZZ == 0 {
			size += 5
			continue
		}
		if i < 10 {
			size += 2
			continue
		}
		if i == 100 {
			size += 4
			continue
		}
		size += 3
	}

	b := make([]byte, size)

	// fizz receives the buffer index to insert
	// 'Fizz', and returns n writen bytes
	fizz := func(l int) int {
		b[l] = 'F'
		b[l+1] = 'i'
		b[l+2] = 'z'
		b[l+3] = 'z'
		return 4
	}

	// buzz receives the buffer index to insert
	// 'Buzz', and returns n writen bytes
	buzz := func(l int) int {
		b[l] = 'B'
		b[l+1] = 'u'
		b[l+2] = 'z'
		b[l+3] = 'z'
		return 4
	}

	// digits converts int (1-100) to ascii
	// and insert it to the buffer index
	digits := func(v, l int) int {
		ret := 0

		t1 := v / 100
		t2 := (v % 100) / 10
		t3 := v % 10

		if t1 != 0 {
			b[l] = byte(zero + t1)
			l++
			ret++
		}
		if t2 != 0 || t1 != 0 {
			b[l] = byte(zero + t2)
			l++
			ret++
		}
		b[l] = byte(zero + t3)
		l++
		ret++
		return ret
	}

	var index int
	for i := start; i <= end; i++ {
		if i%FIZZBUZZ == 0 {
			index += fizz(index)
			index += buzz(index)
			b[index] = newline
			index++
			continue
		}
		if i%BUZZ == 0 {
			index += buzz(index)
			b[index] = newline
			index++
			continue
		}
		if i%FIZZ == 0 {
			index += fizz(index)
			b[index] = newline
			index++
			continue
		}
		index += digits(i, index)
		b[index] = newline
		index++
	}

	return b
}

func main() {
	fmt.Println(string(exercise()))
}
