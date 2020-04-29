package once

// Checker provides test interface of messages
type Checker interface {
	IsSkip(msg interface{}) bool
	MarkAsSent(msg interface{}) error
}
