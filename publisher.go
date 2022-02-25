//
// @project geniusrabbit.com 2015 – 2016, 2019 - 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016, 2019 - 2022
//

package notificationcenter

import (
	"context"

	"go.uber.org/multierr"
)

// Publisher pipeline base declaration
//go:generate mockgen -source $GOFILE -package mocks -destination mocks/publisher.go
type Publisher interface {
	// Publish one or more messages to the pub-service
	Publish(ctx context.Context, messages ...any) error
}

// MultiPublisher wrapper
type MultiPublisher []Publisher

// Publish one or more messages to the banch of pub-services
func (p MultiPublisher) Publish(ctx context.Context, messages ...any) error {
	var errs error
	for _, pub := range p {
		errs = multierr.Append(errs, pub.Publish(ctx, messages...))
	}
	return errs
}

// FuncPublisher provides custom function wrapper for the custom publisher processor
type FuncPublisher func(context.Context, ...any) error

// Publish method call the original custom publisher function
func (f FuncPublisher) Publish(ctx context.Context, messages ...any) error {
	return f(ctx, messages...)
}
