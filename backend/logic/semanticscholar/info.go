package semanticscholar

import (
	"context"
	"learnerai/oapi/semanticscholar/swagger"
	"learnerai/req"
)

// TODO: the result should be cached
func GetPaperInfo(ctx context.Context, r req.PaperInfoFromS2Req) (_ swagger.FullPaper, err error) {
	paperId := toS2PaperId(r.ArxivId)
	opts := swagger.PaperDataApiGetGraphGetPaperOpts{
		Fields: formatFields(r.Fields),
	}
	paper, _, err := s2client.PaperDataApi.GetGraphGetPaper(ctx, paperId, &opts)
	if err != nil {
		return
	}
	return serdeAs[swagger.FullPaper](&paper) // to correct CorpusId type
}
