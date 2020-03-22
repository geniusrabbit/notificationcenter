package notificationcenter

import "bytes"

type multierror struct {
	errs []error
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
