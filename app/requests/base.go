package requests

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseRequest struct{}

// 泛型檢查
func Validate[T any](ctx *gin.Context) (*T, error) {
	var req T
	var err error

	switch ctx.Request.Method {
	case http.MethodGet, http.MethodDelete:
		// GET/DELETE 取 query 參數
		err = ctx.ShouldBindQuery(&req)
	default:
		// 嘗試綁定 JSON body
		err = ctx.ShouldBindJSON(&req)
		if err != nil {
			// 若不是 JSON，再試 Form Data
			if formErr := ctx.ShouldBind(&req); formErr != nil {
				// 返回 Form Data 綁定錯誤
				return nil, fmt.Errorf("Request payload validate fail (Form Data), Err: %v", formErr)
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("Request payload validate fail (JSON/Query), Err: %v", err)
	}

	return &req, nil
}
