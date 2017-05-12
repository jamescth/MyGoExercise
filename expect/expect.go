package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

type cmds []string

func (c cmds) Run() (bytes.Buffer, error) {
	var out bytes.Buffer
	var outerr bytes.Buffer

	fmt.Printf("%v\n", c)
	args := append([]string{"-c"}, c...)
	cmd := exec.Command("/bin/sh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = &out
	cmd.Stderr = &outerr

	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
			}
		}
		log.Printf("execCmd OUT: %s\n", string(out.String()))
		log.Printf("execCmd ERR: %s\n", string(outerr.String()))
		return out, fmt.Errorf("execCmd %v", err)
	}
	return out, nil
}

func main() {
	fmt.Println("vim-go")
}
