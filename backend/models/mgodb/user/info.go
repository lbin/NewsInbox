package user

import (
	"learnerai/golib/app/logger"
	. "learnerai/models/mgodb/common"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserInfo struct {
	OpenId       string    `json:"_id" bson:"_id"`
	UnionId      string    `json:"union_id" bson:"union_id"`
	Uid          string    `json:"uid" bson:"uid"`
	UpdateDate   time.Time `json:"update_date" bson:"update_date"`     //用户更新时间
	FollowAuthor []string  `json:"follow_author" bson:"follow_author"` //关注作者列表
	FollowTag    []string  `json:"follow_tag" bson:"follow_tag"`       //关注标签列表
	FocusArea    []string  `json:"focus_area" bson:"focus_area"`       //关注领域列表
}

type ListUserInfoCondition struct {
	Uid    string `json:"uid"`
	OpenId string `json:"openid"`
	Start  int64  `json:"start"`
	Limit  int64  `json:"limit"`
}

func AddUserInfo(ctx *gin.Context, UserInfo UserInfo) error {
	logger.Infof(ctx, "AddUserInfo [condition:%+v]", UserInfo)
	mClient := DB.Mongo
	var table = CollectionUserInfo
	collection := mClient.Database(GetDataBase()).Collection(table)
	info, err := collection.InsertOne(GetContext(), UserInfo)

	if err != nil {
		logger.Fatalf(ctx, "AddUserInfo [err=%+v]", err)
	} else {
		logger.Infof(ctx, "AddUserInfo [info=%+v]", info)
	}
	return err
}

func ListUserInfo(ctx *gin.Context, condition ListUserInfoCondition) (recordList []*UserInfo) {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionUserInfo)
	filter := bson.M{}
	if condition.Uid != "" {
		filter["uid"] = condition.Uid
	}
	if condition.OpenId != "" {
		filter["_id"] = condition.OpenId
	}
	sort := bson.M{"update_date": -1}
	logger.Infof(ctx, "ListUserInfo Filter:", filter)
	cursor, err := collection.Find(
		GetContext(),
		filter,
		options.Find().SetSort(sort).SetSkip(condition.Start).SetLimit(condition.Limit),
	)
	if err != nil {
		logger.Fatalf(ctx, "ListUserInfo Cursor Failed [err=%+v]", err)
		return
	}
	if err := cursor.Err(); err != nil {
		logger.Fatalf(ctx, "ListUserInfo Cursor Failed [err=%+v]", err)
		return
	}
	defer cursor.Close(GetContext())
	err = cursor.All(GetContext(), &recordList)
	if err != nil {
		logger.Fatalf(ctx, "ListLiveRecord Failed [err=%+v]", err)
	}
	return
}

func UpdateUserInfo(ctx *gin.Context, userInfo UserInfo) error {
	logger.Infof(ctx, "UpdateUserInfo [condition:%+v]", userInfo)
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionUserInfo)
	update := bson.M{}
	update["update_date"] = time.Now()

	if len(userInfo.FollowTag) > 0 {
		update["follow_author"] = userInfo.FollowTag
	}
	if len(userInfo.FollowAuthor) > 0 {
		update["follow_author"] = userInfo.FollowAuthor
	}
	if len(userInfo.FocusArea) > 0 {
		update["focus_area"] = userInfo.FocusArea
	}

	options := options.Update().SetUpsert(true) // 设置 upsert 选项
	filter := bson.M{"uid": userInfo.Uid}
	logger.Infof(ctx, "UpdateUserInfo [filter:%+v] [sort:%+v]", filter, update, options)
	upateRes, err := collection.UpdateOne(GetContext(), filter, bson.M{"$set": update})
	if err != nil {
		logger.Fatalf(ctx, "UpdateUserInfo Failed [result:%+v] [err=%+v]", upateRes, err)
	}
	return err
}

func DeleteUserInfo(ctx *gin.Context, updateTime time.Time) error {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionUserInfo)
	filter := bson.M{
		"update_date": bson.M{
			"$lt": updateTime,
		},
	}
	logger.Infof(ctx, "DeleteUserInfo [filter:%+v]", filter)
	delRes, err := collection.DeleteMany(GetContext(), filter)
	if err != nil {
		logger.Fatalf(ctx, "DeleteUserInfo Failed [result:%+v] [err:%+v]", delRes, err)
	}
	return err
}
