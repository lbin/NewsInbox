package news

import (
	"learnerai/common"
	goLabApp "learnerai/golib/app/request"
	"learnerai/logic/news"
	"learnerai/req"

	"github.com/gin-gonic/gin"
)

func SearchSimilar(ctx *gin.Context) {
	appG := goLabApp.Gin{C: ctx}
	var request req.SearchSimilarReq
	err := ctx.BindJSON(&request)
	if err != nil {
		appG.ResponseError(common.InvalidParams, err.Error())
		return
	}
	data, err := news.SearchSimilar(ctx, request.DocId, request.TopK)
	if err != nil {
		appG.ResponseError(common.ERROR, err.Error())
		return
	}
	appG.ResponseSuccess(data)
}
