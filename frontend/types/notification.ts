/**
 * 通知相關類型定義
 *
 * 包含通知、通知佇列等相關類型
 */

import type { ID, Timestamp, ApiResponse, PaginationParams, PaginatedResponse, NotificationType } from './api'

// ==================== 通知相關類型 ====================

/**
 * 通知
 */
export interface Notification {
  /** 通知 ID */
  id: ID
  /** 用戶 ID */
  user_id: ID
  /** 用戶類型 */
  user_type: 'ADMIN' | 'TEACHER'
  /** 中心 ID (可選) */
  center_id?: ID
  /** 中心名稱 (可選) */
  center_name?: string
  /** 通知標題 */
  title: string
  /** 通知訊息 */
  message: string
  /** 通知類型 */
  type: NotificationType
  /** 關聯資料 ID */
  related_id?: ID
  /** 關聯資料類型 */
  related_type?: string
  /** 是否已讀 */
  is_read: boolean
  /** 讀取時間 */
  read_at?: Timestamp
  /** 建立時間 */
  created_at: Timestamp
}

/**
 * 通知列表查詢參數
 */
export interface NotificationListParams extends PaginationParams {
  /** 通知類型篩選 */
  type?: NotificationType
  /** 是否僅未讀 */
  unread_only?: boolean
  /** 中心 ID */
  center_id?: ID
}

/**
 * 通知列表回應
 */
export type NotificationListResponse = PaginatedResponse<Notification>

/**
 * 未讀通知計數回應
 */
export type UnreadCountResponse = ApiResponse<{
  /** 未讀數量 */
  count: number
}>

/**
 * 標記已讀請求
 */
export interface MarkAsReadRequest {
  /** 通知 ID 清單 (若為空則全部標記已讀) */
  notification_ids?: ID[]
}

/**
 * 標記已讀回應
 */
export type MarkAsReadResponse = ApiResponse<{
  /** 標記數量 */
  count: number
}>

/**
 * 通知設定
 */
export interface NotificationSettings {
  /** 接收新例外通知 */
  receive_exception_notifications: boolean
  /** 接收審核結果通知 */
  receive_approval_notifications: boolean
  /** 接收排課變更通知 */
  receive_schedule_change_notifications: boolean
  /** 接收系統通知 */
  receive_system_notifications: boolean
}

/**
 * 通知設定回應
 */
export type NotificationSettingsResponse = ApiResponse<NotificationSettings>

/**
 * 更新通知設定請求
 */
export interface UpdateNotificationSettingsRequest extends Partial<NotificationSettings> {}

// ==================== 通知佇列相關類型 ====================

/**
 * 佇列統計
 */
export interface QueueStats {
  /** 待處理數量 */
  pending_count: number
  /** 重試中數量 */
  retry_count: number
  /** 已完成數量 */
  completed_count: number
  /** 失敗數量 */
  failed_count: number
  /** 失敗率 (%) */
  failure_rate: number
  /** Redis 是否連線 */
  redis_connected: boolean
  /** Worker 是否運行 */
  worker_running: boolean
}

/**
 * 佇列統計回應
 */
export type QueueStatsResponse = ApiResponse<QueueStats>

/**
 * 通知佇列項目
 */
export interface QueueItem {
  /** 佇列 ID */
  id: ID
  /** 通知類型 */
  type: NotificationType
  /** 目標用戶 ID */
  target_user_id: ID
  /** 目標用戶類型 */
  target_user_type: 'ADMIN' | 'TEACHER'
  /** 標題 */
  title: string
  /** 訊息 */
  message: string
  /** 關聯資料 ID */
  related_id?: ID
  /** 狀態 */
  status: 'PENDING' | 'PROCESSING' | 'COMPLETED' | 'FAILED'
  /** 重試次數 */
  retry_count: number
  /** 錯誤訊息 */
  error_message?: string
  /** 建立時間 */
  created_at: Timestamp
  /** 處理時間 */
  processed_at?: Timestamp
  /** 下次重試時間 */
  next_retry_at?: Timestamp
}

/**
 * 佇列項目列表回應
 */
export type QueueItemListResponse = PaginatedResponse<QueueItem>

/**
 * 重試佇列項目請求
 */
export interface RetryQueueItemRequest {
  /** 佇列項目 ID */
  id: ID
}

/**
 * 重試佇列項目回應
 */
export type RetryQueueItemResponse = ApiResponse<QueueItem>

/**
 * 清除已完成佇列請求
 */
export interface ClearCompletedQueueRequest {
  /** 僅清除 N 天前的資料 */
  older_than_days?: number
}

/**
 * 清除已完成佇列回應
 */
export type ClearCompletedQueueResponse = ApiResponse<{
  /** 清除數量 */
  count: number
}>

// ==================== 通知佇列監控相關類型 ====================

/**
 * 監控儀表板資料
 */
export interface QueueMonitorDashboard {
  /** 佇列統計 */
  queue_stats: QueueStats
  /** 人才庫統計 */
  talent_stats: {
    /** 總人數 */
    total: number
    /** 開放應徵人數 */
    open_hiring: number
    /** 待處理邀請 */
    pending_invites: number
  }
  /** 最近的失敗項目 */
  recent_failures: QueueItem[]
  /** 系統健康狀態 */
  system_health: {
    /** Redis 連線狀態 */
    redis: 'healthy' | 'degraded' | 'unhealthy'
    /** 資料庫連線狀態 */
    database: 'healthy' | 'degraded' | 'unhealthy'
    /** Worker 狀態 */
    worker: 'running' | 'stopped' | 'error'
  }
}

/**
 * 監控儀表板回應
 */
export type QueueMonitorDashboardResponse = ApiResponse<QueueMonitorDashboard>

// ==================== 通知佇列 Worker 相關類型 ====================

/**
 * Worker 配置
 */
export interface WorkerConfig {
  /** 是否啟用 Worker */
  enabled: boolean
  /** 工作執行緒數 */
  worker_count: number
  /** 批次處理大小 */
  batch_size: number
  /** 重試最大次數 */
  max_retries: number
  /** 重試間隔 (秒) */
  retry_interval_seconds: number
}

/**
 * Worker 狀態
 */
export interface WorkerStatus {
  /** 是否運行中 */
  is_running: boolean
  /** 已處理通知數 */
  processed_count: number
  /** 失敗通知數 */
  failed_count: number
  /** 平均處理時間 (毫秒) */
  avg_processing_time_ms: number
  /** 上次執行時間 */
  last_run_at?: Timestamp
  /** 記憶體使用量 */
  memory_usage_mb: number
}

/**
 * Worker 狀態回應
 */
export type WorkerStatusResponse = ApiResponse<WorkerStatus>

/**
 * 更新 Worker 配置請求
 */
export interface UpdateWorkerConfigRequest {
  /** 是否啟用 Worker */
  enabled?: boolean
  /** 工作執行緒數 */
  worker_count?: number
  /** 批次處理大小 */
  batch_size?: number
  /** 重試最大次數 */
  max_retries?: number
  /** 重試間隔 (秒) */
  retry_interval_seconds?: number
}

/**
 * 更新 Worker 配置回應
 */
export type UpdateWorkerConfigResponse = ApiResponse<WorkerConfig>

// ==================== LINE 通知相關類型 ====================

/**
 * LINE 通知記錄
 */
export interface LineNotificationLog {
  /** 記錄 ID */
  id: ID
  /** 通知 ID */
  notification_id: ID
  /** LINE 用戶 ID */
  line_user_id: string
  /** 狀態 */
  status: 'PENDING' | 'SENT' | 'FAILED'
  /** 錯誤碼 */
  error_code?: string
  /** 錯誤訊息 */
  error_message?: string
  /** 發送時間 */
  sent_at?: Timestamp
  /** 建立時間 */
  created_at: Timestamp
}

/**
 * LINE 通知記錄列表回應
 */
export type LineNotificationLogListResponse = PaginatedResponse<LineNotificationLog>

/**
 * 發送 LINE 通知請求
 */
export interface SendLineNotificationRequest {
  /** LINE 用戶 ID */
  line_user_id: string
  /** 通知標題 */
  title: string
  /** 通知訊息 */
  message: string
  /** 通知類型 */
  type: NotificationType
  /** 動作 URL (可選) */
  action_url?: string
}

/**
 * 發送 LINE 通知回應
 */
export type SendLineNotificationResponse = ApiResponse<{
  /** 是否成功發送 */
  success: boolean
  /** 訊息 ID (若成功) */
  message_id?: string
  /** 錯誤訊息 (若失敗) */
  error_message?: string
}>
