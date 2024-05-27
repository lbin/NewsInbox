package paper

import (
	"learnerai/conf"
	"learnerai/golib/app/context"
	"learnerai/golib/app/logger"
	"learnerai/models/milvus/common"
	"testing"
)

func init() {
	logger.Setup("runtime/", "test.log", "debug", "debug", 7)

	conf.GlobalEdgeSetting.Milvus = conf.MilvusConfig{
		Host:     "10.100.3.209:19530",
		User:     "root",
		Password: "Milvus",
	}

	common.Setup()
}

func TestInspectCollection(t *testing.T) {
	ctx := context.GetGinContextWithRequestId()

	c := common.GetClient()
	cr, err := c.GetCollectionStatistics(ctx, collectionName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cr)

	ir, err := c.GetIndexState(ctx, collectionName, fieldEmbedding)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ir)
}
