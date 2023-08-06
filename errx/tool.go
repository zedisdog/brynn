package errx

import (
	"fmt"
)

func Wrap(err error, msg string) error {
	return WrapWithCode(err, msg, 0)
}

func WrapWithCode(err error, msg string, code int) error {
	if err == nil {
		return nil
	}
	e := New(msg).(*Error)
	e.err = err
	e.Code = code
	return e
}

func CodeWrapperWithPrefix(prefix string) func(error, string, int) error {
	return func(err error, s string, code int) error {
		return WrapWithCode(err, fmt.Sprintf("[%s]%s", prefix, s), code)
	}
}
