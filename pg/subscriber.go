//
// @project geniusrabbit.com 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2019
//

package pg

import (
	"time"

	"github.com/geniusrabbit/notificationcenter/subscriber"
	"github.com/lib/pq"
)

type loggerInterface interface {
	Info(params ...interface{})
	Error(params ...interface{})
	Debugf(msg string, params ...interface{})
}

// Subscriber for kafka
type Subscriber struct {
	subscriber.Base

	// Done chanel
	done chan bool

	// Event name of the listening noftification
	eventName string

	// consumer object which receive the messages
	listener *pq.Listener

	// logger interface
	logger loggerInterface
}

// NewSubscriber connection to kafka "group" from list of topics
func NewSubscriber(connect, event string, logger loggerInterface) (*Subscriber, error) {
	if logger == nil {
		logger = nlog()
	}
	return &Subscriber{
		done:      make(chan bool, 1),
		listener:  pq.NewListener(connect, 10*time.Second, time.Minute, nil),
		eventName: event,
		logger:    logger,
	}, nil
}

// SetLogger interface
func (s *Subscriber) SetLogger(logger loggerInterface) {
	s.logger = logger
}

///////////////////////////////////////////////////////////////////////////////
/// Methods
///////////////////////////////////////////////////////////////////////////////

// Listen kafka consumer
func (s *Subscriber) Listen() (err error) {
	if err = s.listener.Listen(s.eventName); err != nil {
		return err
	}
	for {
		select {
		case n := <-s.listener.Notify:
			s.logger.Debugf("Received data from channel [%s]", n.Channel)
			err = s.Handle((*message)(n), true)
		case <-time.After(90 * time.Second):
			s.logger.Info("Received no events for 90 seconds, checking connection")
			err = s.listener.Ping()
		case <-s.done:
			return nil
		}
	}
}

// Close kafka consumer
func (s *Subscriber) Close() (err error) {
	s.done <- true
	return s.CloseAll()
}
