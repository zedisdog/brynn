package reflectx

import (
	"errors"
	"reflect"
	"strings"
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

func GetTag(field reflect.StructField, tags ...string) (content string) {
	for content == "" && len(tags) > 0 {
		tag := tags[0]
		tags = tags[1:]
		content = field.Tag.Get(tag)
	}

	return
}

func FieldByKeyWithOption(v reflect.Value, t reflect.Type, key string, tags ...string) (val reflect.Value, typ reflect.Type, option string, find bool) {
	for i := 0; i < v.NumField(); i++ {
		val = v.Field(i)
		structField := t.Field(i)

		if structField.Anonymous {
			val, typ, option, find = FieldByKeyWithOption(val, val.Type(), key, tags...)
		}
		if find {
			return
		}

		content := GetTag(structField, tags...)

		if content == "" && structField.Name == key {
			typ = val.Type()
			find = true
			return
		} else if content != "" && strings.HasPrefix(content, key+",") {
			typ = val.Type()
			option, _ = strings.CutPrefix(content, key+",")
			find = true
			return
		}
	}

	return
}
