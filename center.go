//
// @project geniusrabbit.com 2016 – 2017, 2019 - 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017, 2019 - 2020
//

package notificationcenter

import (
	"context"
)

// DefaultRegistry is global registry
var DefaultRegistry = NewRegistry()

// PublisherByName returns pub interface by name if exists or Nil otherwise
func PublisherByName(name string) Publisher {
	return DefaultRegistry.Publisher(name)
}

// SubscriberByName returns sub interface by name if exists or Nil otherwise
func SubscriberByName(name string) Subscriber {
	return DefaultRegistry.Subscriber(name)
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
func Register(params ...interface{}) error {
	return DefaultRegistry.Register(params...)
}

// Publish one or more messages to the pub-service
func Publish(ctx context.Context, name string, messages ...interface{}) error {
	return DefaultRegistry.Publish(ctx, name, messages...)
}

// Subscribe new handler on some paticular subscriber interface by name
func Subscribe(ctx context.Context, name string, receiver Receiver) error {
	return DefaultRegistry.Subscribe(ctx, name, receiver)
}

// Listen runs subscribers listen interface
func Listen(ctx context.Context) error {
	return DefaultRegistry.Listen(ctx)
}

// Close notification center
func Close() error {
	return DefaultRegistry.Close()
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
func OnClose() <-chan bool {
	return DefaultRegistry.OnClose()
}
