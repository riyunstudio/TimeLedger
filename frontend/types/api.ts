/**
 * API 通用類型定義
 *
 * 提供所有 API 呼叫的基礎型別系統，包含回應格式、分頁參數、結果類型等
 */

// ==================== 基礎類型別名 ====================

/** 通用 ID 類型 */
export type ID = number

/** 時間戳記 (ISO 8601 格式) */
export type Timestamp = string

/** 日期字串 (YYYY-MM-DD 格式) */
export type DateString = string

/** 分頁排序方向 */
export type SortOrder = 'ASC' | 'DESC'

// ==================== API 回應格式 ====================

/**
 * 統一 API 回應格式
 *
 * 後端所有 API 回應皆使用此格式
 * - code: 錯誤碼，如 "SUCCESS", "SQL_ERROR" (字串類型)
 * - message: 訊息描述
 * - data: 單筆資料
 * - datas: 多筆資料 (部分 API 使用)
 */
export interface ApiResponse<T = unknown> {
  /** 錯誤碼，如 "SUCCESS", "SQL_ERROR" */
  code: string
  /** 訊息描述 */
  message: string
  /** 單筆資料 */
  data?: T
  /** 多筆資料 (部分 API 使用) */
  datas?: T
}

/**
 * 帶分頁的 API 回應
 */
export interface PaginatedApiResponse<T> extends ApiResponse<T[]> {
  pagination: PaginationResult
}

// ==================== 分頁相關類型 ====================

/**
 * 分頁查詢參數
 */
export interface PaginationParams {
  /** 頁碼 (從 1 開始) */
  page?: number
  /** 每頁筆數 (預設 20，最大 100) */
  limit?: number
  /** 排序欄位 */
  sort_by?: string
  /** 排序方向 */
  sort_order?: SortOrder
}

/**
 * 分頁結果
 */
export interface PaginationResult {
  /** 當前頁碼 */
  page: number
  /** 每頁筆數 */
  limit: number
  /** 總筆數 */
  total: number
  /** 總頁數 */
  total_pages: number
  /** 是否有下一頁 */
  has_next: boolean
  /** 是否有上一頁 */
  has_prev: boolean
}

/**
 * 帶分頁的資料回應
 */
export interface PaginatedResponse<T> {
  /** 資料陣列 */
  data: T[]
  /** 分頁資訊 */
  pagination: PaginationResult
}

// ==================== 基礎模型類型 ====================

/**
 * 基礎模型介面
 * 所有具有 id 和時間戳記的模型應繼承此介面
 */
export interface BaseModel {
  /** 主鍵 ID */
  id: ID
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 可追蹤的模型介面
 * 包含軟刪除功能的模型
 */
export interface SoftDeleteModel extends BaseModel {
  /** 刪除時間 (若有值表示已軟刪除) */
  deleted_at?: Timestamp
}

// ==================== 日期範圍類型 ====================

/**
 * 日期範圍
 */
export interface DateRange {
  /** 開始日期 (YYYY-MM-DD) */
  start_date: DateString
  /** 結束日期 (YYYY-MM-DD) */
  end_date: DateString
}

/**
 * 日期時間範圍
 */
export interface DateTimeRange {
  /** 開始時間 (ISO 8601) */
  start_at: Timestamp
  /** 結束時間 (ISO 8601) */
  end_at: Timestamp
}

// ==================== 循環規則類型 ====================

/**
 * 循環頻率類型
 */
export type RecurrenceFrequency = 'NONE' | 'DAILY' | 'WEEKLY' | 'BIWEEKLY' | 'MONTHLY'

/**
 * 循環規則
 *
 * 對應後端 RecurrenceRule 模型
 */
export interface RecurrenceRule {
  /** 循環類型 */
  type: RecurrenceFrequency
  /** 循環間隔 (如每 2 週為 interval: 2) */
  interval: number
  /** 結束日期 (可選，不設限則持續循環) */
  end_date?: DateString
  /** 循環次數上限 (可選) */
  count?: number
  /** 結束日期 ISO 字串 (可選) */
  until?: string
  /** 星期幾循環 (用於 WEEKLY, BIWEEKLY) */
  weekdays?: number[]
}

// ==================== 驗證結果類型 ====================

/**
 * 驗證衝突類型
 */
export type ValidationConflictType = 'OVERLAP' | 'TEACHER_OVERLAP' | 'ROOM_OVERLAP' | 'BUFFER'

/**
 * 驗證衝突詳細資訊
 */
export interface ValidationConflict {
  /** 衝突類型 */
  type: ValidationConflictType
  /** 衝突訊息 */
  message: string
  /** 額外詳細資訊 */
  details?: string
}

/**
 * 驗證結果
 */
export interface ValidationResult {
  /** 是否有效 */
  valid: boolean
  /** 衝突列表 */
  conflicts: ValidationConflict[]
}

// ==================== 狀態枚舉類型 ====================

/**
 * 通用狀態類型
 */
export type CommonStatus = 'ACTIVE' | 'INACTIVE' | 'PENDING' | 'APPROVED' | 'REJECTED' | 'REVOKED'

/**
 * 審核狀態類型
 */
export type ReviewStatus = 'PENDING' | 'APPROVED' | 'REJECTED' | 'REVOKED'

/**
 * 例外申請類型
 */
export type ExceptionType = 'CANCEL' | 'RESCHEDULE' | 'REPLACE_TEACHER'

/**
 * 課程類型
 */
export type CourseType = 'GROUP' | 'INDIVIDUAL' | 'WORKSHOP'

/**
 * 方案等級
 */
export type PlanLevel = 'STARTER' | 'GROWTH' | 'PRO'

/**
 * 用戶類型
 */
export type UserType = 'ADMIN' | 'TEACHER' | 'OWNER' | 'STAFF'

/**
 * 通知類型
 */
export type NotificationType = 'SCHEDULE' | 'EXCEPTION' | 'REVIEW' | 'GENERAL' | 'APPROVAL' | 'CENTER_INVITE'

/**
 * 課表項目類型
 */
export type ScheduleItemType = 'SCHEDULE_RULE' | 'PERSONAL_EVENT' | 'CENTER_SESSION'
