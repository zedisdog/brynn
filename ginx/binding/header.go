package binding

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/zedisdog/brynn/i18n"
	"github.com/zedisdog/brynn/util/reflectx"
	"net/http"
	"reflect"
	"strings"
)

type headerBinding struct{}

func (headerBinding) Name() string {
	return "header"
}

func (hb headerBinding) Bind(req *http.Request, obj any) error {

	if err := hb.parseHeader(req, obj); err != nil {
		return err
	}

	return validate(obj)
}

func (hb headerBinding) parseHeader(req *http.Request, obj any) (err error) {
	destValue := reflect.ValueOf(obj).Elem()
	destType := destValue.Type()
	var errs validator.ValidationErrors

	for i := 0; i < destValue.NumField(); i++ {
		fieldStruct := destType.Field(i)
		fieldValue := destValue.Field(i)
		if fieldStruct.Anonymous {
			errResult := hb.parseHeader(req, fieldValue)
			if errResult != nil {
				errs = append(errs, errResult.(validator.ValidationErrors)...)
			}
			continue
		}

		content := reflectx.GetTag(fieldStruct, "header")
		if content != "" {
			arr := strings.Split(content, ",")
			values := req.Header.Values(arr[0])
			switch len(values) {
			case 0:
				if !reflectx.IsOptional(arr[1:]) {
					errs = append(errs, &SimpleFieldError{
						field: arr[0],
						value: "",
						msg:   i18n.Transf("field [:field] is required", i18n.P{"field": arr[0]}),
					})
					continue
				}
			case 1:
				var v reflect.Value
				v, err = reflectx.ConvertStringTo(reflect.ValueOf(values[0]), fieldValue.Kind())
				if err != nil {
					errs = append(errs, &SimpleFieldError{
						field: arr[0],
						value: "",
						msg:   i18n.Transf("field [:field] type error", i18n.P{"field": arr[0]}),
					})
					continue
				}
				fieldValue.Set(v)
			default:
				if fieldValue.Kind() == reflect.Slice {
					for _, item := range values {
						var v reflect.Value
						v, err = reflectx.ConvertStringTo(reflect.ValueOf(item), fieldValue.Type().Elem().Kind())
						if err != nil {
							errs = append(errs, &SimpleFieldError{
								field: arr[0],
								value: "",
								msg:   i18n.Transf("field [:field] type error", i18n.P{"field": arr[0]}),
							})
							continue
						}
						fieldValue.Set(reflect.Append(fieldValue, v))
					}
				} else {
					panic(errors.New("unsupported"))
				}
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
