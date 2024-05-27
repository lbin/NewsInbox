package paper

import (
	"learnerai/models/mgodb/paper"
	"learnerai/req"

	"github.com/gin-gonic/gin"
)

func GetPaperInfo(ctx *gin.Context, request req.ListPaperInfoReq) (recordList []*paper.PaperInfo) {
	condition := paper.ListPaperInfoCondition{
		Id:            request.Id,
		Category:      request.Category,
		StartTime:     request.StartTime,
		EndTime:       request.EndTime,
		Start:         request.Start,
		Limit:         request.Limit,
		SearchContent: request.SearchContent,
	}
	recordList = paper.ListPaperInfo(ctx, condition)
	return
}

func GetDailyPaperInfo(ctx *gin.Context, request req.ListPaperInfoReq) (recordList []*paper.DailyInfo) {
	condition := paper.ListDailyInfoCondition{
		Category:      request.Category,
		StartTime:     request.StartTime,
		Id:            request.Id,
		EndTime:       request.EndTime,
		Start:         request.Start,
		Limit:         request.Limit,
		SearchContent: request.SearchContent,
	}
	recordList = paper.ListDailyInfo(ctx, condition)
	return
}
