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
	logger.Infof(ctx, "UpdateUserNewsReq:%+v", request)
	if request.Like {
		userNews.Like = []string{request.DocId}
	} else {
		userNews.Like = []string{}
	}
	if request.UnLike {
		userNews.UnLike = []string{request.DocId}
	} else {
		userNews.UnLike = []string{}
	}
	if request.Read {
		userNews.Read = []string{request.DocId}
	} else {
		userNews.Read = []string{}
	}
	if request.Upload {
		userNews.Upload = []string{request.DocId}
	} else {
		userNews.Upload = []string{}
	}
	var err error
	if len(recordList) > 0 {
		userNews.Like, userNews.UnLike = GetLikes(request, recordList[0].Like, recordList[0].UnLike)
		logger.Infof(ctx, "userNewsLike:%+v", userNews.Like)
		// userNews.Read = MergeStringArray([]string{request.DocId}, recordList[0].Read)
		// userNews.Upload = MergeStringArray(userNews.Upload, recordList[0].Upload)
		userNews.Read = GetUpload(request, recordList[0].Read)
		userNews.Upload = GetUpload(request, recordList[0].Upload)
		logger.Infof(ctx, "UpdateUserNewsReq:%+v", userNews)
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

func GetRead(request req.UpdateUserNewsReq, read []string) []string {

	if request.Read {
		read = AddId(request.DocId, read)
	} else {
		read = DeleteId(request.DocId, read)
	}

	return read
}

func GetUpload(request req.UpdateUserNewsReq, upload []string) []string {

	if request.Upload {
		upload = AddId(request.DocId, upload)
	} else {
		upload = DeleteId(request.DocId, upload)
	}

	return upload
}

func GetLikes(request req.UpdateUserNewsReq, likes, unlikes []string) ([]string, []string) {

	if request.Like {
		likes = AddId(request.DocId, likes)
	} else {
		likes = DeleteId(request.DocId, likes)
	}

	if request.UnLike {
		unlikes = AddId(request.DocId, unlikes)
	} else {
		unlikes = DeleteId(request.DocId, unlikes)
	}

	return likes, unlikes
}

func AddId(id string, in1 []string) []string {
	for _, str := range in1 {
		if str == id {
			return in1
		}
	}
	return append(in1, id)
}

func DeleteId(id string, in1 []string) []string {
	for i, str := range in1 {
		if str == id {
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
