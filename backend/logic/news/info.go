package news

import (
	"learnerai/models/mgodb/news"
	"learnerai/req"

	"github.com/gin-gonic/gin"
)

func GetNewsInfo(ctx *gin.Context, request req.ListNewsInfoReq) (recordList []*news.NewsInfo) {
	condition := news.ListNewsInfoCondition{
		Id:            request.Id,
		Tags:          request.Tags,
		StartTime:     request.StartTime,
		EndTime:       request.EndTime,
		Start:         request.Start,
		Limit:         request.Limit,
		SearchContent: request.SearchContent,
	}
	recordList = news.ListNewsInfo(ctx, condition)
	return
}

func GetDailyNewsInfo(ctx *gin.Context, request req.ListNewsInfoReq) (recordList []*news.DailyInfo) {
	condition := news.ListDailyInfoCondition{
		Tags:          request.Tags,
		StartTime:     request.StartTime,
		Id:            request.Id,
		EndTime:       request.EndTime,
		Start:         request.Start,
		Limit:         request.Limit,
		SearchContent: request.SearchContent,
	}
	recordList = news.ListDailyInfo(ctx, condition)
	return
}

func GetDailyUserNewsInfo(ctx *gin.Context, request req.ListUserNewsInfoReq) (recordList []*news.DailyInfo) {
	condition := news.ListDailyUserInfoCondition{
		StartTime: request.StartTime,
		Id:        request.Id,
		EndTime:   request.EndTime,
		Start:     request.Start,
		Limit:     request.Limit,
	}
	recordList = news.ListDailyUserInfo(ctx, condition)
	return
}
