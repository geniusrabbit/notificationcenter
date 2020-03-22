//
// @project geniusrabbit.com 2019 - 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2019 - 2020
//

package pg

import (
	"context"
	"time"

	"github.com/lib/pq"

	nc "github.com/geniusrabbit/notificationcenter"
)

type loggerInterface interface {
	Info(params ...interface{})
	Error(params ...interface{})
	Debugf(msg string, params ...interface{})
}

// Subscriber for kafka
type Subscriber struct {
	nc.ModelSubscriber

	// Done chanel
	done chan bool

	// Event name of the listening noftification
	eventName string

	// consumer object which receive the messages
	listener *pq.Listener

	// logger interface
	logger loggerInterface
}

// NewListenerSubscriber from *pq.Listener object
func NewListenerSubscriber(listener *pq.Listener, event string, logger loggerInterface) (*Subscriber, error) {
	if logger == nil {
		logger = nlog()
	}
	return &Subscriber{
		done:      make(chan bool, 1),
		listener:  listener,
		eventName: event,
		logger:    logger,
	}, nil
}

// MustListenerSubscriber connection to the postgres "event" or panic
func MustListenerSubscriber(listener *pq.Listener, event string, logger loggerInterface) *Subscriber {
	subscriber, err := NewListenerSubscriber(listener, event, logger)
	if err != nil {
		panic(err)
	}
	return subscriber
}

// NewSubscriber connection to postgres "event" from list of topics
func NewSubscriber(connect, event string, logger loggerInterface) (*Subscriber, error) {
	return NewListenerSubscriber(
		pq.NewListener(connect, 10*time.Second, time.Minute, nil),
		event, logger,
	)
}

// MustSubscriber connection to the postgres "event" or panic
func MustSubscriber(connect, event string, logger loggerInterface) *Subscriber {
	subscriber, err := NewSubscriber(connect, event, logger)
	if err != nil {
		panic(err)
	}
	return subscriber
}

// SetLogger interface
func (s *Subscriber) SetLogger(logger loggerInterface) {
	s.logger = logger
}

///////////////////////////////////////////////////////////////////////////////
/// Methods
///////////////////////////////////////////////////////////////////////////////

// PgListener reader function
func (s *Subscriber) PgListener() *pq.Listener {
	return s.listener
}

// Listen kafka consumer
func (s *Subscriber) Listen(ctx context.Context) (err error) {
	if err = s.listener.Listen(s.eventName); err != nil {
		return err
	}
loop:
	for {
		select {
		case n := <-s.listener.Notify:
			s.logger.Debugf("Received data from channel [%s]", n.Channel)
			err = s.ProcessMessage((*message)(n))
		case <-time.After(90 * time.Second):
			s.logger.Info("Received no events for 90 seconds, checking connection")
			err = s.listener.Ping()
		case <-ctx.Done():
			break loop
		case <-s.done:
			break loop
		}
		if err != nil {
			s.logger.Error(err)
		}
	}
	return err
}

// Close kafka consumer
func (s *Subscriber) Close() error {
	s.done <- true
	return nil
}
