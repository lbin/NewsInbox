package req

import "time"

type ListPaperInfoReq struct {
	StartTime     time.Time `json:"start_time"`
	Id            string    `json:"id"`
	Tag           []string  `json:"tag"`
	EndTime       time.Time `json:"end_time"`
	SearchContent string    `json:"search_content"`
	Start         int64     `json:"start"`
	Limit         int64     `json:"limit"`
}

type UnsetDailyInfoField struct {
	Id     string   `json:"id"`
	Fields []string `json:"fields"`
}

type ResetPaperInfoReq struct {
	DocId           string `json:"doc_id"`
	ResetTitleCN    bool   `json:"reset_title_ch"`
	ResetKeywordsCN bool   `json:"reset_keywords_ch"`
	ResetKeywords   bool   `json:"reset_keywords"`
	ResetAbstractCN bool   `json:"reset_abstract_ch"`
}

type SearchSimilarReq struct {
	DocId string `json:"doc_id"`
	TopK  int    `json:"top_k"`
}
