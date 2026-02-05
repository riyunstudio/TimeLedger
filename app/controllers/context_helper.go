package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"timeLedger/global"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
)

// ContextHelper 封裝 gin.Context 的常用操作
// 消除 Controller 中重複的型別轉換程式碼
type ContextHelper struct {
	ctx *gin.Context
}

// NewContextHelper 建立 ContextHelper 實例
func NewContextHelper(ctx *gin.Context) *ContextHelper {
	return &ContextHelper{ctx: ctx}
}

// UserID 取得當前登入使用者 ID
// 回傳 (uint, error) - 如果未登入或取得失敗，回傳 error
func (h *ContextHelper) UserID() (uint, error) {
	uid := h.ctx.GetUint(global.UserIDKey)
	if uid == 0 {
		return 0, fmt.Errorf("user not authenticated")
	}
	return uid, nil
}

// MustUserID 取得當前登入使用者 ID
// 如果未登入，直接回傳 401 錯誤
func (h *ContextHelper) MustUserID() uint {
	uid, err := h.UserID()
	if err != nil {
		h.Unauthorized("Authentication required")
		return 0
	}
	return uid
}

// CenterID 取得中心 ID
// 回傳 (uint, error) - 如果無法取得，回傳 error
func (h *ContextHelper) CenterID() (uint, error) {
	centerID := h.ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		return 0, fmt.Errorf("center not found")
	}
	return centerID, nil
}

// MustCenterID 取得中心 ID
// 如果無法取得，直接回傳 401 錯誤
func (h *ContextHelper) MustCenterID() uint {
	centerID, err := h.CenterID()
	if err != nil {
		h.Unauthorized("Center ID required")
		return 0
	}
	return centerID
}

// UserType 取得使用者類型
// 回傳 (string, bool) - 回傳類型與是否存在
func (h *ContextHelper) UserType() (string, bool) {
	userType, exists := h.ctx.Get(global.UserTypeKey)
	if !exists {
		return "", false
	}
	return userType.(string), true
}

// IsAdmin 檢查是否為管理員
func (h *ContextHelper) IsAdmin() bool {
	userType, _ := h.UserType()
	return userType == "ADMIN" || userType == "OWNER"
}

// IsTeacher 檢查是否為老師
func (h *ContextHelper) IsTeacher() bool {
	userType, _ := h.UserType()
	return userType == "TEACHER"
}

// LineUserID 取得 LINE User ID
func (h *ContextHelper) LineUserID() (string, bool) {
	lineUserID, exists := h.ctx.Get(global.LineUserIDKey)
	if !exists {
		return "", false
	}
	return lineUserID.(string), true
}

// ParamUint 從 URL 參數取得 uint 值
func (h *ContextHelper) ParamUint(key string) (uint, error) {
	val := h.ctx.Param(key)
	if val == "" {
		return 0, fmt.Errorf("parameter '%s' is required", key)
	}
	var result uint
	if _, err := fmt.Sscanf(val, "%d", &result); err != nil {
		return 0, fmt.Errorf("invalid parameter '%s': %w", key, err)
	}
	if result == 0 {
		return 0, fmt.Errorf("parameter '%s' must be greater than 0", key)
	}
	return result, nil
}

// MustParamUint 從 URL 參數取得 uint 值
// 如果失敗，回傳 400 錯誤
func (h *ContextHelper) MustParamUint(key string) uint {
	val, err := h.ParamUint(key)
	if err != nil {
		h.BadRequest(err.Error())
		return 0
	}
	return val
}

// QueryString 取得查詢參數字串
func (h *ContextHelper) QueryString(key string) (string, bool) {
	val := h.ctx.Query(key)
	return val, val != ""
}

// QueryStringOrDefault 取得查詢參數字串，如果不存在則使用預設值
func (h *ContextHelper) QueryStringOrDefault(key, defaultVal string) string {
	val := h.ctx.Query(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// QueryIntOrDefault 取得查詢參數整數值，如果不存在或解析失敗則使用預設值
func (h *ContextHelper) QueryIntOrDefault(key string, defaultVal int) int {
	val := h.ctx.Query(key)
	if val == "" {
		return defaultVal
	}
	result, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return result
}

// QueryUint 從查詢參數取得 uint 值
func (h *ContextHelper) QueryUint(key string) (uint, error) {
	val := h.ctx.Query(key)
	if val == "" {
		return 0, fmt.Errorf("query parameter '%s' is required", key)
	}
	result, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid query parameter '%s': %w", key, err)
	}
	return uint(result), nil
}

// MustQueryUint 從查詢參數取得 uint 值
// 如果失敗，回傳 400 錯誤
func (h *ContextHelper) MustQueryUint(key string) uint {
	val, err := h.QueryUint(key)
	if err != nil {
		h.BadRequest(err.Error())
		return 0
	}
	return val
}

// QueryDate 從查詢參數取得日期 (YYYY-MM-DD 格式)
func (h *ContextHelper) QueryDate(key string) (time.Time, error) {
	val := h.ctx.Query(key)
	if val == "" {
		return time.Time{}, fmt.Errorf("query parameter '%s' is required", key)
	}
	date, err := time.Parse("2006-01-02", val)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format for '%s': %w", key, err)
	}
	return date, nil
}

// MustQueryDate 從查詢參數取得日期
// 如果失敗，回傳 400 錯誤
func (h *ContextHelper) MustQueryDate(key string) time.Time {
	date, err := h.QueryDate(key)
	if err != nil {
		h.BadRequest(err.Error())
		return time.Time{}
	}
	return date
}

// QueryDateRange 取得日期範圍 (from, to)
func (h *ContextHelper) QueryDateRange(fromKey, toKey string) (from, to time.Time, err error) {
	from, err = h.QueryDate(fromKey)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	to, err = h.QueryDate(toKey)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return from, to, nil
}

// MustQueryDateRange 取得日期範圍
// 如果失敗，回傳 400 錯誤
func (h *ContextHelper) MustQueryDateRange(fromKey, toKey string) (from, to time.Time) {
	from, err := h.QueryDate(fromKey)
	if err != nil {
		h.BadRequest(err.Error())
		return time.Time{}, time.Time{}
	}
	to, err = h.QueryDate(toKey)
	if err != nil {
		h.BadRequest(err.Error())
		return time.Time{}, time.Time{}
	}
	return from, to
}

// BindJSON 綁定 JSON 請求體
func (h *ContextHelper) BindJSON(ptr any) error {
	if err := h.ctx.ShouldBindJSON(ptr); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}
	return nil
}

// MustBindJSON 綁定 JSON 請求體
// 如果失敗，回傳 400 錯誤
func (h *ContextHelper) MustBindJSON(ptr any) bool {
	if err := h.BindJSON(ptr); err != nil {
		h.BadRequest(err.Error())
		return false
	}
	return true
}

// GinContext 回傳原始 gin.Context
func (h *ContextHelper) GinContext() *gin.Context {
	return h.ctx
}

// Context 回傳 context.Context
func (h *ContextHelper) Context() interface{} {
	return h.ctx.Request.Context()
}

// ========== 響應方法 ==========

// Success 回傳成功響應
func (h *ContextHelper) Success(data any) {
	h.ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   data,
	})
}

// Created 回傳建立成功響應
func (h *ContextHelper) Created(data any) {
	h.ctx.JSON(http.StatusCreated, global.ApiResponse{
		Code:    0,
		Message: "Created",
		Datas:   data,
	})
}

// NoContent 回傳無內容響應
func (h *ContextHelper) NoContent() {
	h.ctx.Status(http.StatusNoContent)
}

// BadRequest 回傳 400 錯誤
func (h *ContextHelper) BadRequest(message string) {
	h.ctx.JSON(http.StatusBadRequest, global.ApiResponse{
		Code:    global.BAD_REQUEST,
		Message: message,
	})
}

// Unauthorized 回傳 401 錯誤
func (h *ContextHelper) Unauthorized(message string) {
	h.ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
		Code:    global.UNAUTHORIZED,
		Message: message,
	})
}

// Forbidden 回傳 403 錯誤
func (h *ContextHelper) Forbidden(message string) {
	h.ctx.JSON(http.StatusForbidden, global.ApiResponse{
		Code:    global.FORBIDDEN,
		Message: message,
	})
}

// NotFound 回傳 404 錯誤
func (h *ContextHelper) NotFound(message string) {
	h.ctx.JSON(http.StatusNotFound, global.ApiResponse{
		Code:    errInfos.NOT_FOUND,
		Message: message,
	})
}

// Conflict 回傳 409 錯誤
func (h *ContextHelper) Conflict(message string) {
	h.ctx.JSON(http.StatusConflict, global.ApiResponse{
		Code:    409,
		Message: message,
	})
}

// InternalError 回傳 500 錯誤
func (h *ContextHelper) InternalError(message string) {
	h.ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
		Code:    errInfos.SYSTEM_ERROR,
		Message: message,
	})
}

// ErrorWithCode 回傳自訂錯誤碼的錯誤
func (h *ContextHelper) ErrorWithCode(status int, code errInfos.ErrCode, message string) {
	h.ctx.JSON(status, global.ApiResponse{
		Code:    code,
		Message: message,
	})
}

// ErrorWithInfo 回傳使用 errInfos.Res 的錯誤
func (h *ContextHelper) ErrorWithInfo(errInfo *errInfos.Res) {
	// 檢查 errInfo 是否為 nil
	if errInfo == nil {
		// PANIC: 追蹤問題根源
		panic(fmt.Sprintf("FATAL: ErrorWithInfo received nil errInfo! \nPath: %s \nMethod: %s \nStack: %s",
			h.ctx.Request.URL.Path,
			h.ctx.Request.Method,
			string(h.ctx.Request.URL.RawQuery)))
	}

	// 根據錯誤碼決定 HTTP 狀態碼
	status := http.StatusInternalServerError
	switch errInfo.Code {
	case global.UNAUTHORIZED:
		status = http.StatusUnauthorized
	case global.FORBIDDEN:
		status = http.StatusForbidden
	case errInfos.NOT_FOUND:
		status = http.StatusNotFound
	case global.BAD_REQUEST, errInfos.PARAMS_VALIDATE_ERROR:
		status = http.StatusBadRequest
	case errInfos.INVALID_STATUS, errInfos.SCHED_OVERLAP, errInfos.SCHED_BUFFER,
		errInfos.SCHED_RULE_CONFLICT, errInfos.ERR_RESOURCE_LOCKED,
		errInfos.ERR_CONCURRENT_MODIFIED, errInfos.ERR_TX_FAILED:
		status = http.StatusConflict
	}
	h.ctx.JSON(status, global.ApiResponse{
		Code:    errInfo.Code,
		Message: errInfo.Msg,
	})
}

// File 回傳檔案
func (h *ContextHelper) File(data []byte, filename string, contentType string) {
	h.ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	h.ctx.Header("Content-Type", contentType)
	h.ctx.Data(http.StatusOK, contentType, data)
}

// CSV 回傳 CSV 檔案
func (h *ContextHelper) CSV(data []byte, filename string) {
	h.ctx.Header("Content-Type", "text/csv; charset=utf-8")
	h.ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	h.ctx.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}

// ICS 回傳 ICS 檔案
func (h *ContextHelper) ICS(data []byte, filename string) {
	h.ctx.Header("Content-Type", "text/calendar; charset=utf-8")
	h.ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	h.ctx.Data(http.StatusOK, "text/calendar; charset=utf-8", data)
}

// Image 回傳圖片
func (h *ContextHelper) Image(data []byte, contentType string) {
	h.ctx.Header("Content-Type", contentType)
	h.ctx.Data(http.StatusOK, contentType, data)
}

// Redirect 回傳重新導向
func (h *ContextHelper) Redirect(url string) {
	h.ctx.Redirect(http.StatusFound, url)
}
