package notificationcenter

import "bytes"

type multierror struct {
	errs []error
}

// MultiError data type
func MultiError(errors ...error) error {
	return &multierror{errs: errors}
}

func (e *multierror) Error() string {
	var buff bytes.Buffer
	for i, err := range e.errs {
		if i > 0 {
			buff.Write([]byte("\n * "))
		}
		buff.WriteString(err.Error())
	}
	return buff.String()
}
