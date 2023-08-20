package ginx

import (
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
	//TODO
	return
}
