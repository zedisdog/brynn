package ginx

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
)

var _ validator.FieldError = (*FieldError)(nil)

type FieldError struct {
	v              *validator.Validate
	tag            string
	actualTag      string
	ns             string
	structNs       string
	fieldLen       uint8
	structfieldLen uint8
	value          interface{}
	param          string
	kind           reflect.Kind
	typ            reflect.Type
}

// Tag returns the validation tag that failed.
func (fe *FieldError) Tag() string {
	return fe.tag
}

// ActualTag returns the validation tag that failed, even if an
// alias the actual tag within the alias will be returned.
func (fe *FieldError) ActualTag() string {
	return fe.actualTag
}

// Namespace returns the namespace for the field error, with the tag
// name taking precedence over the field's actual name.
func (fe *FieldError) Namespace() string {
	return fe.ns
}

// StructNamespace returns the namespace for the field error, with the field's
// actual name.
func (fe *FieldError) StructNamespace() string {
	return fe.structNs
}

// Field returns the field's name with the tag name taking precedence over the
// field's actual name.
func (fe *FieldError) Field() string {
	return fe.ns[len(fe.ns)-int(fe.fieldLen):]
}

// StructField returns the field's actual name from the struct, when able to determine.
func (fe *FieldError) StructField() string {
	return fe.structNs[len(fe.structNs)-int(fe.structfieldLen):]
}

// Value returns the actual field's value in case needed for creating the error
// message
func (fe *FieldError) Value() interface{} {
	return fe.value
}

// Param returns the param value, in string form for comparison; this will
// also help with generating an error message
func (fe *FieldError) Param() string {
	return fe.param
}

// Kind returns the Field's reflect Kind
func (fe *FieldError) Kind() reflect.Kind {
	return fe.kind
}

// Type returns the Field's reflect Type
func (fe *FieldError) Type() reflect.Type {
	return fe.typ
}

// Error returns the fieldError's error message
func (fe *FieldError) Error() string {
	return fmt.Sprintf("Key: '%s' Error:Field validation for '%s' failed on the '%s' tag", fe.ns, fe.Field(), fe.tag)
}

// Translate returns the FieldError's translated error
// from the provided 'ut.Translator' and registered 'TranslationFunc'
//
// NOTE: if no registered translation can be found, it returns the original
// untranslated error message.
func (fe *FieldError) Translate(ut ut.Translator) string {

	m, ok := fe.v.transTagFunc[ut]
	if !ok {
		return fe.Error()
	}

	fn, ok := m[fe.tag]
	if !ok {
		return fe.Error()
	}

	return fn(ut, fe)
}
