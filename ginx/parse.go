package ginx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/zedisdog/brynn/code"
	"github.com/zedisdog/brynn/codeerr"
	"github.com/zedisdog/brynn/i18n"
	"github.com/zedisdog/brynn/util/reflectx"
	"reflect"
	"strings"
)

func Parse(ctx *gin.Context, container any) (err error) {
	err = ctx.ShouldBindHeader(container)
	if err != nil {
		return
	}
	err = ctx.ShouldBindQuery(container)
	if err != nil {
		return
	}
	switch ContentType(ctx.GetHeader("Content-Type")) {
	case ContentTypeJson:
		err = ctx.ShouldBindJSON(container)
	case ContentTypeForm:
		err = ctx.ShouldBindWith(container, binding.FormPost)
	case ContentTypeMultiPartForm:
		err = ctx.ShouldBindWith(container, binding.FormMultipart)
	default:
		err = codeerr.New(code.BadRequest, i18n.Trans("unsupported Content-Type"))
	}
	if err != nil {
		return
	}

	v := reflect.ValueOf(container).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldStruct := v.Type().Field(i)
		content := reflectx.GetTag(fieldStruct, "form")
		if content == "" {
			continue
		}
		arr := strings.Split(content, ",")

		fieldValue := v.Field(i)
		switch reflectx.BaseType(fieldValue.Type()).Kind() {
		case reflect.Map:
			m := ctx.QueryMap(arr[0])
			if len(m) == 0 && !reflectx.IsOptional(arr[1:]) {
				err = validator.ValidationErrors([]validator.FieldError{
					validator.FieldError(&SimpleFieldError{
						value: "",
						field: arr[0],
						msg: i18n.Transf("field [:field] is required", map[string]any{
							"field": arr[0],
						}),
					}),
				})
				return
			} else if len(m) > 0 {
				var v reflect.Value
				v, err = reflectx.ConvertMapStrStr2MapStrType(reflect.ValueOf(m), fieldValue.Type().Elem())
				if err != nil {
					return
				}

				fieldValue.Set(v)
			}
		case reflect.Slice:
			s := ctx.QueryArray(arr[0])
			if len(s) == 0 && !reflectx.IsOptional(arr[1:]) {
				err = validator.ValidationErrors([]validator.FieldError{
					validator.FieldError(&SimpleFieldError{
						value: s,
						field: arr[0],
						msg: i18n.Transf("field [:field] is required", map[string]any{
							"field": arr[0],
						}),
					}),
				})
				return
			} else if len(s) > 0 {
				var v reflect.Value
				v, err = reflectx.ConvertSliceStr2SliceType(reflect.ValueOf(s), fieldValue.Type().Elem())
				if err != nil {
					return
				}

				fieldValue.Set(v)
			}
		default:
			err = errors.New("unsupported")
		}
	}

	return
}
