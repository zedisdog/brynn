package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/zedisdog/brynn/errx"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func convertTypeOfPtr(tp reflect.Type, target reflect.Value) reflect.Value {
	// keep the original value is a pointer
	if tp.Kind() == reflect.Ptr && target.CanAddr() {
		tp = tp.Elem()
		target = target.Addr()
	}

	for tp.Kind() == reflect.Ptr {
		p := reflect.New(target.Type())
		p.Elem().Set(target)
		target = p
		tp = tp.Elem()
	}

	return target
}

func TestXxx(t *testing.T) {
	a := 1
	fmt.Printf("%p\n", &a)
	res := convertTypeOfPtr(reflect.TypeOf(&a), reflect.ValueOf(a))
	fmt.Printf("%#v\n", res.Interface())
}

func TestXxx2(t *testing.T) {
	j := `{"a": 1, "b": 2, "c": [1,2,3]}`
	var v any
	json.Unmarshal([]byte(j), &v)
	fmt.Printf("%#v\n", v.(map[string]any)["c"])
}

func TestParseHeader(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "https://www.baidu.com", nil)
	r.Header.Add("a", "1")
	r.Header.Add("b", "1")
	r.Header.Add("b", "2")

	ctx := &Context{
		r: r,
	}

	type internal struct {
		A int `header:"a"`
	}
	type testStruct struct {
		A string   `header:"a"`
		B int      `header:"A"`
		C []string `header:"b"`
		D []int    `header:"B"`
		internal
	}
	var tst testStruct
	err := ctx.parseHeader(reflect.ValueOf(&tst).Elem())
	require.Nil(t, err)
	require.Equal(t, "1", tst.A)
	require.Equal(t, 1, tst.B)
	require.Equal(t, []string{"1", "2"}, tst.C)
	require.Equal(t, []int{1, 2}, tst.D)
	require.Equal(t, 1, tst.internal.A)

	type testStruct2 struct {
		A string `header:"c"`
	}
	var tst2 testStruct2
	err = ctx.parseHeader(reflect.ValueOf(&tst2).Elem())
	require.NotNil(t, err)
	require.Equal(t, errx.ValidateError, err.(*errx.Error).Code)
}

func TestParseCookies(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "https://www.baidu.com", nil)
	r.AddCookie(&http.Cookie{Name: "a", Value: "1"})

	ctx := &Context{
		r: r,
	}

	type internal struct {
		A int `cookie:"a"`
	}
	type testStruct struct {
		A string `cookie:"a"`
		B int    `cookie:"a"`
		internal
	}
	var tst testStruct
	err := ctx.parseCookies(reflect.ValueOf(&tst).Elem())
	require.Nil(t, err)
	require.Equal(t, "1", tst.A)
	require.Equal(t, 1, tst.B)
	require.Equal(t, 1, tst.internal.A)

	type testStruct2 struct {
		A string `cookie:"c"`
	}
	var tst2 testStruct2
	err = ctx.parseCookies(reflect.ValueOf(&tst2).Elem())
	require.NotNil(t, err)
	require.Equal(t, errx.ValidateError, err.(*errx.Error).Code)
}

func TestParseForm(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "https://www.baidu.com?a=1", bytes.NewReader([]byte("b=1&c=1&c=2")))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	ctx := &Context{
		r: r,
	}

	type internal struct {
		A []int `form:"c"`
	}
	type testStruct struct {
		A string `form:"b"`
		B int    `form:"a"`
		internal
	}
	var tst testStruct
	err := ctx.parseForm(reflect.ValueOf(&tst).Elem())
	require.Nil(t, err)
	require.Equal(t, "1", tst.A)
	require.Equal(t, 1, tst.B)
	require.Equal(t, []int{1, 2}, tst.internal.A)

	type testStruct2 struct {
		A string `form:"d"`
	}
	var tst2 testStruct2
	err = ctx.parseForm(reflect.ValueOf(&tst2).Elem())
	require.NotNil(t, err)
	require.Equal(t, errx.ValidateError, err.(*errx.Error).Code)
}
