package paper

import (
	"learnerai/common"
	goLabApp "learnerai/golib/app/request"
	"learnerai/logic/paper"
	"learnerai/req"

	"github.com/gin-gonic/gin"
)

func UnsetDailyInfoField(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.UnsetDailyInfoField
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	err = paper.UnsetDailyInfoField(ctx, request)
	if err != nil {
		appG.ResponseError(common.ERROR, err.Error())
		return
	}
	appG.ResponseSuccess(nil)
}
