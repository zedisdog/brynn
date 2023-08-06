package helper

func Ternary[T any](e bool, t T, f T) T {
	if e {
		return t
	} else {
		return f
	}
}
