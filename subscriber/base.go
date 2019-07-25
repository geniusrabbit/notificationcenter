//
// @project geniusrabbit.com 2015 – 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2017, 2019
//

package subscriber

import (
	"errors"
	"io"
	"reflect"

	"github.com/geniusrabbit/notificationcenter"
)

var (
	errHandlerIsNotFound        = errors.New("[subscriber.base] handler is not found")
	errInvalidHandler           = errors.New("[subscribe.base] invalid handler")
	errHandlerAlreadyRegistered = errors.New("[subscribe.base] handler already registered")
)

// Base struct
type Base struct {
	handlers []notificationcenter.Handler
}

// Subscribe new handler
func (s *Base) Subscribe(h notificationcenter.Handler) error {
	if h == nil {
		return errInvalidHandler
	}
	for _, current := range s.handlers {
		if reflect.ValueOf(current).Pointer() == reflect.ValueOf(h).Pointer() {
			return errHandlerAlreadyRegistered
		}
	}
	s.handlers = append(s.handlers, h)
	return nil
}

// Unsubscribe some handler
func (s *Base) Unsubscribe(h notificationcenter.Handler) (err error) {
	if h == nil {
		return errInvalidHandler
	}
	for i, current := range s.handlers {
		if reflect.ValueOf(current).Pointer() == reflect.ValueOf(h).Pointer() {
			s.handlers = append(s.handlers[:i], s.handlers[i+1:]...)
			return nil
		}
	}
	return errHandlerIsNotFound
}

// Handle process the message with list of handlers
func (s *Base) Handle(msg notificationcenter.Message, stopIfError bool) (err error) {
	for _, h := range s.handlers {
		if e := h.Handle(msg); e != nil {
			if err = e; stopIfError {
				break
			}
		}
	}
	return
}

// CloseAll handlers if io.Closer interface is implemented
func (s *Base) CloseAll() (err error) {
	for _, h := range s.handlers {
		if fn, ok := h.(io.Closer); ok {
			err = fn.Close()
		}
		if err != nil {
			break
		}
	}
	return
}
