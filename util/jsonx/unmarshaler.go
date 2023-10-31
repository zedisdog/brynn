package jsonx

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/zedisdog/brynn/util/reflectx"
	"reflect"
	"strings"
)

var (
	ErrMismatched = errors.New("mismatched")
	ErrRequired   = errors.New("required")
)

type JsonUnmarshaler struct {
	tags             []string
	shouldCheckEmpty bool
}

func (j JsonUnmarshaler) Unmarshal(jsonBytes []byte, dest any) (err error) {
	destValue := reflect.ValueOf(dest)
	var src any
	err = json.Unmarshal(jsonBytes, &src)
	if err != nil {
		return
	}
	switch x := src.(type) {
	case map[string]any:
		return j.UnmarshalMap(x, destValue)
	case []any:
		return j.UnmarshalSlice(x, destValue)
	default:
		return errors.New("require src as map[string]any or []any")
	}
}

func (j JsonUnmarshaler) UnmarshalMap(src map[string]any, dest reflect.Value) (err error) {
	switch reflectx.BaseKind(dest) {
	case reflect.Struct:
		err = j.toStruct(src, dest)
	case reflect.Map:
		err = j.toMap(src, dest)
	default:
		err = errors.New("data kind mismatched")
	}

	return
}

func (j JsonUnmarshaler) toMap(src map[string]any, dest reflect.Value) (err error) {
	if reflectx.BaseKind(dest) != reflect.Map {
		return ErrMismatched
	}

	s := reflect.MakeMap(reflectx.BaseType(dest))

	for key, item := range src {
		v := reflect.New(reflectx.BaseType(dest).Elem())
		switch x := item.(type) {
		case []any:
			err = j.UnmarshalSlice(x, v)
		case map[string]any:
			err = j.UnmarshalMap(x, v)
		case float64:
			if reflectx.BaseKind(v) == reflect.Float64 ||
				reflectx.BaseKind(v) == reflect.Float32 ||
				reflectx.BaseKind(v) == reflect.Uint ||
				reflectx.BaseKind(v) == reflect.Uint16 ||
				reflectx.BaseKind(v) == reflect.Uint8 ||
				reflectx.BaseKind(v) == reflect.Uint32 ||
				reflectx.BaseKind(v) == reflect.Uint64 ||
				reflectx.BaseKind(v) == reflect.Int ||
				reflectx.BaseKind(v) == reflect.Int64 ||
				reflectx.BaseKind(v) == reflect.Int32 ||
				reflectx.BaseKind(v) == reflect.Int16 ||
				reflectx.BaseKind(v) == reflect.Int8 {
				err = reflectx.SetValue(reflect.ValueOf(x), v)
			} else {
				err = errors.New("data kind mismatched")
			}
		default:
			err = reflectx.SetValue(reflect.ValueOf(x), v)
		}
		if err != nil {
			return
		}
		s.SetMapIndex(reflect.ValueOf(key), v.Elem())
	}

	err = reflectx.SetValue(s, dest)

	return
}

func (j JsonUnmarshaler) toStruct(src map[string]any, dest reflect.Value) (err error) {
	var errs validator.ValidationErrors
	if reflectx.BaseKind(dest) != reflect.Struct {
		return ErrMismatched
	}

	v := reflect.New(reflectx.BaseType(dest)).Elem()

	vType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldStruct := vType.Field(i)

		if fieldStruct.Anonymous {
			err := j.UnmarshalMap(src, v.Field(i))
			if err != nil {
				switch x := err.(type) {
				case error:
					return x
				case validator.ValidationErrors:
					errs = append(errs, x...)
				}
			}
			continue
		}

		content := reflectx.GetTag(fieldStruct, j.tags...)
		arr := strings.Split(content, ",")

		var err error
		if item, ok := src[arr[0]]; ok {
			switch x := item.(type) {
			case []any:
				if j.shouldCheckEmpty && len(x) == 0 {
					err = ErrRequired
				} else {
					err = j.UnmarshalSlice(x, v.Field(i))
				}
			case map[string]any:
				if j.shouldCheckEmpty && len(x) == 0 {
					err = ErrRequired
				} else {
					err = j.UnmarshalMap(x, v.Field(i))
				}
			case float64:
				if j.shouldCheckEmpty && x == 0 {
					err = ErrRequired
				} else {
					if reflectx.BaseKind(v.Field(i)) == reflect.Float64 ||
						reflectx.BaseKind(v.Field(i)) == reflect.Float32 ||
						reflectx.BaseKind(v.Field(i)) == reflect.Uint ||
						reflectx.BaseKind(v.Field(i)) == reflect.Uint16 ||
						reflectx.BaseKind(v.Field(i)) == reflect.Uint8 ||
						reflectx.BaseKind(v.Field(i)) == reflect.Uint32 ||
						reflectx.BaseKind(v.Field(i)) == reflect.Uint64 ||
						reflectx.BaseKind(v.Field(i)) == reflect.Int ||
						reflectx.BaseKind(v.Field(i)) == reflect.Int64 ||
						reflectx.BaseKind(v.Field(i)) == reflect.Int32 ||
						reflectx.BaseKind(v.Field(i)) == reflect.Int16 ||
						reflectx.BaseKind(v.Field(i)) == reflect.Int8 ||
						reflectx.BaseKind(v.Field(i)) == reflect.Interface {
						err = reflectx.SetValue(reflect.ValueOf(x), v.Field(i))
					} else {
						err = ErrMismatched
					}
				}
			default:
				err = reflectx.SetValue(reflect.ValueOf(x), v.Field(i))
			}
			if err != nil {
				switch x := err.(type) {
				case error:
					return x
				case validator.ValidationErrors:
					errs = append(errs, x...)
				}
			}
		}
	}

	err = reflectx.SetValue(v, dest)
	if err != nil {
		return
	}

	return errs
}

// TODO: 错误类型替换
func (j JsonUnmarshaler) UnmarshalSlice(src []any, dest reflect.Value) (err error) {
	if reflectx.BaseKind(dest) != reflect.Slice {
		err = errors.New("data kind mismatched")
		return
	}

	s := reflect.MakeSlice(reflectx.BaseType(dest), 0, len(src))
	for _, item := range src {
		v := reflect.New(reflectx.BaseType(dest).Elem())
		switch x := item.(type) {
		case []any:
			err = j.UnmarshalSlice(x, v)
		case map[string]any:
			err = j.UnmarshalMap(x, v)
		case float64:
			if reflectx.BaseKind(v) == reflect.Float64 ||
				reflectx.BaseKind(v) == reflect.Float32 ||
				reflectx.BaseKind(v) == reflect.Uint ||
				reflectx.BaseKind(v) == reflect.Uint16 ||
				reflectx.BaseKind(v) == reflect.Uint8 ||
				reflectx.BaseKind(v) == reflect.Uint32 ||
				reflectx.BaseKind(v) == reflect.Uint64 ||
				reflectx.BaseKind(v) == reflect.Int ||
				reflectx.BaseKind(v) == reflect.Int64 ||
				reflectx.BaseKind(v) == reflect.Int32 ||
				reflectx.BaseKind(v) == reflect.Int16 ||
				reflectx.BaseKind(v) == reflect.Int8 {
				err = reflectx.SetValue(reflect.ValueOf(x), v)
			} else {
				err = errors.New("data kind mismatched")
			}
		default:
			err = reflectx.SetValue(reflect.ValueOf(x), v)
		}
		if err != nil {
			return
		}
		s = reflect.Append(s, v.Elem())
	}

	err = reflectx.SetValue(s, dest)

	return
}

func (j JsonUnmarshaler) isOptional(arr []string) bool {
	for _, item := range arr {
		if item == "optional" || item == "omitempty" {
			return true
		}
	}

	return false
}

func (j JsonUnmarshaler) checkEmpty(srcValue reflect.Value, tags []string) (err error) {
	if srcValue.IsZero() && j.shouldCheckEmpty && j.isOptional(tags) {
		return errors.New("required")
	}

	return
}
