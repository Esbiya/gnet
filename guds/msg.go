package guds

import (
	"bytes"
	"encoding/json"
	"github.com/tidwall/gjson"
)

type (
	Data struct {
		gjson.Result
	}
	Reply struct {
		Async  bool
		Status Action
		Body   interface{}
	}
	Message struct {
		length int
		bytes  []byte
		async  bool
		Api    string      `json:"api"`
		Data   interface{} `json:"data,omitempty"`
	}
)

func NewMessage(api string, data interface{}) *Message {
	m := &Message{
		Api:  api,
		Data: data,
	}
	m.bytes = m.Bytes()
	m.length = len(m.bytes)
	return m
}

func (m *Message) reset(async bool, data interface{}) {
	m.async = async
	m.Data = data
	m.bytes = m.Bytes()
	m.length = len(m.bytes)
}

func (m *Message) out() []byte {
	return MergeBytes(IntToBytes(m.length), m.bytes)
}

func (m *Message) Async() bool {
	return m.async
}

func (m *Message) Parse(b []byte) error {
	m.bytes = b
	m.length = len(b)
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	if err := decoder.Decode(&m); err != nil {
		return err
	}
	return nil
}

func (m *Message) Bytes() []byte {
	b, _ := json.Marshal(m)
	return b
}

func (m *Message) Stringify() string {
	return string(m.bytes)
}

func (m *Message) GJson() gjson.Result {
	return gjson.ParseBytes(m.bytes)
}

func (m *Message) ToData() Data {
	return Data{m.GJson()}
}

func (m *Message) Length() int {
	return m.length
}
