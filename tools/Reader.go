package tools

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"time"
)

func ReadUTF(conn net.Conn) []byte {
	l, err := readN(conn, 2)
	if err != nil {
		print(err.Error())
	}
	l1 := binary.BigEndian.Uint16(l)
	data, err := readN(conn, int(l1))
	return data
}

func readN(conn net.Conn, n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func ReadName(conn net.Conn) string {
	var data []byte
	buffer := make([]byte, 1)
	for {
		n, _ := conn.Read(buffer)
		data = append(data, buffer[:n]...)
		if bytes.Contains(data, []byte{116, 0}) {
			d := make([]byte, 1)
			conn.Read(d)
			nl := int(d[0])
			d1 := make([]byte, nl)
			conn.Read(d1)
			return string(d1[:])
		}
	}
	return ""
}

func ReadAll(conn net.Conn) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	data := make([]byte, 0, 4096)
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				// 读取完毕
				break
			}
			// 处理错误
			return data, err
		}
		data = append(data, buf[:n]...)
	}
	return data, nil
}
