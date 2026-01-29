package requests

import "github.com/gin-gonic/gin"

type TalentSearchRequest struct {
	CenterID uint   `json:"-"`
	City     string `form:"city"`
	District string `form:"district"`
	Keyword  string `form:"keyword"`
	Skills   string `form:"skills"`
	Hashtags string `form:"hashtags"`

	// 分頁參數
	Page    int `form:"page"`
	Limit   int `form:"limit"`
	Total   int `json:"total"` // 總筆數（回應用）

	// 排序參數
	SortBy    string `form:"sort_by"`    // name, skills, rating, created_at
	SortOrder string `form:"sort_order"` // asc, desc
}

func ValidateTalentSearch(ctx *gin.Context) (*TalentSearchRequest, error) {
	var req TalentSearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	// 設定預設值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20 // 預設每頁 20 筆
	}
	if req.SortBy == "" {
		req.SortBy = "name" // 預設按姓名排序
	}
	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	return &req, nil
}
