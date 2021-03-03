package guds

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/tidwall/gjson"
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

func Map2Data(m interface{}) Data {
	var b []byte
	switch m.(type) {
	case map[string]interface{}:
		b, _ = json.Marshal(m.(map[string]interface{}))
	case string:
		b = []byte(m.(string))
	case []byte:
		b = m.([]byte)
	case int:
		b = IntToBytes(m.(int))
	}
	return Data{gjson.ParseBytes(b)}
}

func JSON2Bytes(d map[string]interface{}) []byte {
	b, _ := json.Marshal(d)
	return b
}
