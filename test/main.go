package main

import "fmt"

var DEBUG bool

func main() {
	if DEBUG {
		fmt.Println("vim-go")
		return
	}
	fmt.Println("bye")
}
