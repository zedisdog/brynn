package reflectx

import (
	"fmt"
	"reflect"
	"strings"
)

func Unmarshal(src any, dest any, tags ...string) (err error) {

	return
}

func unmarshalMap(src map[string]any, dest reflect.Value, tags ...string) (err error) {
	for i := 0; i<dest.NumField(); i++ {
		content := GetTag(dest.Type().Field(i), tags...)
		if content == "" {
			continue
		}

		arr := strings.Split(content, ",")
		val, ok := src[arr[0]]
		valValue := reflect.ValueOf(val)

		if !ok || valValue.IsZero() {
			if !isOptional(arr[1:]) {
				err = fmt.Errorf("field [%s] is required", arr[0])
				return
			}
			continue
		}
		dest.Field(i).Set(valValue)
	}

	return
}

func isOptional(arr []string) bool {
	for _, item := range arr {
		if item == "optional" {
			return true
		}
	}

	return false
}
