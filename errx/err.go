package errx

const (
	Msg      = "msg"
	Err      = "err"
	File     = "file"
	Line     = "line"
	ErrStack = "errStack"
)

var reserveKeys = []string{
	Msg,
	Err,
	File,
	Line,
	ErrStack,
}
