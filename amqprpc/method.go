package amqprpc

import "fmt"

type ErrorData struct {
	Type    string   `msgpack:"type"`
	Message string   `msgpack:"message"`
	Args    []string `msgpack:"args"`
}

type RPCError struct {
	Err ErrorData `msgpack:"error"`
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("type: %s, message: %s, args: %s", e.Err.Type, e.Err.Message, e.Err.Args)
}

type Method interface {
	GetName() string
	Call(body []byte) (interface{}, *RPCError)
	Setup(serializer Serializer) error
	Cleanup() error
}
