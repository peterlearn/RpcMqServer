package amqprpc

import (
	"github.com/streadway/amqp"
)

type Config struct {
	Dsn               string
	Exchange          string
	IsDurable         bool
	AutoDelete        bool
	PrefetchCount     int
	PrefetchSize      int
	QosGlobal         bool
	Log               Logger
	Serializer        Serializer
	DialConfig        amqp.Config
	ReconnectInterval int
	ClientTimeout     int
}
