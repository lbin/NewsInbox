package router

import (
	"learnerai/controller/common"
	"learnerai/controller/paper"
	s2 "learnerai/controller/semanticscholar"
	"learnerai/controller/user"
	"learnerai/golib/app/middleware"
	"learnerai/golib/app/request"
	"learnerai/logic/auth"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	routerInnerApi := r.Group("/inner/miniapp").Use(request.AddRequestId)
	routerInnerApi.Use()
	{
		routerInnerApi.POST("healthz", common.Health)         //心跳检测
		routerInnerApi.POST("cindexes", common.CreateIndexes) //创建数据库索引
	}
	//业务功能路由
	miniRouterApi := r.Group("/miniapp").Use(request.AddRequestId)
	miniRouterApi.POST("/mini/login", user.GetMiniProgramInfo)
	miniRouterApi.Use(auth.CheckUriToken)
	{
		//用户自身信息接口
		miniRouterApi.POST("/user/info/get", user.GetUserInfoRecord)
		miniRouterApi.POST("/user/info/update", user.UpdateUserInfoRecord)
		//用户与文章的关系接口
		miniRouterApi.POST("/user/paper/get", user.GetUserPaperRecord)
		miniRouterApi.POST("/user/paper/update", user.UpdateUserPaperRecord)
	}

	miniRouterApiWithoutAuth := r.Group("/miniapp").Use(request.AddRequestId)
	miniRouterApiWithoutAuth.Use()
	{
		//文章相关信息接口
		miniRouterApiWithoutAuth.POST("/get/token", user.GetToken)
		miniRouterApiWithoutAuth.POST("/paper/list", paper.GetPaperInfoRecord)
		miniRouterApiWithoutAuth.POST("/daily/list", paper.GetDailyPaperInfoRecord)
		miniRouterApiWithoutAuth.POST("/daily/unset", paper.UnsetDailyInfoField)
		//文章相关元信息接口
		miniRouterApiWithoutAuth.POST("/s2/paper/info", s2.GetPaperInfo)
		miniRouterApiWithoutAuth.POST("/s2/paper/references", s2.GetPaperReference)
		//向量搜索相关接口
		// miniRouterApiWithoutAuth.POST("/vector/search", paper.SearchSimilar)
	}

	return r
}
