package paper

import (
	"learnerai/common"
	"learnerai/logic/paper"
	"learnerai/req"
	"net/http"

	goLabApp "learnerai/golib/app/request"

	"github.com/gin-gonic/gin"
)

func GetPaperInfoRecord(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.ListPaperInfoReq
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data := paper.GetPaperInfo(ctx, request)
	appG.Response(http.StatusOK, common.SUCCESS, data)
}

func GetDailyPaperInfoRecord(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.ListPaperInfoReq
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data := paper.GetDailyPaperInfo(ctx, request)
	appG.Response(http.StatusOK, common.SUCCESS, data)
}