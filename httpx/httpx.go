package httpx

import (
	"encoding/json"
	"errors"
	"github.com/zedisdog/brynn/errx"
	"github.com/zedisdog/brynn/i18n"
	"github.com/zedisdog/brynn/util/reflectx"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
)

var MaxFormSize int64 = 0 * 1024 * 1024

type Context struct {
	r *http.Request
	w http.ResponseWriter
}

func (c *Context) Parse(v any) (err error) {
	destValue := reflect.ValueOf(v)
	if destValue.Kind() != reflect.Ptr {
		err = errx.New(errx.InternalError, "pointer required")
		return
	}
	destValue = destValue.Elem()
	err = c.parseHeader(destValue)
	if err != nil {
		return
	}
	err = c.parseCookies(destValue)
	if err != nil {
		return
	}

	switch ContentType(c.r.Header.Get("Content-Type")) {
	case ContentTypeForm:
		err = c.parseForm(destValue)
	case ContentTypeMultiPartForm:
		err = c.parseMultiPartForm(destValue)
	case ContentTypeJson:
		err = c.parseForm(destValue)
		err = c.parseJson(destValue)
	}
	return
}

func (c *Context) parseHeader(destValue reflect.Value) (err error) {
	destType := destValue.Type()
	for i := 0; i < destValue.NumField(); i++ {
		fieldStruct := destType.Field(i)
		fieldValue := destValue.Field(i)
		if fieldStruct.Anonymous {
			err = c.parseHeader(fieldValue)
			if err != nil {
				return
			}
			continue
		}

		content := reflectx.GetTag(fieldStruct, "header")
		if content != "" {
			arr := strings.Split(content, ",")
			values := c.r.Header.Values(arr[0])
			switch len(values) {
			case 0:
				if !isOptional(arr[1:]) {
					err = errx.New(errx.ValidateError, i18n.Transf("field [:field] is required", i18n.P{"field": arr[0]}))
				}
			case 1:
				var v reflect.Value
				v, err = reflectx.ConvertStringTo(reflect.ValueOf(values[0]), fieldValue.Kind())
				if err != nil {
					return
				}
				fieldValue.Set(v)
			default:
				if fieldValue.Kind() == reflect.Slice {
					for _, item := range values {
						var v reflect.Value
						v, err = reflectx.ConvertStringTo(reflect.ValueOf(item), fieldValue.Type().Elem().Kind())
						if err != nil {
							return
						}
						fieldValue.Set(reflect.Append(fieldValue, v))
					}
				} else {
					panic(errors.New("unsupported"))
				}
			}
		}
	}

	return
}

func (c *Context) parseCookies(destValue reflect.Value) (err error) {
	destType := destValue.Type()
	for i := 0; i < destValue.NumField(); i++ {
		fieldStruct := destType.Field(i)
		fieldValue := destValue.Field(i)
		if fieldStruct.Anonymous {
			err = c.parseCookies(fieldValue)
			if err != nil {
				return
			}
			continue
		}

		content := reflectx.GetTag(fieldStruct, "cookie")
		if content != "" {
			arr := strings.Split(content, ",")
			var cookie *http.Cookie
			cookie, err = c.r.Cookie(arr[0])

			if err != nil {
				if !errors.Is(err, http.ErrNoCookie) {
					return
				} else {
					if !isOptional(arr[1:]) {
						err = errx.New(errx.ValidateError, i18n.Transf("field [:field] is required", i18n.P{"field": arr[0]}))
						return
					}
					err = nil
					continue
				}
			}
			var v reflect.Value
			v, err = reflectx.ConvertStringTo(reflect.ValueOf(cookie.Value), fieldValue.Kind())
			if err != nil {
				return
			}
			fieldValue.Set(v)
		}
	}

	return
}

func (c *Context) parseForm(destValue reflect.Value) (err error) {
	if err = c.r.ParseForm(); err != nil {
		return
	}

	destType := destValue.Type()
	for i := 0; i < destValue.NumField(); i++ {
		fieldStruct := destType.Field(i)
		fieldValue := destValue.Field(i)
		if fieldStruct.Anonymous {
			err = c.parseForm(fieldValue)
			if err != nil {
				return
			}
			continue
		}

		content := reflectx.GetTag(fieldStruct, "form")
		if content != "" {
			arr := strings.Split(content, ",")
			values, ok := c.r.Form[arr[0]]
			if !ok || len(values) == 0 {
				if !isOptional(arr[1:]) {
					err = errx.New(errx.ValidateError, i18n.Transf("field [:field] is required", i18n.P{"field": arr[0]}))
					return
				}
				continue
			}
			switch len(values) {
			case 1:
				var v reflect.Value
				v, err = reflectx.ConvertStringTo(reflect.ValueOf(values[0]), fieldValue.Kind())
				if err != nil {
					return
				}
				fieldValue.Set(v)
			default:
				if fieldValue.Kind() == reflect.Slice {
					for _, item := range values {
						var v reflect.Value
						v, err = reflectx.ConvertStringTo(reflect.ValueOf(item), fieldValue.Type().Elem().Kind())
						if err != nil {
							return
						}
						fieldValue.Set(reflect.Append(fieldValue, v))
					}
				} else {
					panic(errors.New("unsupported"))
				}
			}
		}
	}

	return
}

func (c *Context) parseMultiPartForm(destValue reflect.Value) (err error) {
	if err = c.parseForm(destValue); err != nil {
		return
	}

	if err = c.r.ParseMultipartForm(MaxFormSize); err != nil {
		return
	}

	destType := destValue.Type()
	for i := 0; i < destValue.NumField(); i++ {
		fieldStruct := destType.Field(i)
		fieldValue := destValue.Field(i)
		if fieldStruct.Anonymous {
			err = c.parseMultiPartForm(fieldValue)
			if err != nil {
				return
			}
			continue
		}

		content := reflectx.GetTag(fieldStruct, "file")
		if content != "" {
			arr := strings.Split(content, ",")
			values, ok := c.r.MultipartForm.File[arr[0]]
			if !ok || len(values) == 0 {
				if !isOptional(arr[1:]) {
					err = errx.New(errx.ValidateError, i18n.Transf("field [:field] is required", i18n.P{"field": arr[0]}))
					return
				}
				continue
			}
			switch len(values) {
			case 1:
				if _, ok := fieldValue.Interface().(*multipart.FileHeader); ok {
					fieldValue.Set(reflect.ValueOf(values[0]))
					continue
				}
				var (
					f multipart.File
					c []byte
				)
				f, err = values[0].Open()
				if err != nil {
					return
				}
				c, err = io.ReadAll(f)
				if err != nil {
					return
				}
				switch fieldValue.Interface().(type) {
				case []byte:
					fieldValue.Set(reflect.ValueOf(c))
				case string:
					fieldValue.Set(reflect.ValueOf(string(c)))
				default:
					panic(errors.New("unsupported"))
				}
			default:
				if _, ok := fieldValue.Interface().([]*multipart.FileHeader); ok {
					fieldValue.Set(reflect.ValueOf(values))
					continue
				}

				contents := make([][]byte, 0, len(values))
				for _, item := range values {
					var (
						f multipart.File
						c []byte
					)
					f, err = item.Open()
					if err != nil {
						return
					}
					c, err = io.ReadAll(f)
					if err != nil {
						return
					}
					contents = append(contents, c)
				}
				switch fieldValue.Interface().(type) {
				case [][]byte:
					fieldValue.Set(reflect.ValueOf(contents))
				case []string:
					for _, item := range contents {
						reflect.Append(fieldValue, reflect.ValueOf(string(item)))
					}
				}
			}
		}
	}

	return
}

func (c *Context) parseJson(destValue reflect.Value) (err error) {
	content, err := io.ReadAll(c.r.Body)
	if err != nil {
		return
	}
	req := make(map[string]any)
	err = json.Unmarshal(content, &req)
	if err != nil {
		return
	}
	return reflectx.UnmarshalMap(req, destValue, "json")
}
