//
// @project geniusrabbit.com 2015 – 2016, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016, 2020
//

package notificationcenter

import (
	"context"
	"io"
	"reflect"

	"github.com/pkg/errors"
)

var (
	errInvalidReceiver           = errors.New("invalid receiver value")
	errReceiverAlreadyRegistered = errors.New("receiver already registered")
)

// Subscriber provides methods of working with subscribtion
//go:generate mockgen -source $GOFILE -package mocks -destination mocks/subscriber.go
type Subscriber interface {
	io.Closer

	// Subscribe new receiver to receive messages from the subsribtion
	Subscribe(ctx context.Context, receiver Receiver) error

	// Start processing queue
	Listen(ctx context.Context) error
}

// ModelSubscriber provedes subscibe functionality implementation
type ModelSubscriber struct {
	receivers []Receiver

	// Error handler pointer
	ErrorHandler ErrorHandler

	// Panic handler pointer
	PanicHandler PanicHandler
}

// Subscribe new receiver to receive messages from the subsribtion
func (s *ModelSubscriber) Subscribe(ctx context.Context, receiver Receiver) error {
	if receiver == nil {
		return errInvalidReceiver
	}
	for _, current := range s.receivers {
		if reflect.ValueOf(current).Pointer() == reflect.ValueOf(receiver).Pointer() {
			return errReceiverAlreadyRegistered
		}
	}
	s.receivers = append(s.receivers, receiver)
	return nil
}

// ProcessMessage by all receivers
func (s *ModelSubscriber) ProcessMessage(msg Message) error {
	defer func() {
		if s.PanicHandler == nil {
			return
		}
		if rec := recover(); rec != nil {
			s.PanicHandler(msg, rec)
		}
	}()
	for _, r := range s.receivers {
		if err := r.Receive(msg); err != nil {
			if s.ErrorHandler == nil {
				return err
			}
			s.ErrorHandler(msg, err)
		}
	}
	return nil
}

// Close all receivers if supports io.Closer interface
func (s *ModelSubscriber) Close() error {
	for _, rec := range s.receivers {
		if closer, _ := rec.(io.Closer); closer != nil {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}
