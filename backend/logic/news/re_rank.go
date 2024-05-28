package news

import (
	"errors"
	"fmt"
	"learnerai/golib/app/logger"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

const (
	reRankModelUrl       = "http://halfjourney.xyz:39002/api/re_rank"
	reRankHttpHeaderHost = "embedding.halfjourney.xyz"
)

type computeSimScoreRequest struct {
	Data []any `json:"data"`
}

type computeSimScoreResponse struct {
	Data []numberTable `json:"data"` // list of length 1
}

type varTable[T any] struct {
	Headers  []string `json:"headers"`
	Data     [][]T    `json:"data"`
	Metadata any      `json:"metadata"` // not used
}

type stringTable = varTable[string]
type numberTable = varTable[float64]

type argSort struct {
	sort.Interface
	idx []int
}

func (s argSort) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

// NOTE: input is also mutated!
func ArgSort(v sort.Interface) []int {
	a := argSort{v, make([]int, v.Len())}
	for i := range a.idx {
		a.idx[i] = i
	}
	sort.Sort(a)
	return a.idx
}

func SortSimScore(scores []float64) ([]int, []float64) {
	scp := make([]float64, len(scores))
	copy(scp, scores)
	idx := ArgSort(sort.Reverse(sort.Float64Slice(scp)))
	return idx, scp
}

func ComputeSimScore(ctx *gin.Context, query string, corpus []string) ([]float64, error) {
	reqData := computeSimScoreRequest{
		Data: []any{
			query,
			stringTable{
				Headers: []string{"corpus"},
				Data:    list2column(corpus),
			},
		},
	}
	var respData computeSimScoreResponse

	resp, _, errs := gorequest.New().
		Post(reRankModelUrl).
		AppendHeader("Host", reRankHttpHeaderHost).
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

	return column2list(respData.Data[0].Data), nil
}

func list2column[T any](list []T) [][]T {
	res := make([][]T, len(list))
	for i := range list {
		res[i] = []T{list[i]}
	}
	return res
}

func column2list[T any](table [][]T) []T {
	res := make([]T, len(table))
	for i := range table {
		res[i] = table[i][0]
	}
	return res
}

// curl -H 'content-type: application/json' -d '{"data": ["query", {"data": [["word1"], ["word2"]], "headers": ["corpus"]}]}' 127.0.0.1:7860/api/re_rank
// {"data":[{"headers":["score"],"data":[[-9.025601387023926],[-9.614904403686523]],"metadata":null}],"is_generating":false,"duration":0.024728059768676758,"average_duration":0.33085960149765015}
