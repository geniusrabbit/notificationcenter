//
// @project geniusrabbit.com 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package simple

import (
	"log"

	"github.com/geniusrabbit/notificationcenter/subscriber"
)

type Simple struct {
	subscriber.Base
	chanel chan interface{}
}

func NewSimple() *Simple {
	return &Simple{
		chanel: make(chan interface{}, 1000),
	}
}

// Send messages
func (s *Simple) Send(messages ...interface{}) error {
	for _, m := range messages {
		s.chanel <- m
	}
	return nil
}

// Listen process
func (s *Simple) Listen() error {
	for {
		if m, ok := <-s.chanel; ok {
			if err := s.Handle(m, true); nil != err {
				log.Print(err)
			}
		} else {
			break
		}
	} // end for

	return nil
}

// Close listener
func (s *Simple) Close() error {
	s.Base.CloseAll()
	if nil == s.chanel {
		close(s.chanel)
		s.chanel = nil
	}
	return nil
}
