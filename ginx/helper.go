package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ShouldBind(ctx *gin.Context, s any) (err error) {
	err = httpx.Parse(ctx.Request, s)
	return
}
