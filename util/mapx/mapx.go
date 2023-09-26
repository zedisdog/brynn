package mapx

import (
	"reflect"
)

func ValidateAndAssign(src map[string]any, dest any) (err error) {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Pointer {

	}
	return
}
