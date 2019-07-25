//
// @project geniusrabbit.com 2016, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016, 2019
//

package simple

import (
	"bytes"
	"log"

	"github.com/geniusrabbit/notificationcenter/encoder"
	"github.com/geniusrabbit/notificationcenter/subscriber"
)

// Simple provides permetive chanel implementation of message conveyer
type Simple struct {
	subscriber.Base
	encoder encoder.Encoder
	chanel  chan interface{}
}

// NewSimple event processor
func NewSimple(enc ...encoder.Encoder) *Simple {
	simple := &Simple{
		chanel: make(chan interface{}, 1000),
	}
	if len(enc) > 0 && enc[0] != nil {
		simple.encoder = enc[0]
	} else {
		simple.encoder = encoder.JSON
	}
	return simple
}

// Send messages
func (s *Simple) Send(messages ...interface{}) error {
	for _, m := range messages {
		s.chanel <- m
	}
	return nil
}

// Listen chanel and pricess messages with registered handlers
func (s *Simple) Listen() error {
	for {
		m, ok := <-s.chanel
		if !ok {
			break
		}
		if msg, err := s.messageFrom(m); err != nil {
			log.Print(err)
		} else if err = s.Handle(msg, true); err != nil {
			log.Print(err)
		}
	} // end for
	return nil
}

// Close simple listener
func (s *Simple) Close() error {
	err := s.Base.CloseAll()
	close(s.chanel)
	return err
}

func (s *Simple) messageFrom(m interface{}) (*message, error) {
	var buff bytes.Buffer
	if err := s.encoder(m, &buff); err != nil {
		return nil, err
	}
	return &message{data: buff.Bytes()}, nil
}
