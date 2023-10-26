package error

const (
	Msg      = "msg"
	Err      = "err"
	File     = "file"
	Line     = "line"
	ErrStack = "errStack"
	Code     = "code"
)

var reserveKeys = []string{
	Msg,
	Err,
	File,
	Line,
	ErrStack,
}

func GetField[T any](e error, key string) T {
	switch x := e.(type) {
	case *Error:
		return (*x)[key].(T)
	case Error:
		return x[key].(T)
	default:
		panic("unsupported")
	}
}
