package news

import (
	"learnerai/golib/app/context"
	"learnerai/golib/app/logger"
	model "learnerai/models/mgodb/news"
	"testing"
)

func init() {
	logger.Setup("runtime/", "test.log", "debug", "debug", 7)
}

func TestComputeEmbedding(t *testing.T) {
	vecs, err := computeEmbedding(
		context.GetGinContextWithRequestId(),
		[]model.NewsWithEmbedding{
			{Title: "aaa", Abstract: "ccc"},
			{Title: "bbb", Abstract: "ddd"},
		})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vecs)
	// expected:
	// [[-0.46893066  0.5737397  -0.49394795 ... -0.02208222 -0.3382788 0.9104303 ]
	//  [-0.68060416  0.40745476 -0.95192975 ... -0.45749992 -0.10577605 0.96295744]]
}
