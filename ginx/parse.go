package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/zedisdog/brynn/errx"
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
		err = errx.New(errx.BadRequest, i18n.Trans("unsupported Content-Type"))
	}
	if err != nil {
		return
	}

	v := reflect.ValueOf(container)
	for i := 0; i < v.NumField(); i++ {
		fieldStruct := v.Type().Field(i)
		content := reflectx.GetTag(fieldStruct, "form")
		if content == "" {
			continue
		}
		arr := strings.Split(content, ",")

		fieldValue := v.Field(i)
		if fieldValue.Kind() == reflect.Map {
			m := ctx.QueryMap(arr[0])
			if len(m) == 0 && !reflectx.IsOptional(arr[1:]) {
				err = validator.ValidationErrors([]validator.FieldError{})
			}
		}
		//TODO: there
		m := ctx.QueryMap()

	}

	return
}
