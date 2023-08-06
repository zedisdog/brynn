package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	gin.SetMode("release")
}

func TestShouldBindInJsonBody(t *testing.T) {
	type req struct {
		A int `json:"a"`
	}

	router := gin.Default()
	router.POST("/test", func(c *gin.Context) {
		var r req
		err := ShouldBind(c, &r)
		if err != nil {
			c.String(200, "error:%s", err.Error())
		} else {
			c.String(200, "%d", r.A)
		}
	})

	w := httptest.NewRecorder()
	json := `{"a":1}`
	r, _ := http.NewRequest("POST", "/test", strings.NewReader(json))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)
	require.Equal(t, "1", w.Body.String())
}

func TestShouldBindInQuery(t *testing.T) {
	type Req struct {
		A int `form:"a"`
	}

	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		var r Req
		err := ShouldBind(c, &r)
		if err != nil {
			c.String(200, "error:%s", err.Error())
		} else {
			c.String(200, "%d", r.A)
		}
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test?a=1", strings.NewReader(""))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)
	require.Equal(t, "1", w.Body.String())
}
