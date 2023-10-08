package httpx

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"github.com/zedisdog/brynn/errx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

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

//go:generate go test -bench=Benchmark.+Parse

func BenchmarkSelfParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		j := `{"c": [1,2,3]}`
		r := httptest.NewRequest(http.MethodPost, "https://www.baidu.com?a=1", strings.NewReader(j))
		r.Header.Add("Content-Type", "application/json")
		r.Header.Add("b", "1")

		ctx := Context{
			r: r,
		}

		type request struct {
			A string `form:"a"`
			B string `header:"b"`
			C []int  `json:"c"`
		}
		var req request
		err := ctx.Parse(&req)
		require.Nil(b, err)
	}
}

func BenchmarkGozeroParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		j := `{"c": [1,2,3]}`
		r := httptest.NewRequest(http.MethodPost, "https://www.baidu.com?a=1", strings.NewReader(j))
		r.Header.Add("Content-Type", "application/json")
		r.Header.Add("b", "1")

		type request struct {
			A string `form:"a"`
			B string `header:"b"`
			C []int  `json:"c"`
		}
		var req request
		err := httpx.Parse(r, &req)
		require.Nil(b, err)
	}
}
