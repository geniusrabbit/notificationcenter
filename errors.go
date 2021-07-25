package notificationcenter

import "bytes"

type multiError []error

// MultiError data type
func MultiError(errors ...error) error {
	return multiError(errors)
}

func (e multiError) Error() string {
	if len(e) == 0 {
		return ``
	}
	var buff bytes.Buffer
	for i, err := range e {
		if i > 0 {
			buff.WriteByte('\n')
		}
		buff.Write([]byte(`- `))
		buff.WriteString(err.Error())
	}
	return buff.String()
}

func (e *multiError) Add(err error) {
	if e != nil && err != nil {
		*e = append(*e, err)
	}
}

func (e multiError) AsError() error {
	if len(e) == 0 {
		return nil
	}
	return e
}
