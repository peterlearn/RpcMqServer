package amqprpc

import (
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type callbackFunc func(ch *amqp.Channel) error

type RobustTransport struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	config       *Config
	setupFuncs   []callbackFunc
	cleanupFuncs []callbackFunc
	mu           sync.RWMutex
	interval     time.Duration
	log          Logger
	done         chan struct{}
}

func NewTransport(config *Config) (t *RobustTransport, err error) {
	t = new(RobustTransport)
	t.config = config
	t.log = NewLogWrapper(config.Log)
	t.interval = time.Duration(config.ReconnectInterval) * time.Second
	t.done = make(chan struct{}, 1)

	if err := t.init(); err != nil {
		return nil, err
	}
	go t.handleReconnect()
	t.log.Infof("Created robust transport, dsn: %s", t.config.Dsn)
	return t, nil
}

func (t *RobustTransport) AddSetupFunc(f callbackFunc) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.setupFuncs = append(t.setupFuncs, f)
	return f(t.ch)
}

func (t *RobustTransport) AddCleanupFunc(f callbackFunc) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.cleanupFuncs = append(t.cleanupFuncs, f)
}

func (t *RobustTransport) Close() error {
	t.done <- struct{}{}
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, f := range t.cleanupFuncs {
		if err := f(t.ch); err != nil {
			return err
		}
	}

	if err := t.ch.Close(); err != nil {
		return err
	}

	return t.conn.Close()
}

func (t *RobustTransport) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.log.Debugf("Publish message: %+v, routing key: %s", msg, key)
	return t.ch.Publish(exchange, key, mandatory, immediate, msg)
}

func (t *RobustTransport) handleReconnect() {
	if err := <-t.conn.NotifyClose(make(chan *amqp.Error)); err == nil {
		t.log.Info("Robust transport handleReconnect stopped")
		return
	} else {
		t.log.Errorf("Robust transport received connection error: %+v", err)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	for {
		select {
		case <-t.done:
			t.log.Info("Robust transport handleReconnect stopped")
			return
		default:
			if err := t.recover(); err != nil {
				t.log.Error(err)
				time.Sleep(t.interval)
				continue
			} else {
				t.log.Info("Robust connection was restored")
				go t.handleReconnect()
				return
			}
		}
	}
}

func (t *RobustTransport) recover() error {
	if err := t.init(); err != nil {
		return err
	}

	for _, f := range t.setupFuncs {
		if err := f(t.ch); err != nil {
			return err
		}
	}

	return nil
}

func (t *RobustTransport) init() (err error) {
	if t.conn, err = amqp.DialConfig(t.config.Dsn, t.config.DialConfig); err != nil {
		return err
	}

	if t.ch, err = t.conn.Channel(); err != nil {
		return err
	}

	if err = t.ch.Qos(t.config.PrefetchCount,
		t.config.PrefetchSize,
		t.config.QosGlobal); err != nil {
		return err
	}

	if t.config.Exchange != "" {
		if err = t.ch.ExchangeDeclare(
			t.config.Exchange,
			amqp.ExchangeDirect,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return err
		}
	}
	return nil
}
