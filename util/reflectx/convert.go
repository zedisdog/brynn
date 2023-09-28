package reflectx

import (
	"errors"
	"reflect"
	"strconv"
)

func ConvertTo[T any](src reflect.Value) (result T, err error) {
	r, err := ConvertToType(src, reflect.ValueOf(result).Kind())
	if err != nil {
		return
	}
	result = r.Interface().(T)
	return
}

func ConvertToType(src reflect.Value, kind reflect.Kind) (result reflect.Value, err error) {
	var r reflect.Value
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r, err = ConvertIntTo(src, kind)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r, err = ConvertUintTo(src, kind)
	case reflect.Bool:
		r, err = ConvertBoolTo(src, kind)
	case reflect.String:
		r, err = ConvertStringTo(src, kind)
	case reflect.Float32, reflect.Float64:
		r, err = ConvertFloatTo(src, kind)
	default:
		err = errors.New("unsupported type")
		return
	}

	if err != nil {
		return
	}

	result = r

	return
}

func ConvertBoolTo(src reflect.Value, kind reflect.Kind) (result reflect.Value, err error) {
	s := src.Interface().(bool)
	var r int8
	if s {
		r = 1
	} else {
		r = 0
	}
	switch kind {
	case reflect.Int:
		result = reflect.ValueOf(int(r))
	case reflect.Int8:
		result = reflect.ValueOf(r)
	case reflect.Int16:
		result = reflect.ValueOf(int16(r))
	case reflect.Int32:
		result = reflect.ValueOf(int32(r))
	case reflect.Int64:
		result = reflect.ValueOf(int64(r))
	case reflect.Uint64:
		result = reflect.ValueOf(uint64(r))
	case reflect.Uint32:
		result = reflect.ValueOf(uint32(r))
	case reflect.Uint16:
		result = reflect.ValueOf(uint16(r))
	case reflect.Uint8:
		result = reflect.ValueOf(uint8(r))
	case reflect.Uint:
		result = reflect.ValueOf(uint(r))
	case reflect.Float32:
		result = reflect.ValueOf(float32(r))
	case reflect.Float64:
		result = reflect.ValueOf(float64(r))
	case reflect.Bool:
		result = reflect.ValueOf(s)
	case reflect.String:
		if s {
			result = reflect.ValueOf("true")
		} else {
			result = reflect.ValueOf("false")
		}
	default:
		err = errors.New("unsupported type")
	}

	return
}

func ConvertUintTo(src reflect.Value, kind reflect.Kind) (result reflect.Value, err error) {
	r := src.Uint()
	switch kind {
	case reflect.Int:
		result = reflect.ValueOf(int(r))
	case reflect.Int8:
		result = reflect.ValueOf(int8(r))
	case reflect.Int16:
		result = reflect.ValueOf(int16(r))
	case reflect.Int32:
		result = reflect.ValueOf(int32(r))
	case reflect.Int64:
		result = reflect.ValueOf(int64(r))
	case reflect.Uint64:
		result = reflect.ValueOf(r)
	case reflect.Uint32:
		result = reflect.ValueOf(uint32(r))
	case reflect.Uint16:
		result = reflect.ValueOf(uint16(r))
	case reflect.Uint8:
		result = reflect.ValueOf(uint8(r))
	case reflect.Uint:
		result = reflect.ValueOf(uint(r))
	case reflect.Float32:
		result = reflect.ValueOf(float32(r))
	case reflect.Float64:
		result = reflect.ValueOf(float64(r))
	case reflect.Bool:
		result = reflect.ValueOf(r != 0)
	case reflect.String:
		result = reflect.ValueOf(strconv.FormatUint(r, 10))
	default:
		err = errors.New("unsupported type")
	}

	return
}

func ConvertIntTo(src reflect.Value, kind reflect.Kind) (result reflect.Value, err error) {
	r := src.Int()
	switch kind {
	case reflect.Int:
		result = reflect.ValueOf(int(r))
	case reflect.Int8:
		result = reflect.ValueOf(int8(r))
	case reflect.Int16:
		result = reflect.ValueOf(int16(r))
	case reflect.Int32:
		result = reflect.ValueOf(int32(r))
	case reflect.Int64:
		result = reflect.ValueOf(r)
	case reflect.Uint64:
		result = reflect.ValueOf(uint64(r))
	case reflect.Uint32:
		result = reflect.ValueOf(uint32(r))
	case reflect.Uint16:
		result = reflect.ValueOf(uint16(r))
	case reflect.Uint8:
		result = reflect.ValueOf(uint8(r))
	case reflect.Uint:
		result = reflect.ValueOf(uint(r))
	case reflect.Float32:
		result = reflect.ValueOf(float32(r))
	case reflect.Float64:
		result = reflect.ValueOf(float64(r))
	case reflect.Bool:
		result = reflect.ValueOf(r != 0)
	case reflect.String:
		result = reflect.ValueOf(strconv.FormatInt(r, 10))
	default:
		err = errors.New("unsupported type")
	}

	return
}

func ConvertFloatTo(src reflect.Value, kind reflect.Kind) (result reflect.Value, err error) {
	r := src.Float()
	switch kind {
	case reflect.Int:
		result = reflect.ValueOf(int(r))
	case reflect.Int8:
		result = reflect.ValueOf(int8(r))
	case reflect.Int16:
		result = reflect.ValueOf(int16(r))
	case reflect.Int32:
		result = reflect.ValueOf(int32(r))
	case reflect.Int64:
		result = reflect.ValueOf(int64(r))
	case reflect.Uint64:
		result = reflect.ValueOf(uint64(r))
	case reflect.Uint32:
		result = reflect.ValueOf(uint32(r))
	case reflect.Uint16:
		result = reflect.ValueOf(uint16(r))
	case reflect.Uint8:
		result = reflect.ValueOf(uint8(r))
	case reflect.Uint:
		result = reflect.ValueOf(uint(r))
	case reflect.Float32:
		result = reflect.ValueOf(float32(r))
	case reflect.Float64:
		result = reflect.ValueOf(r)
	case reflect.Bool:
		result = reflect.ValueOf(r != 0)
	case reflect.String:
		result = reflect.ValueOf(strconv.FormatFloat(r, 'g', -1, 10))
	default:
		err = errors.New("unsupported type")
	}

	return
}

func ConvertStringTo(src reflect.Value, kind reflect.Kind) (result reflect.Value, err error) {
	r := src.String()
	switch kind {
	case reflect.Int:
		var i int
		i, err = strconv.Atoi(r)
		result = reflect.ValueOf(i)
	case reflect.Int8:
		var i int64
		i, err = strconv.ParseInt(r, 10, 8)
		result = reflect.ValueOf(int8(i))
	case reflect.Int16:
		var i int64
		i, err = strconv.ParseInt(r, 10, 16)
		result = reflect.ValueOf(int16(i))
	case reflect.Int32:
		var i int64
		i, err = strconv.ParseInt(r, 10, 32)
		result = reflect.ValueOf(int32(i))
	case reflect.Int64:
		var i int64
		i, err = strconv.ParseInt(r, 10, 64)
		result = reflect.ValueOf(i)
	case reflect.Uint64:
		var i uint64
		i, err = strconv.ParseUint(r, 10, 64)
		result = reflect.ValueOf(i)
	case reflect.Uint32:
		var i uint64
		i, err = strconv.ParseUint(r, 10, 32)
		result = reflect.ValueOf(uint32(i))
	case reflect.Uint16:
		var i uint64
		i, err = strconv.ParseUint(r, 10, 16)
		result = reflect.ValueOf(uint16(i))
	case reflect.Uint8:
		var i uint64
		i, err = strconv.ParseUint(r, 10, 8)
		result = reflect.ValueOf(uint8(i))
	case reflect.Uint:
		var i uint64
		i, err = strconv.ParseUint(r, 10, 64)
		result = reflect.ValueOf(uint(i))
	case reflect.Float32:
		var i float64
		i, err = strconv.ParseFloat(r, 32)
		result = reflect.ValueOf(float32(i))
	case reflect.Float64:
		var i float64
		i, err = strconv.ParseFloat(r, 64)
		result = reflect.ValueOf(i)
	case reflect.Bool:
		if r == "" || r == "false" || r == "False" {
			result = reflect.ValueOf(false)
		} else {
			result = reflect.ValueOf(true)
		}
	case reflect.String:
		result = reflect.ValueOf(r)
	default:
		err = errors.New("unsupported type")
	}

	return
}
