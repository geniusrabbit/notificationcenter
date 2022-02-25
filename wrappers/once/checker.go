package once

// Checker provides test interface of messages
type Checker interface {
	IsSkip(msg any) bool
	MarkAsSent(msg any) error
}
