package mapx

import (
	"reflect"
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

// func Unmarshal(m map[string][]any, v any) (err error) {
// 	value := reflect.ValueOf(v)
// }

func Set(value reflect.Value, v any) (err error) {
	vv := reflect.ValueOf(v)
	// fmt.Printf("%+v\n", value.Type().Name())
	if value.Kind() == reflect.Invalid || value.Kind() == vv.Kind() {
		value.Set(vv)
	}

	return
}
