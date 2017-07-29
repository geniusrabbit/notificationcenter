//
// @project GeniusRabbit 2016 – 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017
//

package nats

import (
	"github.com/nats-io/nats"

	"github.com/geniusrabbit/notificationcenter/subscriber"
)

// Subscriber for NATS queue
type Subscriber struct {
	subscriber.Base
	group      string
	topics     []string
	conn       *nats.Conn
	closeEvent chan bool
}

// NewSubscriber object
func NewSubscriber(url, group string, topics []string, options ...nats.Option) (*Subscriber, error) {
	var conn, err = nats.Connect(url, options...)
	if nil != err || nil == conn {
		return nil, err
	}

	return &Subscriber{
		group:      group,
		topics:     topics,
		conn:       conn,
		closeEvent: make(chan bool, 1),
	}, nil
}

// MustNewSubscriber object
func MustNewSubscriber(url, group string, topics []string, options ...nats.Option) *Subscriber {
	var sub, err = NewSubscriber(url, group, topics, options...)
	if nil != err || nil == sub {
		panic(err)
	}
	return sub
}

// Listen kafka consumer
func (s *Subscriber) Listen() (_ error) {
	for _, topic := range s.topics {
		if s.group != "" {
			s.conn.QueueSubscribe(topic, s.group, s.message)
		} else {
			s.conn.Subscribe(topic, s.message)
		}
	}
	<-s.closeEvent
	return
}

// message execute
func (s *Subscriber) message(m *nats.Msg) {
	s.Handle(m.Data, false)
}

// Close nats client
func (s *Subscriber) Close() error {
	if nil != s.conn {
		s.conn.Close()
		s.conn = nil
		s.closeEvent <- true
	}
	return nil
}
