//
// @project geniusrabbit.com 2015 – 2016, 2019 - 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016, 2019 - 2020
//

package notificationcenter

import "context"

// Publisher pipeline base declaration
//go:generate mockgen -source $GOFILE -package mocks -destination mocks/publisher.go
type Publisher interface {
	// Publish one or more messages to the pub-service
	Publish(ctx context.Context, messages ...interface{}) error
}
