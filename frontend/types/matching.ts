/**
 * 智慧媒合相關類型定義
 *
 * 包含人才庫搜尋、智慧媒合、替代時段建議等相關類型
 */

import type { ID, Timestamp, DateString, ApiResponse, PaginationParams, PaginatedResponse } from './api'

// ==================== 智慧媒合相關類型 ====================

/**
 * 智慧媒合搜尋請求
 */
export interface SmartMatchingRequest {
  /** 開課 ID */
  offering_id: ID
  /** 目標日期 */
  target_date: DateString
  /** 目標開始時間 */
  target_start_time: string
  /** 目標結束時間 */
  target_end_time: string
  /** 需要的技能 (可選) */
  required_skills?: string[]
  /** 偏好的縣市 */
  preferred_city?: string
  /** 最低評分 */
  min_rating?: number
  /** 最大結果數 */
  limit?: number
}

/**
 * 媒合分數詳細資訊
 */
export interface MatchScoreDetail {
  /** 媒合分數 (0-100) */
  score: number
  /** 可用性分數 (0-40) */
  availability_score: number
  /** 內部評價分數 (0-40) */
  evaluation_score: number
  /** 技能/地區匹配分數 (0-20) */
  match_score: number
}

/**
 * 智慧媒合結果
 */
export interface SmartMatchingResult {
  /** 教師 ID */
  teacher_id: ID
  /** 教師名稱 */
  teacher_name: string
  /** 教師頭像 */
  avatar_url?: string
  /** 媒合分數 */
  match_score: number
  /** 分數詳細 */
  score_detail: MatchScoreDetail
  /** 技能匹配列表 */
  matched_skills: string[]
  /** 平均評分 */
  rating: number
  /** 可用性狀態 */
  availability: 'AVAILABLE' | 'BUSY' | 'UNAVAILABLE'
  /** 衝突說明 (若有) */
  conflict_note?: string
  /** 距離 (公里) */
  distance_km?: number
}

/**
 * 智慧媒合回應
 */
export type SmartMatchingResponse = ApiResponse<SmartMatchingResult[]>

/**
 * 智慧媒合搜尋建議
 */
export interface MatchingSuggestion {
  /** 建議類型 */
  type: 'SKILL' | 'TIME' | 'ROOM'
  /** 建議內容 */
  message: string
  /** 建議值 */
  value: string | string[]
}

/**
 * 搜尋建議回應
 */
export type MatchingSuggestionsResponse = ApiResponse<MatchingSuggestion[]>

// ==================== 替代時段建議相關類型 ====================

/**
 * 替代時段建議請求
 */
export interface AlternativeSlotsRequest {
  /** 開課 ID */
  offering_id: ID
  /** 原定日期 */
  original_date: DateString
  /** 原定開始時間 */
  original_start_time: string
  /** 原定結束時間 */
  original_end_time: string
  /** 搜尋天數範圍 */
  days_range?: number
  /** 偏好的教師 ID 清單 (可選) */
  preferred_teacher_ids?: ID[]
}

/**
 * 替代時段
 */
export interface AlternativeSlot {
  /** 建議日期 */
  date: DateString
  /** 開始時間 */
  start_time: string
  /** 結束時間 */
  end_time: string
  /** 可用教師數量 */
  available_teacher_count: number
  /** 衝突說明 */
  conflict_note?: string
  /** 媒合分數 */
  match_score?: number
}

/**
 * 替代時段回應
 */
export type AlternativeSlotsResponse = ApiResponse<AlternativeSlot[]>

// ==================== 人才庫相關類型 ====================

/**
 * 人才搜尋請求
 */
export interface TalentSearchRequest {
  /** 搜尋關鍵字 */
  keyword?: string
  /** 技能分類 */
  category?: string
  /** 技能名稱 */
  skill_name?: string
  /** 縣市 */
  city?: string
  /** 區域 */
  district?: string
  /** 是否僅顯示開放應徵 */
  open_to_hiring_only?: boolean
  /** 最低評分 */
  min_rating?: number
  /** 狀態篩選 */
  status?: 'ACTIVE' | 'ALL'
}

/**
 * 人才庫統計
 */
export interface TalentPoolStats {
  /** 總人數 */
  total_count: number
  /** 開放應徵人數 */
  open_hiring_count: number
  /** 已加入中心人數 */
  member_count: number
  /** 平均評分 */
  average_rating: number
  /** 月變化量 */
  monthly_change: number
  /** 月趨勢 (過去 7 個月) */
  monthly_trend: number[]
  /** 待處理邀請數 */
  pending_invites: number
  /** 已接受邀請數 */
  accepted_invites: number
  /** 已拒絕邀請數 */
  declined_invites: number
  /** 縣市分布 */
  city_distribution: Array<{
    name: string
    count: number
  }>
  /** 熱門技能 */
  top_skills: Array<{
    name: string
    count: number
  }>
}

/**
 * 人才庫統計回應
 */
export type TalentPoolStatsResponse = ApiResponse<TalentPoolStats>

/**
 * 人才卡片資料
 */
export interface TalentCard {
  /** 教師 ID */
  id: ID
  /** 名稱 */
  name: string
  /** 頭像 URL */
  avatar_url?: string
  /** 縣市 */
  city?: string
  /** 區域 */
  district?: string
  /** 技能列表 */
  skills: Array<{
    category: string
    skill_name: string
  }>
  /** 證照數量 */
  certificate_count: number
  /** 平均評分 */
  rating: number
  /** 評價數量 */
  review_count: number
  /** 是否開放應徵 */
  is_open_to_hiring: boolean
  /** 是否已加入中心 */
  is_member: boolean
  /** 簡介 */
  bio?: string
  /** 個人標籤 */
  personal_hashtags?: string[]
  /** 公開聯絡資訊 */
  public_contact_info?: string
  /** 證照詳細列表 */
  certificates?: Array<{
    id: ID
    name: string
    issuer?: string
    obtained_at?: string
    expiry_date?: string
  }>
}

/**
 * 人才搜尋結果
 */
export interface TalentSearchResponseData {
  /** 人才清單 */
  talents: TalentCard[]
  /** 分頁資訊 */
  pagination: PaginationResult
}

export type TalentSearchResponse = ApiResponse<TalentSearchResponseData>

/**
 * 人才搜尋篩選器
 */
export interface TalentFilterPanel {
  /** 搜尋關鍵字 */
  keyword?: string
  /** 技能分類 */
  category?: string
  /** 技能名稱 */
  skill_name?: string
  /** 縣市 */
  city?: string
  /** 區域 */
  district?: string
  /** 評分範圍 */
  rating_range?: [number, number]
  /** 是否僅顯示開放應徵 */
  open_to_hiring_only?: boolean
  /** 排序方式 */
  sort_by?: 'rating' | 'name' | 'recent'
  /** 排序方向 */
  sort_order?: 'ASC' | 'DESC'
}

// ==================== 人才邀請相關類型 ====================

/**
 * 邀請人才請求
 */
export interface InviteTalentRequest {
  /** 教師 ID 清單 */
  teacher_ids: ID[]
  /** 邀請訊息 */
  message?: string
}

/**
 * 邀請結果
 */
export interface InviteResult {
  /** 教師 ID */
  teacher_id: ID
  /** 邀請 Token */
  token?: string
  /** 狀態 */
  status: 'PENDING' | 'SKIPPED' | 'FAILED'
  /** 失敗原因 (若有) */
  reason?: string
}

/**
 * 邀請人才回應
 */
export type InviteTalentResponse = ApiResponse<{
  /** 成功邀請數 */
  success_count: number
  /** 失敗數量 */
  failed_count: number
  /** 失敗的教師 ID 清單 */
  failed_ids: ID[]
  /** 邀請結果清單 */
  invitations: InviteResult[]
  /** 訊息 */
  message: string
}>

/**
 * 取消邀請請求
 */
export interface CancelInviteRequest {
  /** 邀請 Token */
  token: string
}

// ==================== 教師課表查詢相關類型 ====================

/**
 * 教師課表查詢請求
 */
export interface TeacherScheduleQueryRequest {
  /** 教師 ID */
  teacher_id: ID
  /** 開始日期 */
  start_date: DateString
  /** 結束日期 */
  end_date: DateString
  /** 包含個人行程 */
  include_personal_events?: boolean
}

/**
 * 教師課表項目
 */
export interface TeacherScheduleItem {
  /** 項目類型 */
  type: 'SCHEDULE_RULE' | 'PERSONAL_EVENT'
  /** ID */
  id: ID | string
  /** 標題 */
  title: string
  /** 日期 */
  date: DateString
  /** 開始時間 */
  start_time: string
  /** 結束時間 */
  end_time: string
  /** 中心名稱 (若是課堂) */
  center_name?: string
  /** 狀態 */
  status?: string
}

/**
 * 教師課表回應
 */
export type TeacherScheduleResponse = ApiResponse<TeacherScheduleItem[]>

// ==================== 快速篩選標籤相關類型 ====================

/**
 * 快速篩選標籤
 */
export interface QuickFilterTag {
  /** 標籤 ID */
  id: string
  /** 標籤名稱 */
  name: string
  /** 標籤類型 */
  type: 'skill' | 'city' | 'rating' | 'custom'
  /** 計數 */
  count?: number
  /** 是否已選中 */
  is_selected?: boolean
}

/**
 * 技能分布統計
 */
export interface SkillDistribution {
  /** 技能名稱 */
  skill_name: string
  /** 技能分類 */
  category: string
  /** 人數 */
  count: number
}

/**
 * 技能分布回應
 */
export type SkillDistributionResponse = ApiResponse<SkillDistribution[]>

// ==================== 搜尋歷史相關類型 ====================

/**
 * 搜尋歷史項目
 */
export interface RecentSearch {
  /** 搜尋關鍵字 */
  keyword: string
  /** 搜尋時間 */
  searched_at: Timestamp
  /** 結果數量 */
  result_count: number
}

/**
 * 搜尋歷史回應
 */
export type RecentSearchesResponse = ApiResponse<RecentSearch[]>

// ==================== 比較模式相關類型 ====================

/**
 * 比較模式教師資料
 */
export interface CompareModeTeacher {
  /** 教師 ID */
  id: ID
  /** 名稱 */
  name: string
  /** 頭像 */
  avatar_url?: string
  /** 評分 */
  rating: number
  /** 技能列表 */
  skills: string[]
  /** 可用日期清單 */
  available_dates: DateString[]
  /** 總媒合分數 */
  total_score: number
}

/**
 * 比較模式回應
 */
export type CompareModeResponse = ApiResponse<CompareModeTeacher[]>
