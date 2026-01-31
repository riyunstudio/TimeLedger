/**
 * 類型定義匯出中心
 *
 * 統一匯出所有模組的類型定義，建議直接從此檔案匯入
 *
 * @example
 * ```typescript
 * import type { Teacher, TeacherSkill, ApiResponse } from '~/types'
 * ```
 */

// ==================== API 通用類型 ====================
export * from './api'

// ==================== 管理員相關類型 ====================
export * from './admin'

// ==================== 教師相關類型 ====================
export * from './teacher'

// ==================== 中心相關類型 ====================
export * from './center'

// ==================== 排課相關類型 ====================
export * from './scheduling'

// ==================== 智慧媒合相關類型 ====================
export * from './matching'

// ==================== 通知相關類型 ====================
export * from './notification'

// ==================== 相容性匯出 (舊版 API Response 格式) ====================

/**
 * @deprecated 請使用 ./api.ts 中的 ApiResponse
 * 此類型僅為向後相容保留
 */
export interface LegacyApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

/**
 * @deprecated 請使用 ./scheduling.ts 中的 ScheduleCell
 */
export interface ScheduleCell {
  date: string
  time: string
  items: (ScheduleRule | PersonalEvent)[]
  has_conflict: boolean
}

/**
 * @deprecated 請使用 ./matching.ts 中的 MatchScore
 */
export interface MatchScore {
  teacher_id: number
  teacher_name: string
  match_score: number
  skill_match: number
  rating: number
  notes?: string
}

// ==================== 常用組合類型 ====================

/**
 * 教師課表項目 (用於前端顯示)
 */
export interface TeacherScheduleItem {
  id: string
  type: string
  title: string
  date: string
  start_time: string
  end_time: string
  room_id: number
  teacher_id?: number
  center_id: number
  center_name?: string
  status: string
  rule_id?: number
  data?: unknown
  is_cross_day_part?: boolean
}

/**
 * 課表移動請求
 */
export interface MoveScheduleItem {
  item_id: number
  item_type: 'SCHEDULE_RULE' | 'PERSONAL_EVENT' | 'CENTER_SESSION'
  center_id: number
  new_date: string
  new_start_time: string
  new_end_time: string
  update_mode?: 'SINGLE' | 'FUTURE' | 'ALL'
}

/**
 * 課表驗證衝突
 */
export interface ScheduleValidationConflict {
  type: 'OVERLAP' | 'TEACHER_OVERLAP' | 'ROOM_OVERLAP' | 'BUFFER'
  message: string
  current_gap_minutes?: number
  required_buffer_minutes?: number
  previous_session?: {
    id: number
    course_name: string
    end_at: string
  }
  can_override?: boolean
}

/**
 * 課表驗證結果
 */
export interface ScheduleValidation {
  valid: boolean
  conflicts: ScheduleValidationConflict[]
}
