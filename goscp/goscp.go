// go_scp
// https://gist.github.com/jedy/3357393
package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
)

const privateKey = `content of id_rsa`

func main() {
	signer, _ := ssh.ParsePrivateKey([]byte(privateKey))
	clientConfig := &ssh.ClientConfig{
		User: "james",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	client, err := ssh.Dial("tcp", "127.0.0.1:22", clientConfig)
	if err != nil {
		log.Fatal("Failed to dial: ", err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err.Error())
	}
	defer session.Close()

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		content := "123456789\n"
		fmt.Fprintln(w, "D0755", 0, "testdir") // mkdir
		fmt.Fprintln(w, "C0644", len(content), "testfile1")
		fmt.Fprint(w, content)
		fmt.Fprint(w, "\x00") // transfer end with \x00
		fmt.Fprintln(w, "C0644", len(content), "testfile2")
		fmt.Fprint(w, content)
		fmt.Fprint(w, "\x00")
	}()

	// http://stackoverflow.com/questions/19043557/creating-a-file-in-ssh-client-for-golang
	// this is 'remote scp' in 'sink' or 'to' mode.  The scp client relies
	// on another copy of scp being available on the remote
	// scp -t being 'to'/sink mode,
	// scp -f being 'from'/source mode
	// read this how scp protocol works:
	// https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works
	if err := session.Run("/usr/bin/scp -tr ./"); err != nil {
		log.Fatal("Failed to run: ", err.Error())
	}
}
