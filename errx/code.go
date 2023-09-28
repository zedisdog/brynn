package errx

import (
	"fmt"
	"github.com/zedisdog/brynn/i18n"
)

type Code int32

func (c Code) Message() (msg string) {
	msg, ok := FindMsg(c)
	if !ok {
		return i18n.Transf("unknown error code [:code]", i18n.P{"code": c})
	}

	return
}

const (
	InternalError Code = 500
	ValidateError Code = 422
)

var message = map[Code]string{
	InternalError: i18n.Trans("internal error"),
}

func RegisterCode(code int32, msg string) {
	if _, ok := message[Code(code)]; ok {
		panic(New(InternalError, fmt.Sprintf("code [%d] is already exists.", code)))
	}
	message[Code(code)] = msg
}

func FindMsg(code Code) (msg string, exists bool) {
	msg, exists = message[code]
	return
}
