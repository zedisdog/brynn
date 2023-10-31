package jsonx

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
)

var _ validator.FieldError = (*SimpleFieldError)(nil)

type SimpleFieldError struct {
	field string
	value any
	msg   string
}

// Tag returns the validation tag that failed.
func (fe *SimpleFieldError) Tag() string {
	panic("not implement")
}

// ActualTag returns the validation tag that failed, even if an
// alias the actual tag within the alias will be returned.
func (fe *SimpleFieldError) ActualTag() string {
	panic("not implement")
}

// Namespace returns the namespace for the field error, with the tag
// name taking precedence over the field's actual name.
func (fe *SimpleFieldError) Namespace() string {
	panic("not implement")
}

// StructNamespace returns the namespace for the field error, with the field's
// actual name.
func (fe *SimpleFieldError) StructNamespace() string {
	panic("not implement")
}

// Field returns the field's name with the tag name taking precedence over the
// field's actual name.
func (fe *SimpleFieldError) Field() string {
	return fe.field
}

// StructField returns the field's actual name from the struct, when able to determine.
func (fe *SimpleFieldError) StructField() string {
	panic("not implement")
}

// Value returns the actual field's value in case needed for creating the error
// message
func (fe *SimpleFieldError) Value() interface{} {
	return fe.value
}

// Param returns the param value, in string form for comparison; this will
// also help with generating an error message
func (fe *SimpleFieldError) Param() string {
	panic("not implement")
}

// Kind returns the Field's reflect Kind
func (fe *SimpleFieldError) Kind() reflect.Kind {
	panic("not implement")
}

// Type returns the Field's reflect Type
func (fe *SimpleFieldError) Type() reflect.Type {
	panic("not implement")
}

// Error returns the fieldError's error message
func (fe *SimpleFieldError) Error() string {
	//return fmt.Sprintf("Key: '%s' Errorx:Field validation for '%s' failed on the '%s' tag", fe.ns, fe.Field(), fe.tag)
	return fe.msg
}

// Translate returns the FieldError's translated error
// from the provided 'ut.Translator' and registered 'TranslationFunc'
//
// NOTE: if no registered translation can be found, it returns the original
// untranslated error message.
func (fe *SimpleFieldError) Translate(ut ut.Translator) string {
	panic("not implement")
}
