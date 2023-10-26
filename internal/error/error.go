package error

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
)

func NewWithSkip(msg string, skip int) error {
	e := make(map[string]any, 10)
	e[Msg] = msg

	_, file, line, _ := runtime.Caller(1 + skip)
	e[File] = file
	e[Line] = line

	e[ErrStack] = fmt.Sprintf(
		"%s:%d:%s\n",
		file,
		line,
		msg,
	)

	return Error(e)
}

func WrapWithSkip(err error, message string, skip int) error {
	e := NewWithSkip(message, skip+1)
	tmp := e.(Error)
	tmp[Err] = err

	switch x := err.(type) {
	case *Error:
		tmp[ErrStack] = fmt.Sprintf(
			"%s:%d:%s\n%s",
			tmp[File],
			tmp[Line],
			tmp[Msg],
			x.ErrStack(),
		)
	case Error:
		tmp[ErrStack] = fmt.Sprintf(
			"%s:%d:%s\n%s",
			tmp[File],
			tmp[Line],
			tmp[Msg],
			x.ErrStack(),
		)
	default:
		tmp[ErrStack] = fmt.Sprintf(
			"%s:%d:%s\n%s",
			tmp[File],
			tmp[Line],
			tmp[Msg],
			x.(error).Error(),
		)
	}

	return tmp
}

type Error map[string]any

func (e Error) Error() string {
	msg := e[Msg].(string)
	if sub, ok := e[Err]; ok {
		switch x := sub.(type) {
		case *Error:
			msg += ":" + (*x)[Msg].(string)
		case Error:
			msg += ":" + x[Msg].(string)
		default:
			msg += ":" + sub.(error).Error()
		}
	}

	return msg
}

func (e Error) ErrStack() string {
	return e[ErrStack].(string)
}

type Fields map[string]any

func (e Error) Set(fields Fields) {
	for key, value := range fields {
		for _, item := range reserveKeys {
			if key == item {
				panic(errors.New("reserveKey used"))
			}
		}
		e[key] = value
	}
}

func (e Error) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		switch {
		case s.Flag('+'):
			_, _ = s.Write([]byte(fmt.Sprintf("%s:%d:%s\n", e[File], e[Line], e[Msg])))
		case s.Flag('#'):
			_, _ = s.Write([]byte(e[ErrStack].(string)))
		default:
			if e[Err] != nil {
				_, _ = s.Write([]byte(fmt.Sprintf("%s]<=[%v\n", e[Msg], e[Err])))
			} else {
				_, _ = s.Write([]byte(e[Msg].(string) + "\n"))
			}
		}
	}
}

func (e Error) GRPCStatus() *status.Status {
	var code codes.Code
	if e[Code] != nil {
		code = codes.Code(e[Code].(int))
	} else {
		code = codes.Unknown
	}
	s := status.New(code, e.Error())

	if len(e) > 0 {
		var (
			err error
		)
		m := &Map{}
		m.Fields = Map2Pb(e)
		s, err = s.WithDetails(m)
		if err != nil {
			panic(fmt.Errorf("build status failed: %w", err))
		}
	}

	return s
}
