//
// @project GeniusRabbit 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package nats

import (
	"github.com/geniusrabbit/notificationcenter/subscriber"
	"github.com/nats-io/nats"
)

// Subscriber for NATS queue
type Subscriber struct {
	subscriber.Base
	topics []string
	conn   *nats.Conn
}

// NewSubscriber object
func NewSubscriber(topics []string, url string, options ...nats.Option) (*Subscriber, error) {
	var conn, err = nats.Connect(url, options...)
	if nil != err || nil == conn {
		return nil, err
	}

	return &Log{topics: topics, conn: conn}, nil
}

// MustNewSubscriber object
func MustNewSubscriber(topics []string, url string, options ...nats.Option) *Subscriber {
	var sub, err = NewSubscriber(topics, url, options...)
	if nil != err || nil == sub {
		panic(err)
	}
	return sub
}

// Listen kafka consumer
func (s *Subscriber) Listen() (_ error) {
	for _, t := range s.topics {
		s.conn.Subscribe(t, s.message)
	}
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
	}
}
