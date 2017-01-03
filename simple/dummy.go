//
// @project GeniusRabbit 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package simple

import (
	"fmt"
	"log"
	"strings"

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
func (e *Dummy) Send(msg ...interface{}) error {
	log.Println("[DUMMY]", strings.TrimSpace(fmt.Sprintln(msg...)))
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
