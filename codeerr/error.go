package codeerr

import (
	"fmt"
	"github.com/zedisdog/brynn/code"
	"github.com/zedisdog/brynn/errx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const CodeField = "code"

type ErrorWithCode struct {
	errx.Errorx
}

// New e error
func New(code int, message string) error {
	e := errx.NewWithSkip(message, 1).(errx.Errorx)
	e[CodeField] = code
	return &ErrorWithCode{
		Errorx: e,
	}
}

// NewCode new an error with given code.
func NewCode(errCode int) error {
	e := errx.NewWithSkip(code.CodeMsg(errCode), 1).(errx.Errorx)
	e[CodeField] = errCode
	return &ErrorWithCode{
		Errorx: e,
	}
}

// NewMsg new an error with given error message.
func NewMsg(errMsg string) error {
	e := errx.NewWithSkip(errMsg, 1).(errx.Errorx)
	e[CodeField] = code.InternalError
	return &ErrorWithCode{
		Errorx: e,
	}
}

func Wrap(err error, code int, message string) error {
	e := errx.WrapWithSkip(err, message, 1).(errx.Errorx)
	e[CodeField] = code
	return &ErrorWithCode{
		Errorx: e,
	}
}

func WrapCode(err error, c int) error {
	e := errx.WrapWithSkip(err, code.CodeMsg(c), 1).(errx.Errorx)
	e[CodeField] = c
	return &ErrorWithCode{
		Errorx: e,
	}
}

func WrapMsg(err error, message string) error {
	e := errx.WrapWithSkip(err, message, 1).(errx.Errorx)
	e[CodeField] = code.InternalError
	return &ErrorWithCode{
		Errorx: e,
	}
}

func (e ErrorWithCode) GRPCStatus() *status.Status {
	var code codes.Code
	if e.Errorx[CodeField] != nil {
		code = codes.Code(e.Errorx[CodeField].(int))
	} else {
		code = codes.Unknown
	}
	s := status.New(code, e.Error())

	if len(e.Errorx) > 0 {
		var (
			err error
		)
		m := &errx.Map{}
		m.Fields = errx.Map2Pb(e)
		s, err = s.WithDetails(m)
		if err != nil {
			panic(fmt.Errorf("build status failed: %w", err))
		}
	}

	return s
}
