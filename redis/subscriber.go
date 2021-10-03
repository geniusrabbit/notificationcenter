package redis

import (
	"context"

	"github.com/go-redis/redis"

	nc "github.com/geniusrabbit/notificationcenter"
)

// Subscriber for Redis channels
type Subscriber struct {
	nc.ModelSubscriber
	ctx context.Context

	// List of subscibed channels
	channels []string

	// connection to the NATS stream server
	cli *redis.Client
	sub *redis.PubSub

	// Close connection event queue
	closeEvent chan bool

	// logger interface
	logger loggerInterface
}

// NewSubscriber creates new subscriber object
func NewSubscriber(options ...Option) (*Subscriber, error) {
	var opts Options
	for _, opt := range options {
		opt(&opts)
	}
	sub := &Subscriber{
		ctx: context.Background(),
		ModelSubscriber: nc.ModelSubscriber{
			ErrorHandler: opts.ErrorHandler,
			PanicHandler: opts.PanicHandler,
		},
		cli:        opts.connect(),
		channels:   opts.channels(),
		closeEvent: make(chan bool, 1),
		logger:     opts.logger(),
	}
	go sub.subscribe(sub.ctx)
	return sub, nil
}

// MustNewSubscriber creates new subscriber object
func MustNewSubscriber(options ...Option) *Subscriber {
	sub, err := NewSubscriber(options...)
	if err != nil || sub == nil {
		panic(err)
	}
	return sub
}

// Listen kafka consumer
func (s *Subscriber) Listen(ctx context.Context) error {
	select {
	case <-ctx.Done():
	case <-s.closeEvent:
	}
	return nil
}

func (s *Subscriber) subscribe(ctx context.Context) {
	s.sub = s.cli.Subscribe(s.channels...)
	for msg := range s.sub.Channel() {
		s.message(ctx, msg)
	}
}

// message execute
func (s *Subscriber) message(ctx context.Context, m *redis.Message) {
	if err := s.ProcessMessage(messageFromRedis(ctx, m)); err != nil {
		s.logger.Error(err)
	}
}

// Close nstream client
func (s *Subscriber) Close() error {
	err := s.cli.Close()
	if err == nil {
		err = s.sub.Close()
	}
	s.closeEvent <- true
	return err
}
