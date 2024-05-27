package user

import (
	"learnerai/golib/app/logger"
	. "learnerai/models/mgodb/common"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserPaper struct {
	Uid        string    `json:"_id" bson:"_id"`
	UpdateDate time.Time `json:"update_date" bson:"update_date"`
	Like       []string  `json:"like" bson:"like"`     //用户喜欢的paper id
	UnLike     []string  `json:"unlike" bson:"unlike"` //用户不喜欢的paper id
	Read       []string  `json:"read" bson:"read"`     //用户已经阅读的paper id
}

type ListUserPaperCondition struct {
	Uid   string `json:"uid"`
	Start int64  `json:"start"`
	Limit int64  `json:"limit"`
}

func AddUserPaper(ctx *gin.Context, UserPaper UserPaper) error {
	logger.Infof(ctx, "AddUserPaper [condition:%+v]", UserPaper)
	mClient := DB.Mongo
	var table = CollectionPaperInfo
	collection := mClient.Database(GetDataBase()).Collection(table)
	info, err := collection.InsertOne(GetContext(), UserPaper)

	if err != nil {
		logger.Fatalf(ctx, "AddUserPaper [err=%+v]", err)
	} else {
		logger.Infof(ctx, "AddUserPaper [info=%+v]", info)
	}
	return err
}

func ListUserPaper(ctx *gin.Context, condition ListUserPaperCondition) (recordList []*UserPaper) {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionPaperInfo)
	filter := bson.M{}
	if condition.Uid != "" {
		filter["_id"] = condition.Uid
	}
	sort := bson.M{"update_date": -1}

	logger.Infof(ctx, "ListUserPaper Filter:", filter)
	cursor, err := collection.Find(
		GetContext(),
		filter,
		options.Find().SetSort(sort).SetSkip(condition.Start).SetLimit(condition.Limit),
	)
	if err != nil {
		logger.Fatalf(ctx, "ListUserPaper Cursor Failed [err=%+v]", err)
		return
	}
	if err := cursor.Err(); err != nil {
		logger.Fatalf(ctx, "ListUserPaper Cursor Failed [err=%+v]", err)
		return
	}
	defer cursor.Close(GetContext())
	err = cursor.All(GetContext(), &recordList)
	if err != nil {
		logger.Fatalf(ctx, "ListLiveRecord Failed [err=%+v]", err)
	}
	return
}

func UpdateUserPaper(ctx *gin.Context, UserPaper UserPaper) error {
	logger.Infof(ctx, "UpdateUserPaper [condition:%+v]", UserPaper)
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionPaperInfo)
	update := bson.M{}
	update["update_date"] = time.Now()

	/*
		Like       []string  `json:"like" bson:"like"`     //用户喜欢的paper id
		UnLike     []string  `json:"unlike" bson:"unlike"` //用户不喜欢的paper id
		Read       []string  `json:"read" bson:"read"`     //用户已经阅读的paper id
	*/
	if len(UserPaper.Like) > 0 {
		update["like"] = UserPaper.Like
	}
	if len(UserPaper.UnLike) > 0 {
		update["unlike"] = UserPaper.UnLike
	}
	if len(UserPaper.Read) > 0 {
		update["read"] = UserPaper.Read
	}
	options := options.Update().SetUpsert(true) // 设置 upsert 选项
	filter := bson.M{"_id": UserPaper.Uid}
	logger.Infof(ctx, "UpdateUserPaper [filter:%+v] [sort:%+v]", filter, update, options)
	upateRes, err := collection.UpdateOne(GetContext(), filter, bson.M{"$set": update})
	if err != nil {
		logger.Fatalf(ctx, "UpdateUserPaper Failed [result:%+v] [err=%+v]", upateRes, err)
	}
	return err
}

func DeleteUserPaper(ctx *gin.Context, updateTime time.Time) error {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionPaperInfo)
	filter := bson.M{
		"update_date": bson.M{
			"$lt": updateTime,
		},
	}
	logger.Infof(ctx, "DeleteUserPaper [filter:%+v]", filter)
	delRes, err := collection.DeleteMany(GetContext(), filter)
	if err != nil {
		logger.Fatalf(ctx, "DeleteUserPaper Failed [result:%+v] [err:%+v]", delRes, err)
	}
	return err
}
