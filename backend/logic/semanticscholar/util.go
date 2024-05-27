package semanticscholar

import (
	"bytes"
	"encoding/json"
	"learnerai/oapi/semanticscholar/swagger"
	"strings"

	"github.com/antihax/optional"
)

var s2client *swagger.APIClient

func init() {
	conf := swagger.NewConfiguration()
	// conf.Host = "api.semanticscholar.org"
	// conf.Scheme = "https"
	conf.BasePath = "https://api.semanticscholar.org/graph/v1" // this is a hack, Host and Scheme don't work
	s2client = swagger.NewAPIClient(conf)
}

func toS2PaperId(arxivId string) string {
	return "ARXIV:" + getShortId(arxivId)
}

func getShortId(arxivId string) string {
	dud, shortId, found := strings.Cut(arxivId, "arxiv.org/abs/")
	if !found {
		shortId = dud
	}
	return shortId
}

func formatFields(fields []string) optional.String {
	if len(fields) == 0 {
		return optional.EmptyString()
	}
	return optional.NewString(strings.Join(fields, ","))
}

func serdeAs[T any](src any) (dst T, err error) {
	b, err := json.Marshal(src)
	if err != nil {
		return
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber() // since we are handling any fields, keep them safe. NOTABLY CorpusId is int64 (invalid json if pedantic).
	err = d.Decode(&dst)
	return
}

// most external ids are just string, but some (notably CorpusId!) are number (parsed into json.Number, which is string).
// works with:
// swagger.AllOfReferenceCitedPaper
// swagger.FullPaper
func CastExternalIds(m *interface{}) map[string]any {
	if m != nil {
		return (*m).(map[string]any)
	}
	return nil
}
