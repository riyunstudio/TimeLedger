package requests

import "github.com/gin-gonic/gin"

type TalentSearchRequest struct {
	CenterID uint   `json:"-"`
	City     string `form:"city"`
	District string `form:"district"`
	Keyword  string `form:"keyword"`
	Skills   string `form:"skills"`
	Hashtags string `form:"hashtags"`
}

func ValidateTalentSearch(ctx *gin.Context) (*TalentSearchRequest, error) {
	var req TalentSearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
