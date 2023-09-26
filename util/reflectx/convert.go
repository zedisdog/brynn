package reflectx

import (
	"errors"
	"reflect"
	"strconv"
)

func ConvertTo[T any](src reflect.Value) (result T, err error) {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ConvertIntTo[T](src)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ConvertUintTo[T](src)
	case reflect.Bool:
		return ConvertBoolTo[T](src)
	case reflect.String:
		return ConvertStringTo[T](src)
	}

	err = errors.New("unsupported type")
	return
}

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
	case float32:
		result = any(float32(r)).(T)
	case float64:
		result = any(float64(r)).(T)
	case bool:
		result = any(s).(T)
	case string:
		if s {
			result = any("true").(T)
		} else {
			result = any("false").(T)
		}
	default:
		err = errors.New("unsupported type")
	}

	return
}

func ConvertUintTo[T any](src reflect.Value) (result T, err error) {
	r := src.Uint()
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
		result = any(r).(T)
	case uint32:
		result = any(uint32(r)).(T)
	case uint16:
		result = any(uint16(r)).(T)
	case uint8:
		result = any(uint8(r)).(T)
	case uint:
		result = any(uint(r)).(T)
	case float32:
		result = any(float32(r)).(T)
	case float64:
		result = any(float64(r)).(T)
	case bool:
		result = any(r != 0).(T)
	case string:
		result = any(strconv.FormatUint(r, 10)).(T)
	default:
		err = errors.New("unsupported type")
	}

	return
}

func ConvertIntTo[T any](src reflect.Value) (result T, err error) {
	r := src.Int()
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
		result = any(r).(T)
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
	case float32:
		result = any(float32(r)).(T)
	case float64:
		result = any(float64(r)).(T)
	case bool:
		result = any(r != 0).(T)
	case string:
		result = any(strconv.FormatInt(r, 10)).(T)
	default:
		err = errors.New("unsupported type")
	}

	return
}

func ConvertFloatTo[T any](src reflect.Value) (result T, err error) {
	r := src.Float()
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
	case float32:
		result = any(float32(r)).(T)
	case float64:
		result = any(r).(T)
	case bool:
		result = any(r != 0).(T)
	case string:
		result = any(strconv.FormatFloat(r, 'g', -1, 10)).(T)
	default:
		err = errors.New("unsupported type")
	}

	return
}

func ConvertStringTo[T any](src reflect.Value) (result T, err error) {
	r := src.String()
	switch any(result).(type) {
	case int:
		var i int
		i, err = strconv.Atoi(r)
		result = any(i).(T)
	case int8:
		var i int64
		i, err = strconv.ParseInt(r, 10, 8)
		result = any(int8(i)).(T)
	case int16:
		var i int64
		i, err = strconv.ParseInt(r, 10, 16)
		result = any(int16(i)).(T)
	case int32:
		var i int64
		i, err = strconv.ParseInt(r, 10, 32)
		result = any(int32(i)).(T)
	case int64:
		var i int64
		i, err = strconv.ParseInt(r, 10, 64)
		result = any(i).(T)
	case uint64:
		var i uint64
		i, err = strconv.ParseUint(r, 10, 64)
		result = any(i).(T)
	case uint32:
		var i uint64
		i, err = strconv.ParseUint(r, 10, 32)
		result = any(uint32(i)).(T)
	case uint16:
		var i uint64
		i, err = strconv.ParseUint(r, 10, 16)
		result = any(uint16(i)).(T)
	case uint8:
		var i uint64
		i, err = strconv.ParseUint(r, 10, 8)
		result = any(uint8(i)).(T)
	case uint:
		var i uint64
		i, err = strconv.ParseUint(r, 10, 64)
		result = any(uint(i)).(T)
	case float32:
		var i float64
		i, err = strconv.ParseFloat(r, 32)
		result = any(float32(i)).(T)
	case float64:
		var i float64
		i, err = strconv.ParseFloat(r, 64)
		result = any(i).(T)
	case bool:
		if r == "" || r == "false" || r == "False" {
			result = any(false).(T)
		} else {
			result = any(true).(T)
		}
	case string:
		result = any(r).(T)
	default:
		err = errors.New("unsupported type")
	}

	return
}
