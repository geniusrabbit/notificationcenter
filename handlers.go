package notificationcenter

// ErrorHandler type to process error value
type ErrorHandler func(msg Message, err error)

// PanicHandler type to process panic action
type PanicHandler func(msg Message, recoverData interface{})
