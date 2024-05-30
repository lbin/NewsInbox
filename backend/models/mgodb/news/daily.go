package news

import (
	"learnerai/golib/app/logger"
	. "learnerai/models/mgodb/common"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DailyInfo struct {
	NewsId    string    `json:"_id" bson:"_id"`
	Title     string    `json:"title" bson:"title"`
	Summary   string    `json:"summary" bson:"summary"`
	Class     []string  `json:"class" bson:"class"`
	KeyPoints string    `json:"key_points" bson:"key_points"`
	Tags      []string  `json:"tags" bson:"tags"`
	Products  []string  `json:"products" bson:"products"`
	Teams     []string  `json:"teams" bson:"teams"`
	URL       string    `json:"url" bson:"url"`
	Sender    string    `json:"sender" bson:"sender"`
	AddPerson string    `json:"add_person" bson:"add_person"`
	Published time.Time `json:"published" bson:"published"`
	// RawContent string    `json:"raw_content" bson:"raw_content"`
}

type ListDailyInfoCondition struct {
	SearchContent string    `json:"search_content"`
	Id            string    `json:"id"`
	Tags          []string  `json:"tags"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Start         int64     `json:"start"`
	Limit         int64     `json:"limit"`
}

type ListDailyUserInfoCondition struct {
	Id        []string  `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Start     int64     `json:"start"`
	Limit     int64     `json:"limit"`
}

// func AddDailyInfo(ctx *gin.Context, DailyInfo DailyInfo) error {
// 	logger.Infof(ctx, "AddDailyInfo [condition:%+v]", CollectionNewsDaily)
// 	mClient := DB.Mongo
// 	var table = CollectionKaggleSnapshot
// 	collection := mClient.Database(GetDataBase()).Collection(table)
// 	info, err := collection.InsertOne(GetContext(), DailyInfo)

// 	if err != nil {
// 		logger.Fatalf(ctx, "AddDailyInfo [err=%+v]", err)
// 	} else {
// 		logger.Infof(ctx, "AddDailyInfo [info=%+v]", info)
// 	}
// 	return err
// }

func ListDailyInfo(ctx *gin.Context, condition ListDailyInfoCondition) (recordList []*DailyInfo) {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionNewsDaily)

	filter := bson.M{}
	filter["published"] = bson.M{"$gte": condition.StartTime, "$lte": condition.EndTime}
	if len(condition.Tags) > 0 {
		filter["categories"] = bson.M{"$in": condition.Tags}
	}
	if condition.SearchContent != "" {
		regex := primitive.Regex{Pattern: condition.SearchContent, Options: "i"} // 不区分大小写
		filter["title"] = bson.M{"$regex": regex}
	}
	if condition.Id != "" {
		filter["_id"] = condition.Id
		delete(filter, "published")
	}
	sort := bson.M{"published": -1}

	logger.Infof(ctx, "ListDailyInfo Filter:", filter)
	cursor, err := collection.Find(
		GetContext(),
		filter,
		options.Find().SetSort(sort).SetSkip(condition.Start).SetLimit(condition.Limit),
	)
	if err != nil {
		logger.Fatalf(ctx, "ListDailyInfo Cursor Failed [err=%+v]", err)
		return
	}
	if err := cursor.Err(); err != nil {
		logger.Fatalf(ctx, "ListDailyInfo Cursor Failed [err=%+v]", err)
		return
	}
	defer cursor.Close(GetContext())
	err = cursor.All(GetContext(), &recordList)
	if err != nil {
		logger.Fatalf(ctx, "ListLiveRecord Failed [err=%+v]", err)
	}
	return
}

func ListDailyUserInfo(ctx *gin.Context, condition ListDailyUserInfoCondition) (recordList []*DailyInfo) {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionNewsDaily)

	filter := bson.M{}
	filter["published"] = bson.M{"$gte": condition.StartTime, "$lte": condition.EndTime}

	if len(condition.Id) > 0 {
		filter["_id"] = bson.M{"$in": condition.Id}
		delete(filter, "published")
	}
	sort := bson.M{"published": -1}

	logger.Infof(ctx, "ListDailyInfo Filter:", filter)
	cursor, err := collection.Find(
		GetContext(),
		filter,
		options.Find().SetSort(sort).SetSkip(condition.Start).SetLimit(condition.Limit),
	)
	if err != nil {
		logger.Fatalf(ctx, "ListDailyInfo Cursor Failed [err=%+v]", err)
		return
	}
	if err := cursor.Err(); err != nil {
		logger.Fatalf(ctx, "ListDailyInfo Cursor Failed [err=%+v]", err)
		return
	}
	defer cursor.Close(GetContext())
	err = cursor.All(GetContext(), &recordList)
	if err != nil {
		logger.Fatalf(ctx, "ListLiveRecord Failed [err=%+v]", err)
	}
	return
}

func DeleteDailyInfo(ctx *gin.Context, updateTime time.Time) error {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionNewsDaily)
	filter := bson.M{
		"update_date": bson.M{
			"$lt": updateTime,
		},
	}
	logger.Infof(ctx, "DeleteDailyInfo [filter:%+v]", filter)
	delRes, err := collection.DeleteMany(GetContext(), filter)
	if err != nil {
		logger.Fatalf(ctx, "DeleteDailyInfo Failed [result:%+v] [err:%+v]", delRes, err)
	}
	return err
}

func UnsetDailyInfoField(ctx *gin.Context, id string, fields []string) error {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionNewsDaily)
	unset := make(bson.M)
	for _, f := range fields {
		unset[f] = ""
	}
	update := bson.M{"$unset": unset}
	logger.Infof(ctx, "UnsetDailyInfoField [update:%+v]", update)
	unsetRes, err := collection.UpdateByID(GetContext(), id, update)
	if err != nil {
		logger.Errorf(ctx, "UnsetDailyInfoField Failed [result:%+v] [err=%+v]", unsetRes, err)
	}
	return err

}

type IncNewsCondition struct {
	DocId   string `json:"doc_id"`
	ReadCnt int    `json:"read_cnt"`
	LikeCnt int    `json:"like_cnt"`
}

func IncNewsCount(ctx *gin.Context, incInfo IncNewsCondition) error {
	logger.Infof(ctx, "UpdatenewsInfo [condition:%+v]", incInfo)
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionNewsDaily)

	//更新文档计数
	inc := bson.M{}
	if incInfo.LikeCnt > 0 {
		inc["like_cnt"] = incInfo.LikeCnt
	}
	if incInfo.ReadCnt > 0 {
		inc["read_cnt"] = incInfo.ReadCnt
	}
	update := bson.M{}
	update["$inc"] = inc

	set := bson.M{}
	set["update_time"] = time.Now()
	update["$set"] = set

	options := options.Update().SetUpsert(true) // 设置 upsert 选项
	filter := bson.M{"_id": incInfo.DocId}
	logger.Infof(ctx, "UpdatenewsInfo [filter:%+v] [sort:%+v]", filter, update, options)
	upateRes, err := collection.UpdateOne(GetContext(), filter, update)
	if err != nil {
		logger.Fatalf(ctx, "UpdatenewsInfo Failed [result:%+v] [err=%+v]", upateRes, err)
	}
	return err
}
