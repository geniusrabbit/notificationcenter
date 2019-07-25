package subscriber

import (
	"testing"

	"github.com/geniusrabbit/notificationcenter"
)

type (
	fnkWrapper = notificationcenter.FuncHandler
	message    = notificationcenter.Message
)

func Test_Subscriber(t *testing.T) {
	var (
		handler = fnkWrapper(func(msg message) error { return nil })
		tests   = []struct {
			handler           notificationcenter.Handler
			subscribeResult   error
			unsubscribeResult error
		}{
			{
				handler:           fnkWrapper(func(msg message) error { return nil }),
				subscribeResult:   nil,
				unsubscribeResult: nil,
			},
			{
				handler:           handler,
				subscribeResult:   nil,
				unsubscribeResult: nil,
			},
			{
				handler:           handler,
				subscribeResult:   errHandlerAlreadyRegistered,
				unsubscribeResult: errHandlerIsNotFound,
			},
			{
				handler:           nil,
				subscribeResult:   errInvalidHandler,
				unsubscribeResult: errInvalidHandler,
			},
		}
		base Base
	)

	// Test subscription
	for _, test := range tests {
		if err := base.Subscribe(test.handler); err != test.subscribeResult {
			t.Errorf(`invalid handler subsribtion: "%v" must be "%v"`, err, test.subscribeResult)
		}
	}

	// Test unsubscribtion
	for _, test := range tests {
		if err := base.Unsubscribe(test.handler); err != test.unsubscribeResult {
			t.Errorf(`invalid handler unsubsribtion: "%v" must be "%v"`, err, test.unsubscribeResult)
		}
	}
}
