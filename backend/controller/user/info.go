package user

import (
	"learnerai/common"
	user "learnerai/logic/user"
	"learnerai/req"
	"net/http"

	goLabApp "learnerai/golib/app/request"

	"github.com/gin-gonic/gin"
)

func GetUserInfoRecord(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.ListUserReq
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data := user.GetUserInfo(ctx, request)
	appG.Response(http.StatusOK, common.SUCCESS, data)
}

func UpdateUserInfoRecord(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.UpdateUserInfoReq
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data := user.UpdateUserInfo(ctx, request)
	appG.Response(http.StatusOK, common.SUCCESS, data)
}

func GetMiniProgramInfo(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.GetMiniProgramCode
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data, err := user.GetUserMiniProgramInfo(ctx, request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	appG.Response(http.StatusOK, common.SUCCESS, data)
}

func GetToken(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.GetTokenReq
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data := user.GetToken(ctx, request.Uid)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	appG.Response(http.StatusOK, common.SUCCESS, data)
}
