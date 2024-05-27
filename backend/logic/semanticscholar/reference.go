package semanticscholar

import (
	"context"
	"learnerai/oapi/semanticscholar/swagger"
	"learnerai/req"

	"github.com/antihax/optional"
)

// TODO: the result should be cached
func GetPaperReference(ctx context.Context, r req.PaperReferencesFromS2Req) (refs []swagger.AllOfReferenceCitedPaper, err error) {
	paperId := toS2PaperId(r.ArxivId)

	opts := swagger.PaperDataApiGetGraphGetPaperReferencesOpts{
		Offset: optional.NewInt32(0), // volatile
		Limit:  optional.NewInt32(100),
		Fields: formatFields(r.Fields),
	}
	for {
		r, _, err := s2client.PaperDataApi.GetGraphGetPaperReferences(ctx, paperId, &opts)
		if err != nil {
			return nil, err
		}
		if r.Next == 0 {
			// `next` is marked as required, but document says in title "absent if no more data exists".
			// so in current codegen, it is deserialized as 0 as golang handles it.
			// however other codegen may enforce the requirement check so API call will fail miserably.
			break
		}
		for _, data := range r.Data {
			if data.CitedPaper != nil {
				// *data.CitedPaper is map[string]interface{} because of flawed codegen
				v, err := serdeAs[swagger.AllOfReferenceCitedPaper](data.CitedPaper)
				if err != nil {
					return nil, err
				}
				refs = append(refs, v)
			}
		}
		opts.Offset = optional.NewInt32(r.Next)
	}
	return
}
