package reflectx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Unmarshal(src any, dest any, tags ...string) (err error) {

	return
}

func baseType(v reflect.Type) (t reflect.Type) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v
}

func unmarshalSlice(src []any, dest reflect.Value, tags ...string) (result reflect.Value, err error) {
	if dest.Kind() != reflect.Slice {
		err = errors.New("data structure mismatched")
		return
	}

	elementType := dest.Type().Elem()
	result = reflect.MakeSlice(dest.Type(), 0, len(src))
	for _, item := range src {
		switch x := item.(type) {
		case map[string]any:
			v := reflect.New(baseType(elementType))
			err = unmarshalMap(x, v.Elem(), tags...)
			if err != nil {
				return
			}
			result = reflect.Append(result, convertTypeOfPtr(dest.Type().Elem(), v.Elem()))
		default:
			result = reflect.Append(result, convertTypeOfPtr(dest.Type().Elem(), reflect.ValueOf(item)))
		}
	}

	return
}

// TODO: unmarshalSlice
func unmarshalMap(src map[string]any, dest reflect.Value, tags ...string) (err error) {
	if dest.Kind() != reflect.Struct {
		err = errors.New("data structure mismatched")
		return
	}
	for i := 0; i < dest.NumField(); i++ {
		if dest.Type().Field(i).Anonymous {
			err = unmarshalMap(src, dest.Field(i), tags...)
			if err != nil {
				return
			}
			continue
		}

		content := GetTag(dest.Type().Field(i), tags...)
		if content == "" {
			continue
		}

		arr := strings.Split(content, ",")
		val, ok := src[arr[0]]

		if !ok {
			if !isOptional(arr[1:]) {
				err = fmt.Errorf("field [%s] is required", arr[0])
				return
			}
			continue
		}

		switch x := val.(type) {
		case map[string]any:
			err = unmarshalMap(x, dest.Field(i), tags...)
			if err != nil {
				return
			}
		case []any:
			var res reflect.Value
			res, err = unmarshalSlice(x, dest.Field(i))
			if err != nil {
				return
			}
			dest.Field(i).Set(res)
		default:
			valValue := reflect.ValueOf(val)
			if valValue.IsZero() {
				if !isOptional(arr[1:]) {
					err = fmt.Errorf("field [%s] is required", arr[0])
					return
				}
				continue
			}

			dest.Field(i).Set(convertTypeOfPtr(dest.Field(i).Type(), valValue))
		}
	}

	return
}

func baseKind(value reflect.Value) reflect.Kind {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return value.Kind()
}

func convertTypeOfPtr(tp reflect.Type, target reflect.Value) reflect.Value {
	// keep the original value is a pointer
	if tp.Kind() == reflect.Ptr && target.CanAddr() {
		tp = tp.Elem()
		target = target.Addr()
	}

	for tp.Kind() == reflect.Ptr {
		p := reflect.New(target.Type())
		p.Elem().Set(target)
		target = p
		tp = tp.Elem()
	}

	return target
}

func isOptional(arr []string) bool {
	for _, item := range arr {
		if item == "optional" {
			return true
		}
	}

	return false
}
