package req

type ListUserReq struct {
	Code  string `json:"code"`
	Start int64  `json:"start"`
	Limit int64  `json:"limit"`
}

type UpdateUserInfoReq struct {
	Code         string   `json:"code"`
	FollowAuthor []string `json:"follow_author"` //关注作者列表
	FollowTag    []string `json:"follow_tag"`    //关注标签列表
	FocusArea    []string `json:"focus_area"`    //关注领域列表
}

type UpdateUserNewsReq struct {
	Code   string `json:"code"`
	DocId  string `json:"doc_id"`
	Like   bool   `json:"like"`   //用户喜欢的news id
	UnLike bool   `json:"unlike"` //用户不喜欢的news id
	Read   bool   `json:"read"`   //用户已经阅读的news id
}

type GetMiniProgramCode struct {
	Code string `json:"code"`
}

type GetMiniProgramCodeRes struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
}

type GetTokenReq struct {
	Uid string `json:"uid"`
}
