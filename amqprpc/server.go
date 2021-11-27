package amqprpc

import (
	"sync"

	"github.com/streadway/amqp"
)

type Server struct {
	transport  *RobustTransport
	serializer Serializer
	queues     map[string]amqp.Queue
	registries []*Registry
	mu         sync.Mutex
	config     *Config
	log        Logger
}

func NewServer(config *Config) (srv *Server, err error) {
	srv = new(Server)
	srv.config = config
	srv.log = NewLogWrapper(config.Log)
	srv.queues = make(map[string]amqp.Queue)

	if config.Serializer != nil {
		srv.serializer = config.Serializer
	} else {
		srv.serializer = DefaultSerializer
	}

	if srv.transport, err = NewTransport(config); err != nil {
		return nil, err
	}

	return srv, nil
}

func (s *Server) Setup() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transport.AddCleanupFunc(s.cleanup)
	return s.transport.AddSetupFunc(s.setup)
}

func (s *Server) setup(ch *amqp.Channel) error {
	for _, r := range s.registries {
		for name, meth := range r.GetMethods() {
			if err := s.regMethod(ch, name, meth); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Server) cleanup(ch *amqp.Channel) error {
	for _, r := range s.registries {
		for name, meth := range r.GetMethods() {
			if err := s.unregMethod(name, meth); err != nil {
				return err
			}
		}
	}

	for name := range s.queues {
		if err := ch.QueueUnbind(name, "", s.config.Exchange, nil); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.log.Info("Close rpc server")
	return s.transport.Close()
}

func (s *Server) AddRegistry(reg *Registry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.registries = append(s.registries, reg)
}

func (s *Server) regMethod(ch *amqp.Channel, name string, meth Method) error {
	if err := meth.Setup(s.serializer); err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		//fmt.Sprintf("%s.%s", s.config.Exchange, name), // name
		Six007rpcApi,
		s.config.IsDurable,                            // durable
		s.config.AutoDelete,                           // delete when unused
		false,                                         // exclusive
		false,                                         // no-wait
		//amqp.Table{
		//	"x-dead-letter-exchange": s.config.Exchange,
		//}, // arguments
		nil,
	)

	if err != nil {
		return err
	}

	s.queues[q.Name] = q

	if err := ch.QueueBind(
		q.Name,
		Six007rpcApi,
		//"",
		s.config.Exchange,
		false,
		//amqp.Table{},
		nil,
	); err != nil {
		return err
	}

	return s.consume(ch, q, meth)
}

func (s *Server) unregMethod(name string, meth Method) error {
	if err := meth.Cleanup(); err != nil {
		return err
	}
	return nil
}

func (s *Server) messageHadler(msg amqp.Delivery, meth Method) {
	if msg.Type != RPCCallMessageType {
		s.log.Warnf("Invalid message type: %s", msg.Type)
		return
	}

	result, rpcErr := meth.Call(msg.Body)

	var (
		msgType string
		body    []byte
		err     error
	)

	if rpcErr != nil {
		if body, err = s.serializer.Marshal(rpcErr); err != nil {
			s.log.Error(err)
			return
		}
		msgType = RPCErrorMessageType
	} else {
		if body, err = s.serializer.Marshal(result); err != nil {
			s.log.Error(err)
			return
		}
		msgType = RPCResultMessageType
	}

	if err := s.transport.Publish(
		"",          // exchange
		msg.ReplyTo, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   s.serializer.GetContentType(),
			CorrelationId: msg.CorrelationId,
			Body:          body,
			Type:          msgType,
		}); err != nil {
		s.log.Errorf("Failed to publish a message: %+v, rpc method name: %s", msg, meth.GetName())
		msg.Reject(false)
		return

	}
	msg.Ack(false)
}

func (s *Server) consume(ch *amqp.Channel, q amqp.Queue, meth Method) error {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		return err
	}

	go func(msgs <-chan amqp.Delivery) {
		s.log.Infof("Started consumer: %s", q.Name)
		for {
			msg, ok := <-msgs
			if !ok {
				s.log.Infof("Stopped consumer: %s", q.Name)
				return
			}
			go s.messageHadler(msg, meth)
		}
	}(msgs)

	return nil
}
