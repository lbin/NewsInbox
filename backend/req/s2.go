package req

type PaperInfoFromS2Req struct {
	ArxivId string   `json:"arxiv_id"`
	Fields  []string `json:"fields,omitempty"`
}

type PaperReferencesFromS2Req struct {
	ArxivId string   `json:"arxiv_id"`
	Fields  []string `json:"fields,omitempty"`
}
