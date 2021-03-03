package guds

import (
	"bytes"
	"encoding/binary"
)

// 大端序字节转 int
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

// int 转大端序字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func MergeBytes(b1 []byte, b2 ...[]byte) []byte {
	var buffer bytes.Buffer
	buffer.Write(b1)
	for _, b := range b2 {
		buffer.Write(b)
	}
	return buffer.Bytes()
}
