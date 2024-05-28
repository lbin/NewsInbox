package user

import (
	"learnerai/golib/app/context"
	"learnerai/golib/app/logger"
	"learnerai/models/mgodb/news"
	"learnerai/models/mgodb/user"
	"learnerai/req"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUserNews(ctx *gin.Context, request req.ListUserReq) (recordList []*user.UserNews) {
	condition := user.ListUserNewsCondition{
		Start: request.Start,
		Limit: request.Limit,
		Uid:   context.GetUserID(ctx),
	}
	if condition.Limit == 0 {
		condition.Limit = 1
	}
	recordList = user.ListUserNews(ctx, condition)
	return
}

func UpdateUserNews(ctx *gin.Context, request req.UpdateUserNewsReq) user.UserNews {
	var userNews user.UserNews
	//1.查询用户信息，如果存在记录，在原有信息上更新；如果不存在记录，插入记录更新 todo 此处可能存在并发问题
	condition := user.ListUserNewsCondition{
		Start: 0,
		Limit: 1,
		Uid:   context.GetUserID(ctx),
	}
	recordList := user.ListUserNews(ctx, condition)
	userNews.Uid = context.GetUserID(ctx)
	userNews.UpdateDate = time.Now()
	if request.Like {
		userNews.Like = []string{request.DocId}
	}
	if request.UnLike {
		userNews.UnLike = []string{request.DocId}
	}
	if request.Read {
		userNews.Read = []string{request.DocId}
	}
	var err error
	if len(recordList) > 0 {
		userNews.Like, userNews.UnLike = GetLikes(request, recordList[0].Like, recordList[0].UnLike)
		userNews.Read = MergeStringArray([]string{request.DocId}, recordList[0].Read)
		err = user.UpdateUserNews(ctx, userNews)
	} else {
		err = user.AddUserNews(ctx, userNews)
	}

	if err != nil {
		logger.Fatalf(ctx, "UpdateUserNews Failed [result:%+v] [err=%+v]", userNews, err)
	}
	//更新文章阅读数，更新文章喜欢数
	incNewsCnt(ctx, request)
	return userNews
}

func GetLikes(request req.UpdateUserNewsReq, likes, unlikes []string) ([]string, []string) {

	if request.Like {
		likes = append(likes, request.DocId)
		//delete unlike doc id
		unlikes = DeleteId(request.DocId, unlikes)
	}
	if request.UnLike {
		unlikes = append(unlikes, request.DocId)
		//delete like doc id
		likes = DeleteId(request.DocId, likes)
	}

	return likes, unlikes
}

func DeleteId(id string, in1 []string) []string {
	//被删除的id可能不在数组中
	for i, str := range in1 {
		if str == in1[i] {
			return append(in1[:i], in1[i+1:]...)
		}
	}
	return in1
}

func incNewsCnt(ctx *gin.Context, request req.UpdateUserNewsReq) {
	var incInfo news.IncNewsCondition
	incInfo.DocId = request.DocId
	incInfo.ReadCnt = 1
	if request.Like {
		incInfo.LikeCnt = 1
	}
	err := news.IncNewsCount(ctx, incInfo)
	if err != nil {
		logger.Fatalf(ctx, "UpdatenewsInfo Failed [result:%+v] [err=%+v]", request, err)
	}
}
