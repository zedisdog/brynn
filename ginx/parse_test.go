package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func MakeRequest() *http.Request {
	j := `{"c": [1,2,3]}`
	r := httptest.NewRequest(http.MethodPost, "https://www.baidu.com/?a=1&filters[sort]=1", strings.NewReader(j))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("b", "1")
	return r
}

func TestParse(t *testing.T) {
	type request struct {
		A int               `form:"a"`
		B int               `header:"b"`
		C []int             `json:"c"`
		D map[string]string `form:"filters"`
	}
	var req request

	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		err := Parse(c, &req)
		require.Nil(t, err)
		require.Equal(t, 1, req.B)
		require.Equal(t, 1, req.A)
		require.Equal(t, []int{1, 2, 3}, req.C)
		fmt.Printf("%+v\n", req.D)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, MakeRequest())
}

func BenchmarkParse(b *testing.B) {
	type request struct {
		A int   `form:"a"`
		B int   `header:"b"`
		C []int `json:"c"`
	}
	var req request
	for i := 0; i < b.N; i++ {
		r := gin.Default()
		r.POST("/", func(c *gin.Context) {
			err := Parse(c, &req)
			require.Nil(b, err)
			require.Equal(b, 1, req.B)
			require.Equal(b, 1, req.A)
			require.Equal(b, []int{1, 2, 3}, req.C)
		})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, MakeRequest())
	}
}
