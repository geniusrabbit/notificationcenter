//
// @project geniusrabbit.com 2015 – 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2017
//

package subscriber

import (
	"io"
	"reflect"

	"github.com/geniusrabbit/notificationcenter"
)

// Base struct
type Base struct {
	handlers []notificationcenter.Handler
}

// Subscribe new handler
func (s *Base) Subscribe(h notificationcenter.Handler) error {
	if nil == s.handlers {
		s.handlers = make([]notificationcenter.Handler, 0, 10)
	} else {
		for _, H := range s.handlers {
			if reflect.ValueOf(H).Pointer() == reflect.ValueOf(h).Pointer() {
				return nil
			}
		} // end for
	}
	s.handlers = append(s.handlers, h)
	return nil
}

// Unsubscribe some handler
func (s *Base) Unsubscribe(h notificationcenter.Handler) error {
	if nil != s.handlers {
		for i, H := range s.handlers {
			if reflect.ValueOf(H).Pointer() == reflect.ValueOf(h).Pointer() {
				s.handlers = append(s.handlers[i:], s.handlers[:i+1]...)
				break
			}
		} // end for
	}
	return nil
}

// Handle all handlers
func (s *Base) Handle(item interface{}, stopIfError bool) (err error) {
	if nil == s.handlers {
		return
	}

	for _, h := range s.handlers {
		if e := h.Handle(item); nil != e {
			err = e
			if stopIfError {
				break
			}
		}
	}
	return
}

// CloseAll handlers
func (s *Base) CloseAll() (err error) {
	if nil != s.handlers {
		for _, h := range s.handlers {
			if fn, ok := h.(io.Closer); ok {
				err = fn.Close()
			}
			if nil != err {
				break
			}
		}
	}
	return
}
