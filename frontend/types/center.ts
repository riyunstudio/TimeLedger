/**
 * 中心相關類型定義
 *
 * 包含中心、課程、教室、方案等相關類型
 */

import type { ID, Timestamp, ApiResponse, PaginationParams, PaginatedResponse, PlanLevel } from './api'

// ==================== 中心基礎類型 ====================

/**
 * 中心設定
 */
export interface CenterSettings {
  /** 是否允許公開註冊 */
  allow_public_register: boolean
  /** 預設語言 */
  default_language: string
}

/**
 * 中心
 */
export interface Center {
  /** 中心 ID */
  id: ID
  /** 中心名稱 */
  name: string
  /** 方案等級 */
  plan_level: PlanLevel
  /** 中心設定 */
  settings: CenterSettings
  /** 聯絡電話 */
  phone?: string
  /** 聯絡郵箱 */
  email?: string
  /** 地址 */
  address?: string
  /** 縣市 */
  city?: string
  /** 區域 */
  district?: string
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 中心列表查詢參數
 */
export interface CenterListParams extends PaginationParams {
  /** 搜尋關鍵字 */
  keyword?: string
  /** 縣市篩選 */
  city?: string
  /** 方案等級篩選 */
  plan_level?: PlanLevel
}

/**
 * 中心列表項目
 */
export interface CenterListItem {
  /** 中心 ID */
  id: ID
  /** 中心名稱 */
  name: string
  /** 方案等級 */
  plan_level: PlanLevel
  /** 縣市 */
  city?: string
  /** 區域 */
  district?: string
  /** 教師數量 */
  teacher_count: number
  /** 建立時間 */
  created_at: Timestamp
}

/**
 * 中心列表回應
 */
export type CenterListResponse = PaginatedResponse<CenterListItem>

// ==================== 課程相關類型 ====================

/**
 * 課程
 */
export interface Course {
  /** 課程 ID */
  id: ID
  /** 所屬中心 ID */
  center_id: ID
  /** 課程名稱 */
  name: string
  /** 課程說明 */
  description?: string
  /** 教師緩衝時間 (分鐘) */
  teacher_buffer_min: number
  /** 教室緩衝時間 (分鐘) */
  room_buffer_min: number
  /** 課程類型 */
  course_type?: string
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 課程列表查詢參數
 */
export interface CourseListParams extends PaginationParams {
  /** 搜尋關鍵字 */
  keyword?: string
}

/**
 * 課程列表回應
 */
export type CourseListResponse = PaginatedResponse<Course>

/**
 * 新增課程請求
 */
export interface CreateCourseRequest {
  /** 課程名稱 */
  name: string
  /** 課程說明 */
  description?: string
  /** 教師緩衝時間 (分鐘) */
  teacher_buffer_min?: number
  /** 教室緩衝時間 (分鐘) */
  room_buffer_min?: number
  /** 課程類型 */
  course_type?: string
}

/**
 * 更新課程請求
 */
export interface UpdateCourseRequest extends Partial<CreateCourseRequest> {}

// ==================== 開課相關類型 ====================

/**
 * 開課 (Offering)
 */
export interface Offering {
  /** 開課 ID */
  id: ID
  /** 所屬中心 ID */
  center_id: ID
  /** 關聯課程 ID */
  course_id: ID
  /** 課程名稱 (關聯取得) */
  course_name?: string
  /** 課程時長（分鐘，從關聯課程取得） */
  course_duration?: number
  /** 預設教室 ID */
  default_room_id?: ID
  /** 預設教師 ID */
  default_teacher_id?: ID
  /** 預設開始時間 (HH:MM) */
  default_start_time?: string
  /** 預設結束時間 (HH:MM) */
  default_end_time?: string
  /** 是否允許緩衝覆寫 */
  allow_buffer_override: boolean
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 開課列表回應
 */
export type OfferingListResponse = PaginatedResponse<Offering>

/**
 * 新增開課請求
 */
export interface CreateOfferingRequest {
  /** 課程 ID */
  course_id: ID
  /** 預設教室 ID */
  default_room_id?: ID
  /** 預設教師 ID */
  default_teacher_id?: ID
  /** 是否允許緩衝覆寫 */
  allow_buffer_override?: boolean
}

// ==================== 教室相關類型 ====================

/**
 * 教室
 */
export interface Room {
  /** 教室 ID */
  id: ID
  /** 所屬中心 ID */
  center_id: ID
  /** 教室名稱 */
  name: string
  /** 容納人數 */
  capacity: number
  /** 設備說明 */
  equipment?: string
  /** 清潔時間 (分鐘) */
  cleaning_time?: number
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 教室列表回應
 */
export type RoomListResponse = PaginatedResponse<Room>

/**
 * 新增教室請求
 */
export interface CreateRoomRequest {
  /** 教室名稱 */
  name: string
  /** 容納人數 */
  capacity: number
  /** 設備說明 */
  equipment?: string
  /** 清潔時間 (分鐘) */
  cleaning_time?: number
}

/**
 * 更新教室請求
 */
export interface UpdateRoomRequest extends Partial<CreateRoomRequest> {}

// ==================== 中心會員相關類型 ====================

/**
 * 中心會員關聯
 */
export interface CenterMembership {
  /** 關聯 ID */
  id: ID
  /** 中心 ID */
  center_id: ID
  /** 中心名稱 */
  center_name?: string
  /** 教師 ID */
  teacher_id: ID
  /** 角色 */
  role?: string
  /** 狀態 */
  status: 'ACTIVE' | 'INACTIVE' | 'INVITED'
  /** 加入時間 */
  joined_at?: Timestamp
}

/**
 * 中心會員列表回應
 */
export type CenterMembershipListResponse = PaginatedResponse<CenterMembership>

// ==================== 中心設定相關類型 ====================

/**
 * 中心設定更新請求
 */
export interface UpdateCenterSettingsRequest {
  /** 中心名稱 */
  name?: string
  /** 聯絡電話 */
  phone?: string
  /** 聯絡郵箱 */
  email?: string
  /** 地址 */
  address?: string
  /** 縣市 */
  city?: string
  /** 區域 */
  district?: string
  /** 設定 */
  settings?: Partial<CenterSettings>
}

// ==================== 中心統計相關類型 ====================

/**
 * 中心統計資訊
 */
export interface CenterStats {
  /** 中心 ID */
  center_id: ID
  /** 教師數量 */
  teacher_count: number
  /** 課程數量 */
  course_count: number
  /** 開課數量 */
  offering_count: number
  /** 教室數量 */
  room_count: number
  /** 本月課堂數 */
  monthly_session_count: number
  /** 待審核例外數 */
  pending_exception_count: number
}

// ==================== 國定假日相關類型 ====================

/**
 * 國定假日
 */
export interface Holiday {
  /** 假日 ID */
  id: ID
  /** 所屬中心 ID */
  center_id: ID
  /** 假日名稱 */
  name: string
  /** 假日日期 */
  date: string
  /** 是否為補班日 */
  is_makeup_day: boolean
  /** 建立時間 */
  created_at: Timestamp
}

/**
 * 國定假日列表回應
 */
export type HolidayListResponse = PaginatedResponse<Holiday>

/**
 * 新增國定假日請求
 */
export interface CreateHolidayRequest {
  /** 假日名稱 */
  name: string
  /** 假日日期 (YYYY-MM-DD) */
  date: string
  /** 是否為補班日 */
  is_makeup_day?: boolean
}

/**
 * 批量新增國定假日請求
 */
export interface BulkCreateHolidaysRequest {
  /** 假日清單 */
  holidays: CreateHolidayRequest[]
}

// ==================== 資源相關類型 ====================

/**
 * 資源類型
 */
export type ResourceType = 'DOCUMENT' | 'VIDEO' | 'IMAGE' | 'AUDIO' | 'LINK'

/**
 * 資源
 */
export interface Resource {
  /** 資源 ID */
  id: ID
  /** 所屬中心 ID */
  center_id: ID
  /** 資源名稱 */
  name: string
  /** 資源類型 */
  resource_type: ResourceType
  /** 資源 URL */
  url: string
  /** 說明 */
  description?: string
  /** 建立者 ID */
  created_by: ID
  /** 建立者名稱 */
  created_by_name?: string
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 資源列表回應
 */
export type ResourceListResponse = PaginatedResponse<Resource>

/**
 * 新增資源請求
 */
export interface CreateResourceRequest {
  /** 資源名稱 */
  name: string
  /** 資源類型 */
  resource_type: ResourceType
  /** 資源 URL */
  url: string
  /** 說明 */
  description?: string
}

// ==================== 範本相關類型 ====================

/**
 * 課表範本
 */
export interface ScheduleTemplate {
  /** 範本 ID */
  id: ID
  /** 所屬中心 ID */
  center_id: ID
  /** 範本名稱 */
  name: string
  /** 範本說明 */
  description?: string
  /** 範本內容 (JSON) */
  content: Record<string, unknown>
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 課表範本列表回應
 */
export type TemplateListResponse = PaginatedResponse<ScheduleTemplate>

/**
 * 新增範本請求
 */
export interface CreateTemplateRequest {
  /** 範本名稱 */
  name: string
  /** 範本說明 */
  description?: string
  /** 範本內容 */
  content: Record<string, unknown>
}
