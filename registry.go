package notificationcenter

import (
	"context"
	"io"
	"log"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

// Error list...
var (
	ErrInvalidRegisterParameter     = errors.New(`invalid register parameter`)
	ErrUndefinedPublisherInterface  = errors.New(`undefined publisher interface`)
	ErrUndefinedSubscriberInterface = errors.New(`undefined subscriber interface`)
	ErrInterfaceAlreadySubscribed   = errors.New("[notificationcenter] interface already subscribed")
)

// Registry provides functionality of access to pub/sub
// interfaces by string names.
type Registry struct {
	mx              sync.RWMutex
	closeEventCount int32
	closeEvent      chan bool
	publishers      map[string]Publisher
	subscribers     map[string]Subscriber
}

// NewRegistry init new registry object
func NewRegistry() *Registry {
	return &Registry{
		closeEvent:  make(chan bool),
		publishers:  map[string]Publisher{},
		subscribers: map[string]Subscriber{},
	}
}

// Publisher returns pub interface by name if exists or Nil otherwise
func (r *Registry) Publisher(name string) Publisher {
	r.mx.RLock()
	defer r.mx.RUnlock()
	return r.publishers[name]
}

// Subscriber returns sub interface by name if exists or Nil otherwise
func (r *Registry) Subscriber(name string) Subscriber {
	r.mx.RLock()
	defer r.mx.RUnlock()
	return r.subscribers[name]
}

// Register one or more Publisher or Subscriber services.
// As input parameters must be order of parameters {Name, interface}
//
// Example:
// ```
//   nc.Register(
//     "events", kafka.MustNewSubscriber(),
//     "notifications", nats.MustNewSubscriber(),
//   )
// ```
func (r *Registry) Register(params ...interface{}) error {
	var name string
	r.mx.Lock()
	defer r.mx.Unlock()
	for i, p := range params {
		switch v := p.(type) {
		case string:
			if name != "" {
				return errors.Wrap(ErrInvalidRegisterParameter, "pubsub-name")
			}
			name = v
		case Publisher:
			if _, ok := r.publishers[name]; ok {
				return errors.Wrap(ErrInterfaceAlreadySubscribed, name)
			}
			r.publishers[name] = v
			name = ""
		case Subscriber:
			if _, ok := r.subscribers[name]; ok {
				return errors.Wrap(ErrInterfaceAlreadySubscribed, name)
			}
			r.subscribers[name] = v
			name = ""
		default:
			return errors.Wrapf(ErrInvalidRegisterParameter, "%d:%t", i, v)
		}
	}
	return nil
}

// Publish one or more messages to the pub-service
func (r *Registry) Publish(ctx context.Context, name string, messages ...interface{}) error {
	pub := r.Publisher(name)
	if pub == nil {
		return errors.Wrapf(ErrUndefinedPublisherInterface, name)
	}
	return pub.Publish(ctx, messages...)
}

// Subscribe new handler on some paticular subscriber interface by name
func (r *Registry) Subscribe(ctx context.Context, name string, receiver Receiver) error {
	if sub := r.Subscriber(name); sub != nil {
		return sub.Subscribe(ctx, receiver)
	}
	return errors.Wrap(ErrUndefinedSubscriberInterface, name)
}

// Listen runs subscribers listen interface
func (r *Registry) Listen(ctx context.Context) (err error) {
	var w sync.WaitGroup
	r.mx.RLock()
	for name, sub := range r.subscribers {
		w.Add(1)
		go func(ctx context.Context, name string, sub Subscriber) {
			defer w.Done()
			if err = sub.Listen(ctx); err != nil {
				log.Println(name, err)
			}
		}(ctx, name, sub)
	}
	r.mx.RUnlock()
	w.Wait()
	return err
}

// Close notification center
func (r *Registry) Close() error {
	var errors []error
	r.mx.Lock()
	defer r.mx.Unlock()
	for _, pub := range r.publishers {
		if cl, ok := pub.(io.Closer); ok {
			if err := cl.Close(); err != nil {
				errors = append(errors, err)
			}
		}
	}
	for _, sub := range r.subscribers {
		if cl, ok := sub.(io.Closer); ok {
			if err := cl.Close(); err != nil {
				errors = append(errors, err)
			}
		}
	}
	for i := int32(0); i < atomic.LoadInt32(&r.closeEventCount); i++ {
		r.closeEvent <- true
	}
	if len(errors) > 0 {
		return MultiError(errors...)
	}
	return nil
}

// OnClose event will be executed only after closing all interfaces
//
// Usecases in the application makes subsribing for the finishing event very convinient
//
// ```go
// func myDatabaseObserver() {
//   <- notificationcenter.OnClose()
//   // ... Do something
// }
// ```
func (r *Registry) OnClose() <-chan bool {
	atomic.AddInt32(&r.closeEventCount, 1)
	return (<-chan bool)(r.closeEvent)
}
