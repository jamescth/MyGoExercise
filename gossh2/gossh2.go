/*
 https://www.digitalocean.com/community/tutorials/how-to-set-up-ssh-keys--2

 1. ssh-keygen -t rsa
 2. ssh-copy-id user@ip
    or
	cat ~/.ssh/id_rsa.pub | ssh user@ip "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"
*/
package main

import (
	"fmt"
	"os"
	_ "time"

	_ "golang.org/x/crypto/ssh"
)

var (
	keys []string
)

func main() {
	// cmd := os.Args[1]
	// hosts := os.Args[2]
	// results := make(chan string, 10)
	// timeout := time.After(5 * time.Second)

	// init the structure w/ the configuration for ssh packet
	/*
		config := &ssh.ClientConfig{
			User: os.Getenv("hoj9"),
			Auth: []ssh.ClientAuth{makeKeyring()},
		}
	*/
	keys = []string{os.Getenv("HOME") + "/.ssh/id_rsa", os.Getenv("HOME") + "/.ssh/id_dsa", os.Getenv("HOME") + "/.ssh/id_ecdsa"}

	fmt.Println(keys)
}
