package news

import (
	"learnerai/common"
	"learnerai/logic/news"
	"learnerai/req"
	"net/http"

	goLabApp "learnerai/golib/app/request"

	"github.com/gin-gonic/gin"
)

func GetNewsInfoRecord(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.ListNewsInfoReq
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data := news.GetNewsInfo(ctx, request)
	appG.Response(http.StatusOK, common.SUCCESS, data)
}

func GetDailyNewsInfoRecord(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.ListNewsInfoReq
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data := news.GetDailyNewsInfo(ctx, request)
	appG.Response(http.StatusOK, common.SUCCESS, data)
}
