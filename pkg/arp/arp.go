package arp

const (
	// ArpRequest ARP请求
	ArpRequest = 0x01
	// ArpReply ARP应答
	ArpReply    = 0x02
	RArpRequest = 0x03
	RArpReply   = 0x04
)

// ArpPacket  ARP数据包
type ArpPacket struct {
	DstMac [6]byte // 目的地址
	SrcMac [6]byte // 来源地址
	Frame  uint16  // 长度或类型

	HwType     uint16   // 硬件类型
	ProtoType  uint16   // 协议类型
	HwLen      byte     // 硬件大小
	ProtoLen   byte     // 协议大小
	Op         uint16   // 操作
	ArpSrcMac  [6]byte  // 源硬件地址
	ArpSrcIp   [4]byte  // 源IPv4地址
	ArpDstMac  [6]byte  // 目标硬件地址
	ArpDstIp   [4]byte  // 目标IPv4地址
	ArpPadding [18]byte // 填充字段
	//ArpFCS     [4]byte
}
