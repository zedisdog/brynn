package reflectx

import (
	"errors"
	"reflect"
)

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
