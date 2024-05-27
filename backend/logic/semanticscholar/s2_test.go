package semanticscholar_test

import (
	"context"
	s2 "learnerai/logic/semanticscholar"
	"learnerai/oapi/semanticscholar/swagger"
	"learnerai/req"
	"testing"
)

func TestS2Reference(t *testing.T) {
	r, err := s2.GetPaperReference(context.Background(), req.PaperReferencesFromS2Req{
		ArxivId: "2307.09288",
		Fields:  []string{"paperId", "externalIds", "title"},
	})
	assertSwaggerError(t, err)

	t.Logf("%#v", r)
	for _, v := range r {
		m := s2.CastExternalIds(v.ExternalIds)
		t.Logf("%+v", m)
	}
}

func assertSwaggerError(t *testing.T, err error) {
	if err != nil {
		t.Fatal("request failed:", err)
		if ge, ok := err.(swagger.GenericSwaggerError); ok {
			t.Errorf("%+v\n", ge.Model())
			t.Error(string(ge.Body()))
		}
	}
}

func TestS2Paper(t *testing.T) {
	p, err := s2.GetPaperInfo(context.Background(), req.PaperInfoFromS2Req{
		ArxivId: "1905.10044",
		Fields:  []string{"paperId", "externalIds", "title"},
	})
	assertSwaggerError(t, err)

	t.Logf("%+v", p)
	t.Logf("%+v", s2.CastExternalIds(p.ExternalIds))
}
