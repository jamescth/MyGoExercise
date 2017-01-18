package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	//"time"
)

const numBuffer = 100

var EndRnStr = []byte("END\r\n")

var freeBufCh chan struct{}

func main() {
	/*
		f, err := os.Create(time.Now().Format("2006-01-02T150405.pprof"))
		if err != nil {
			panic(err)
		}

		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			trace.Stop()
			f.Close()
			os.Exit(1)
		}()

		if err := trace.Start(f); err != nil {
			panic(err)
		}
	*/

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

	cnt := 0
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		//err = conn.(*net.TCPConn).SetReadBuffer(6 * 1024 * 1024)
		if err != nil {
			panic(err)
		}

		go worker(conn, cnt)
		cnt++
		fmt.Println("conn", cnt)
	}

}

func worker(conn net.Conn, index int) {
	r := bufio.NewReaderSize(conn, 1*1024*1024)

	for {
		//_, err := r.ReadSlice('\n')
		//_, err = r.ReadSlice('\n')
		buf := <-freeBufCh
		_, err := r.ReadSlice('\n')
		_, err = r.ReadSlice('\n')
		freeBufCh <- buf

		if err != nil {
			break
		}
		conn.Write(EndRnStr)
	}
}
