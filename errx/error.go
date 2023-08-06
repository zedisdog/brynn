package errx

import (
	"errors"
	"fmt"
	"github.com/zedisdog/brynn/util/stringsx"
	"runtime"
	"strconv"
	"strings"
)

type Error struct {
	Code    int
	Message string
	err     error
	file    string
	line    int
	Detail  interface{}
}

func (e Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s<=%s", e.Message, e.err.Error())
	} else {
		return e.Message
	}
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		switch {
		case s.Flag('+'):
			_, _ = s.Write([]byte(fmt.Sprintf("%s:%d:%s\n", e.file, e.line, e.Message)))
		case s.Flag('#'):
			strArr := []string{fmt.Sprintf("%+v", e)}

			err := errors.Unwrap(e)
			for {
				if err == nil {
					break
				}

				if _, ok := err.(*Error); ok {
					strArr = append(strArr, fmt.Sprintf("%+v", err))
				} else {
					strArr = append(strArr, err.Error()+"\n")
				}

				err = errors.Unwrap(err)
			}
			_, _ = s.Write([]byte(strings.Join(strArr, "")))
		default:
			_, _ = s.Write([]byte(e.Error()))
		}
	}
}

func WithDetail(detail interface{}) func(err *Error) {
	return func(err *Error) {
		err.Detail = detail
	}
}

// New auto determine the caller and return error with msg.
func New(msg string, opts ...func(err *Error)) error {
	// get call stack, and parse caller in it.
	buf := make([]byte, 10240)
	runtime.Stack(buf, false)
	arr := strings.Split(string(buf), "\n")

	var index int
	for index = 1; index < len(arr); index += 2 {
		if !stringsx.ContainersAny(arr[index], []string{"zedisdog/brynn/errx.New", "zedisdog/brynn/errx.Wrap", "zedisdog/brynn/errx.Make"}) {
			break
		}
	}
	location := strings.Split(strings.Split(strings.Trim(arr[index+1], "\t"), " ")[0], ":")
	line, err := strconv.Atoi(location[1])
	if err != nil {
		panic(err)
	}

	err = &Error{
		Message: msg,
		file:    location[0],
		line:    line,
	}

	for _, set := range opts {
		set(err.(*Error))
	}

	return err
}

func NewWithCode(code int, msg string, opts ...func(err *Error)) error {
	err := New(msg, opts...).(*Error)
	err.Code = code
	return err
}
