package request

import (
	"context"
	"learnerai/golib/common"

	"github.com/gin-gonic/gin"
)

type ErrorWithCode interface {
	Code() int32
}

func JsonController[Request any, Response any](handler func(context.Context, Request) (Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appG := Gin{C: ctx}
		var rin Request
		err := ctx.BindJSON(&rin)
		if err != nil {
			appG.ResponseError(common.InvalidParams, err.Error())
			return
		}
		rout, err := handler(ctx, rin)
		if err != nil {
			code := int32(common.ERROR)
			if ec, ok := err.(ErrorWithCode); ok {
				code = ec.Code()
			}
			appG.ResponseError(code, err.Error())
			return
		}
		appG.ResponseSuccess(rout)
	}
}
