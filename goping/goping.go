/*
 * http://my.oschina.net/ybusad/blog/300155
 * http://blog.csdn.net/gophers/article/details/21481447
 */
package main

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)

	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

const usage = `
	Usage: goping <url>
	Example: goping www.google.com
`

func GetLocalIP() string {
	ifaces, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range ifaces {
		// check the address type and if it is not a loopback
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func main() {
	arg_num := len(os.Args)
	if arg_num < 2 {
		fmt.Println(usage)
		return
	}
	lip := GetLocalIP()

	var (
		icmp ICMP
		// source ip can be local or eth0 ip
		// laddr    net.IPAddr = net.IPAddr{IP: net.ParseIP("0.0.0.0")}
		laddr    net.IPAddr = net.IPAddr{IP: net.ParseIP(lip)}
		raddr, _            = net.ResolveIPAddr("ip", os.Args[1])
	)

	conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	icmp.Type = 8
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 0
	icmp.SequenceNum = 0

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = CheckSum(buffer.Bytes())
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, icmp)

	fmt.Printf("PING %s (%s)\n", lip, raddr.String())

	recv := make([]byte, 1024)

	statistic := list.New()
	sended_packets := 0

	for i := 4; i > 0; i-- {

		if _, err := conn.Write(buffer.Bytes()); err != nil {
			fmt.Println(err.Error())
			return
		}
		sended_packets++
		t_start := time.Now()

		conn.SetReadDeadline((time.Now().Add(time.Second * 5)))
		_, err := conn.Read(recv)

		if err != nil {
			fmt.Println("Timeout")
			continue
		}

		t_end := time.Now()

		dur := t_end.Sub(t_start).Nanoseconds() / 1e6

		fmt.Printf("from %s: time = %dms\n", raddr.String(), dur)

		statistic.PushBack(dur)

		//for i := 0; i < recvsize; i++ {
		//	if i%16 == 0 {
		//		fmt.Println("")
		//	}
		//	fmt.Printf("%.2x ", recv[i])
		//}
		//fmt.Println("")

	}

	defer func() {
		//
		var min, max, sum int64
		if statistic.Len() == 0 {
			min, max, sum = 0, 0, 0
		} else {
			min, max, sum = statistic.Front().Value.(int64), statistic.Front().Value.(int64), int64(0)
		}

		for v := statistic.Front(); v != nil; v = v.Next() {

			val := v.Value.(int64)

			switch {
			case val < min:
				min = val
			case val > max:
				max = val
			}

			sum = sum + val
		}
		recved, losted := statistic.Len(), sended_packets-statistic.Len()
		fmt.Printf("%s ping statistics\n %d packets transmitted, %d received, %d packet loss (%.1f%% loss)，\nrtt min/max/avg = %d/%d/%.0f ms\n",
			raddr.String(),
			sended_packets, recved, losted, float32(losted)/float32(sended_packets)*100,
			min, max, float32(sum)/float32(recved),
		)
	}()
}
