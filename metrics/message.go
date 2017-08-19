//
// @project geniusrabbit::notificationcenter 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package metrics

import "github.com/demdxx/gocast"

// Message type
const (
	MessageTypeUndefined = iota
	MessageTypeIncrement
	MessageTypeGauge
	MessageTypeTiming
	MessageTypeCount
	MessageTypeUnique
)

// Message optimised for metric use
type Message struct {
	Name  string            `json:"name,omitempty"`
	Type  int               `json:"type,omitempty"`
	Tags  map[string]string `json:"tags,omitempty"`
	Value interface{}       `json:"value,omitempty"`
}

// ValueInt type
func (m *Message) ValueInt() int {
	return gocast.ToInt(m.Value)
}

// ValueString type
func (m *Message) ValueString() string {
	return gocast.ToString(m.Value)
}
