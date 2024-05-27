package common_test

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

func TestConnection(t *testing.T) {
	ctx := context.GetGinContextWithRequestId()

	c := common.GetClient()

	dbs, err := c.ListDatabases(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("server databases: %#v", dbs)
}
