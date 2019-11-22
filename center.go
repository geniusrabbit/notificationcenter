//
// @project geniusrabbit.com 2016 – 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017, 2019
//

package notificationcenter

import (
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"
)

var (
	streamers       = map[string]Streamer{}
	subscribers     = map[string]Subscriber{}
	closeEvent      = make(chan bool, 10)
	closeEventCount int32
)

// StreamByName returns streamer interface by codename if exists
func StreamByName(name string) Streamer {
	l, _ := streamers[name]
	return l
}

// SubscriberByName interface
func SubscriberByName(name string) Subscriber {
	s, _ := subscribers[name]
	return s
}

// Register the new streamer or subscriber interface to future using
// Object can implement both interface and process all events from start to end
func Register(name string, elem ...interface{}) error {
	for _, el := range elem {
		var registered bool

		if logger, ok := el.(Streamer); ok {
			if _, ok := streamers[name]; ok {
				return ErrInterfaceAlreadySubscribed
			}
			streamers[name] = logger
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

// Unregister exist streamer or subscriber interface from the global storage
func Unregister(elem ...interface{}) error {
	for _, el := range elem {
		for k, i := range streamers {
			if i == el {
				delete(streamers, k)
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

// UnregisterAllByName exist streamer or subscriber interface from the global storage
func UnregisterAllByName(name string, _streamers, _subscribers bool) {
	if i, ok := streamers[name]; ok && _streamers {
		if cl, ok := i.(io.Closer); ok {
			cl.Close()
		}
		delete(streamers, name)
	}
	if i, ok := subscribers[name]; ok && _subscribers {
		if cl, ok := i.(io.Closer); ok {
			cl.Close()
		}
		delete(subscribers, name)
	}
}

// Send message object to the stream by name
func Send(name string, msg ...interface{}) error {
	if l, ok := streamers[name]; ok {
		return l.Send(msg...)
	}
	return ErrInvalidObject
}

// Subscribe new handler on some paticular subscriber interface by name
func Subscribe(name string, h Handler) error {
	if sub, _ := subscribers[name]; sub != nil {
		return sub.Subscribe(h)
	}
	return fmt.Errorf("Undefined subscriber with name: %s", name)
}

// Unsubscribe some paticular handler interface from subscriber with the *name*
func Unsubscribe(name string, h Handler) error {
	if sub, _ := subscribers[name]; nil != sub {
		return sub.Unsubscribe(h)
	}
	return fmt.Errorf("Undefined subscriber with name: %s", name)
}

// Listen runs subscribers listen interface
func Listen() (err error) {
	if subscribers == nil {
		return nil
	}
	var w sync.WaitGroup
	for name, sub := range subscribers {
		w.Add(1)
		go func(name string, sub Subscriber) {
			if err = sub.Listen(); err != nil {
				log.Println(name, err)
			}
			w.Done()
		}(name, sub)
	}
	w.Wait()
	return err
}

// Close notification center
func Close() {
	for _, i := range streamers {
		if cl, ok := i.(io.Closer); ok {
			cl.Close()
		}
	}
	for _, i := range subscribers {
		if cl, ok := i.(io.Closer); ok {
			cl.Close()
		}
	}
	for i := int32(0); i < atomic.LoadInt32(&closeEventCount); i++ {
		closeEvent <- true
	}
}

// OnClose event will be executed only after closing all interfaces
//
// Usecases in the application makes subsribing for the finishing event very convinient
// ```go
// func myDatabaseObserver() {
//   <- notificationcenter.OnClose()
//   // ... Do something
// }
// ```
func OnClose() <-chan bool {
	atomic.AddInt32(&closeEventCount, 1)
	return (<-chan bool)(closeEvent)
}
