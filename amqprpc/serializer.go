package amqprpc

import (
	"github.com/vmihailenco/msgpack/v5"
)

type Serializer interface {
	GetContentType() string
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

var DefaultSerializer = new(MsgPackSerializer)

type MsgPackSerializer struct{}

func (s *MsgPackSerializer) GetContentType() string {
	return "application/x-msgpack"
}

func (s *MsgPackSerializer) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (s *MsgPackSerializer) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}
