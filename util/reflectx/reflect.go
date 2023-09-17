package reflectx

import (
	"errors"
	"reflect"
)

// func CopyValue(srcV reflect.Value, srcT reflect.Type, destV reflect.Value, destT reflect.Type) (err error) {
// 	if destT.Kind() == reflect.Pointer {
// 		pointer = true
// 		destT = destT.Elem()
// 		destV = destV.Elem()

// 		if !destV.IsValid() {

// 		}
// 	}

// 	if srcT.Name() != destT.Name() {
// 		return errors.New("type not same")
// 	}

// 	println(pointer)
// 	return
// }

func NewTypeToValue(t reflect.Type, v reflect.Value) (err error) {
	if t.Kind() != reflect.Pointer || v.Elem().IsValid() {
		return errors.New("not supported")
	}

	v.Set(reflect.New(t))

	return
}
