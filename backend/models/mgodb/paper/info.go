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

type PaperInfo struct {
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

/*
*   {
    _id: '1903.05041',
    submitter: 'Yuval Pinter',
    authors: 'Yuval Pinter, Marc Marone, Jacob Eisenstein',
    title: 'Character Eyes: Seeing Language through Character-Level Taggers',
    comments: null,
    'journal-ref': null,
    doi: null,
    'report-no': null,
    categories: [ 'cs.CL' ],
    license: 'http://arxiv.org/licenses/nonexclusive-distrib/1.0/',
    abstract: '  Character-level models have been used extensively in recent years in NLP\n' +
      'tasks as both supplements and replacements for closed-vocabulary token-level\n' +
      'word representations. In one popular architecture, character-level LSTMs are\n' +
      'used to feed token representations into a sequence tagger predicting\n' +
      'token-level annotations such as part-of-speech (POS) tags. In this work, we\n' +
      'examine the behavior of POS taggers across languages from the perspective of\n' +
      'individual hidden units within the character LSTM. We aggregate the behavior of\n' +
      'these units into language-level metrics which quantify the challenges that\n' +
      'taggers face on languages with different morphological properties, and identify\n' +
      'links between synthesis and affixation preference and emergent behavior of the\n' +
      'hidden tagger layer. In a comparative experiment, we show how modifying the\n' +
      'balance between forward and backward hidden units affects model arrangement and\n' +
      'performance in these types of languages.\n',
    versions: [ { version: 'v1', created: ISODate('2019-03-12T16:42:39.000Z') } ],
    update_date: ISODate('2019-03-13T00:00:00.000Z'),
    authors_parsed: [
      [ 'Pinter', 'Yuval', '' ],
      [ 'Marone', 'Marc', '' ],
      [ 'Eisenstein', 'Jacob', '' ]
    ]
  }
*/

type ListPaperInfoCondition struct {
	Id            string    `json:"id"`
	SearchContent string    `json:"search_content"`
	Category      []string  `json:"category"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Start         int64     `json:"start"`
	Limit         int64     `json:"limit"`
}

func (c ListPaperInfoCondition) asFilter() bson.M {
	filter := bson.M{}

	if c.Id != "" {
		filter["_id"] = c.Id
		return filter // shortcut
	}

	filter["update_date"] = bson.M{"$gte": c.StartTime, "$lte": c.EndTime}
	if len(c.Category) == 1 {
		filter["categories"] = c.Category[0]
	} else if len(c.Category) > 1 {
		filter["categories"] = bson.M{"$in": c.Category}
	}
	if c.SearchContent != "" {
		regex := primitive.Regex{Pattern: c.SearchContent, Options: "i"} // 不区分大小写
		filter["title"] = bson.M{"$regex": regex}
	}
	return filter
}

func AddPaperInfo(ctx *gin.Context, paperinfo PaperInfo) error {
	logger.Infof(ctx, "AddPaperInfo [condition:%+v]", paperinfo)
	mClient := DB.Mongo
	var table = CollectionKaggleSnapshot
	collection := mClient.Database(GetDataBase()).Collection(table)
	info, err := collection.InsertOne(GetContext(), paperinfo)

	if err != nil {
		logger.Fatalf(ctx, "AddPaperInfo [err=%+v]", err)
	} else {
		logger.Infof(ctx, "AddPaperInfo [info=%+v]", info)
	}
	return err
}

func ListPaperInfo(ctx *gin.Context, condition ListPaperInfoCondition) (recordList []*PaperInfo) {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionKaggleSnapshot)
	filter := condition.asFilter()
	sort := bson.M{"update_date": -1}

	logger.Infof(ctx, "ListPaperInfo Filter:", filter)
	cursor, err := collection.Find(
		GetContext(),
		filter,
		options.Find().SetSort(sort).SetSkip(condition.Start).SetLimit(condition.Limit),
	)
	if err != nil {
		logger.Fatalf(ctx, "ListPaperInfo Cursor Failed [err=%+v]", err)
		return
	}
	if err := cursor.Err(); err != nil {
		logger.Fatalf(ctx, "ListPaperInfo Cursor Failed [err=%+v]", err)
		return
	}
	defer cursor.Close(GetContext())
	err = cursor.All(GetContext(), &recordList)
	if err != nil {
		logger.Fatalf(ctx, "ListLiveRecord Failed [err=%+v]", err)
	}
	return
}

func DeletePaperInfo(ctx *gin.Context, updateTime time.Time) error {
	mClient := DB.Mongo
	collection := mClient.Database(GetDataBase()).Collection(CollectionKaggleSnapshot)
	filter := bson.M{
		"update_date": bson.M{
			"$lt": updateTime,
		},
	}
	logger.Infof(ctx, "DeletePaperInfo [filter:%+v]", filter)
	delRes, err := collection.DeleteMany(GetContext(), filter)
	if err != nil {
		logger.Fatalf(ctx, "DeletePaperInfo Failed [result:%+v] [err:%+v]", delRes, err)
	}
	return err
}
