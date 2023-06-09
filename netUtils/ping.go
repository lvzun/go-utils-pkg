package netUtils

import (
	"bytes"
	"encoding/binary"
	"fmt"

	//"fmt"
	"net"
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

func Ping(addr string) (int, error) {
	return PingTimeOut(addr, 300)
}

func PingTimeOut(addr string, milliscond int) (int, error) {
	var (
		icmp     ICMP
		laddr    = net.IPAddr{IP: net.ParseIP("0.0.0.0")} //***IP地址改成你自己的网段***
		raddr, _ = net.ResolveIPAddr("ip", addr)
	)

	//如果你要使用网络层的其他协议还可以设置成 ip:ospf、ip:arp 等
	conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
	if err != nil {
		//fmt.Println(err.Error())
		return 0, err
	}
	defer conn.Close()

	//开始填充数据包
	icmp.Type = 8 //8->echo message  0->reply message
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 0
	icmp.SequenceNum = 0

	var (
		buffer bytes.Buffer
	)
	//先在buffer中写入icmp数据报求去校验和
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = CheckSum(buffer.Bytes())
	//然后清空buffer并把求完校验和的icmp数据报写入其中准备发送
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, icmp)

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		//fmt.Println(err.Error())
		return 0, err
	}

	//读取结果
	t_start := time.Now()

	conn.SetReadDeadline(time.Now().Add(time.Duration(milliscond) * time.Millisecond))
	recv := make([]byte, 1024)
	_, err = conn.Read(recv)
	if err != nil {
		//fmt.Println("请求超时")
		return 0, err
	}

	t_end := time.Now()

	dur := t_end.Sub(t_start).Nanoseconds() / 1e6

	//fmt.Printf("来自 %s 的回复: 时间 = %dms\n", raddr.String(), dur)

	return int(dur), nil
}

func TelNet(ip, port string) bool {
	address := fmt.Sprintf("%s:%s", ip, port)
	fmt.Println("address:", address)
	_, err := net.DialTimeout("tcp", address, time.Second*2)
	return err == nil
}
