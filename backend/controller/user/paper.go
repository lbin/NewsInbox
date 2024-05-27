package user

import (
	"learnerai/common"
	user "learnerai/logic/user"
	"learnerai/req"
	"net/http"

	goLabApp "learnerai/golib/app/request"

	"github.com/gin-gonic/gin"
)

func GetUserPaperRecord(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.ListUserReq
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data := user.GetUserPaper(ctx, request)
	appG.Response(http.StatusOK, common.SUCCESS, data)
}

func UpdateUserPaperRecord(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.UpdateUserPaperReq
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data := user.UpdateUserPaper(ctx, request)
	appG.Response(http.StatusOK, common.SUCCESS, data)
}
