package user

import (
	"encoding/json"
	"io/ioutil"
	"learnerai/golib/app/logger"
	"learnerai/logic/auth"
	"learnerai/models/mgodb/user"
	"learnerai/req"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat"
)

/*
WxMiniConf:

	AppId: wx67ffeb75fc13b242
	Secret: 2c810016ffb7c2ce97af3f436c4d204e

参考文档：https://silenceper.com/wechat/miniprogram/auth.html
*/
func GetUserMiniProgramInfo(ctx *gin.Context, request req.GetMiniProgramCode) (res req.GetMiniProgramCodeRes, err error) {
	//替换为您的小程序配置信息
	config := &wechat.Config{
		AppID:     "wx67ffeb75fc13b242",
		AppSecret: "2c810016ffb7c2ce97af3f436c4d204e",
	}
	wc := wechat.NewWechat(config)

	//使用code获取用户信息
	mini := wc.GetMiniProgram()
	result, err := mini.Code2Session(request.Code)
	if err != nil {
		logger.Fatalf(ctx, "UpdateUserInfo Failed [result:%+v] [err=%+v]", result, err)
		return
	}
	logger.Infof(ctx, "UpdateUserInfo Failed [result:%+v]", result)
	condition := user.ListUserInfoCondition{
		Start:  0,
		Limit:  1,
		OpenId: result.OpenID,
	}

	//查询用户信息，如果不存在，重新插入一条
	recordList := user.ListUserInfo(ctx, condition)
	var uid string
	if len(recordList) == 0 {
		var userInfo user.UserInfo
		userInfo.FocusArea = []string{}
		userInfo.Uid = GenUidId()
		uid = userInfo.Uid
		userInfo.OpenId = result.OpenID
		userInfo.UnionId = result.UnionID
		userInfo.FollowAuthor = []string{}
		userInfo.FollowTag = []string{}
		err = user.AddUserInfo(ctx, userInfo)
		if err != nil {
			logger.Fatalf(ctx, "UpdateUserInfo Failed [result:%+v] [err=%+v]", userInfo, err)
			return
		}

		// // TODO 问题，这里新建表，后续请求写不进去，但是重新再登录一次就可以写进去
		// var userNews user.UserNews
		// userNews.Uid = userInfo.Uid
		// userNews.UpdateDate = time.Now()
		// userNews.Like = []string{}
		// userNews.UnLike = []string{}
		// userNews.Read = []string{}
		// userNews.Upload = []string{}
		// err = user.AddUserNews(ctx, userNews)
		// if err != nil {
		// 	logger.Fatalf(ctx, "UpdateUserNews Failed [result:%+v] [err=%+v]", userNews, err)
		// 	return
		// }
		// time.Sleep(1 * time.Second) // Add a delay of 1 second

	} else {
		uid = recordList[0].Uid
	}

	//生成用户token，用于登录校验
	var tokenInfo auth.TokenInfo
	tokenInfo.Uid = uid
	res.Token = auth.CreateToken(auth.SignToken, tokenInfo)

	//获取access token，返回前端小程序
	res.AccessToken = GetAccessToken(ctx)
	return
}

func GenUidId() string {
	return strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(rand.Int63(), 10)
}

func GetToken(ctx *gin.Context, uid string) (res req.GetMiniProgramCodeRes) {
	res.AccessToken = GetAccessToken(ctx)
	//生成用户token，用于登录校验
	var tokenInfo auth.TokenInfo
	tokenInfo.Uid = uid
	res.Token = auth.CreateToken(auth.SignToken, tokenInfo)
	condition := user.ListUserInfoCondition{
		Start: 0,
		Limit: 1,
		Uid:   uid,
	}
	//查询用户信息，如果不存在，重新插入一条
	recordList := user.ListUserInfo(ctx, condition)
	logger.Infof(ctx, "UpdateUserInfo Failed [result:%+v]", recordList)
	return
}

const BaseAuthUrl = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx67ffeb75fc13b242&secret=2c810016ffb7c2ce97af3f436c4d204e"

type accessTokenInfo struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func GetAccessToken(ctx *gin.Context) (accessToken string) {
	response, err := http.Get(BaseAuthUrl)
	if err != nil {
		logger.Fatalf(ctx, "err = [%+v]", err)
		return ""
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Fatalf(ctx, "err = [%+v]", err)
		return ""
	}
	info := accessTokenInfo{}
	_ = json.Unmarshal(data, &info)
	logger.Infof(ctx, "tencent access token return [%+v]", string(data))
	return info.AccessToken
}
