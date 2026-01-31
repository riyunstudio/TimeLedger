/**
 * 排課驗證 API Composable
 *
 * 提供排課驗證相關的 API 呼叫功能，
 * 包含時段重疊檢查、緩衝時間檢查和完整驗證
 */

import type { ID, DateString, ApiResponse } from '~/types/api'
import type { ValidationResult, ValidationConflict } from '~/types/api'

// ==================== 驗證請求類型 ====================

/**
 * 驗證衝突類型
 */
export type ValidationConflictType =
  | 'OVERLAP'      // 硬性衝突（時段重疊）
  | 'TEACHER_OVERLAP' // 教師時段重疊
  | 'ROOM_OVERLAP'     // 教室時段重疊
  | 'TEACHER_BUFFER'   // 教師緩衝時間不足
  | 'ROOM_BUFFER'      // 教室緩衝時間不足

/**
 * 衝突詳細資訊
 */
export interface ValidationConflictDetail {
  /** 衝突類型 */
  type: ValidationConflictType
  /** 衝突訊息 */
  message: string
  /** 當前間隔分鐘數（緩衝衝突時） */
  current_gap_minutes?: number
  /** 所需緩衝分鐘數（緩衝衝突時） */
  required_buffer_minutes?: number
  /** 上一堂課資訊 */
  previous_session?: {
    id: ID
    course_name: string
    end_at: string
  }
  /** 是否可覆寫 */
  can_override?: boolean
}

/**
 * 基礎驗證請求
 */
export interface BaseValidationRequest {
  /** 中心 ID */
  center_id: ID
  /** 教師 ID（可選，某些檢查可能不需要） */
  teacher_id?: ID
  /** 教室 ID */
  room_id: ID
  /** 日期（YYYY-MM-DD） */
  date: DateString
  /** 開始時間（HH:mm） */
  start_time: string
  /** 結束時間（HH:mm） */
  end_time: string
}

/**
 * 檢查時段重疊請求
 */
export interface CheckOverlapRequest extends BaseValidationRequest {
  /** 規則 ID（更新時傳入，避免與自身衝突） */
  rule_id?: ID
}

/**
 * 檢查教師緩衝時間請求
 */
export interface CheckTeacherBufferRequest extends BaseValidationRequest {
  /** 規則 ID（更新時傳入） */
  rule_id?: ID
}

/**
 * 檢查教室緩衝時間請求
 */
export interface CheckRoomBufferRequest extends BaseValidationRequest {
  /** 規則 ID（更新時傳入） */
  rule_id?: ID
}

/**
 * 完整驗證請求
 */
export interface FullValidationRequest extends BaseValidationRequest {
  /** 規則 ID（更新時傳入） */
  rule_id?: ID
  /** 是否覆寫緩衝衝突（僅當允許時） */
  override_buffer_conflict?: boolean
}

// ==================== 驗證回應類型 ====================

/**
 * 衝突檢查回應
 */
export interface ConflictCheckResponse {
  /** 是否有衝突 */
  has_conflict: boolean
  /** 衝突列表 */
  conflicts: ValidationConflictDetail[]
}

/**
 * 緩衝檢查回應
 */
export interface BufferCheckResponse {
  /** 是否有緩衝衝突 */
  has_buffer_conflict: boolean
  /** 衝突列表 */
  conflicts: ValidationConflictDetail[]
  /** 是否可覆寫 */
  can_override: boolean
}

/**
 * 完整驗證回應
 */
export interface ScheduleValidationResponse {
  /** 是否有效（無衝突） */
  valid: boolean
  /** 所有衝突列表 */
  conflicts: ValidationConflictDetail[]
  /** 硬衝突列表 */
  hard_conflicts: ValidationConflictDetail[]
  /** 緩衝衝突列表 */
  buffer_conflicts: ValidationConflictDetail[]
}

// ==================== Composable ====================

export const useSchedulingValidation = () => {
  const api = useApi()

  /**
   * 檢查時段重疊
   *
   * 驗證新時段是否與現有課程重疊
   *
   * @param request - 檢查重疊請求參數
   * @returns 衝突檢查回應
   */
  async function checkOverlap(request: CheckOverlapRequest): Promise<ConflictCheckResponse> {
    try {
      const response = await api.post<ConflictCheckResponse>(
        '/admin/scheduling/check-overlap',
        request
      )
      return response
    } catch (error) {
      console.error('檢查時段重疊失敗:', error)
      throw error
    }
  }

  /**
   * 檢查教師緩衝時間
   *
   * 驗證教師上一堂課與新時段的間隔是否足夠
   *
   * @param request - 檢查緩衝時間請求參數
   * @returns 緩衝檢查回應
   */
  async function checkTeacherBuffer(request: CheckTeacherBufferRequest): Promise<BufferCheckResponse> {
    try {
      const response = await api.post<BufferCheckResponse>(
        '/admin/scheduling/check-teacher-buffer',
        request
      )
      return response
    } catch (error) {
      console.error('檢查教師緩衝時間失敗:', error)
      throw error
    }
  }

  /**
   * 檢查教室緩衝時間
   *
   * 驗證教室使用間隔是否足夠
   *
   * @param request - 檢查緩衝時間請求參數
   * @returns 緩衝檢查回應
   */
  async function checkRoomBuffer(request: CheckRoomBufferRequest): Promise<BufferCheckResponse> {
    try {
      const response = await api.post<BufferCheckResponse>(
        '/admin/scheduling/check-room-buffer',
        request
      )
      return response
    } catch (error) {
      console.error('檢查教室緩衝時間失敗:', error)
      throw error
    }
  }

  /**
   * 完整驗證
   *
   * 一次檢查所有規則（重疊、教師緩衝、教室緩衝）
   *
   * @param request - 完整驗證請求參數
   * @returns 完整驗證回應
   */
  async function validateSchedule(request: FullValidationRequest): Promise<ScheduleValidationResponse> {
    try {
      const response = await api.post<ScheduleValidationResponse>(
        '/admin/scheduling/validate',
        request
      )
      return response
    } catch (error) {
      console.error('排課驗證失敗:', error)
      throw error
    }
  }

  /**
   * 快速驗證（組合多個檢查）
   *
   * 在表單提交前進行快速驗證
   *
   * @param request - 基礎驗證請求參數
   * @returns 完整驗證回應
   */
  async function quickValidate(request: BaseValidationRequest & { rule_id?: ID }): Promise<ScheduleValidationResponse> {
    // 同時發起多個驗證請求
    const [overlapResult, teacherBufferResult, roomBufferResult] = await Promise.allSettled([
      checkOverlap(request as CheckOverlapRequest),
      request.teacher_id ? checkTeacherBuffer(request as CheckTeacherBufferRequest) : null,
      checkRoomBuffer(request as CheckRoomBufferRequest),
    ])

    // 收集所有衝突
    const allConflicts: ValidationConflictDetail[] = []
    const hardConflicts: ValidationConflictDetail[] = []
    const bufferConflicts: ValidationConflictDetail[] = []

    // 處理時段重疊檢查結果
    if (overlapResult.status === 'fulfilled' && overlapResult.value.has_conflict) {
      allConflicts.push(...overlapResult.value.conflicts)
      hardConflicts.push(...overlapResult.value.conflicts.filter(
        c => c.type === 'OVERLAP' || c.type === 'TEACHER_OVERLAP' || c.type === 'ROOM_OVERLAP'
      ))
    }

    // 處理教師緩衝檢查結果
    if (teacherBufferResult?.status === 'fulfilled' && teacherBufferResult.value.has_buffer_conflict) {
      allConflicts.push(...teacherBufferResult.value.conflicts)
      bufferConflicts.push(...teacherBufferResult.value.conflicts)
    }

    // 處理教室緩衝檢查結果
    if (roomBufferResult?.status === 'fulfilled' && roomBufferResult.value.has_buffer_conflict) {
      allConflicts.push(...roomBufferResult.value.conflicts)
      bufferConflicts.push(...roomBufferResult.value.conflicts)
    }

    return {
      valid: allConflicts.length === 0,
      conflicts: allConflicts,
      hard_conflicts: hardConflicts,
      buffer_conflicts: bufferConflicts,
    }
  }

  /**
   * 格式化衝突訊息
   *
   * 將衝突物件轉換為易讀的訊息陣列
   *
   * @param conflicts - 衝突列表
   * @returns 格式化後的訊息陣列
   */
  function formatConflictMessages(conflicts: ValidationConflictDetail[]): string[] {
    return conflicts.map(conflict => {
      switch (conflict.type) {
        case 'OVERLAP':
        case 'TEACHER_OVERLAP':
        case 'ROOM_OVERLAP':
          return `${conflict.message}`

        case 'TEACHER_BUFFER':
        case 'ROOM_BUFFER':
          if (conflict.current_gap_minutes !== undefined && conflict.required_buffer_minutes !== undefined) {
            const gap = conflict.current_gap_minutes
            const required = conflict.required_buffer_minutes
            return `${conflict.message}（目前間隔 ${gap} 分鐘，需間隔 ${required} 分鐘）`
          }
          return conflict.message

        default:
          return conflict.message
      }
    })
  }

  /**
   * 判斷是否為硬衝突
   *
   * 硬衝突（時段重疊）必須處理，無法直接覆寫
   *
   * @param conflicts - 衝突列表
   * @returns 是否存在硬衝突
   */
  function hasHardConflicts(conflicts: ValidationConflictDetail[]): boolean {
    return conflicts.some(conflict =>
      conflict.type === 'OVERLAP' ||
      conflict.type === 'TEACHER_OVERLAP' ||
      conflict.type === 'ROOM_OVERLAP'
    )
  }

  /**
   * 判斷是否可覆寫
   *
   * 檢查所有衝突是否都可覆寫
   *
   * @param conflicts - 衝突列表
   * @returns 是否可覆寫
   */
  function canOverrideAll(conflicts: ValidationConflictDetail[]): boolean {
    return conflicts.length > 0 && conflicts.every(conflict => conflict.can_override)
  }

  return {
    // API 方法
    checkOverlap,
    checkTeacherBuffer,
    checkRoomBuffer,
    validateSchedule,
    quickValidate,

    // 工具函數
    formatConflictMessages,
    hasHardConflicts,
    canOverrideAll,
  }
}

// ==================== 類型匯出 ====================

export type {
  ValidationConflictDetail,
  CheckOverlapRequest,
  CheckTeacherBufferRequest,
  CheckRoomBufferRequest,
  FullValidationRequest,
  ConflictCheckResponse,
  BufferCheckResponse,
  ScheduleValidationResponse,
}
