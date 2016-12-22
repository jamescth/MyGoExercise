package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

const numBuffer = 128

var EndRnStr = []byte("END\r\n")

var freeBufCh chan struct{}

func main() {

	var portOpt = flag.Int("port", 11211, "port on which to listen for connections")
	flag.Parse()

	freeBufCh = make(chan struct{}, numBuffer)
	for i := 0; i < numBuffer; i++ {
		freeBufCh <- struct{}{}
	}

	fmt.Println("Port :", *portOpt)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *portOpt))

	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		//err = conn.(*net.TCPConn).SetReadBuffer(128 * 1024)
		//if err != nil {
		//panic(err)
		//}

		go worker(conn)

	}

}

func worker(conn net.Conn) {
	r := bufio.NewReader(conn)
	for {
		//buf := <-freeBufCh
		_, err := r.ReadSlice('\n')
		_, err = r.ReadSlice('\n')
		if err != nil {
			break
		}
		//freeBufCh <- buf
		conn.Write(EndRnStr)
	}
}
