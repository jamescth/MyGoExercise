package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	fmt.Printf("%s: PID %d\n", os.Args[0], syscall.Getpid())

	cmd := &exec.Cmd{
		Path: os.Args[1],
		Args: []string{"link"},
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID,
	}
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
	cmd.Wait()
}
