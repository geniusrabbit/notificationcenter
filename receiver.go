//
// @project geniusrabbit.com 2015 – 2016, 2019 - 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016, 2019 - 2022
//

package notificationcenter

// Receiver describe interface of message processing
//go:generate mockgen -source $GOFILE -package mocks -destination mocks/receiver.go
type Receiver interface {
	Receive(msg Message) error
}

// FuncReceiver implements Receiver interface for a single function
type FuncReceiver func(msg Message) error

// Receive message from sub-service to process it with function
func (f FuncReceiver) Receive(msg Message) error {
	return f(msg)
}
