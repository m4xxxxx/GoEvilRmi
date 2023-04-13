package tools

import (
	"bytes"
)

func WriteUTFparse(data []byte) []byte {
	var bs bytes.Buffer
	port1 := len(data)
	p1 := byte(port1 >> 8)   // 高八位
	p2 := byte(port1 & 0xff) // 低八位
	bs.WriteByte(p1)
	bs.WriteByte(p2)
	bs.Write(data)
	return bs.Bytes()
}

func Writeshortparse(port1 int) []byte {
	var bs bytes.Buffer
	p1 := byte(port1 >> 8)   // 高八位
	p2 := byte(port1 & 0xff) // 低八位
	bs.WriteByte(p1)
	bs.WriteByte(p2)
	return bs.Bytes()
}
