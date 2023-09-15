package ginx

import (
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/gin-gonic/gin"
)

func NewContext(c *gin.Context) Context {
	return &context{
		Context: c,
	}
}

type Context interface {
	Parse(container any) (err error)
	Abort()
	AbortWithStatus(code int)
	Param(key string) string
	Query(key string) (value string)
}

type context struct {
	*gin.Context
}

func (c *context) Parse(container any) (err error) {
	t := reflect.TypeOf(container)
	v := reflect.ValueOf(container)

	err = c.parseHeader(v, t)
	if err != nil {
		return
	}

	return
}

func (c *context) parseCookie(v reflect.Value, t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		vField := v.Field(i)
		if !vField.CanSet() {
			vField = reflect.NewAt(vField.Type(), unsafe.Pointer(vField.UnsafeAddr())).Elem()
		}

		fieldName := tField.Tag.Get("cookie")
		for _, item := range c.Request.Cookies() {
			if item.Name == fieldName {
				switch vField.Interface().(type) {
				case string:
					vField.Set(reflect.ValueOf(item.Value))
				case int, int32: //TODO: 基础类型转换
					value, err := strconv.Atoi(item.Value)
					if err != nil {
						panic(err)
					}
					vField.Set(reflect.ValueOf(value))
				}
				break
			}
		}
	}

	return
}

func (c *context) parseHeader(v reflect.Value, t reflect.Type) (err error) {
	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		vField := v.Field(i)
		if !vField.CanSet() {
			vField = reflect.NewAt(vField.Type(), unsafe.Pointer(vField.UnsafeAddr())).Elem()
		}

		fieldName := tField.Tag.Get("header")
		switch vField.Kind() {
		case reflect.String:
			vField.Set(reflect.ValueOf(strings.Join(c.Request.Header.Values(fieldName), ",")))
		case reflect.Int:

		}

		if vField.Kind() == reflect.Slice && vField.Type().Elem().Kind() == reflect.String {
			vField.Set(reflect.ValueOf(c.Request.Header.Values(fieldName)))
		} else if vField.Kind() == reflect.String {
			vField.Set(reflect.ValueOf(c.Request.Header.Get(fieldName)))
		}
	}

	return
}
