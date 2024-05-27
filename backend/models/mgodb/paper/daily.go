package paper

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
	DocId           string     `json:"_id" bson:"_id"`
	EntryId         string     `json:"entry_id" bson:"entry_id"`
	Submitter       string     `json:"submitter" bson:"submitter"`
	Authors         string     `json:"authors" bson:"authors"`
	Title           string     `json:"title" bson:"title"`
	TitleCN         string     `json:"title_ch" bson:"title_ch"`
	Comments        string     `json:"comment" bson:"comment"`
	JournalRef      string     `json:"journal_ref" bson:"journal_ref"`
	Doi             string     `json:"doi" bson:"doi"`
	ReportNo        string     `json:"report-no" bson:"report-no"`
	PrimaryCategory string     `json:"primary_category" bson:"primary_category"`
	Categories      []string   `json:"categories" bson:"categories"`
	PdfUrl          string     `json:"pdf_url" bson:"pdf_url"`
	Keywords        string     `json:"keywords" bson:"keywords"`
	KeywordsCN      string     `json:"keywords_ch" bson:"keywords_ch"`
	Abstract        string     `json:"abstract" bson:"abstract"`
	AbstractCN      string     `json:"abstract_ch" bson:"abstract_ch"`
	UpdateDate      time.Time  `json:"updated" bson:"updated"`
	Published       time.Time  `json:"published" bson:"published"`
	AuthorsParsed   [][]string `json:"authors_parsed" bson:"authors_parsed"`
	LikeCnt         int        `json:"like_cnt" bson:"like_cnt"`
	ReadCnt         int        `json:"read_cnt" bson:"read_cnt"`
	FReadCnt        int        `json:"fread_cnt" bson:"fread_cnt"`
}

type ListDailyInfoCondition struct {
	SearchContent string    `json:"search_content"`
	Id            string    `json:"id"`
	Category      []string  `json:"category"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Start         int64     `json:"start"`
	Limit         int64     `json:"limit"`
}

func AddDailyInfo(ctx *gin.Context, DailyInfo DailyInfo) error {
	logger.Infof(ctx, "AddDailyInfo [condition:%+v]", CollectionNewsDaily)
	mClient := DB.Mongo
	var table = CollectionKaggleSnapshot
	collection := mClient.Database(GetDataBase()).Collection(table)
	info, err := collection.InsertOne(GetContext(), DailyInfo)

	if err != nil {
		logger.Fatalf(ctx, "AddDailyInfo [err=%+v]", err)
	} else {
		logger.Infof(ctx, "AddDailyInfo [info=%+v]", info)
	}
	return err
}

func ListDailyInfo(ctx *gin.Context, condition ListDailyInfoCondition) (recordList []*DailyInfo) {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionNewsDaily)
	

	filter := bson.M{}
	filter["updated"] = bson.M{"$gte": condition.StartTime, "$lte": condition.EndTime}
	if len(condition.Category) > 0 {
		filter["categories"] = bson.M{"$in": condition.Category}
	}
	if condition.SearchContent != "" {
		regex := primitive.Regex{Pattern: condition.SearchContent, Options: "i"} // 不区分大小写
		filter["title"] = bson.M{"$regex": regex}
	}
	if condition.Id != "" {
		filter["_id"] = condition.Id
		delete(filter, "updated")
	}
	sort := bson.M{"updated": -1}

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

type IncPaperCondition struct {
	DocId   string `json:"doc_id"`
	ReadCnt int    `json:"read_cnt"`
	LikeCnt int    `json:"like_cnt"`
}

func IncPaperCount(ctx *gin.Context, incInfo IncPaperCondition) error {
	logger.Infof(ctx, "UpdatepaperInfo [condition:%+v]", incInfo)
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
	logger.Infof(ctx, "UpdatepaperInfo [filter:%+v] [sort:%+v]", filter, update, options)
	upateRes, err := collection.UpdateOne(GetContext(), filter, update)
	if err != nil {
		logger.Fatalf(ctx, "UpdatepaperInfo Failed [result:%+v] [err=%+v]", upateRes, err)
	}
	return err
}
