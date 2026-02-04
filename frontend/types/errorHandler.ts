/**
 * 錯誤處理類型定義
 *
 * 集中管理錯誤處理相關的類型定義
 * 避免在 composables/ 目錄中重複匯出導致 Nuxt 警告
 */

import type { ApiResponse } from '~/types/api'

/**
 * 錯誤處理的額外選項
 */
export interface ErrorHandlerOptions {
  /** 錯誤標題（可自訂） */
  title?: string
  /** 是否顯示為 Toast（預設為 Alert） */
  asToast?: boolean
  /** Toast 類型（asToast 為 true 時有效） */
  toastType?: 'success' | 'error' | 'warning' | 'info'
  /** 額外的動作按鈕 */
  actions?: ErrorAction[]
  /** 錯誤發生後的回調 */
  onError?: (error: ApiResponse) => void
  /** 成功時的回調 */
  onSuccess?: (data: any) => void
}

/**
 * 錯誤動作按鈕
 */
export interface ErrorAction {
  label: string
  style?: 'primary' | 'warning' | 'critical' | 'secondary'
  action: () => void
}
