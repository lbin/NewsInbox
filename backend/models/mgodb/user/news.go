package user

import (
	"learnerai/golib/app/logger"
	. "learnerai/models/mgodb/common"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserNews struct {
	Uid        string    `json:"_id" bson:"_id"`
	UpdateDate time.Time `json:"update_date" bson:"update_date"`
	Like       []string  `json:"like" bson:"like"`     //用户喜欢的news id
	UnLike     []string  `json:"unlike" bson:"unlike"` //用户不喜欢的news id
	Read       []string  `json:"read" bson:"read"`     //用户已经阅读的news id
}

type ListUserNewsCondition struct {
	Uid   string `json:"uid"`
	Start int64  `json:"start"`
	Limit int64  `json:"limit"`
}

func AddUserNews(ctx *gin.Context, UserNews UserNews) error {
	logger.Infof(ctx, "AddUserNews [condition:%+v]", UserNews)
	mClient := DB.Mongo
	var table = CollectionNewsInfo
	collection := mClient.Database(GetDataBase()).Collection(table)
	info, err := collection.InsertOne(GetContext(), UserNews)

	if err != nil {
		logger.Fatalf(ctx, "AddUserNews [err=%+v]", err)
	} else {
		logger.Infof(ctx, "AddUserNews [info=%+v]", info)
	}
	return err
}

func ListUserNews(ctx *gin.Context, condition ListUserNewsCondition) (recordList []*UserNews) {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionNewsInfo)
	filter := bson.M{}
	if condition.Uid != "" {
		filter["_id"] = condition.Uid
	}
	sort := bson.M{"update_date": -1}

	logger.Infof(ctx, "ListUserNews Filter:", filter)
	cursor, err := collection.Find(
		GetContext(),
		filter,
		options.Find().SetSort(sort).SetSkip(condition.Start).SetLimit(condition.Limit),
	)
	if err != nil {
		logger.Fatalf(ctx, "ListUserNews Cursor Failed [err=%+v]", err)
		return
	}
	if err := cursor.Err(); err != nil {
		logger.Fatalf(ctx, "ListUserNews Cursor Failed [err=%+v]", err)
		return
	}
	defer cursor.Close(GetContext())
	err = cursor.All(GetContext(), &recordList)
	if err != nil {
		logger.Fatalf(ctx, "ListLiveRecord Failed [err=%+v]", err)
	}
	return
}

func UpdateUserNews(ctx *gin.Context, UserNews UserNews) error {
	logger.Infof(ctx, "UpdateUserNews [condition:%+v]", UserNews)
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionNewsInfo)
	update := bson.M{}
	update["update_date"] = time.Now()

	/*
		Like       []string  `json:"like" bson:"like"`     //用户喜欢的news id
		UnLike     []string  `json:"unlike" bson:"unlike"` //用户不喜欢的news id
		Read       []string  `json:"read" bson:"read"`     //用户已经阅读的news id
	*/
	if len(UserNews.Like) > 0 {
		update["like"] = UserNews.Like
	}
	if len(UserNews.UnLike) > 0 {
		update["unlike"] = UserNews.UnLike
	}
	if len(UserNews.Read) > 0 {
		update["read"] = UserNews.Read
	}
	options := options.Update().SetUpsert(true) // 设置 upsert 选项
	filter := bson.M{"_id": UserNews.Uid}
	logger.Infof(ctx, "UpdateUserNews [filter:%+v] [sort:%+v]", filter, update, options)
	upateRes, err := collection.UpdateOne(GetContext(), filter, bson.M{"$set": update})
	if err != nil {
		logger.Fatalf(ctx, "UpdateUserNews Failed [result:%+v] [err=%+v]", upateRes, err)
	}
	return err
}

func DeleteUserNews(ctx *gin.Context, updateTime time.Time) error {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionNewsInfo)
	filter := bson.M{
		"update_date": bson.M{
			"$lt": updateTime,
		},
	}
	logger.Infof(ctx, "DeleteUserNews [filter:%+v]", filter)
	delRes, err := collection.DeleteMany(GetContext(), filter)
	if err != nil {
		logger.Fatalf(ctx, "DeleteUserNews Failed [result:%+v] [err:%+v]", delRes, err)
	}
	return err
}
