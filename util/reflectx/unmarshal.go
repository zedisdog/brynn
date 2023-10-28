package reflectx

import (
	"errors"
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
	if dest.Kind() != reflect.Ptr || dest.IsNil() {
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

	if BaseKind(v) == reflect.Float64 {
		if BaseType(dest.Type()).Kind() == reflect.Interface {
			dest.Elem().Set(v)
		} else {
			var err error
			v, err = ConvertFloatTo(v, BaseType(dest.Type()).Kind())
			if err != nil {
				panic(err)
			}
			SetValue(v.Interface(), dest)
		}
	} else {
		dest.Elem().Set(v)
	}
}

func UnmarshalSlice(src []any, dest reflect.Value, tags ...string) (err error) {
	if BaseType(dest.Type()).Kind() != reflect.Slice {
		err = errors.New("data structure mismatched")
		return
	}

	s := reflect.MakeSlice(BaseTypeByValue(dest), 0, len(src))
	for _, item := range src {
		v := reflect.New(BaseTypeByValue(dest).Elem())
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
		case float64:
			if BaseType(v.Type()).Kind() == reflect.Float64 ||
				BaseType(v.Type()).Kind() == reflect.Float32 ||
				BaseType(v.Type()).Kind() == reflect.Uint ||
				BaseType(v.Type()).Kind() == reflect.Uint16 ||
				BaseType(v.Type()).Kind() == reflect.Uint8 ||
				BaseType(v.Type()).Kind() == reflect.Uint32 ||
				BaseType(v.Type()).Kind() == reflect.Uint64 ||
				BaseType(v.Type()).Kind() == reflect.Int ||
				BaseType(v.Type()).Kind() == reflect.Int64 ||
				BaseType(v.Type()).Kind() == reflect.Int32 ||
				BaseType(v.Type()).Kind() == reflect.Int16 ||
				BaseType(v.Type()).Kind() == reflect.Int8 {
				SetValue(x, v)
			} else {
				err = errors.New("data structure mismatched")
				return
			}
		default:
			SetValue(x, v)
		}
		s = reflect.Append(s, v.Elem())
	}

	SetValue(s.Interface(), dest)

	return
}

func UnmarshalMap(src map[string]any, dest reflect.Value, tags ...string) (err error) {
	if BaseType(dest.Type()).Kind() != reflect.Struct {
		err = errors.New("data structure mismatched")
		return
	}

	v := reflect.New(BaseType(dest.Type())).Elem()

	vType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldStruct := vType.Field(i)

		if fieldStruct.Anonymous {
			err = UnmarshalMap(src, v.Field(i), tags...)
			if err != nil {
				return
			}
			continue
		}

		content := GetTag(fieldStruct, tags...)
		arr := strings.Split(content, ",")
		if item, ok := src[arr[0]]; ok {
			switch x := item.(type) {
			case []any:
				err = UnmarshalSlice(x, v.Field(i), tags...)
			case map[string]any:
				err = UnmarshalMap(x, v.Field(i), tags...)
			case float64:
				if BaseType(v.Field(i).Type()).Kind() == reflect.Float64 ||
					BaseType(v.Field(i).Type()).Kind() == reflect.Float32 ||
					BaseType(v.Field(i).Type()).Kind() == reflect.Uint ||
					BaseType(v.Field(i).Type()).Kind() == reflect.Uint16 ||
					BaseType(v.Field(i).Type()).Kind() == reflect.Uint8 ||
					BaseType(v.Field(i).Type()).Kind() == reflect.Uint32 ||
					BaseType(v.Field(i).Type()).Kind() == reflect.Uint64 ||
					BaseType(v.Field(i).Type()).Kind() == reflect.Int ||
					BaseType(v.Field(i).Type()).Kind() == reflect.Int64 ||
					BaseType(v.Field(i).Type()).Kind() == reflect.Int32 ||
					BaseType(v.Field(i).Type()).Kind() == reflect.Int16 ||
					BaseType(v.Field(i).Type()).Kind() == reflect.Int8 ||
					BaseType(v.Field(i).Type()).Kind() == reflect.Interface {
					SetValue(x, v.Field(i))
				} else {
					err = errors.New("data structure mismatched")
					return
				}
			default:
				SetValue(x, v.Field(i))
			}
			if err != nil {
				return
			}
		}
	}

	SetValue(v.Interface(), dest)

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
