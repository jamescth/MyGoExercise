package ssh

import (
	"fmt"
	_ "strings"
	"testing"
)

/*
var (
	usr    = "root"
	passwd = "abc123"
	host   = "artemis13.datadomain.com"
	cmd1   = "ls"
	cmd2   = "whoami"
)
*/
var (
	usr    = "james"
	passwd = "james"
	host   = "192.168.180.224"
	cmds   = []string{"ls", "whoami"}
	neg    = []string{"/home/james/golang/src/gossh/ssh/while.py",
		"kill `pidof python`"}
)

func TestSSH(t *testing.T) {
	myssh := SSH_Config{
		Usr:    usr,
		Passwd: passwd,
		Host:   host,
	}

	c, err := myssh.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	for _, ctx := range cmds {
		// fmt.Println(idx, ctx)
		str, err := myssh.Run_Session(c, ctx)
		if err != nil {
			t.Fatal(ctx, err)
		}
		fmt.Println("CMD:", ctx)
		fmt.Println(str)
	}
}

func TestNeg(t *testing.T) {
	myssh := SSH_Config{
		Usr:    usr,
		Passwd: passwd,
		Host:   host,
	}

	c, err := myssh.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	fmt.Println("Negative test for while loop:", neg[0])
	str, err := myssh.Run_Session(c, neg[0])
	if err == nil {
		t.Fatal("Test should fail due to too long but didn't")
	}
	fmt.Println("expected error:", err)
	fmt.Println("clean up:", neg[1])
	str, err = myssh.Run_Session(c, neg[1])
	if err != nil {
		t.Fatal(neg[1], err)
	}
	fmt.Println(str)
}
