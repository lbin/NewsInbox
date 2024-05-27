package common

import (
	"github.com/gin-gonic/gin"
	"learnerai/golib/app/request"
	"net/http"
	. "learnerai/common"
	mongoCommon "learnerai/models/mgodb/common"
)

func CreateIndexes(ctx *gin.Context) {
	appG := request.Gin{C: ctx}
	mongoCommon.CreateIndex(ctx)
	appG.Response(http.StatusOK, SUCCESS, nil)
}

//心跳检测【k8s】
func Health(ctx *gin.Context) {
	appG := request.Gin{C: ctx}
	appG.Response(http.StatusOK, SUCCESS, nil)
	return
}
