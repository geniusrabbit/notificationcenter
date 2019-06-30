//
// @project GeniusRabbit 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package nats

import (
	"github.com/geniusrabbit/notificationcenter/subscriber"
	nstream "github.com/nats-io/stan.go"
)

// Subscriber for NATS queue
type Subscriber struct {
	subscriber.Base
	group      string
	topics     []string
	options    []nstream.Option
	conn       *nstream.Conn
	closeEvent chan bool
}

// NewSubscriber object
func NewSubscriber(url, clusterID, clientID, group string, topics []string, subOptions []nstream.Option, options ...nstream.Option) (*Subscriber, error) {
	var conn, err = nstream.Connect(clusterID, clientID, url, options...)
	if err != nil || conn == nil {
		return nil, err
	}

	return &Subscriber{
		group:      group,
		topics:     topics,
		options:    subOptions,
		conn:       conn,
		closeEvent: make(chan bool, 1),
	}, nil
}

// MustNewSubscriber object
func MustNewSubscriber(url, clusterID, clientID, group string, topics []string, subOptions []nstream.Option, options ...nstream.Option) *Subscriber {
	var sub, err = NewSubscriber(url, clusterID, clientID, group, topics, subOptions, options...)
	if err != nil || sub == nil {
		panic(err)
	}
	return sub
}

// Listen kafka consumer
func (s *Subscriber) Listen() (_ error) {
	for _, topic := range s.topics {
		if s.group != "" {
			s.conn.QueueSubscribe(topic, s.group, s.message, s.options...)
		} else {
			s.conn.Subscribe(topic, s.message, s.options...)
		}
	}
	<-s.closeEvent
	return
}

// message execute
func (s *Subscriber) message(m *nstream.Msg) {
	s.Handle(m.Data, false)
}

// Close nstream client
func (s *Subscriber) Close() error {
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
		s.closeEvent <- true
	}
	return nil
}
