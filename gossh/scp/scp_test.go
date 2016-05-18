package scp

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func getAgent() (agent.Agent, error) {
	agentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	return agent.NewClient(agentConn), err
}

func TestCopyPath(t *testing.T) {
	f, _ := ioutil.TempFile("", "")
	fmt.Fprintln(f, "hello world")
	f.Close()

	defer os.Remove(f.Name())
	defer os.Remove(f.Name() + "-copy")

	agent, err := getAgent()
	if err != nil {
		t.Fatal("Failed to connecto to SSH_AUTH_SOCK:", err)
	}

	client, err := ssh.Dial("tcp", "localhost:22", &ssh.ClientConfig{
		User: os.Getenv("USER"),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agent.Signers),
		},
	})

	if err != nil {
		t.Fatal("Fail to dial:", err)
	}

	session, err := client.NewSession()
	if err != nil {
		t.Fatal("Fail to create new session:", err)
	}

	dest := f.Name() + "-copy"
	err = CopyPath(f.Name(), dest, session)

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Fatal("No file or directory:", dest)
	}
}
