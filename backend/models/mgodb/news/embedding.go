package news

import (
	"learnerai/golib/app/logger"
	"learnerai/models/mgodb/common"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type NewsWithEmbedding struct {
	DocId     string                     `json:"_id" bson:"_id"` // DailyInfo has v1 part, NewsInfo doesn't.
	Title     string                     `json:"title" bson:"title"`
	Abstract  string                     `json:"abstract" bson:"abstract"`
	Embedding map[string]EmbeddingVector `json:"embedding" bson:"embedding"` // indexed by model name.
}

type EmbeddingVector struct {
	Base64     string    `json:"base64" bson:"base64"`
	Dim        int32     `json:"dim" bson:"dim"`
	Dtype      string    `json:"dtype" bson:"dtype"` // numpy dtype, usually float32.
	CreateTime time.Time `json:"create_time" bson:"create_time"`
}

func FindNewsWithEmbeddingById(ctx *gin.Context, collectionName string, arxivId string) *NewsWithEmbedding {
	collection := common.GetCollection(collectionName)
	wctx := common.GetContext()

	res := collection.FindOne(wctx, bson.M{"_id": arxivId})
	if err := res.Err(); err != nil {
		logger.Errorf(ctx, "FindNewsWithEmbedding failed [err=%+v]", err)
		return nil
	}
	var p NewsWithEmbedding
	err := res.Decode(&p)
	if err != nil {
		logger.Errorf(ctx, "FindNewsWithEmbedding decode failed [err=%+v]", err)
		return nil
	}
	return &p
}
