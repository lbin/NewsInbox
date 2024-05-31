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

type NewsInfo struct {
	DocId         string     `json:"_id" bson:"_id"`
	Submitter     string     `json:"submitter" bson:"submitter"`
	Authors       string     `json:"authors" bson:"authors"`
	Title         string     `json:"title" bson:"title"`
	Comments      string     `json:"comments" bson:"comments"`
	Doi           string     `json:"doi" bson:"doi"`
	ReportNo      string     `json:"report-no" bson:"report-no"`
	Categories    []string   `json:"categories" bson:"categories"`
	Versions      []Version  `json:"versions" bson:"versions"`
	Abstract      string     `json:"abstract" bson:"abstract"`
	UpdateDate    time.Time  `json:"update_date" bson:"update_date"`
	AuthorsParsed [][]string `json:"authors_parsed" bson:"authors_parsed"`
}

type Version struct {
	Pversion string    `json:"version" bson:"version"`
	Created  time.Time `json:"created" bson:"created"`
}

type ListNewsInfoCondition struct {
	Id            string    `json:"id"`
	SearchContent string    `json:"search_content"`
	Tags          []string  `json:"tags"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Start         int64     `json:"start"`
	Limit         int64     `json:"limit"`
}

func (c ListNewsInfoCondition) asFilter() bson.M {
	filter := bson.M{}

	if c.Id != "" {
		filter["_id"] = c.Id
		return filter // shortcut
	}

	filter["update_date"] = bson.M{"$gte": c.StartTime, "$lte": c.EndTime}
	if len(c.Tags) == 1 {
		filter["categories"] = c.Tags[0]
	} else if len(c.Tags) > 1 {
		filter["categories"] = bson.M{"$in": c.Tags}
	}
	if c.SearchContent != "" {
		regex := primitive.Regex{Pattern: c.SearchContent, Options: "i"} // 不区分大小写
		filter["title"] = bson.M{"$regex": regex}
	}
	return filter
}

func AddNewsInfo(ctx *gin.Context, newsinfo NewsInfo) error {
	logger.Infof(ctx, "AddNewsInfo [condition:%+v]", newsinfo)
	mClient := DB.Mongo
	var table = CollectionNewsDaily
	collection := mClient.Database(GetDataBase()).Collection(table)
	info, err := collection.InsertOne(GetContext(), newsinfo)

	if err != nil {
		logger.Fatalf(ctx, "AddNewsInfo [err=%+v]", err)
	} else {
		logger.Infof(ctx, "AddNewsInfo [info=%+v]", info)
	}
	return err
}

func ListNewsInfo(ctx *gin.Context, condition ListNewsInfoCondition) (recordList []*NewsInfo) {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionNewsDaily)
	filter := condition.asFilter()
	sort := bson.M{"update_date": -1}

	logger.Infof(ctx, "ListNewsInfo Filter:", filter)
	cursor, err := collection.Find(
		GetContext(),
		filter,
		options.Find().SetSort(sort).SetSkip(condition.Start).SetLimit(condition.Limit),
	)
	if err != nil {
		logger.Fatalf(ctx, "ListNewsInfo Cursor Failed [err=%+v]", err)
		return
	}
	if err := cursor.Err(); err != nil {
		logger.Fatalf(ctx, "ListNewsInfo Cursor Failed [err=%+v]", err)
		return
	}
	defer cursor.Close(GetContext())
	err = cursor.All(GetContext(), &recordList)
	if err != nil {
		logger.Fatalf(ctx, "ListLiveRecord Failed [err=%+v]", err)
	}
	return
}

func DeleteNewsInfo(ctx *gin.Context, updateTime time.Time) error {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionNewsDaily)
	filter := bson.M{
		"update_date": bson.M{
			"$lt": updateTime,
		},
	}
	logger.Infof(ctx, "DeleteNewsInfo [filter:%+v]", filter)
	delRes, err := collection.DeleteMany(GetContext(), filter)
	if err != nil {
		logger.Fatalf(ctx, "DeleteNewsInfo Failed [result:%+v] [err:%+v]", delRes, err)
	}
	return err
}
