package user

import (
	"learnerai/golib/app/context"
	"learnerai/golib/app/logger"
	"learnerai/models/mgodb/paper"
	"learnerai/models/mgodb/user"
	"learnerai/req"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUserPaper(ctx *gin.Context, request req.ListUserReq) (recordList []*user.UserPaper) {
	condition := user.ListUserPaperCondition{
		Start: request.Start,
		Limit: request.Limit,
		Uid:   context.GetUserID(ctx),
	}
	if condition.Limit == 0 {
		condition.Limit = 1
	}
	recordList = user.ListUserPaper(ctx, condition)
	return
}

func UpdateUserPaper(ctx *gin.Context, request req.UpdateUserPaperReq) user.UserPaper {
	var userPaper user.UserPaper
	//1.查询用户信息，如果存在记录，在原有信息上更新；如果不存在记录，插入记录更新 todo 此处可能存在并发问题
	condition := user.ListUserPaperCondition{
		Start: 0,
		Limit: 1,
		Uid:   context.GetUserID(ctx),
	}
	recordList := user.ListUserPaper(ctx, condition)
	userPaper.Uid = context.GetUserID(ctx)
	userPaper.UpdateDate = time.Now()
	if request.Like {
		userPaper.Like = []string{request.DocId}
	}
	if request.UnLike {
		userPaper.UnLike = []string{request.DocId}
	}
	if request.Read {
		userPaper.Read = []string{request.DocId}
	}
	var err error
	if len(recordList) > 0 {
		userPaper.Like, userPaper.UnLike = GetLikes(request, recordList[0].Like, recordList[0].UnLike)
		userPaper.Read = MergeStringArray([]string{request.DocId}, recordList[0].Read)
		err = user.UpdateUserPaper(ctx, userPaper)
	} else {
		err = user.AddUserPaper(ctx, userPaper)
	}

	if err != nil {
		logger.Fatalf(ctx, "UpdateUserPaper Failed [result:%+v] [err=%+v]", userPaper, err)
	}
	//更新文章阅读数，更新文章喜欢数
	incPaperCnt(ctx, request)
	return userPaper
}

func GetLikes(request req.UpdateUserPaperReq, likes, unlikes []string) ([]string, []string) {

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

func incPaperCnt(ctx *gin.Context, request req.UpdateUserPaperReq) {
	var incInfo paper.IncPaperCondition
	incInfo.DocId = request.DocId
	incInfo.ReadCnt = 1
	if request.Like {
		incInfo.LikeCnt = 1
	}
	err := paper.IncPaperCount(ctx, incInfo)
	if err != nil {
		logger.Fatalf(ctx, "UpdatepaperInfo Failed [result:%+v] [err=%+v]", request, err)
	}
}
