package code

import "github.com/zedisdog/brynn/i18n"

const (
	InternalError = 500
	ValidateError = 422
	BadRequest    = 400
)

var codes = map[int]string{
	InternalError: i18n.Trans("internal error"),
	ValidateError: i18n.Trans("validate error"),
	BadRequest:    i18n.Trans("bad request"),
}

func CodeMsg(code int) string {
	return codes[code]
}
