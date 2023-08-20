package ginx

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"net/url"
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
