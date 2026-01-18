package requests

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BaseRequest struct{}

func Validate[T any](ctx *gin.Context) (*T, error) {
	var req T

	switch ctx.Request.Method {
	case http.MethodGet:
		if err := ctx.ShouldBindQuery(&req); err != nil {
			return nil, fmt.Errorf("request payload validate fail (Query), Err: %v", err)
		}
		return &req, nil

	case http.MethodDelete:
		// JSON -> Query
		jsonErr := ctx.ShouldBindBodyWithJSON(&req)
		if jsonErr == nil {
			return &req, nil
		}
		if err := ctx.ShouldBindQuery(&req); err != nil {
			return nil, fmt.Errorf("request payload validate fail (JSON/Query), JSON Err: %v, Query Err: %v", jsonErr, err)
		}
		return &req, nil

	default:
		// Content-Type: application/json => 只 bind JSON，不做 fallback（避免二次 bind 出 EOF）
		if isJSONRequest(ctx) {
			if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
				return nil, fmt.Errorf("request payload validate fail (JSON), Err: %v", err)
			}
			return &req, nil
		}

		// 非 JSON => bind Form（或其他 gin 自動 binder）
		if err := ctx.ShouldBind(&req); err != nil {
			return nil, fmt.Errorf("request payload validate fail (Form), Err: %v", err)
		}
		return &req, nil
	}
}

func isJSONRequest(ctx *gin.Context) bool {
	ct := strings.ToLower(ctx.GetHeader("Content-Type"))
	return strings.Contains(ct, "application/json")
}

func ValidateURI[T any](ctx *gin.Context) (*T, error) {
	var req T
	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, fmt.Errorf("request payload validate fail (URI), Err: %v", err)
	}
	return &req, nil
}
