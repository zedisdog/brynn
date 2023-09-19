package mapx

import (
	"errors"
	"reflect"
	"strings"

	"github.com/zedisdog/brynn/util/slicex"
)

func Merge[K comparable, V any](m1 map[K]V, m2 map[K]V) (result map[K]V) {
	result = make(map[K]V, len(m1)+len(m2))
	if m1 == nil && m2 == nil {
		return
	} else if m1 == nil {
		return m2
	} else if m2 == nil {
		return m1
	}

	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		result[k] = v
	}

	return
}

const (
	keyName = "input"
)

var (
	ErrNotMatch = errors.New("type not match")
)

func Unmarshal(src map[string]any, dest any) (err error) {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Pointer || !destValue.Elem().IsValid() {
		return errors.New("only support pointer")
	}
	destValue = destValue.Elem()
	destType := reflect.TypeOf(dest).Elem()

	return unmarshal(src, destValue, destType)
}

func unmarshal(src map[string]any, dValue reflect.Value, dType reflect.Type) (err error) {
	for i := 0; i < dValue.NumField(); i++ {
		fieldValue := dValue.Field(i)
		fieldType := dType.Field(i)

		if fieldType.Anonymous {
			unmarshal(src, fieldValue, fieldValue.Type())
			continue
		}

		var key string
		var options []string
		if tag := fieldType.Tag.Get(keyName); tag != "" {
			arr := strings.Split(tag, ",")
			key = arr[0]
			options = arr[1:]
		} else {
			key = fieldType.Name
		}

		if v, ok := src[key]; ok {
			srcValue := reflect.ValueOf(v)
			if fieldValue.Kind() == srcValue.Kind() {
				fieldValue.Set(srcValue)
			} else {
				return ErrNotMatch
			}
		} else if !slicex.Containers("optional", options) {
			return errors.New("field <" + key + "> is required")
		}
	}

	return
}

// func Set(value reflect.Value, v any) (err error) {
// 	vv := reflect.ValueOf(v)
// 	// fmt.Printf("%+v\n", value.Type().Name())
// 	if value.Kind() == reflect.Invalid || value.Kind() == vv.Kind() {
// 		value.Set(vv)
// 	}

// 	return
// }
