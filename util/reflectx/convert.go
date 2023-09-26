package reflectx

import (
	"github.com/zedisdog/brynn/errx"
	"github.com/zedisdog/brynn/i18n"
	"reflect"
)

// func ConvertTo(src reflect.Value, kind reflect.Kind) (reflect.Value, error) {

// }
// TODO: 所有转换所有
func ConvertBoolTo[T any](src reflect.Value) (result T, err error) {
	s := src.Interface().(bool)
	var r int8
	if s {
		r = 1
	} else {
		r = 0
	}
	switch any(result).(type) {
	case int:
		result = any(int(r)).(T)
	case int8:
		result = any(int8(r)).(T)
	case int16:
		result = any(int16(r)).(T)
	case int32:
		result = any(int32(r)).(T)
	case int64:
		result = any(int64(r)).(T)
	case uint64:
		result = any(uint64(r)).(T)
	case uint32:
		result = any(uint32(r)).(T)
	case uint16:
		result = any(uint16(r)).(T)
	case uint8:
		result = any(uint8(r)).(T)
	case uint:
		result = any(uint(r)).(T)
	case bool:
		result = any(s).(T)
	case string:
		if s {
			result = any("true").(T)
		} else {
			result = any("false").(T)
		}
	default:
		err = errx.New(errx.InternalError, i18n.Trans("unsupported type"))
	}

	return
}
