package requests

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

// InitValidators 註冊自訂驗證器到 Gin's validator engine（在程式啟動時呼叫）
func InitValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date_format", DateFormatValidator)
		v.RegisterValidation("time_format", TimeFormatValidator)
	}
}

// DateFormatValidator 日期格式驗證 (YYYY-MM-DD)
var DateFormatValidator = func(fl validator.FieldLevel) bool {
	dateStr, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// TimeFormatValidator 時間格式驗證 (HH:MM 或 HH:MM:SS)
var TimeFormatValidator = func(fl validator.FieldLevel) bool {
	timeStr, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	_, err := time.Parse("15:04", timeStr)
	if err != nil {
		_, err = time.Parse("15:04:05", timeStr)
	}
	return err == nil
}
