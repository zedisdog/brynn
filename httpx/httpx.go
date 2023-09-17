package httpx

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/zedisdog/brynn/util/mapx"
)

var MaxFormSize int64 = 0 * 1024 * 1024

type Context struct {
	r *http.Request
	w http.ResponseWriter
}

func (c *Context) Parse(v any) (err error) {
	var values map[string][]any

	values = mapx.Merge(values, c.readHeader())
	values = mapx.Merge(values, c.readCookies())

	contentType := ContentType(c.r.Header.Get("Content-Type"))
	switch contentType {
	case ContentTypeMultiPartForm:
		var form map[string][]any
		form, err = c.readMultiForm()
		if err != nil {
			return
		}
		values = mapx.Merge(values, form)
	default:
		var form map[string][]any
		form, err = c.readForm()
		if err != nil {
			return
		}
		values = mapx.Merge(values, form)
	}

	if contentType == ContentTypeJson {
		var result map[string][]any
		result, err = c.readJson()
		values = mapx.Merge(values, result)

	}
	return
}

func (c *Context) readHeader() (result map[string][]any) {
	result = make(map[string][]any, len(c.r.Header))
	for k, values := range c.r.Header {
		result[k] = gconv.SliceAny(values)
	}

	return
}

func (c *Context) readCookies() (result map[string][]any) {
	cookies := c.r.Cookies()
	result = make(map[string][]any, len(cookies))
	for _, cookie := range cookies {
		result[cookie.Name] = []any{cookie.Value}
	}

	return
}

func (c *Context) readForm() (result map[string][]any, err error) {
	if err = c.r.ParseForm(); err != nil {
		return
	}
	result = make(map[string][]any, len(c.r.Form))
	for k, v := range c.r.Form {
		result[k] = gconv.SliceAny(v)
	}

	return
}

func (c *Context) readMultiForm() (result map[string][]any, err error) {
	if err = c.r.ParseMultipartForm(MaxFormSize); err != nil {
		return
	}
	result = make(map[string][]any, len(c.r.Form)+len(c.r.MultipartForm.File))
	for k, v := range c.r.Form {
		result[k] = gconv.SliceAny(v)
	}

	for k, files := range c.r.MultipartForm.File {
		result[k] = gconv.SliceAny(files)
	}

	return
}

func (c *Context) readJson() (result map[string][]any, err error) {
	content, err := io.ReadAll(c.r.Body)
	if err != nil {
		return
	}
	var v any
	err = json.Unmarshal(content, &v)
	if err != nil {
		return
	}

	result = make(map[string][]any, 1)
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Slice {
		result["body"] = v.([]any)
	}

	if value.Kind() == reflect.Map {
		for k, v := range v.(map[string]any) {
			result[k] = []any{v}
		}
	}

	return
}
