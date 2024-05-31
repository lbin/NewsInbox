package req

type NewsInfoFromS2Req struct {
	ArxivId string   `json:"arxiv_id"`
	Fields  []string `json:"fields,omitempty"`
}

type NewsReferencesFromS2Req struct {
	ArxivId string   `json:"arxiv_id"`
	Fields  []string `json:"fields,omitempty"`
}
