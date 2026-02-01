package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"

	"github.com/redis/go-redis/v9"
)

// CacheKeyPrefix Redis Key 前綴
const CacheKeyPrefix = "timeledger"

// 預設快取時間
const (
	CacheDurationShort  = 5 * time.Minute  // 5分鐘 - 頻繁變動的資料
	CacheDurationMedium = 30 * time.Minute // 30分鐘 - 一般資料
	CacheDurationLong   = 24 * time.Hour   // 24小時 - 幾乎不變的資料
)

// CacheService 快取服務
type CacheService struct {
	BaseService
	app    *app.App
	redis  *redis.Client
}

// NewCacheService 建立 Cache Service
func NewCacheService(app *app.App) *CacheService {
	baseSvc := NewBaseService(app, "CacheService")
	return &CacheService{
		BaseService: *baseSvc,
		app:         app,
		redis:       app.Redis.DB0,
	}
}

// buildKey 組合快取 Key
func (s *CacheService) buildKey(category string, keys ...string) string {
	result := fmt.Sprintf("%s:%s", CacheKeyPrefix, category)
	for _, key := range keys {
		result += fmt.Sprintf(":%s", key)
	}
	return result
}

// Get 取得快取資料
func (s *CacheService) Get(ctx context.Context, category string, keys ...string) (string, error) {
	key := s.buildKey(category, keys...)
	return s.redis.Get(ctx, key).Result()
}

// GetBytes 取得快取資料（位元組）
func (s *CacheService) GetBytes(ctx context.Context, category string, keys ...string) ([]byte, error) {
	key := s.buildKey(category, keys...)
	data, err := s.redis.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return data, err
}

// GetJSON 取得快取資料並解析為 JSON
func (s *CacheService) GetJSON(ctx context.Context, dest interface{}, category string, keys ...string) error {
	data, err := s.GetBytes(ctx, category, keys...)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("cache miss")
	}
	return json.Unmarshal(data, dest)
}

// Set 設定快取（使用預設時間）
func (s *CacheService) Set(ctx context.Context, category string, value interface{}, keys ...string) error {
	return s.SetWithTTL(ctx, CacheDurationMedium, category, value, keys...)
}

// SetWithTTL 設定快取（自訂時間）
func (s *CacheService) SetWithTTL(ctx context.Context, ttl time.Duration, category string, value interface{}, keys ...string) error {
	key := s.buildKey(category, keys...)

	var data []byte
	var err error

	// 根據類型序列化
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
	}

	if err := s.redis.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete 刪除快取
func (s *CacheService) Delete(ctx context.Context, category string, keys ...string) error {
	key := s.buildKey(category, keys...)
	return s.redis.Del(ctx, key).Err()
}

// DeleteByPattern 依 Pattern 刪除快取
func (s *CacheService) DeleteByPattern(ctx context.Context, category string, pattern string) error {
	// 組合前綴和 pattern，使用通配符匹配
	fullPattern := fmt.Sprintf("%s:%s:%s", CacheKeyPrefix, category, pattern)
	s.Logger.Debug("deleting cache by pattern", "pattern", fullPattern)
	keys, err := s.redis.Keys(ctx, fullPattern).Result()
	if err != nil {
		s.Logger.Error("failed to get keys for deletion", "error", err)
		return err
	}
	s.Logger.Debug("found keys to delete", "keys", keys, "count", len(keys))
	if len(keys) > 0 {
		delErr := s.redis.Del(ctx, keys...).Err()
		if delErr != nil {
			s.Logger.Error("failed to delete keys", "error", delErr, "keys", keys)
		}
		return delErr
	}
	return nil
}

// Exists 檢查是否存在
func (s *CacheService) Exists(ctx context.Context, category string, keys ...string) (bool, error) {
	key := s.buildKey(category, keys...)
	result, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Increment 遞增計數器
func (s *CacheService) Increment(ctx context.Context, category string, keys ...string) error {
	key := s.buildKey(category, keys...)
	return s.redis.Incr(ctx, key).Err()
}

// Decrement 遞減計數器
func (s *CacheService) Decrement(ctx context.Context, category string, keys ...string) error {
	key := s.buildKey(category, keys...)
	return s.redis.Decr(ctx, key).Err()
}

// SetNX 僅在不存在時設定（用於防止快取擊穿）
func (s *CacheService) SetNX(ctx context.Context, ttl time.Duration, category string, value interface{}, keys ...string) (bool, error) {
	key := s.buildKey(category, keys...)

	var data []byte
	var err error

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return false, fmt.Errorf("failed to marshal value: %w", err)
		}
	}

	return s.redis.SetNX(ctx, key, data, ttl).Result()
}

// GetOrSet 取得或設定（常用模式）
func (s *CacheService) GetOrSet(ctx context.Context, ttl time.Duration, category string, fetchFunc func() (interface{}, error), keys ...string) (interface{}, error) {
	// 先嘗試從快取取得
	var result interface{}
	err := s.GetJSON(ctx, &result, category, keys...)
	if err == nil {
		return result, nil
	}

	// 快取未命中，從來源取得
	data, err := fetchFunc()
	if err != nil {
		return nil, err
	}

	// 存入快取（非同步，不影響主要流程）
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.SetWithTTL(ctx, ttl, category, data, keys...)
	}()

	return data, nil
}

// InvalidateCategory 使整個類別的快取失效
func (s *CacheService) InvalidateCategory(ctx context.Context, category string) error {
	return s.DeleteByPattern(ctx, category, "*")
}

// IsHealthy 檢查 Redis 連線
func (s *CacheService) IsHealthy(ctx context.Context) bool {
	return s.redis.Ping(ctx).Err() == nil
}

// =====================================================
// 特定業務的快取方法
// =====================================================

// CacheCategory 常量
const (
	CacheCategoryTeacher   = "teacher"
	CacheCategorySchedule  = "schedule"
	CacheCategoryException = "exception"
	CacheCategoryCenter    = "center"
	CacheCategoryOffering  = "offering"
	CacheCategoryCourse    = "course"
	CacheCategoryRoom      = "room"
	CacheCategoryGeo       = "geo"
	CacheCategoryUser      = "user"
	CacheCategoryStats     = "stats"
)

// GetTeacherProfile 取得教師資料（快取）
func (s *CacheService) GetTeacherProfile(ctx context.Context, teacherID uint) (*TeacherProfile, error) {
	var profile TeacherProfile
	err := s.GetJSON(ctx, &profile, CacheCategoryTeacher, fmt.Sprintf("profile:%d", teacherID))
	return &profile, err
}

// SetTeacherProfile 設定教師資料（快取）
func (s *CacheService) SetTeacherProfile(ctx context.Context, teacherID uint, profile *TeacherProfile) error {
	return s.SetWithTTL(ctx, CacheDurationMedium, CacheCategoryTeacher, profile, fmt.Sprintf("profile:%d", teacherID))
}

// InvalidateTeacher 使教師快取失效
func (s *CacheService) InvalidateTeacher(ctx context.Context, teacherID uint) error {
	return s.Delete(ctx, CacheCategoryTeacher, fmt.Sprintf("profile:%d", teacherID))
}

// GetSchedule 取得課表（快取）
func (s *CacheService) GetSchedule(ctx context.Context, teacherID uint, weekStart string) (*ScheduleData, error) {
	var schedule ScheduleData
	err := s.GetJSON(ctx, &schedule, CacheCategorySchedule, fmt.Sprintf("%d:%s", teacherID, weekStart))
	return &schedule, err
}

// SetSchedule 設定課表（快取）
func (s *CacheService) SetSchedule(ctx context.Context, teacherID uint, weekStart string, schedule *ScheduleData) error {
	return s.SetWithTTL(ctx, CacheDurationShort, CacheCategorySchedule, schedule, fmt.Sprintf("%d:%s", teacherID, weekStart))
}

// InvalidateSchedule 使課表快取失效
func (s *CacheService) InvalidateSchedule(ctx context.Context, teacherID uint) error {
	return s.DeleteByPattern(ctx, CacheCategorySchedule, fmt.Sprintf("%d:*", teacherID))
}

// GetExceptionStats 取得例外統計（快取）
func (s *CacheService) GetExceptionStats(ctx context.Context, centerID uint) (map[string]int, error) {
	var stats map[string]int
	err := s.GetJSON(ctx, &stats, CacheCategoryStats, fmt.Sprintf("exceptions:%d", centerID))
	return stats, err
}

// SetExceptionStats 設定例外統計（快取）
func (s *CacheService) SetExceptionStats(ctx context.Context, centerID uint, stats map[string]int) error {
	return s.SetWithTTL(ctx, CacheDurationShort, CacheCategoryStats, stats, fmt.Sprintf("exceptions:%d", centerID))
}

// =====================================================
// 輔助類型定義
// =====================================================

// TeacherProfile 教師基本資料
type TeacherProfile struct {
	ID             uint     `json:"id"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	AvatarURL      string   `json:"avatar_url"`
	Bio            string   `json:"bio"`
	City           string   `json:"city"`
	District       string   `json:"district"`
	IsOpenToHiring bool     `json:"is_open_to_hiring"`
	CenterCount    int      `json:"center_count"`
	Skills         []string `json:"skills"`
	Certificates   int      `json:"certificates"`
}

// ScheduleData 課表資料
type ScheduleData struct {
	Days          []ScheduleDay `json:"days"`
	WeekStart     string        `json:"week_start"`
	WeekEnd       string        `json:"week_end"`
	TotalHours    float64       `json:"total_hours"`
	TotalSessions int           `json:"total_sessions"`
}

// ScheduleDay 課表某天
type ScheduleDay struct {
	Date    string         `json:"date"`
	Weekday string         `json:"weekday"`
	Items   []ScheduleItem `json:"items"`
}

// ScheduleItem 課表項目
type ScheduleItem struct {
	ID         uint   `json:"id"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Date       string `json:"date"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Duration   int    `json:"duration"`
	RoomID     uint   `json:"room_id"`
	RoomName   string `json:"room_name"`
	CenterID   uint   `json:"center_id"`
	CenterName string `json:"center_name"`
	Color      string `json:"color"`
	Status     string `json:"status"`
}

// =====================================================
// Center 快取方法
// =====================================================

// GetCenterSettings 取得中心設置（快取）
func (s *CacheService) GetCenterSettings(ctx context.Context, centerID uint) (*models.CenterSettings, error) {
	var settings models.CenterSettings
	err := s.GetJSON(ctx, &settings, CacheCategoryCenter, fmt.Sprintf("settings:%d", centerID))
	return &settings, err
}

// SetCenterSettings 設定中心設置（快取）
func (s *CacheService) SetCenterSettings(ctx context.Context, centerID uint, settings *models.CenterSettings) error {
	return s.SetWithTTL(ctx, CacheDurationLong, CacheCategoryCenter, settings, fmt.Sprintf("settings:%d", centerID))
}

// InvalidateCenterSettings 使中心設置快取失效
func (s *CacheService) InvalidateCenterSettings(ctx context.Context, centerID uint) error {
	return s.Delete(ctx, CacheCategoryCenter, fmt.Sprintf("settings:%d", centerID))
}

// GetCenterBasic 取得中心基本資訊（快取）
func (s *CacheService) GetCenterBasic(ctx context.Context, centerID uint) (*CenterBasicInfo, error) {
	var info CenterBasicInfo
	err := s.GetJSON(ctx, &info, CacheCategoryCenter, fmt.Sprintf("basic:%d", centerID))
	return &info, err
}

// SetCenterBasic 設定中心基本資訊（快取）
func (s *CacheService) SetCenterBasic(ctx context.Context, centerID uint, info *CenterBasicInfo) error {
	return s.SetWithTTL(ctx, CacheDurationMedium, CacheCategoryCenter, info, fmt.Sprintf("basic:%d", centerID))
}

// InvalidateCenterBasic 使中心基本資訊快取失效
func (s *CacheService) InvalidateCenterBasic(ctx context.Context, centerID uint) error {
	return s.Delete(ctx, CacheCategoryCenter, fmt.Sprintf("basic:%d", centerID))
}

// CenterBasicInfo 中心基本資訊（用於快取）
type CenterBasicInfo struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	PlanLevel string `json:"plan_level"`
	IsActive  bool   `json:"is_active"`
}

// =====================================================
// Course 快取方法
// =====================================================

// GetCourseList 取得課程列表（快取）
func (s *CacheService) GetCourseList(ctx context.Context, centerID uint) ([]CourseCacheItem, error) {
	var courses []CourseCacheItem
	err := s.GetJSON(ctx, &courses, CacheCategoryCourse, fmt.Sprintf("list:%d", centerID))
	
	// 快取不存在或解析失敗時，返回空切片而非錯誤
	if err != nil {
		s.Logger.Debug("course cache miss or parse error", "center_id", centerID, "error", err)
		return []CourseCacheItem{}, nil
	}
	
	return courses, nil
}

// SetCourseList 設定課程列表（快取）
func (s *CacheService) SetCourseList(ctx context.Context, centerID uint, courses []CourseCacheItem) error {
	s.Logger.Debug("caching course list", "center_id", centerID, "count", len(courses))
	return s.SetWithTTL(ctx, CacheDurationMedium, CacheCategoryCourse, courses, fmt.Sprintf("list:%d", centerID))
}

// InvalidateCourseList 使課程列表快取失效
func (s *CacheService) InvalidateCourseList(ctx context.Context, centerID uint) error {
	pattern := fmt.Sprintf("list:%d", centerID)
	s.Logger.Debug("invalidating course list cache", "center_id", centerID, "pattern", pattern)
	return s.DeleteByPattern(ctx, CacheCategoryCourse, pattern)
}

// CourseCacheItem 課程快取項目
type CourseCacheItem struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	DefaultDuration  int    `json:"default_duration"`
	ColorHex         string `json:"color_hex"`
	RoomBufferMin    int    `json:"room_buffer_min"`
	TeacherBufferMin int    `json:"teacher_buffer_min"`
	IsActive         bool   `json:"is_active"`
}

// =====================================================
// Room 快取方法
// =====================================================

// GetRoomList 取得教室列表（快取）
func (s *CacheService) GetRoomList(ctx context.Context, centerID uint) ([]RoomCacheItem, error) {
	var rooms []RoomCacheItem
	err := s.GetJSON(ctx, &rooms, CacheCategoryRoom, fmt.Sprintf("list:%d", centerID))
	return rooms, err
}

// SetRoomList 設定教室列表（快取）
func (s *CacheService) SetRoomList(ctx context.Context, centerID uint, rooms []RoomCacheItem) error {
	s.Logger.Debug("caching room list", "center_id", centerID, "count", len(rooms))
	return s.SetWithTTL(ctx, CacheDurationMedium, CacheCategoryRoom, rooms, fmt.Sprintf("list:%d", centerID))
}

// InvalidateRoomList 使教室列表快取失效
func (s *CacheService) InvalidateRoomList(ctx context.Context, centerID uint) error {
	pattern := fmt.Sprintf("list:%d", centerID)
	s.Logger.Debug("invalidating room list cache", "center_id", centerID, "pattern", pattern)
	return s.DeleteByPattern(ctx, CacheCategoryRoom, pattern)
}

// RoomCacheItem 教室快取項目
type RoomCacheItem struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	IsActive bool   `json:"is_active"`
}

// =====================================================
// Today Schedule 快取方法
// =====================================================

// GetTodaySchedule 取得今日排課（快取）
func (s *CacheService) GetTodaySchedule(ctx context.Context, centerID uint, date string) (*TodayScheduleCache, error) {
	var schedule TodayScheduleCache
	err := s.GetJSON(ctx, &schedule, CacheCategorySchedule, fmt.Sprintf("today:%d:%s", centerID, date))
	return &schedule, err
}

// SetTodaySchedule 設定今日排課（快取）
func (s *CacheService) SetTodaySchedule(ctx context.Context, centerID uint, date string, schedule *TodayScheduleCache) error {
	return s.SetWithTTL(ctx, CacheDurationShort, CacheCategorySchedule, schedule, fmt.Sprintf("today:%d:%s", centerID, date))
}

// InvalidateTodaySchedule 使今日排課快取失效
func (s *CacheService) InvalidateTodaySchedule(ctx context.Context, centerID uint, date string) error {
	return s.Delete(ctx, CacheCategorySchedule, fmt.Sprintf("today:%d:%s", centerID, date))
}

// InvalidateAllTodaySchedules 使某中心所有今日排課快取失效
func (s *CacheService) InvalidateAllTodaySchedules(ctx context.Context, centerID uint) error {
	return s.DeleteByPattern(ctx, CacheCategorySchedule, fmt.Sprintf("today:%d:*", centerID))
}

// TodayScheduleCache 今日排課快取結構
type TodayScheduleCache struct {
	Date       string              `json:"date"`
	CenterID   uint                `json:"center_id"`
	Sessions   []SessionCacheItem  `json:"sessions"`
	TotalCount int                 `json:"total_count"`
	CachedAt   time.Time           `json:"cached_at"`
}

// SessionCacheItem 課堂快取項目
type SessionCacheItem struct {
	ID         uint   `json:"id"`
	RuleID     uint   `json:"rule_id"`
	CourseName string `json:"course_name"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	RoomID     uint   `json:"room_id"`
	RoomName   string `json:"room_name"`
	TeacherID  uint   `json:"teacher_id"`
	TeacherName string `json:"teacher_name"`
	Status     string `json:"status"`
	Color      string `json:"color"`
}

// =====================================================
// Teacher Profile 快取方法
// =====================================================

// 注意：TeacherProfile 快取方法已在前面的區塊定義

// =====================================================
// Exception Stats 快取方法
// =====================================================

// 注意：ExceptionStats 快取方法已在前面的區塊定義
