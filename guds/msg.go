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
		length     int
		bytes      []byte
		bodyLength int
		bodyBytes  []byte
		async      bool
		Api        string      `json:"api"`
		Body       interface{} `json:"body,omitempty"`
	}
)

func (m *Message) BodyBytes() []byte {
	b, _ := json.Marshal(m.Body)
	return b
}

func (m *Message) BodyLength() int {
	return m.bodyLength
}

func (m *Message) BodyStringify() string {
	return string(m.bodyBytes)
}

func NewMessage(api string, data interface{}) *Message {
	m := &Message{
		Api:  api,
		Body: data,
	}
	m.bytes = m.Bytes()
	m.length = len(m.bytes)
	m.bodyBytes = m.BodyBytes()
	m.bodyLength = len(m.bodyBytes)
	return m
}

func (m *Message) reset(async bool, body interface{}) {
	m.async = async
	m.Body = body
	m.bytes = m.Bytes()
	m.length = len(m.bytes)
	m.bodyBytes = m.BodyBytes()
	m.bodyLength = len(m.bodyBytes)
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
	m.bodyBytes = m.BodyBytes()
	m.bodyLength = len(m.bodyBytes)
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
	return gjson.ParseBytes(m.bodyBytes)
}

func (m *Message) ToData() Data {
	return Data{m.GJson()}
}

func (m *Message) Length() int {
	return m.length
}
