/*
 Providing a traceroute package. This uses PacketConn example as reference
 https://godoc.org/golang.org/x/net/icmp#PacketConn

 https://github.com/aeden/traceroute uses syscall.Socket() implementation.
 syscall.
*/
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// golang.org/x/net/internal/iana defines ProtocolICMP.
// However, after go1.5, it's not allowed to access any internal files
// So, define a new one here
// https://godoc.org/golang.org/x/net/internal/iana
const (
	ProtocolICMP = 1
	TargetAddr   = "10.25.130.49"
)

//func traceroute(dst string) {
func main() {
	ttl := 1

	for {
		// start := time.Now()

		// set up the icmp receiver
		c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err != nil {
			log.Fatalf("listen err, %s", err)
		}

		if p := c.IPv4PacketConn(); p != nil {
			if err := p.SetTTL(ttl); err != nil {
				log.Fatal(err)
			}
		}
		wm := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  1,
				Data: []byte("HELLO-R-U-There"),
			},
		}

		wb, err := wm.Marshal(nil)
		if err != nil {
			log.Fatal(err)
		}

		start := time.Now()
		if _, err := c.WriteTo(wb, &net.IPAddr{IP: net.ParseIP(TargetAddr)}); err != nil {
			log.Fatalf("WriteTo err, %s", err)
		}

		rb := make([]byte, 1500)
		n, peer, err := c.ReadFrom(rb)

		delta := time.Since(start)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(delta, n, peer)
		rm, err := icmp.ParseMessage(ProtocolICMP, rb[:n])
		if err != nil {
			log.Fatal(err)
		}

		switch rm.Type {
		case ipv4.ICMPTypeEchoReply:
			// if we get Echo Reply, we successfully ping the target, retrun
			log.Printf("got reflection from %v", peer)
			return
		default:
			// keep increament ttl if we cant ping
			ttl++
			//log.Printf("got %+v; want echo reply", rm)
		}
		c.Close()
	} // for loop
}
