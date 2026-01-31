/**
 * 管理員相關類型定義
 *
 * 包含管理員用戶、認證、權限等相關類型
 */

import type { ID, Timestamp, ApiResponse, PaginationParams, PaginatedResponse, SortOrder } from './api'

// ==================== 用戶基礎類型 ====================

/**
 * 用戶基礎介面
 */
export interface User {
  /** 用戶 ID */
  id: ID
  /** 電子郵件 */
  email?: string
  /** 用戶名稱 */
  name: string
  /** 頭像 URL */
  avatar_url?: string
  /** 個人簡介 */
  bio?: string
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 公開聯絡資訊
 */
export interface PublicContactInfo {
  /** Instagram */
  instagram?: string
  /** YouTube */
  youtube?: string
  /** 個人網站 */
  website?: string
  /** 其他 */
  other?: string
}

// ==================== 管理員類型 ====================

/**
 * 管理員角色類型
 */
export type AdminRole = 'ADMIN' | 'OWNER' | 'STAFF'

/**
 * 管理員用戶
 */
export interface AdminUser extends User {
  /** 用戶類型 */
  user_type: AdminRole
  /** 角色 (如 "owner", "admin", "staff") */
  role?: string
  /** 所屬中心 ID */
  center_id?: ID
  /** LINE 用戶 ID (用於接收通知) */
  line_user_id?: string
  /** 是否已啟用 LINE 通知 */
  line_notify_enabled?: boolean
  /** LINE 綁定時間 */
  line_bound_at?: Timestamp
}

// ==================== 管理員清單回應 ====================

/**
 * 管理員清單查詢參數
 */
export interface AdminListParams extends PaginationParams {
  /** 搜尋關鍵字 */
  keyword?: string
  /** 角色篩選 */
  role?: AdminRole
  /** 狀態篩選 */
  status?: string
}

/**
 * 管理員清單項目
 */
export interface AdminListItem {
  /** 管理員 ID */
  id: ID
  /** 名稱 */
  name: string
  /** 電子郵件 */
  email: string
  /** 角色 */
  role: AdminRole
  /** 頭像 URL */
  avatar_url?: string
  /** 狀態 */
  status: 'ACTIVE' | 'INACTIVE'
  /** 建立時間 */
  created_at: Timestamp
  /** 最後登入時間 */
  last_login_at?: Timestamp
}

/**
 * 管理員清單回應
 */
export type AdminListResponse = PaginatedResponse<AdminListItem>

// ==================== 認證相關類型 ====================

/**
 * 登入請求
 */
export interface LoginRequest {
  /** 電子郵件 */
  email: string
  /** 密碼 */
  password: string
}

/**
 * 登入回應
 */
export interface AuthResponseData {
  /** JWT Token */
  token: string
  /** 刷新 Token */
  refresh_token?: string
  /** 用戶資訊 */
  user: AdminUser
}

/**
 * 登入 API 回應
 */
export type AuthResponse = ApiResponse<AuthResponseData>

/**
 * 管理員資料更新請求
 */
export interface UpdateAdminRequest {
  /** 名稱 */
  name?: string
  /** 電子郵件 */
  email?: string
  /** 頭像 URL */
  avatar_url?: string
  /** 角色 (僅 OWNER 可修改) */
  role?: AdminRole
}

// ==================== LINE 綁定相關類型 ====================

/**
 * LINE 綁定狀態
 */
export interface LineBindingStatus {
  /** 是否已綁定 */
  is_bound: boolean
  /** 綁定時間 */
  bound_at?: Timestamp
  /** 是否接收通知 */
  notify_enabled: boolean
}

/**
 * LINE 綁定驗證碼回應
 */
export interface LineBindCodeResponse {
  /** 驗證碼 */
  code: string
  /** 過期時間 */
  expires_at: Timestamp
  /** QR Code 圖片 URL */
  qr_code_url?: string
}

/**
 * LINE 通知設定
 */
export interface LineNotificationSettings {
  /** 接收新例外通知 */
  receive_exception_notifications: boolean
  /** 接收審核結果通知 */
  receive_approval_notifications: boolean
}

/**
 * LINE 綁定相關 API 回應
 */
export type LineBindingResponse = ApiResponse<LineBindingStatus>
export type LineBindCodeApiResponse = ApiResponse<LineBindCodeResponse>
export type LineNotificationSettingsResponse = ApiResponse<LineNotificationSettings>

// ==================== 中心管理員相關 ====================

/**
 * 中心管理員關聯
 */
export interface CenterAdmin {
  /** 關聯 ID */
  id: ID
  /** 中心 ID */
  center_id: ID
  /** 中心名稱 */
  center_name: string
  /** 管理員 ID */
  admin_id: ID
  /** 管理員名稱 */
  admin_name: string
  /** 角色 */
  role: AdminRole
  /** 狀態 */
  status: 'ACTIVE' | 'INACTIVE'
  /** 建立時間 */
  created_at: Timestamp
}
