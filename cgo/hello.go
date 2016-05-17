package main

// #include <stdio.h>
// #include <errno.h>
// typedef int (*intFunc) ();
//
// int
// bridge_int_func(intFunc f)
// {
//		return f();
// }
//
// int fortytwo()
// {
//	    return 42;
// }
/*
#include <stdio.h>
#include <stdlib.h>

void hello() {
	printf("hello world\n");
}

void myprint(char* s) {
    printf("%s", s);
}
*/
import "C"
import (
	"bufio"
	_ "flag"
	"fmt"
	"os"
	"unsafe"
)

func main() {
	/*
			var (
			dir = flag.String("proc", "", "/proc path")
		)

		flag.Parse()


		f := C.intFunc(C.fortytwo)
		fmt.Println(int(C.bridge_int_func(f)))
		// Output: 42
		Example()
	*/
	C.hello()

	/*
		if *dir != "" {
			read_proc(*dir)
		}
	*/
}

func Example() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func read_proc(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	r := bufio.NewReader(file)

	s, e := Readln(r)
	for e == nil {
		fmt.Println(s)
		s, e = Readln(r)
	}
	return nil
}
