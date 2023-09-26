package conv

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ConvertTo[T any](a any) (result T, err error) {
	switch x := a.(type) {
	case bool:
		result = BoolTo[T](x)
	case int:
		result = IntTo[T](x)
	case int8:
		result = IntTo[T](x)
	case int16:
		result = IntTo[T](x)
	case int32:
		result = IntTo[T](x)
	case int64:
		result = IntTo[T](x)
	case uint:
		result = IntTo[T](x)
	case uint8:
		result = IntTo[T](x)
	case uint16:
		result = IntTo[T](x)
	case uint32:
		result = IntTo[T](x)
	case uint64:
		result = IntTo[T](x)
	case float32:
		result = FloatTo[T](x)
	case float64:
		result = FloatTo[T](x)
	case string:
		result, err = StrTo[T](x)
	default:
		value := reflect.ValueOf(a)
		switch value.Kind() {
		case reflect.Int:
			println("ok")
		}
	}

	return
}

type intBase interface {
	~int | ~int8 | ~int16 | ~int64 | ~int32 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func IntTo[T any, S intBase](v S) (result T) {
	switch any(result).(type) {
	case bool:
		result = any(0 != v).(T)
	case int:
		result = any(int(v)).(T)
	case int8:
		result = any(int8(v)).(T)
	case int16:
		result = any(int16(v)).(T)
	case int32:
		result = any(int32(v)).(T)
	case int64:
		result = any(int64(v)).(T)
	case uint:
		result = any(uint(v)).(T)
	case uint8:
		result = any(uint8(v)).(T)
	case uint16:
		result = any(uint16(v)).(T)
	case uint32:
		result = any(uint32(v)).(T)
	case uint64:
		result = any(uint64(v)).(T)
	case float32:
		result = any(float32(v)).(T)
	case float64:
		result = any(float64(v)).(T)
	case string:
		result = any(fmt.Sprintf("%d", v)).(T)
	}

	return
}

type floatBase interface {
	~float32 | ~float64
}

func FloatTo[T any, S floatBase](v S) (result T) {
	switch any(result).(type) {
	case bool:
		result = any(0 != v).(T)
	case int:
		result = any(int(v)).(T)
	case int8:
		result = any(int8(v)).(T)
	case int16:
		result = any(int16(v)).(T)
	case int32:
		result = any(int32(v)).(T)
	case int64:
		result = any(int64(v)).(T)
	case uint:
		result = any(uint(v)).(T)
	case uint8:
		result = any(uint8(v)).(T)
	case uint16:
		result = any(uint16(v)).(T)
	case uint32:
		result = any(uint32(v)).(T)
	case uint64:
		result = any(uint64(v)).(T)
	case float32:
		result = any(float32(v)).(T)
	case float64:
		result = any(float64(v)).(T)
	case string:
		result = any(fmt.Sprintf("%f", v)).(T)
	}

	return
}

func BoolTo[T any](v bool) (result T) {
	switch any(result).(type) {
	case bool:
		result = any(v).(T)
	case int:
		if v {
			result = any(1).(T)
		} else {
			result = any(0).(T)
		}
	case int8:
		if v {
			result = any(int8(1)).(T)
		} else {
			result = any(int8(1)).(T)
		}
	case int16:
		if v {
			result = any(int16(1)).(T)
		} else {
			result = any(int16(1)).(T)
		}
	case int32:
		if v {
			result = any(int32(1)).(T)
		} else {
			result = any(int32(1)).(T)
		}
	case int64:
		if v {
			result = any(int64(1)).(T)
		} else {
			result = any(int64(1)).(T)
		}
	case uint:
		if v {
			result = any(uint(1)).(T)
		} else {
			result = any(uint(1)).(T)
		}
	case uint8:
		if v {
			result = any(uint8(1)).(T)
		} else {
			result = any(uint8(1)).(T)
		}
	case uint16:
		if v {
			result = any(uint16(1)).(T)
		} else {
			result = any(uint16(1)).(T)
		}
	case uint32:
		if v {
			result = any(uint32(1)).(T)
		} else {
			result = any(uint32(1)).(T)
		}
	case uint64:
		if v {
			result = any(uint64(1)).(T)
		} else {
			result = any(uint64(1)).(T)
		}
	case float32:
		if v {
			result = any(float32(1)).(T)
		} else {
			result = any(float32(1)).(T)
		}
	case float64:
		if v {
			result = any(float64(1)).(T)
		} else {
			result = any(float64(1)).(T)
		}
	case string:
		if v {
			result = any("true").(T)
		} else {
			result = any("false").(T)
		}
	}

	return
}

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
