package interval

import (
	"context"
	"time"

	nc "github.com/geniusrabbit/notificationcenter"
)

type interval struct {
	nc.ModelSubscriber

	timeInterval time.Duration
	ticker       *time.Ticker
	msgFnk       func() interface{}
}

// NewSubscriber with interval message generation
func NewSubscriber(timeInterval time.Duration, options ...Option) nc.Subscriber {
	var opts Options
	for _, opt := range options {
		opt(&opts)
	}
	if opts.Handler == nil {
		opts.Handler = func() interface{} { return struct{}{} }
	}
	return &interval{
		ModelSubscriber: nc.ModelSubscriber{
			ErrorHandler: opts.ErrorHandler,
			PanicHandler: opts.PanicHandler,
		},
		timeInterval: timeInterval,
		msgFnk:       opts.Handler,
	}
}

// Listen kafka consumer
func (s *interval) Listen(ctx context.Context) error {
	s.ticker = time.NewTicker(s.timeInterval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-s.ticker.C:
			if err := s.ProcessMessage(s.message()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *interval) message() *message {
	return &message{s.msgFnk()}
}

// Close nstream client
func (s *interval) Close() error {
	s.ticker.Stop()
	return s.ModelSubscriber.Close()
}
