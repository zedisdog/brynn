package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestNormal(t *testing.T) {
	u, err := url.Parse("http://www.baidu.com?a=1&a=2")
	require.Nil(t, err)
	r := http.Request{
		URL: u,
	}

	err = r.ParseForm()
	require.Nil(t, err)

	fmt.Printf("%+v\n", r.Form)

	type res struct {
		A []int `form:"a"`
	}
	var req res

	err = httpx.Parse(&r, &req)
	require.Nil(t, err)

	fmt.Printf("%+v\n", req)
}

func TestContext_ParseHeader(t *testing.T) {
	{
		type testStruct struct {
			A string `header:"a"`
		}

		testReq := httptest.NewRequest(http.MethodGet, "http://www.baidu.com", nil)
		testReq.Header.Set("a", "123")
		c := &context{
			Context: &gin.Context{
				Request: testReq,
			},
		}

		var data testStruct

		err := c.parseHeader(reflect.ValueOf(&data).Elem(), reflect.TypeOf(&data).Elem())
		require.Nil(t, err)
		require.Equal(t, "123", data.A)
	}

	{
		type testStruct struct {
			A []string `header:"a"`
		}

		testReq := httptest.NewRequest(http.MethodGet, "http://www.baidu.com", nil)
		testReq.Header.Add("a", "123")
		testReq.Header.Add("a", "321")
		c := &context{
			Context: &gin.Context{
				Request: testReq,
			},
		}

		var data testStruct

		err := c.parseHeader(reflect.ValueOf(&data).Elem(), reflect.TypeOf(&data).Elem())
		require.Nil(t, err)
		require.Equal(t, "123", data.A[0])
		require.Equal(t, "321", data.A[1])
	}

	{
		type testStruct struct {
			a []string `header:"a"`
		}

		testReq := httptest.NewRequest(http.MethodGet, "http://www.baidu.com", nil)
		testReq.Header.Add("a", "123")
		testReq.Header.Add("a", "321")
		c := &context{
			Context: &gin.Context{
				Request: testReq,
			},
		}

		var data testStruct

		err := c.parseHeader(reflect.ValueOf(&data).Elem(), reflect.TypeOf(&data).Elem())
		require.Nil(t, err)
		require.Equal(t, "123", data.a[0])
		require.Equal(t, "321", data.a[1])
	}
}
