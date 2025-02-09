package user

import (
	"learnerai/common"
	user "learnerai/logic/user"
	"learnerai/req"
	"net/http"

	goLabApp "learnerai/golib/app/request"

	"github.com/gin-gonic/gin"
)

func GetUserNewsRecord(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.ListUserReq
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data := user.GetUserNews(ctx, request)
	appG.Response(http.StatusOK, common.SUCCESS, data)
}

func UpdateUserNewsRecord(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.UpdateUserNewsReq
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data := user.UpdateUserNews(ctx, request)
	appG.Response(http.StatusOK, common.SUCCESS, data)
}
