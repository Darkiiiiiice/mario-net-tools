package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net-example/pkg/arp"
	"net-example/pkg/net"
	"syscall"
)

func main() {
	// 通过 Linux 原始套接字完成数据包的发送
	recvFd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, net.Htons16(syscall.ETH_P_ARP))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(recvFd)
	//syscall.Socketpair()

	var count = 0

	fmt.Println("======================== Recv")
	for  {
		fmt.Printf("============= Recv count: %d\n", count)

		buf := make([]byte, 60)
		n, fromAddr, err := syscall.Recvfrom(recvFd, buf, 0)
		if err != nil {
			log.Fatalln(err)
		}

		packet := new(arp.ArpPacket)
		buffer := bytes.NewBuffer(buf)
		if err := binary.Read(buffer, binary.BigEndian, packet); err != nil {
			log.Fatalln(err)
		}
		if packet.Op == arp.ArpReply {
			fmt.Println(buf)
			fmt.Println(n)
			fmt.Println(fromAddr)
			any := fromAddr.(*syscall.SockaddrLinklayer)
			fmt.Println(any)
			fmt.Printf("%+v\n",packet)
		}



		count++
	}

	fmt.Println("============= Recv end")
	// 关闭套接字
	if err := syscall.Close(recvFd); err != nil {
		log.Fatalln(err)
	}
}
