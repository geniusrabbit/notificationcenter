//
// @project geniusrabbit.com 2016 – 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017
//

package notificationcenter

import (
	"fmt"
	"io"
	"sync"
)

var (
	loggers         = map[string]Logger{}
	subscribers     = map[string]Subscriber{}
	closeEvent      = make(chan bool, 1000)
	closeEventCount int
)

// LoggerByName interface
func LoggerByName(name string) Logger {
	l, _ := loggers[name]
	return l
}

// SubscriberByName interface
func SubscriberByName(name string) Subscriber {
	s, _ := subscribers[name]
	return s
}

// Register logger/subscriber
func Register(name string, elem ...interface{}) error {
	for _, el := range elem {
		var registered bool

		if logger, ok := el.(Logger); ok {
			if _, ok := loggers[name]; ok {
				return ErrInterfaceAlreadySubscribed
			}
			loggers[name] = logger
			registered = true
		}

		if subscriber, ok := el.(Subscriber); ok {
			if _, ok := subscribers[name]; ok {
				return ErrInterfaceAlreadySubscribed
			}
			subscribers[name] = subscriber
			registered = true
		}

		if !registered {
			return ErrInvalidObject
		}
	}
	return nil
}

// Unregister logger/subscriber object
func Unregister(elem ...interface{}) error {
	for _, el := range elem {
		for k, i := range loggers {
			if i == el {
				delete(loggers, k)
			}
		}
		for k, i := range subscribers {
			if i == el {
				delete(subscribers, k)
			}
		}
	}
	return nil
}

// UnregisterAllByName logger/subscriber
func UnregisterAllByName(name string, _loggers, _subscribers bool) {
	if i, ok := loggers[name]; ok && _loggers {
		if cl, ok := i.(io.Closer); ok {
			cl.Close()
		}
		delete(loggers, name)
	}
	if i, ok := subscribers[name]; ok && _subscribers {
		if cl, ok := i.(io.Closer); ok {
			cl.Close()
		}
		delete(subscribers, name)
	}
}

// Send message
func Send(name string, msg ...interface{}) error {
	if l, ok := loggers[name]; ok {
		return l.Send(msg...)
	}
	return ErrInvalidObject
}

// Subscribe new handler
func Subscribe(name string, h Handler) error {
	if sub, _ := subscribers[name]; nil != sub {
		return sub.Subscribe(h)
	}
	return fmt.Errorf("Undefined subscriber with name: %s", name)
}

// Unsubscribe this handler by ptr
func Unsubscribe(name string, h Handler) error {
	if sub, _ := subscribers[name]; nil != sub {
		return sub.Unsubscribe(h)
	}
	return fmt.Errorf("Undefined subscriber with name: %s", name)
}

// Listen subscribers
func Listen() (err error) {
	if nil == subscribers {
		return nil
	}

	var w sync.WaitGroup

	for _, sub := range subscribers {
		w.Add(1)
		go func() {
			err = sub.Listen()
			w.Done()
		}()
	}

	w.Wait()
	return
}

// Close notification center
func Close() {
	for _, i := range loggers {
		if cl, ok := i.(io.Closer); ok {
			cl.Close()
		}
	}

	for _, i := range subscribers {
		if cl, ok := i.(io.Closer); ok {
			cl.Close()
		}
	}

	for i := 0; i < closeEventCount; i++ {
		closeEvent <- true
	}
}

// OnClose event
func OnClose() <-chan bool {
	closeEventCount++
	return (<-chan bool)(closeEvent)
}
