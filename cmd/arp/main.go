package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"net-example/pkg/arp"
	net2 "net-example/pkg/net"
	"sync"
	"syscall"
)

const (
	// ArpRequest ARP请求
	ArpRequest = 0x01
	// ArpReply ARP应答
	ArpReply    = 0x02
	RArpRequest = 0x03
	RArpReply   = 0x04
)


func main() {
	// 获取网卡接口
	eno1, err := net.InterfaceByName("eno1")
	if err != nil {
		log.Fatalln(err)
	}

	// 目标Mac地址
	dstMac, err := net.ParseMAC("ff:ff:ff:ff:ff:ff")
	if err != nil {
		log.Fatalln(err)
	}
	// 来源Mac地址
	//srcMac, err := net.ParseMAC("AA:BB:CC:DD:EE:FF")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	// 构造arp数据包
	packet := new(arp.ArpPacket)
	// 构造头部信息
	copy(packet.DstMac[:], dstMac)
	copy(packet.SrcMac[:], eno1.HardwareAddr)
	packet.Frame = syscall.ETH_P_ARP

	// 构造ARP类型
	packet.HwType = 0x01
	packet.ProtoType = 0x0800
	packet.HwLen = 0x06
	packet.ProtoLen = 0x04
	packet.Op = ArpRequest

	// 构造来源ARP地址信息
	srcIp := net.ParseIP("192.168.199.153")
	copy(packet.ArpSrcMac[:], eno1.HardwareAddr)
	copy(packet.ArpSrcIp[:], srcIp.To4())

	// 构造目的ARP地址信息
	dstIp := net.ParseIP("192.168.199.101")
	//copy(packet.ArpDstMac[:], dstMac)
	copy(packet.ArpDstIp[:], dstIp.To4())

	// 将Packet结构转为bytes
	buffer := bytes.NewBuffer(make([]byte, 0))
	if err := binary.Write(buffer, binary.BigEndian, packet); err != nil {
		log.Fatalln(err)
	}

	var group = sync.WaitGroup{}

	group.Add(1)
	go func() {
		// 通过 Linux 原始套接字完成数据包的发送
		sockfd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, net2.Htons16(syscall.ETH_P_ARP))
		if err != nil {
			log.Fatalln(err)
		}
		var count = 0
		linklayer := new(syscall.SockaddrLinklayer)
		linklayer.Ifindex = eno1.Index
		for count < 10 {
			// 发送数据
			if err = syscall.Sendto(sockfd, buffer.Bytes(), 0, linklayer); err != nil {
				log.Fatalln(err)
			}

			//buf := make([]byte, 2048)
			//recvfrom, from, err := syscall.Recvfrom(sockfd, buf, 0)
			//if err != nil {
			//	log.Fatalln(err)
			//}
			//
			//fmt.Println(buf)
			//fmt.Println(recvfrom)
			//fmt.Println(from)

			count++
		}
		// 关闭套接字
		if err := syscall.Close(sockfd); err != nil {
			log.Fatalln(err)
		}
		group.Done()
	}()

	//group.Add(1)
	//go func() {
	//	// 通过 Linux 原始套接字完成数据包的发送
	//	recvFd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//
	//	var count = 0
	//
	//	fmt.Println("======================== Recv")
	//	for count < 5{
	//		fmt.Printf("============= Recv count: %d\n", count)
	//
	//		buf := make([]byte, 0)
	//		recvfrom, from, err := syscall.Recvfrom(recvFd, buf, 0)
	//		if err != nil {
	//			log.Fatalln(err)
	//		}
	//
	//		fmt.Println(buf)
	//		fmt.Println(recvfrom)
	//		fmt.Println(from)
	//
	//		count++
	//	}
	//
	//	fmt.Println("============= Recv end")
	//	// 关闭套接字
	//	if err := syscall.Close(recvFd); err != nil {
	//		log.Fatalln(err)
	//	}
	//	group.Done()
	//}()

	group.Wait()

}

//0000   ff ff ff ff ff ff 94 d9 b3 20 7c f3 08 06 00 01   ......... |.....
//0010   08 00 06 04 00 01 94 d9 b3 20 7c f3 c0 a8 c7 e3   ......... |.....
//0020   00 00 00 00 00 00 c0 a8 c7 65 00 00 00 00 00 00   .........e......
//0030   00 00 00 00 00 00 00 00 00 00 00 00               ............
//
//0000   ff ff ff ff ff ff 40 b0 76 81 ad 3f 08 06 00 01   ......@.v..?....
//0010   08 00 06 04 00 01 40 b0 76 81 ad 3f c0 a8 c7 99   ......@.v..?....
//0020   00 00 00 00 00 00 c0 a8 c7 65 00 00 00 00 00 00   .........e......
//0030   00 00..
