package httpx

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/zedisdog/brynn/errx"
	"mime"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestXxx(t *testing.T) {
	types, _ := mime.ExtensionsByType("application/x-www-form-urlencoded")
	fmt.Printf("%+v\n", types)
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
