//
// @project GeniusRabbit 2016 - 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2017
//

package simple

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/geniusrabbit/notificationcenter"
)

// Dummy struct
type Dummy struct {
}

// NewDummy object
func NewDummy() *Dummy {
	return &Dummy{}
}

// Send message
func (e *Dummy) Send(messages ...interface{}) (err error) {
	for _, m := range messages {
		switch msg := m.(type) {
		case map[string]func() error:
			for bucket, handler := range msg {
				t := time.Now()
				err = handler()
				log.Println("[DUMMY]", bucket, time.Since(t), err)
			}
		default:
			log.Println("[DUMMY]", strings.TrimSpace(fmt.Sprintln(messages...)))
		}
	}
	return nil
}

// Subscribe new handler
func (e *Dummy) Subscribe(h notificationcenter.Handler) error {
	return nil
}

// Unsubscribe this handler by ptr
func (e *Dummy) Unsubscribe(h notificationcenter.Handler) error {
	return nil
}

// Listen process
func (e *Dummy) Listen() error {
	return nil
}

// Close proxy handler
func (e *Dummy) Close() error {
	return nil
}
