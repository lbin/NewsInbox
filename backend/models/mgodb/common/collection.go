package common

import (
	"context"

	"learnerai/golib/app/logger"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DATABASE = "newsinbox"

// 核心配置信息
const CollectionNewsDaily = "t_news_daily"
const CollectionUserInfo = "t_user_info"
const CollectionNewsInfo = "t_user_news"

//todo 注意事项
//1.新加表是否需要同步，如果需要同步必须加入 MapSyncType2Table 中
//2.为 "community_id",	"update_time", 加上索引
//3.新加表业务中需要增加 import export 函数，支持数据导入导出
//4.所有的update操作都需要更新update时间戳

const DefaultMiniLimit = 5
const DefaultListLimit = 10
const DefaultNormalLimit = 100
const DefaultBiggestLimit = 1000

type MongoDbIndex struct {
	CollectionName string           `json:"collection_name"` //表名
	Keys           map[string]int32 `json:"keys"`            //索引项
	IsUnique       bool             `json:"is_unique"`       //是否唯一索引
	IndexName      string           `json:"index_name"`      //是否唯一索引
}

var table2Indexes = []MongoDbIndex{
	{
		CollectionName: CollectionNewsDaily,
		Keys: map[string]int32{
			"title": 1,
		},
		IsUnique:  true,
		IndexName: "t_news_daily_index_1",
	},
	{
		CollectionName: CollectionNewsDaily,
		Keys: map[string]int32{
			"updated":    1,
			"categories": 1,
		},
		IsUnique:  false,
		IndexName: "t_news_daily_index_2",
	},
	{
		CollectionName: CollectionUserInfo,
		Keys: map[string]int32{
			"update_date": -1,
		},
		IsUnique:  false,
		IndexName: "t_user_info_1",
	},
	{
		CollectionName: CollectionNewsInfo,
		Keys: map[string]int32{
			"update_date": -1,
		},
		IsUnique:  false,
		IndexName: "t_user_paper_1",
	},
}

func GetDataBase() string {
	return DATABASE
}

func GetSyncDataBase() string {
	return DATABASE
}

func GetEdgeDataBase() string {
	return DATABASE
}

func GetCollection(name string) *mongo.Collection {
	return DB.Mongo.Database(GetDataBase()).Collection(name)
}

func CreateIndex(ctx *gin.Context) {
	//创建业务主索引
	executeCreateIndex(ctx, table2Indexes)
}

func executeCreateIndex(ctx *gin.Context, indexes []MongoDbIndex) {
	mClient := DB.Mongo
	for _, index := range indexes {
		collection := mClient.Database(GetDataBase()).Collection(index.CollectionName)
		//配置索引参数
		keys := bson.D{}
		for key, value := range index.Keys {
			singleKey := bson.E{key, value}
			keys = append(keys, singleKey)
		}

		indexName, err := collection.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    keys,
				Options: options.Index().SetUnique(index.IsUnique).SetName(index.IndexName),
			},
		)
		if err != nil {
			logger.Fatalf(ctx, "create indexes error [%s] [%s] [%+v] [%s] ", index.CollectionName, index.IndexName, keys, err)
			continue
		}
		logger.Infof(ctx, "create indexes name [%s] [%s] [%+v]", index.CollectionName, indexName, keys)
	}
}
