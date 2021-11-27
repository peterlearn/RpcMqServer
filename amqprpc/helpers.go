package amqprpc

import (
	"github.com/google/uuid"
)

func makeCorrelationId() string {
	return uuid.NewString()
}

func makeRandomQueueName() string {
	return "queue-" + uuid.NewString()
}
