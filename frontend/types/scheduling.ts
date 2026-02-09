/**
 * 排課相關類型定義
 *
 * 包含課表規則、例、外、個人行程、課堂筆記等相關類型
 */

import type {
  ID,
  Timestamp,
  DateString,
  ApiResponse,
  PaginationParams,
  PaginatedResponse,
  DateRange,
  RecurrenceRule,
  ValidationResult,
  ScheduleItemType,
  ExceptionType,
  ReviewStatus,
} from './api'

// ==================== 課表規則相關類型 ====================

/**
 * 課表規則 (ScheduleRule)
 */
export interface ScheduleRule {
  /** 規則 ID */
  id: ID
  /** 所屬中心 ID */
  center_id: ID
  /** 關聯開課 ID */
  offering_id: ID
  /** 開課名稱 */
  offering_name?: string
  /** 關聯課程 ID */
  course_id?: ID
  /** 課程名稱 */
  course_name?: string
  /** 指定教師 ID (可選) */
  teacher_id?: ID
  /** 教師名稱 */
  teacher_name?: string
  /** 教室 ID */
  room_id: ID
  /** 教室名稱 */
  room_name?: string
  /** 星期幾 (0-6, 0 為週日) */
  weekday: number
  /** 開始時間 (HH:mm) */
  start_time: string
  /** 結束時間 (HH:mm) */
  end_time: string
  /** 有效範圍 */
  effective_range: DateRange
  /** 循環規則 */
  recurrence_rule?: RecurrenceRule
  /** 是否已鎖定 */
  is_locked?: boolean
  /** 鎖定原因 */
  lock_reason?: string
  /** 例外列表 */
  exceptions?: ScheduleException[]
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 課表規則列表查詢參數
 */
export interface ScheduleRuleListParams extends PaginationParams {
  /** 中心 ID */
  center_id?: ID
  /** 教師 ID */
  teacher_id?: ID
  /** 教室 ID */
  room_id?: ID
  /** 開課 ID */
  offering_id?: ID
  /** 星期幾 */
  weekday?: number
  /** 開始日期 */
  start_date?: DateString
  /** 結束日期 */
  end_date?: DateString
  /** 搜尋關鍵字 */
  keyword?: string
}

/**
 * 課表規則列表回應
 */
export type ScheduleRuleListResponse = PaginatedResponse<ScheduleRule>

/**
 * 新增課表規則請求
 */
export interface CreateScheduleRuleRequest {
  /** 開課 ID */
  offering_id: ID
  /** 指定教師 ID (可選) */
  teacher_id?: ID
  /** 教室 ID */
  room_id: ID
  /** 星期幾 (0-6) */
  weekday: number
  /** 開始時間 (HH:mm) */
  start_time: string
  /** 結束時間 (HH:mm) */
  end_time: string
  /** 有效開始日期 */
  effective_start_date: DateString
  /** 有效結束日期 */
  effective_end_date: DateString
  /** 循環規則 (可選) */
  recurrence_rule?: RecurrenceRule
}

/**
 * 更新課表規則請求
 */
export interface UpdateScheduleRuleRequest {
  /** 指定教師 ID (可選) */
  teacher_id?: ID | null
  /** 教室 ID */
  room_id?: ID
  /** 開始時間 (HH:mm) */
  start_time?: string
  /** 結束時間 (HH:mm) */
  end_time?: string
  /** 有效開始日期 */
  effective_start_date?: DateString
  /** 有效結束日期 */
  effective_end_date?: DateString
}

// ==================== 例外申請相關類型 ====================

/**
 * 例外申請 (ScheduleException)
 */
export interface ScheduleException {
  /** 例外 ID */
  id: ID
  /** 所屬中心 ID */
  center_id: ID
  /** 關聯規則 ID */
  rule_id: ID
  /** 教師 ID */
  teacher_id: ID
  /** 教師名稱 */
  teacher_name?: string
  /** 原始日期 */
  original_date: DateString
  /** 例外類型 */
  type: ExceptionType
  /** 狀態 */
  status: ReviewStatus
  /** 新開始時間 (可選) */
  new_start_at?: Timestamp
  /** 新結束時間 (可選) */
  new_end_at?: Timestamp
  /** 替換教師 ID (可選) */
  new_teacher_id?: ID
  /** 替換教師名稱 (可選) */
  new_teacher_name?: string
  /** 原因 */
  reason: string
  /** 審核者 ID */
  reviewed_by?: ID
  /** 審核者名稱 */
  reviewed_by_name?: string
  /** 審核時間 */
  reviewed_at?: Timestamp
  /** 審核意見 */
  review_note?: string
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 例外申請列表查詢參數
 */
export interface ExceptionListParams extends PaginationParams {
  /** 中心 ID */
  center_id?: ID
  /** 教師 ID */
  teacher_id?: ID
  /** 狀態篩選 */
  status?: ReviewStatus
  /** 例外類型篩選 */
  type?: ExceptionType
  /** 開始日期 */
  start_date?: DateString
  /** 結束日期 */
  end_date?: DateString
}

/**
 * 例外申請列表回應
 */
export type ExceptionListResponse = PaginatedResponse<ScheduleException>

/**
 * 新增例外申請請求
 */
export interface CreateExceptionRequest {
  /** 中心 ID */
  center_id: ID
  /** 規則 ID */
  rule_id: ID
  /** 原始日期 */
  original_date: DateString
  /** 例外類型 */
  type: ExceptionType
  /** 新開始時間 (可選) */
  new_start_at?: string
  /** 新結束時間 (可選) */
  new_end_at?: string
  /** 替換教師 ID (可選) */
  new_teacher_id?: ID
  /** 原因 */
  reason: string
}

/**
 * 審核例外請求
 */
export interface ReviewExceptionRequest {
  /** 審核動作 */
  action: 'APPROVE' | 'REJECT'
  /** 審核意見 */
  note?: string
  /** 是否覆寫緩衝衝突 (僅核准時有效) */
  override_buffer_conflict?: boolean
}

// ==================== 個人行程相關類型 ====================

/**
 * 個人行程 (PersonalEvent)
 */
export interface PersonalEvent {
  /** 行程 ID (展開後為 string: "originalId_date") */
  id: ID | string
  /** 原始 ID (用於 API 呼叫) */
  originalId?: ID
  /** 所屬教師 ID */
  teacher_id: ID
  /** 行程標題 */
  title: string
  /** 開始時間 */
  start_at: Timestamp
  /** 結束時間 */
  end_at: Timestamp
  /** 是否為全天行程 */
  is_all_day?: boolean
  /** 循環規則 */
  recurrence_rule?: RecurrenceRule
  /** 顏色 (十六進位) */
  color: string
  /** 備註 */
  notes?: string
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 個人行程列表回應
 */
export type PersonalEventListResponse = ApiResponse<PersonalEvent[]>

/**
 * 新增個人行程請求
 */
export interface CreatePersonalEventRequest {
  /** 行程標題 */
  title: string
  /** 開始時間 */
  start_at: string
  /** 結束時間 */
  end_at: string
  /** 是否為全天行程 */
  is_all_day?: boolean
  /** 顏色 */
  color_hex?: string
  /** 備註 */
  notes?: string
  /** 循環規則 */
  recurrence_rule?: RecurrenceRule
}

/**
 * 更新個人行程請求
 */
export interface UpdatePersonalEventRequest extends Partial<CreatePersonalEventRequest> {
  /** 更新模式 */
  update_mode?: 'SINGLE' | 'FUTURE' | 'ALL'
}

// ==================== 課堂筆記相關類型 ====================

/**
 * 課堂筆記 (SessionNote)
 */
export interface SessionNote {
  /** 筆記 ID */
  id: ID
  /** 所屬中心 ID */
  center_id: ID
  /** 關聯規則 ID */
  rule_id: ID
  /** 課程日期 */
  session_date: DateString
  /** 課堂筆記內容 */
  content: string
  /** 課前準備 */
  prep_note: string
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 課堂筆記回應
 */
export interface SessionNoteResponse {
  /** 筆記內容 */
  note: SessionNote
  /** 是否為新筆記 */
  is_new: boolean
}

/**
 * 課堂筆記 API 回應
 */
export type SessionNoteApiResponse = ApiResponse<SessionNoteResponse>

/**
 * 新增/更新課堂筆記請求
 */
export interface SaveSessionNoteRequest {
  /** 規則 ID */
  rule_id: ID
  /** 課程日期 */
  session_date: DateString
  /** 課堂筆記內容 */
  content: string
  /** 課前準備 */
  prep_note: string
}

// ==================== 課表檢視相關類型 ====================

/**
 * 課表項目
 */
export interface ScheduleItem {
  /** 項目類型 */
  type: ScheduleItemType
  /** 項目 ID */
  id: ID | string
  /** 標題 */
  title: string
  /** 開始時間 */
  start_time: string
  /** 結束時間 */
  end_time: string
  /** 日期 */
  date: string
  /** 顏色 */
  color?: string
  /** 狀態 */
  status?: string
  /** 中心名稱 */
  center_name?: string
  /** 原始資料 */
  data?: ScheduleRule | PersonalEvent
  /** 教室 ID */
  room_id?: ID
  /** 教師 ID */
  teacher_id?: ID
  /** 中心 ID */
  center_id?: ID
  /** 規則 ID (用於關聯課堂筆記) */
  rule_id?: ID
  /** 是否為跨日課程的一部分 */
  is_cross_day_part?: boolean
}

/**
 * 每日課表
 */
export interface DaySchedule {
  /** 日期 */
  date: DateString
  /** 星期幾 (0-6) */
  day_of_week: number
  /** 課表項目列表 */
  items: ScheduleItem[]
}

/**
 * 每週課表
 */
export interface WeekSchedule {
  /** 週開始日期 */
  week_start: DateString
  /** 週結束日期 */
  week_end: DateString
  /** 每日課表 */
  days: DaySchedule[]
}

/**
 * 每日課表摘要
 */
export interface TodaySummary {
  /** 日期 */
  date: DateString
  /** 星期幾 */
  day_of_week: number
  /** 今日課程數 */
  session_count: number
  /** 今日總時數 */
  total_hours: number
  /** 第一堂課時間 */
  first_session_time?: string
  /** 最後一堂課時間 */
  last_session_time?: string
}

// ==================== 衝突檢測相關類型 ====================

/**
 * 衝突類型
 */
export type ConflictType = 'OVERLAP' | 'TEACHER_OVERLAP' | 'ROOM_OVERLAP' | 'TEACHER_BUFFER' | 'ROOM_BUFFER'

/**
 * 衝突詳細資訊
 */
export interface ConflictDetail {
  /** 衝突類型 */
  type: ConflictType
  /** 衝突訊息 */
  message: string
  /** 當前間隔分鐘數 */
  current_gap_minutes?: number
  /** 所需緩衝分鐘數 */
  required_buffer_minutes?: number
  /** 上一堂課資訊 */
  previous_session?: {
    id: ID
    course_name: string
    end_at: Timestamp
  }
  /** 是否可覆寫 */
  can_override?: boolean
}

/**
 * 課表驗證回應
 */
export type ScheduleValidationResponse = ApiResponse<ValidationResult>

/**
 * 衝突檢測請求
 */
export interface CheckConflictRequest {
  /** 中心 ID */
  center_id: ID
  /** 規則 ID (更新時傳入) */
  rule_id?: ID
  /** 教師 ID */
  teacher_id?: ID
  /** 教室 ID */
  room_id: ID
  /** 日期 */
  date: DateString
  /** 開始時間 */
  start_time: string
  /** 結束時間 */
  end_time: string
  /** 是否覆寫緩衝衝突 */
  override_buffer_conflict?: boolean
}

// ==================== 循環編輯相關類型 ====================

/**
 * 循環編輯模式
 */
export type RecurrenceEditMode = 'SINGLE' | 'FUTURE' | 'ALL'

/**
 * 循環編輯影響預覽
 */
export interface RecurrenceEditImpact {
  /** 受影響的課程數量 */
  affected_count: number
  /** 受影響的日期清單 */
  affected_dates: DateString[]
  /** 將被刪除的例外數量 */
  cancelled_exceptions: number
  /** 將被更新的例外數量 */
  updated_exceptions: number
}

/**
 * 循環編輯請求
 */
export interface EditRecurringRequest {
  /** 規則 ID */
  rule_id: ID
  /** 編輯日期 */
  edit_date: DateString
  /** 編輯模式 */
  mode: RecurrenceEditMode
  /** 新的開始時間 (可選) */
  new_start_time?: string
  /** 新的結束時間 (可選) */
  new_end_time?: string
  /** 新的教室 ID (可選) */
  new_room_id?: ID
  /** 新的教師 ID (可選) */
  new_teacher_id?: ID
}

// ==================== 課表移動相關類型 ====================

/**
 * 移動課表項目請求
 */
export interface MoveScheduleItemRequest {
  /** 項目 ID */
  item_id: ID
  /** 項目類型 */
  item_type: ScheduleItemType
  /** 中心 ID */
  center_id: ID
  /** 新日期 */
  new_date: DateString
  /** 新開始時間 */
  new_start_time: string
  /** 新結束時間 */
  new_end_time: string
  /** 更新模式 */
  update_mode?: RecurrenceEditMode
}

/**
 * 課表移動回應
 */
export type MoveScheduleItemResponse = ApiResponse<{
  /** 移動後的課表規則 */
  rule?: ScheduleRule
  /** 移動後的個人行程 */
  event?: PersonalEvent
}>

// ==================== 資源佔用相關類型 ====================

/**
 * 學期期間 (Term)
 */
export interface Term {
  /** 學期 ID */
  id: ID
  /** 所屬中心 ID */
  center_id: ID
  /** 學期名稱 */
  name: string
  /** 開始日期 (YYYY-MM-DD) */
  start_date: DateString
  /** 結束日期 (YYYY-MM-DD) */
  end_date: DateString
  /** 建立時間 */
  created_at?: Timestamp
  /** 更新時間 */
  updated_at?: Timestamp
}

/**
 * 佔用規則查詢參數
 */
export interface OccupancyRulesParams {
  /** 學期 ID */
  term_id: ID
  /** 老師 ID (可選) */
  teacher_id?: ID
  /** 教室 ID (可選) */
  room_id?: ID
}

/**
 * 佔用規則 (用於週曆視圖)
 *
 * 此類型用於資源佔用表頁面，簡化了 ScheduleRule 以提高效能
 */
export interface OccupancyRule {
  /** 規則 ID */
  id: ID
  /** 關聯開課 ID */
  offering_id: ID
  /** 開課名稱 */
  offering_name?: string
  /** 關聯課程 ID */
  course_id?: ID
  /** 課程名稱 */
  course_name?: string
  /** 指定教師 ID (可選) */
  teacher_id?: ID
  /** 教師名稱 */
  teacher_name?: string
  /** 教室 ID */
  room_id: ID
  /** 教室名稱 */
  room_name?: string
  /** 星期幾 (0-6, 0 為週日) */
  weekday: number
  /** 開始時間 (HH:mm) */
  start_time: string
  /** 結束時間 (HH:mm) */
  end_time: string
  /** 有效範圍 */
  effective_range: DateRange
}

/**
 * 佔用規則列表回應
 */
export type OccupancyRulesResponse = ApiResponse<OccupancyRule[]>

/**
 * 規則複製請求
 */
export interface CopyRulesRequest {
  /** 來源學期 ID */
  source_term_id: ID
  /** 目標學期 ID */
  target_term_id: ID
  /** 要複製的規則 ID 列表 */
  rule_ids: ID[]
}

/**
 * 規則複製回應
 */
export type CopyRulesResponse = ApiResponse<{
  /** 成功複製的規則數量 */
  copied_count: number
  /** 失敗的規則 ID 列表 */
  failed_ids?: ID[]
  /** 訊息 */
  message?: string
}>

// ==================== 矩陣視圖相關類型 (Phase 2) ====================

/**
 * 矩陣視圖日期範圍
 */
export interface MatrixDateRange {
  /** 開始日期 (YYYY-MM-DD) */
  start_date: string
  /** 結束日期 (YYYY-MM-DD) */
  end_date: string
}

/**
 * 矩陣視圖項目（課程場次）
 *
 * 包含後端已計算好的 CSS 定位屬性
 */
export interface MatrixItem {
  /** 項目 ID（展開後的課程場次 ID） */
  id: number
  /** 原始規則 ID */
  rule_id: number
  /** 課程名稱 */
  title: string
  /** 日期 (YYYY-MM-DD) */
  date: string
  /** 開始時間 (HH:mm) */
  start_time: string
  /** 結束時間 (HH:mm) */
  end_time: string
  /** 開始小時 (用於定位) */
  start_hour: number
  /** 開始分鐘 */
  start_minute: number
  /** 持續分鐘數 */
  duration: number
  /** CSS top 百分比 (0-100)，後端已計算 */
  top_offset: number
  /** CSS height 百分比 (相對於時段)，後端已計算 */
  height_percent: number
  /** 開課 ID */
  offering_id: number
  /** 開課名稱 */
  offering_name: string
  /** 教師 ID */
  teacher_id?: number
  /** 教師名稱 */
  teacher_name: string
  /** 教室 ID */
  room_id: number
  /** 教室名稱 */
  room_name: string
  /** 是否為假日 */
  is_holiday: boolean
  /** 是否有例外 */
  has_exception: boolean
  /** 例外類型 */
  exception_type?: string
  /** 是否為停課 */
  is_suspended: boolean
  /** 課程顏色 */
  color?: string
}

/**
 * 矩陣視圖資源（老師或教室）
 */
export interface MatrixResource {
  /** 資源 ID */
  id: number
  /** 資源名稱 */
  name: string
  /** 資源類型：'teacher' | 'room' */
  type: 'teacher' | 'room'
  /** 該資源的課程項目 */
  items: MatrixItem[]
}

/**
 * 矩陣視圖回應（BFF 模式）
 *
 * 後端直接回傳前端可直接渲染的矩陣結構
 */
export interface MatrixViewResponse {
  /** 橫軸時段，如 [9, 10, 11, ...] */
  time_slots: number[]
  /** 縱軸資源（老師或教室） */
  resources: MatrixResource[]
  /** 日期範圍 */
  date_range: MatrixDateRange
}

/**
 * 矩陣視圖查詢參數
 */
export interface MatrixViewParams {
  /** 開始日期 (YYYY-MM-DD) */
  start_date: string
  /** 結束日期 (YYYY-MM-DD) */
  end_date: string
  /** 查詢類型：'teacher' | 'room' | 'all' */
  type: 'teacher' | 'room' | 'all'
  /** 是否包含停課 (預設 true) */
  include_suspended?: boolean
  /** 指定資源 ID（逗號分隔，可選） */
  resource_ids?: string
}

/**
 * 矩陣視圖 API 回應類型
 */
export type MatrixViewApiResponse = ApiResponse<MatrixViewResponse>

/**
 * 矩陣視圖週曆卡片（用於前端渲染）
 */
export interface MatrixScheduleCard {
  /** 唯一鍵值 */
  key: string
  /** 矩陣項目資料 */
  item: MatrixItem
  /** CSS 樣式（直接使用後端計算的值） */
  style: {
    top: string
    height: string
    left?: string
    width?: string
  }
  /** 是否為重疊課程 */
  isOverlapping: boolean
  /** 重疊數量 */
  overlapCount: number
  /** 是否為重疊中的第一個 */
  isFirstInOverlap: boolean
}

/**
 * 矩陣視圖週曆格子
 */
export interface MatrixWeekCell {
  /** 星期幾 (1-7，週一到週日) */
  weekday: number
  /** 日期 (YYYY-MM-DD) */
  date: string
  /** 該格子的課程列表 */
  items: MatrixScheduleCard[]
  /** 是否為假日 */
  isHoliday: boolean
}
