package notificationcenter

import (
	"context"
	"reflect"

	"github.com/geniusrabbit/notificationcenter/decoder"
	"github.com/pkg/errors"
)

var errInvalidReturnType = errors.New("invalid return types")

// ReceiverFrom converts income handler type to Receiver interface
func ReceiverFrom(handler any) Receiver {
	switch h := handler.(type) {
	case Receiver:
		return h
	case func(msg Message) error:
		return FuncReceiver(h)
	default:
		return ExtFuncReceiver(h)
	}
}

var (
	errorType   = reflect.TypeOf((*error)(nil)).Elem()
	contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
	msgType     = reflect.TypeOf((*Message)(nil)).Elem()
)

// ExtFuncReceiver wraps function argument with arbitrary input data type
func ExtFuncReceiver(f any, decs ...decoder.Decoder) Receiver {
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		panic("argument must be a function")
	}
	dec := decoder.JSON
	if len(decs) > 0 && decs[0] != nil {
		dec = decs[0]
	}
	var (
		ft        = fv.Type()
		argMapper = make([]func(Message) (reflect.Value, error), 0, ft.NumIn())
		retMapper = make([]func(reflect.Value) error, 0, ft.NumOut())
	)
	for i := 0; i < ft.NumIn(); i++ {
		inType := ft.In(i)
		switch inType {
		case contextType:
			argMapper = append(argMapper, func(msg Message) (reflect.Value, error) {
				return reflect.ValueOf(msg.Context()), nil
			})
		case msgType:
			argMapper = append(argMapper, func(msg Message) (reflect.Value, error) {
				return reflect.ValueOf(msg), nil
			})
		default:
			argMapper = append(argMapper, func(msg Message) (reflect.Value, error) {
				newValue, newValueI := newValue(inType)
				err := dec(msg.Body(), newValueI)
				return newValue, err
			})
		}
	}
	for i := 0; i < ft.NumOut(); i++ {
		outType := ft.Out(i)
		switch outType {
		case errorType:
			retMapper = append(retMapper, func(v reflect.Value) error {
				if v.IsNil() {
					return nil
				}
				return v.Interface().(error)
			})
		default:
			panic(errInvalidReturnType)
		}
	}
	return FuncReceiver(func(msg Message) error {
		args := make([]reflect.Value, 0, len(argMapper))
		for _, fm := range argMapper {
			arg, err := fm(msg)
			if err != nil {
				return err
			}
			args = append(args, arg)
		}
		retVals := fv.Call(args)
		for i, fr := range retMapper {
			if err := fr(retVals[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

func newValue(t reflect.Type) (reflect.Value, any) {
	if t.Kind() == reflect.Ptr {
		return newValue(t.Elem())
	}
	v := reflect.New(t)
	i := v.Interface()
	return v, i
}
