package errx

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/status"
	"runtime"
)

func New(msg string) error {
	return NewWithSkip(msg, 1)
}

func Wrap(err error, msg string) error {
	return WrapWithSkip(err, msg, 1)
}

// NewFromStatus new an error from status.Status
func NewFromStatus(s *status.Status) error {
	e := NewWithSkip(s.Message(), 1)
	tmp := e.(Errorx)
	if statusDetails := s.Details(); len(statusDetails) > 0 {
		if details, ok := statusDetails[0].(*Map); ok {
			d := PbMap2MapStrAny(details)
			for key, value := range d {
				tmp[key] = value
			}
		}
	}
	return tmp
}

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

	return Errorx(e)
}

func WrapWithSkip(err error, message string, skip int) error {
	e := NewWithSkip(message, skip+1)
	tmp := e.(Errorx)
	tmp[Err] = err

	switch x := err.(type) {
	case *Errorx:
		tmp[ErrStack] = fmt.Sprintf(
			"%s:%d:%s\n%s",
			tmp[File],
			tmp[Line],
			tmp[Msg],
			x.ErrStack(),
		)
	case Errorx:
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

type Errorx map[string]any

func (e Errorx) Error() string {
	msg := e[Msg].(string)
	if sub, ok := e[Err]; ok {
		switch x := sub.(type) {
		case *Errorx:
			msg += ":" + (*x)[Msg].(string)
		case Errorx:
			msg += ":" + x[Msg].(string)
		default:
			msg += ":" + sub.(error).Error()
		}
	}

	return msg
}

func (e Errorx) ErrStack() string {
	return e[ErrStack].(string)
}

type Fields map[string]any

func (e Errorx) Set(fields Fields) {
	for key, value := range fields {
		for _, item := range reserveKeys {
			if key == item {
				panic(errors.New("reserveKey used"))
			}
		}
		e[key] = value
	}
}

func (e Errorx) Get(field string) (value any, ok bool) {
	value, ok = e[field]
	return
}

func (e Errorx) Format(s fmt.State, c rune) {
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
