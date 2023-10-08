package reflectx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func baseType(v reflect.Type) (t reflect.Type) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v
}

// UnmarshalSlice 解析[]any到其他变量中
func UnmarshalSlice(src []any, dest reflect.Value, tags ...string) (result reflect.Value, err error) {
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
			err = UnmarshalMap(x, v.Elem(), tags...)
			if err != nil {
				return
			}
			result = reflect.Append(result, convertTypeOfPtr(dest.Type().Elem(), v.Elem()))
		case float64:
			var v reflect.Value
			v, err = ConvertFloatTo(reflect.ValueOf(x), elementType.Kind())
			if err != nil {
				return
			}
			result = reflect.Append(result, convertTypeOfPtr(dest.Type().Elem(), v))
		default:
			result = reflect.Append(result, convertTypeOfPtr(dest.Type().Elem(), reflect.ValueOf(item)))
		}
	}

	return
}

// UnmarshalMap 解析map[string]any到结构体中
func UnmarshalMap(src map[string]any, dest reflect.Value, tags ...string) (err error) {
	if dest.Kind() != reflect.Struct {
		err = errors.New("data structure mismatched")
		return
	}
	for i := 0; i < dest.NumField(); i++ {
		if dest.Type().Field(i).Anonymous {
			err = UnmarshalMap(src, dest.Field(i), tags...)
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
			if !IsOptional(arr[1:]) {
				err = fmt.Errorf("field [%s] is required", arr[0])
				return
			}
			continue
		}

		switch x := val.(type) {
		case map[string]any:
			if dest.Field(i).Kind() == reflect.Ptr {
				v := reflect.New(baseType(dest.Field(i).Type()))
				err = UnmarshalMap(x, v.Elem(), tags...)
				if err != nil {
					return
				}
				dest.Field(i).Set(convertTypeOfPtr(dest.Field(i).Type(), v.Elem()))
			} else {
				err = UnmarshalMap(x, dest.Field(i), tags...)
				if err != nil {
					return
				}
			}
		case []any:
			var res reflect.Value
			res, err = UnmarshalSlice(x, dest.Field(i), tags...)
			if err != nil {
				return
			}
			dest.Field(i).Set(res)
		case float64:
			var v reflect.Value
			v, err = ConvertFloatTo(reflect.ValueOf(x), baseType(dest.Field(i).Type()).Kind())
			if err != nil {
				return
			}
			dest.Field(i).Set(convertTypeOfPtr(dest.Field(i).Type(), v))
		default:
			valValue := reflect.ValueOf(val)
			if valValue.IsZero() {
				if !IsOptional(arr[1:]) {
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

// convertTypeOfPtr 自动处理指针的赋值, 从gozero抄来的
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

func IsOptional(arr []string) bool {
	for _, item := range arr {
		if item == "optional" || item == "omitempty" {
			return true
		}
	}

	return false
}
