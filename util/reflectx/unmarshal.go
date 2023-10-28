package reflectx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func BaseType(v reflect.Type) (t reflect.Type) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v
}

func Unmarshal(src any, dest reflect.Value, tags ...string) (err error) {
	switch x := src.(type) {
	case map[string]any:
		return UnmarshalMap(x, dest, tags...)
	case []any:
		return UnmarshalSlice(x, dest, tags...)
	default:
		panic(errors.New("unsupported"))
	}
}

func SetValue(src any, dest reflect.Value) {
	if dest.Kind() != reflect.Ptr {
		if dest.CanAddr() {
			dest = dest.Addr()
		} else {
			panic(errors.New("value must be pointer or canAddr"))
		}
	}

	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	tp := dest.Elem().Type()
	v := srcValue
	for tp.Kind() == reflect.Ptr {
		n := reflect.New(v.Type())
		n.Elem().Set(v)
		v = n
		tp = tp.Elem()
	}

	dest.Elem().Set(v)
}

func UnmarshalSlice(src []any, dest reflect.Value, tags ...string) (err error) {
	if BaseKind(dest) != reflect.Slice {
		err = errors.New("data structure mismatched")
		return
	}

	s := reflect.MakeSlice(BaseTypeByValue(dest), 0, len(src))
	for _, item := range src {
		v := reflect.New(BaseTypeByValue(dest))
		switch x := item.(type) {
		case []any:
			err = UnmarshalSlice(x, v, tags...)
			if err != nil {
				return
			}
		case map[string]any:
			err = UnmarshalMap(x, v, tags...)
			if err != nil {
				return
			}
		default:
			SetValue(x, v)
		}
		s = reflect.Append(s, v)
	}

	SetValue(s.Interface().([]any), dest)

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
				v := reflect.New(BaseType(dest.Field(i).Type()))
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
			err = UnmarshalSlice(x, dest.Field(i), tags...)
			if err != nil {
				return
			}
			dest.Field(i).Set(res)
		case float64:
			var v reflect.Value
			v, err = ConvertFloatTo(reflect.ValueOf(x), BaseType(dest.Field(i).Type()).Kind())
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

func BaseKind(v reflect.Value) reflect.Kind {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v.Kind()
}

func BaseTypeByValue(v reflect.Value) reflect.Type {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v.Type()
}
