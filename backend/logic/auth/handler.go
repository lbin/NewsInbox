package auth

import (
	"learnerai/common"
	"learnerai/golib/app/context"
	"learnerai/golib/app/logger"
	"learnerai/golib/app/request"

	"github.com/gin-gonic/gin"
)

func CheckUriToken(ctx *gin.Context) {
	appG := request.Gin{C: ctx}

	token := context.GetToken(ctx)
	if token == "" {
		appG.ResponseError(common.ErrorCheckTokenFail, nil)
		ctx.Abort()
		return
	}

	userInfo, err := ParseToken(SignToken, token)
	if err != nil {
		appG.ResponseError(common.ErrorCheckTokenFail, nil)
		ctx.Abort()
		return
	}

	// 优先存储的内容
	context.SetToken(ctx, token)
	//电话号码加密
	context.SetUserID(ctx, userInfo.Uid)
	logger.Infof(ctx, "CheckUriToken succeeded, token: %v, info: %+v", token, userInfo)
	ctx.Next()
}
