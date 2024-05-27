package user

import (
	"learnerai/golib/app/context"
	"learnerai/golib/app/logger"
	"learnerai/models/mgodb/user"
	"learnerai/req"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(ctx *gin.Context, request req.ListUserReq) (recordList []*user.UserInfo) {
	condition := user.ListUserInfoCondition{
		Start: request.Start,
		Limit: request.Limit,
		Uid:   context.GetUserID(ctx),
	}
	if condition.Limit == 0 {
		condition.Limit = 1
	}
	recordList = user.ListUserInfo(ctx, condition)
	return
}

func UpdateUserInfo(ctx *gin.Context, request req.UpdateUserInfoReq) (userInfo user.UserInfo) {

	//1.查询用户信息，如果存在记录，在原有信息上更新；如果不存在记录，插入记录更新 todo 此处可能存在并发问题
	condition := user.ListUserInfoCondition{
		Start: 0,
		Limit: 1,
		Uid:   context.GetUserID(ctx),
	}
	recordList := user.ListUserInfo(ctx, condition)

	var err error
	userInfo.Uid = context.GetUserID(ctx)
	userInfo.UpdateDate = time.Now()
	if len(recordList) > 0 {
		if len(request.FocusArea) > 0 {
			userInfo.FocusArea = MergeStringArray(request.FocusArea, recordList[0].FocusArea)
		}
		if len(request.FollowAuthor) > 0 {
			userInfo.FollowAuthor = MergeStringArray(request.FollowAuthor, recordList[0].FollowAuthor)
		}
		if len(request.FollowTag) > 0 {
			userInfo.FollowTag = MergeStringArray(request.FollowTag, recordList[0].FollowTag)
		}
		err = user.UpdateUserInfo(ctx, userInfo)
	}
	if err != nil {
		logger.Fatalf(ctx, "UpdateUserInfo Failed [result:%+v] [err=%+v]", userInfo, err)
	}
	return userInfo
}

func MergeStringArray(in1, in2 []string) []string {
	var mapStr = make(map[string]bool, 0)
	for _, str := range in1 {
		mapStr[str] = true
	}

	var ans []string
	ans = append(ans, in1...)
	for _, str := range in2 {
		_, ok := mapStr[str]
		if !ok {
			ans = append(ans, str)
		}
	}
	return ans
}
