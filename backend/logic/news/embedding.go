package news

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"learnerai/golib/app/logger"
	"learnerai/models/mgodb/common"
	mgodb "learnerai/models/mgodb/news"
	milvus "learnerai/models/milvus/news"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

const (
	embeddingModelName      = "allenai-specter" // should be a valid mongo key name
	embeddingDataType       = "float32"
	embeddingByteSize       = 4
	embeddingModelUrl       = "http://halfjourney.xyz:39002/api/compute_embedding"
	embeddingHttpHeaderHost = "embedding.halfjourney.xyz"
)

type computeEmbeddingRequest struct {
	Data    [][]string `json:"data"` // the request is batched, each field is a list.
	Batched bool       `json:"batched"`
}

type computeEmbeddingResponse struct {
	Data [][]mgodb.EmbeddingVector `json:"data"` // the request is batched so result is list of fields, but in this case length 1.
}

func toF32Vector(emb mgodb.EmbeddingVector) ([]float32, error) {
	if emb.Dtype != embeddingDataType {
		return nil, fmt.Errorf("unexpected dtype %q, should be %q", emb.Dtype, embeddingDataType)
	}
	b, err := base64.StdEncoding.DecodeString(emb.Base64)
	if err != nil {
		return nil, err
	}
	expectedLen := embeddingByteSize * emb.Dim
	if len(b) != int(expectedLen) {
		return nil, fmt.Errorf("unexpected binary length %d, should be %d", len(b), expectedLen)
	}
	vf32 := make([]float32, emb.Dim)
	for j := 0; j < int(emb.Dim); j++ {
		vf32[j] = f32le(b[j*4:])
	}
	return vf32, nil
}

func f32le(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(b))
}

func computeEmbedding(ctx *gin.Context, docs []mgodb.NewsWithEmbedding) ([]mgodb.EmbeddingVector, error) {
	n := len(docs)
	titles := make([]string, n)
	abstracts := make([]string, n)
	for i := range docs {
		titles[i] = docs[i].Title
		abstracts[i] = docs[i].Abstract
	}

	reqData := computeEmbeddingRequest{
		Data:    [][]string{titles, abstracts},
		Batched: true,
	}
	var respData computeEmbeddingResponse

	resp, _, errs := gorequest.New().
		Post(embeddingModelUrl).
		AppendHeader("Host", embeddingHttpHeaderHost).
		Type(gorequest.TypeJSON).
		SendStruct(reqData).
		EndStruct(&respData)

	err := errors.Join(errs...)
	if err != nil {
		logger.Errorf(ctx, "http request failed [err=%+v]", err)
		return nil, err
	}

	logger.Debugf(ctx, "http response detail: %+v", resp)

	if resp.StatusCode != http.StatusOK {
		logger.Errorf(ctx, "http request is not OK: %q", resp.Status)
		return nil, fmt.Errorf("http request failed with code %v", resp.StatusCode)
	}

	embs := respData.Data[0]
	return embs, nil
}

func ComputeEmbeddingFromText(ctx *gin.Context, title, abstract string) ([]float32, error) {
	v, err := computeEmbedding(ctx, []mgodb.NewsWithEmbedding{{Title: title, Abstract: abstract}})
	if err != nil {
		logger.Errorf(ctx, "ComputeEmbedding from text failed [err=%+v]", err)
		return nil, err
	}

	if len(v) == 0 {
		err = fmt.Errorf("empty result list")
		logger.Errorf(ctx, "ComputeEmbedding from text failed [err=%+v]", err)
		return nil, err
	}

	return toF32Vector(v[0])
}

func SearchByEmbedding(ctx *gin.Context, embedding []float32, topK int) (milvus.ScoredNewsRecords, error) {
	results, err := milvus.SearchNewsRecords(
		ctx,
		milvus.NewsRecords{{Embedding: embedding}}, // single search
		topK,
		"",   // no filter
		0, 0, // no offset or limit
	)
	if err != nil {
		logger.Errorf(ctx, "SearchNewsRecords by embedding failed [err=%+v]", err)
		return nil, err
	}
	if len(results) == 0 {
		err = fmt.Errorf("empty result list")
		logger.Errorf(ctx, "SearchNewsRecords by embedding failed [err=%+v]", err)
		return nil, err
	}
	return results[0], err
}

type ScoredResult struct {
	DocId          string  `json:"doc_id"`
	EmbeddingScore float64 `json:"embedding_score"`
	ReRankScore    float64 `json:"re_rank_score"`
}

func SearchSimilar(ctx *gin.Context, arxivId string, topK int) ([]ScoredResult, error) {
	collName := common.CollectionNewsDaily
	if strings.Contains(arxivId, "v") {
		collName = common.CollectionNewsDaily
	}
	pp := mgodb.FindNewsWithEmbeddingById(ctx, collName, arxivId)
	if pp == nil {
		err := fmt.Errorf("document %q not found in %q", arxivId, collName)
		return nil, err
	}
	emb, ok := pp.Embedding[embeddingModelName]

	var vec []float32
	if ok {
		v, err := toF32Vector(emb)
		if err != nil {
			logger.Errorf(ctx, "decode embedding from database failed [err=%+v]", err)
			return nil, err
		}
		vec = v
	} else {
		v, err := ComputeEmbeddingFromText(ctx, pp.Title, pp.Abstract)
		if err != nil {
			logger.Errorf(ctx, "compute embedding from text failed [err=%+v]", err)
			return nil, err
		}
		vec = v
	}

	res, err := SearchByEmbedding(ctx, vec, topK*10)
	if err != nil {
		logger.Errorf(ctx, "search by embedding failed [err=%+v]", err)
		return nil, err
	}

	// re-rank
	query := joinTextForReRank(*pp)
	corpus := make([]string, len(res))
	for i := range res {
		pr := mgodb.FindNewsWithEmbeddingById(ctx, common.CollectionNewsDaily, res[i].ArxivId)
		// TODO when fill milvus, also carry title and abstract.
		if pr != nil {
			corpus[i] = joinTextForReRank(*pr)
		}
	}
	ss, err := ComputeSimScore(ctx, query, corpus)
	if err != nil {
		logger.Errorf(ctx, "compute sim score for re-rank failed [err=%+v]", err)
		return nil, err
	}
	idx, ss := SortSimScore(ss)

	sres := make([]ScoredResult, 0, topK)
	for i := 0; i < topK && i < len(idx); i++ {
		sres = append(sres, ScoredResult{
			DocId:          res[idx[i]].ArxivId,
			EmbeddingScore: float64(res[idx[i]].Score),
			ReRankScore:    ss[i],
		})
	}

	return sres, nil
}

func joinTextForReRank(pp mgodb.NewsWithEmbedding) string {
	return "Title: " + pp.Title + "\n" + "Abstract: " + pp.Abstract
}
