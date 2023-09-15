package conv

import (
	"fmt"
	"strconv"
	"strings"
)

func StrTo[T any](str string) (result T, err error) {
	switch any(result).(type) {
	case bool:
		switch strings.ToLower(str) {
		case "0", "false":
			result = any(false).(T)
		default:
			result = any(true).(T)
		}
	case int:
		var intVal int64
		intVal, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			err = fmt.Errorf("the value %q cannot be parsed as int", str)
		}
		result = any(int(intVal)).(T)
	case int8:
		var intVal int64
		intVal, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			err = fmt.Errorf("the value %q cannot be parsed as int", str)
		}
		result = any(int8(intVal)).(T)
	case int16:
		var intVal int64
		intVal, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			err = fmt.Errorf("the value %q cannot be parsed as int", str)
		}
		result = any(int16(intVal)).(T)
	case int32:
		var intVal int64
		intVal, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			err = fmt.Errorf("the value %q cannot be parsed as int", str)
		}
		result = any(int32(intVal)).(T)
	case int64:
		var intVal int64
		intVal, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			err = fmt.Errorf("the value %q cannot be parsed as int", str)
		}
		result = any(intVal).(T)
	case uint:
		var uintVal uint64
		uintVal, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			err = fmt.Errorf("the value %q cannot be parsed as int", str)
		}
		result = any(uint(uintVal)).(T)
	case uint8:
		var uintVal uint64
		uintVal, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			err = fmt.Errorf("the value %q cannot be parsed as int", str)
		}
		result = any(uint8(uintVal)).(T)
	case uint16:
		var uintVal uint64
		uintVal, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			err = fmt.Errorf("the value %q cannot be parsed as int", str)
		}
		result = any(uint16(uintVal)).(T)
	case uint32:
		var uintVal uint64
		uintVal, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			err = fmt.Errorf("the value %q cannot be parsed as int", str)
		}
		result = any(uint32(uintVal)).(T)
	case uint64:
		var uintVal uint64
		uintVal, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			err = fmt.Errorf("the value %q cannot be parsed as int", str)
		}
		result = any(uintVal).(T)
	case float32:
		var f64Val float64
		f64Val, err = strconv.ParseFloat(str, 64)
		if err != nil {
			err = fmt.Errorf("the value %q cannot be parsed as float", str)
		}
		result = any(float32(f64Val)).(T)
	case float64:
		var f64Val float64
		f64Val, err = strconv.ParseFloat(str, 64)
		if err != nil {
			err = fmt.Errorf("the value %q cannot be parsed as float", str)
		}
		result = any(f64Val).(T)
	case string:
		result = any(str).(T)
	default:
		err = fmt.Errorf("Unsupported Type ")
	}

	return
}
