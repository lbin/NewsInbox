package news

import (
	"fmt"
	"learnerai/golib/app/logger"
	"learnerai/models/milvus/common"

	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	collectionName = "allenai_specter" // letter, number and underscore
	partitionName  = ""                // use default
	fieldArxivId   = "arxiv_id"
	fieldEmbedding = "embedding"
	fieldTitle     = "title"
	fieldAbstract  = "abstract"

	// indexed by IVF flat, with search parameters
	ivfFlatMetricType     = entity.COSINE
	ivfFlatIndexNumList   = 1024
	ivfFlatSearchNumProbe = 10
)

type NewsRecord struct {
	ArxivId   string    `json:"arxiv_id"`
	Embedding []float32 `json:"embedding"`
	Title     string    `json:"title"`
	Abstract  string    `json:"abstract"`
}

func (r NewsRecord) floatVector() entity.FloatVector {
	return r.Embedding
}

func (r NewsRecord) vector() entity.Vector {
	return r.floatVector()
}

type NewsRecords []NewsRecord

func (s NewsRecords) vectors() []entity.Vector {
	r := make([]entity.Vector, len(s))
	for i := range s {
		r[i] = s[i].vector()
	}
	return r
}

func getColumnAs[T entity.Column](r client.ResultSet, name string) (t T) {
	c := r.GetColumn(name)
	if c == nil {
		return
	}
	v, ok := c.(T)
	if !ok {
		return
	}
	return v
}

func fromResultSet(r client.ResultSet) (NewsRecords, error) {
	idCol := getColumnAs[*entity.ColumnVarChar](r, fieldArxivId)
	if idCol == nil {
		return nil, fmt.Errorf("failed to parse ids")
	}
	s := make(NewsRecords, idCol.Len())
	for i := range s {
		s[i].ArxivId = idCol.Data()[i]
	}

	if col := getColumnAs[*entity.ColumnFloatVector](r, fieldEmbedding); col != nil {
		for i := range s {
			s[i].Embedding = col.Data()[i]
		}
	}
	if col := getColumnAs[*entity.ColumnString](r, fieldTitle); col != nil {
		for i := range s {
			s[i].Title = col.Data()[i]
		}
	}
	if col := getColumnAs[*entity.ColumnString](r, fieldAbstract); col != nil {
		for i := range s {
			s[i].Abstract = col.Data()[i]
		}
	}

	return s, nil
}

func GetNewsRecords(ctx *gin.Context, arxivIds []string, start, limit int64) (NewsRecords, error) {
	c := common.GetClient()
	wctx, cancel := common.GetContext(ctx)
	defer cancel()

	err := c.LoadCollection(wctx, collectionName, false)
	if err != nil {
		logger.Errorf(ctx, "failed to load collection into memory [err=%+v]", err)
		return nil, err
	}

	idCol := entity.NewColumnVarChar(fieldArxivId, arxivIds)

	r, err := c.QueryByPks(
		wctx,
		collectionName, nil,
		idCol,
		[]string{fieldArxivId, fieldEmbedding}, // projection
		defaultSearchQueryOption(start, limit),
	)
	if err != nil {
		logger.Errorf(ctx, "failed to conduct query by primary key [err=%+v]", err)
		return nil, err
	}

	return fromResultSet(r)
}

func DeleteNewsRecords(ctx *gin.Context, arxivIds []string) error {
	c := common.GetClient()
	wctx, cancel := common.GetContext(ctx)
	defer cancel()

	ids := entity.NewColumnVarChar(fieldArxivId, arxivIds)

	err := c.DeleteByPks(wctx, collectionName, partitionName, ids)
	if err != nil {
		logger.Errorf(ctx, "delete news record failed [err=%+v]", err)
		return err
	}

	return nil
}

type ScoredNewsRecord struct {
	NewsRecord
	Score float32 `json:"score"`
}

type ScoredNewsRecords []ScoredNewsRecord

func fromSearchResult(r client.SearchResult) (ScoredNewsRecords, error) {
	prs, err := fromResultSet(r.Fields)
	if err != nil {
		return nil, err
	}

	sprs := make([]ScoredNewsRecord, r.ResultCount)
	for i := 0; i < r.ResultCount; i++ {
		sprs[i].NewsRecord = prs[i]
		sprs[i].Score = r.Scores[i]
	}
	return sprs, nil
}

func defaultSearchQueryOption(start, limit int64) client.SearchQueryOptionFunc {
	if limit == 0 {
		limit = 16383 - start
	}
	return func(option *client.SearchQueryOption) {
		option.Offset = start
		option.Limit = limit
		option.ConsistencyLevel = entity.ClEventually // relaxed
	}
}

// A vector similarity search in Milvus calculates the distance
// between query vector(s) and vectors in the collection with specified similarity metrics,
// and returns the most similar results.
// Handles batch.
// Empty string means no filter applied.
// start + limit < 16384. 0 for max limit.
func SearchNewsRecords(ctx *gin.Context, inps NewsRecords, topK int, filter string, start, limit int64) ([]ScoredNewsRecords, error) {
	c := common.GetClient()
	wctx, cancel := common.GetContext(ctx)
	defer cancel()

	err := c.LoadCollection(wctx, collectionName, false)
	if err != nil {
		logger.Errorf(ctx, "failed to load collection into memory [err=%+v]", err)
		return nil, err
	}

	sp, err := entity.NewIndexIvfFlatSearchParam(ivfFlatSearchNumProbe)
	if err != nil {
		logger.Errorf(ctx, "failed to create index search param [err=%+v]", err)
		return nil, err
	}

	r, err := c.Search(
		wctx,
		collectionName, nil,
		filter,
		[]string{fieldArxivId, fieldEmbedding}, // projection // TODO add title, abstract
		inps.vectors(), fieldEmbedding, ivfFlatMetricType,
		topK, sp,
		defaultSearchQueryOption(start, limit),
	)
	if err != nil {
		logger.Errorf(ctx, "failed to conduct search [err=%+v]", err)
		return nil, err
	}

	var rlist []ScoredNewsRecords
	for _, re := range r {
		s, err := fromSearchResult(re)
		if err != nil {
			logger.Errorf(ctx, "failed to parse search result [err=%+v]", err)
			return nil, err
		}
		rlist = append(rlist, s)
	}

	return rlist, nil
}

// A query retrieves vectors via scalar filtering based on boolean expression.
func QueryNewsRecords(ctx *gin.Context, expr string, start, limit int64) (NewsRecords, error) {
	c := common.GetClient()
	wctx, cancel := common.GetContext(ctx)
	defer cancel()

	err := c.LoadCollection(wctx, collectionName, false)
	if err != nil {
		logger.Errorf(ctx, "failed to load collection into memory [err=%+v]", err)
		return nil, err
	}

	r, err := c.Query(
		wctx,
		collectionName, nil,
		expr,
		[]string{fieldArxivId, fieldEmbedding}, // projection // TODO add title, abstract
		defaultSearchQueryOption(start, limit),
	)
	if err != nil {
		logger.Errorf(ctx, "failed to conduct query [err=%+v]", err)
		return nil, err
	}

	return fromResultSet(r)
}
