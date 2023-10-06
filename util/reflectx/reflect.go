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

func GetTag(field reflect.StructField, tags ...string) (content string) {
	for content == "" && len(tags) > 0 {
		tag := tags[0]
		tags = tags[1:]
		content = field.Tag.Get(tag)
	}

	return
}

//func Assign(src any, dest reflect.Value, tags ...string) (err error) {
//	switch src.(type) {
//	case map[string]any:
//		for i := 0; i < dest.NumField(); i++ {
//			fieldStruct := dest.Type().Field(i)
//			if fieldStruct.Anonymous {
//				err = Assign(src, dest.Field(i), tags...)
//				if err != nil {
//					return
//				}
//			}
//
//			content := GetTag(dest.Type().Field(i), tags...)
//			arr := strings.Split(content, ",")
//			if content == "" {
//				continue
//			}
//			fieldValue := dest.Field(i)
//			if fieldValue.IsZero() {
//				if !isOptional(arr[1:]) {
//					err = errx.New(errx.ValidateError, i18n.Transf("field [:field] is required", i18n.P{"field": arr[0]}))
//					return
//				}
//				continue
//			}
//		}
//	}
//}
//
//func isOptional(arr []string) bool {
//	for _, item := range arr {
//		if item == "optional" {
//			return true
//		}
//	}
//
//	return false
//}

//func FieldByKeyWithOption(v reflect.Value, t reflect.Type, key string, tags ...string) (val reflect.Value, typ reflect.Type, option string, find bool) {
//	for i := 0; i < v.NumField(); i++ {
//		val = v.Field(i)
//		structField := t.Field(i)
//
//		if structField.Anonymous {
//			val, typ, option, find = FieldByKeyWithOption(val, val.Type(), key, tags...)
//		}
//		if find {
//			return
//		}
//
//		content := GetTag(structField, tags...)
//
//		if content == "" && structField.Name == key {
//			typ = val.Type()
//			find = true
//			return
//		} else if content != "" && strings.HasPrefix(content, key+",") {
//			typ = val.Type()
//			option, _ = strings.CutPrefix(content, key+",")
//			find = true
//			return
//		}
//	}
//
//	return
//}
