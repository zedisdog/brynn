package errx

import (
	"fmt"
	"runtime"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func wrap(err error, code Code, message string, skip int) error {
	e := NewWithPos(code, message, skip+1)
	tmp := e.(*Error)
	tmp.err = err
	if oriErr, ok := err.(*Error); ok && oriErr.Details["stack"] != "" {
		tmp.Details["stack"] = fmt.Sprintf(
			"%s%s",
			tmp.Details["stack"],
			oriErr.Details["stack"],
		)
	} else {
		tmp.Details["stack"] = fmt.Sprintf(
			"%s%s",
			tmp.Details["stack"],
			strings.Trim(err.Error(), "\n")+"\n",
		)
	}
	return e
}

func Wrap(err error, code Code, message string) error {
	return wrap(err, code, message, 1)
}

func WrapCode(err error, code Code) error {
	return wrap(err, code, code.Message(), 1)
}

func WrapMsg(err error, message string) error {
	return wrap(err, InternalError, message, 1)
}

// NewFromStatus new an error from status.Status
func NewFromStatus(s *status.Status) error {
	e := NewWithPos(Code(s.Code()), s.Message(), 1)
	if d := s.Details(); len(d) > 0 {
		if details, ok := d[0].(*Map); ok {
			var err error
			e.(*Error).Details, err = PbMap2MapStrAny(details)
			if err != nil {
				panic(WrapMsg(err, "convert grpc status failed"))
			}
		}
	}
	return e
}

// New e error
func New(code Code, message string) error {
	return NewWithPos(code, message, 1)
}

// NewCode new an error with given code.
// if code has been predefined, it'll find error message auto.
// if not, it'll use the zero value of string.
func NewCode(errCode Code) error {
	return NewWithPos(errCode, errCode.Message(), 1)
}

// NewMsg new an error with given error message.
// it'll fill field Code with UNKNOWN_ERROR.
func NewMsg(errMsg string) error {
	return NewWithPos(InternalError, errMsg, 1)
}

// NewWithPos new an error with caller info(file and line).
func NewWithPos(code Code, message string, skip int) error {
	e := &Error{
		Code:    code,
		Msg:     message,
		Details: make(map[string]any),
	}
	_, file, line, _ := runtime.Caller(1 + skip)
	e.Details["file"] = file
	e.Details["line"] = line
	e.Details["stack"] = fmt.Sprintf(
		"%s:%d:%s\n",
		file,
		line,
		message,
	)

	return e
}

type ErrorDetails map[string]any

type Error struct {
	Code    Code
	Msg     string
	Details ErrorDetails
	err     error
}

func (e *Error) WithDetails(details ErrorDetails) error {
	ori := e.Details
	e.Details = details
	e.Details["file"] = ori["file"]
	e.Details["line"] = ori["line"]
	e.Details["stack"] = ori["stack"]
	return e
}

func (e Error) Error() string {
	return e.Msg
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) GRPCStatus() *status.Status {
	s := status.New(codes.Code(e.Code), e.Error())

	if len(e.Details) > 0 {
		var (
			err error
		)
		if err != nil { //如果构建status报错，就直接使用报错信息构建status
			return WrapMsg(err, "build status failed").(*Error).GRPCStatus()
		}
		m := &Map{}
		m.Fields, err = Map2Pb(e.Details)
		if err != nil { //如果构建status报错，就直接使用报错信息构建status
			return WrapMsg(err, "build status failed").(*Error).GRPCStatus()
		}
		s, err = s.WithDetails(m)
		if err != nil { //如果构建status报错，就直接使用报错信息构建status
			return WrapMsg(err, "build status failed").(*Error).GRPCStatus()
		}
	}

	return s
}

// Format 实现Format接口来在打印error时展示更详细的信息,记录了调用栈
func (e Error) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		switch {
		case s.Flag('+'):
			_, _ = s.Write([]byte(fmt.Sprintf("%s:%d:%s\n", e.Details["file"], e.Details["line"], e.Msg)))
		case s.Flag('#'):
			_, _ = s.Write([]byte(e.Details["stack"].(string)))
		default:
			if e.err != nil {
				_, _ = s.Write([]byte(fmt.Sprintf("%s]<=[%v", e.Msg, e.err)))
			} else {
				_, _ = s.Write([]byte(e.Msg))
			}
		}
	}
}
