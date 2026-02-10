/**
 * 錯誤處理工具
 *
 * 提供統一的錯誤處理機制，整合錯誤碼對照表與現有的 alert/toast 系統
 * @deprecated 請使用 composables/useErrorHandler.ts
 */

import { ref } from 'vue'
import {
  ERROR_MESSAGES,
  isSuccessCode,
  isPermissionError,
  isUnauthorizedError,
  isValidationError,
  type ErrorCode,
  NUMERIC_ERROR_CODE_MAP,
} from '~/constants/errorCodes'
import type { ApiResponse } from '~/types/api'
import type { ErrorHandlerOptions } from '~/types/errorHandler'

// ==================== 錯誤類型定義 ====================

/**
 * API 錯誤介面
 */
export interface ApiError {
  /** 錯誤碼 */
  code: string
  /** 錯誤訊息 */
  message: string
  /** 額外資料 */
  data?: unknown
}

/**
 * 網路錯誤介面
 */
export interface NetworkError {
  /** HTTP 狀態碼 */
  status?: number
  /** 狀態文字 */
  statusText?: string
  /** 錯誤訊息 */
  message: string
  /** 原始錯誤 */
  originalError?: Error
}

// ==================== 全域錯誤狀態 ====================

const errorToastRef = ref<{
  show: (message: string, duration?: number) => void
} | null>(null)

/**
 * 註冊錯誤提示組件
 */
export function registerErrorToast(component: { show: (message: string, duration?: number) => void }) {
  errorToastRef.value = component
}

// ==================== 輔助函數 ====================

/**
 * 根據錯誤碼取得使用者友善訊息
 * 優先級順序：
 * 1. ERROR_MESSAGES[映射後的字串錯誤碼]（優先使用本地化中文訊息）
 * 2. 後端原始 error.message（僅當沒有對應的本地化訊息時使用）
 * 3. 預設錯誤訊息
 *
 * 注意：後端可能返回英文訊息，前端應該始終顯示本地化的中文訊息
 * 以確保使用者看到一致且在地化的錯誤說明
 */
function getErrorMessage(error: ApiError | NetworkError): string {
  if ('code' in error && error.code) {
    const code = error.code
    const originalCode = String(code)

    // 嘗試映射數字錯誤碼到字串錯誤碼
    let mappedCode: string | undefined
    if (typeof code === 'number' && NUMERIC_ERROR_CODE_MAP[code]) {
      mappedCode = NUMERIC_ERROR_CODE_MAP[code]
    } else if (typeof code === 'string') {
      mappedCode = NUMERIC_ERROR_CODE_MAP[code] || code
    } else {
      mappedCode = originalCode
    }

    // 優先使用本地化的中文訊息（從 ERROR_MESSAGES）
    // 這確保使用者始終看到中文的錯誤說明，而非後端的英文訊息
    if (mappedCode && ERROR_MESSAGES[mappedCode]) {
      return ERROR_MESSAGES[mappedCode]
    }

    // 使用原始錯誤碼查找
    if (ERROR_MESSAGES[originalCode]) {
      return ERROR_MESSAGES[originalCode]
    }

    // 沒有對應的本地化訊息時，才使用後端提供的訊息
    if (error.message) {
      return error.message
    }

    return '發生未知錯誤'
  }
  return error.message || '發生未知錯誤'
}

/**
 * 取得 HTTP 狀態碼對應的處理方式
 */
function getDefaultHandler(statusCode: number) {
  // 預設處理器包裝，將 options 傳入閉包的 options
  const makeHandler = (handler: (error: ApiError | NetworkError) => void): ErrorHandlerFn => {
    return (error: ApiError | NetworkError, options: ErrorHandlerOptions) => {
      handler(error)
    }
  }

  const handlers: Record<number, ErrorHandlerFn> = {
    401: makeHandler((error) => {
      redirectToLogin()
      showErrorAlert('請先登入再進行操作', '未授權')
    }),
    403: makeHandler((error) => {
      showErrorAlert('您沒有權限執行此操作', '禁止存取')
    }),
    404: makeHandler((error) => {
      showErrorAlert('找不到請求的資源', '404')
    }),
    409: makeHandler((error) => {
      // 衝突錯誤，根據錯誤訊息提供更詳細的描述
      let message = '操作與現有資料衝突，請檢查後重試'
      const errorMessage = 'message' in error ? error.message : ''

      // 檢測是否為課表相關的衝突
      const scheduleRelatedKeywords = ['時段', '課程', '重疊', 'schedule', 'session', 'room', 'teacher', '教師', '教室']
      const isScheduleRelated = scheduleRelatedKeywords.some(keyword =>
        errorMessage.toLowerCase().includes(keyword.toLowerCase())
      )

      if (isScheduleRelated) {
        message = '時段與現有課程重疊，請檢查排行程'
      } else if (errorMessage.includes('email') || errorMessage.includes('信箱')) {
        message = '此電子郵件已被註冊，請使用其他信箱'
      } else if (errorMessage.includes('已存在') || errorMessage.includes('already exists')) {
        message = '資料已存在，請勿重複建立'
      }

      showErrorAlert(message, '衝突錯誤')
    }),
    422: makeHandler((error) => {
      showErrorAlert('輸入資料驗證失敗，請檢查後重試', '驗證錯誤')
    }),
    429: makeHandler((error) => {
      showErrorAlert('請求過於頻繁，請稍後再試', '速率限制')
    }),
    500: makeHandler((error) => {
      showErrorAlert('系統錯誤，請稍後再試', '伺服器錯誤')
    }),
  }
  return handlers[statusCode] || defaultErrorHandler
}

/**
 * 預設錯誤處理函數
 */
function defaultErrorHandler(error: ApiError | NetworkError, options: ErrorHandlerOptions) {
  const message = getErrorMessage(error)
  const title = options.title || getErrorTitle(error)

  if (options.showAlert !== false) {
    showErrorAlert(message, title)
  }
}

/**
 * 取得錯誤標題
 */
function getErrorTitle(error: ApiError | NetworkError): string {
  if ('code' in error && error.code) {
    if (isPermissionError(error.code)) {
      return '權限錯誤'
    }
    if (isValidationError(error.code)) {
      return '驗證錯誤'
    }
    if (isUnauthorizedError(error.code)) {
      return '未授權'
    }
  }
  return '操作失敗'
}

/**
 * 顯示錯誤提示
 */
function showErrorAlert(message: string, title?: string) {
  // 使用全域 alert 系統
  if (typeof window !== 'undefined' && (window as any).$alert) {
    ; (window as any).$alert({
      message,
      title: title || '操作失敗',
      type: 'error',
    })
  } else {
    // Fallback 到原生 alert
    alert(`${title || '錯誤'}: ${message}`)
  }
}

/**
 * 導向登入頁（401 錯誤處理專用）
 */
function redirectToLogin() {
  if (typeof window !== 'undefined') {
    // 清除所有登入資訊
    clearAllAuthData()

    // 根據目前路徑判斷登入頁面
    const currentPath = window.location.pathname
    let loginPath = '/teacher/login' // 預設老師登入頁

    // 判斷應該前往哪個登入頁
    if (currentPath.startsWith('/admin')) {
      loginPath = '/admin/login'
    }

    // 只有在不已經是登入頁的情況下才導向
    if (!window.location.pathname.includes('/login')) {
      window.location.href = `${loginPath}?redirect=${encodeURIComponent(currentPath)}`
    }
  }
}

/**
 * 清除所有登入相關資料
 */
function clearAllAuthData() {
  if (typeof localStorage !== 'undefined') {
    // 清除所有可能的 token storage keys
    const tokenKeys = ['token', 'auth_token', 'admin_token', 'teacher_token', 'refresh_token']
    tokenKeys.forEach(key => {
      localStorage.removeItem(key)
    })
  }

  if (typeof sessionStorage !== 'undefined') {
    sessionStorage.clear()
  }
}

/**
 * 記錄錯誤到監控服務
 */
function logError(error: ApiError | NetworkError, context?: Record<string, unknown>) {
  const errorData = {
    code: 'code' in error ? error.code : `HTTP_${error.status}`,
    message: error.message,
    timestamp: new Date().toISOString(),
    url: typeof window !== 'undefined' ? window.location.href : '',
    userAgent: typeof navigator !== 'undefined' ? navigator.userAgent : '',
    ...context,
  }

  // 控制台輸出
  console.error('[Error Handler]', errorData)

  // 未來可以整合監控服務，如 Sentry
  // if (typeof window !== 'undefined' && (window as any).Sentry) {
  //   ;(window as any).Sentry.captureException(error)
  // }
}

// ==================== 主要錯誤處理類別 ====================

/**
 * 錯誤處理器
 *
 * 提供靜態方法來處理各種類型的錯誤
 */
export class ErrorHandler {
  /**
   * 處理 API 錯誤回應
   */
  static handleApiError(error: ApiError, options: ErrorHandlerOptions = {}): void {
    const mergedOptions = {
      showAlert: true,
      logError: true,
      redirectOnUnauthorized: true,
      ...options,
    }

    // 檢查是否為成功回應
    if (isSuccessCode(error.code)) {
      return
    }

    // 自訂處理器優先
    if (mergedOptions.onCustomHandler) {
      mergedOptions.onCustomHandler(error)
      return
    }

    // 記錄錯誤
    if (mergedOptions.logError) {
      logError(error, mergedOptions.context)
    }

    // 權限錯誤特殊處理
    if (isUnauthorizedError(error.code)) {
      const handler = getDefaultHandler(401, mergedOptions)
      handler(error)
      return
    }

    // 預設處理
    defaultErrorHandler(error, mergedOptions)
  }

  /**
   * 處理網路錯誤
   */
  static handleNetworkError(error: NetworkError, options: ErrorHandlerOptions = {}): void {
    const mergedOptions = {
      showAlert: true,
      logError: true,
      redirectOnUnauthorized: true,
      ...options,
    }

    // 自訂處理器優先
    if (mergedOptions.onCustomHandler) {
      mergedOptions.onCustomHandler(error)
      return
    }

    // 記錄錯誤
    if (mergedOptions.logError) {
      logError(error, mergedOptions.context)
    }

    // 根據 HTTP 狀態碼處理
    if (error.status) {
      const handler = getDefaultHandler(error.status)
      handler(error, mergedOptions)
      return
    }

    // 網路錯誤預設處理
    if (mergedOptions.showAlert !== false) {
      showErrorAlert('網路連線錯誤，請檢查網路連線後重試', '網路錯誤')
    }
  }

  /**
   * 處理未知錯誤
   */
  static handleUnknownError(error: unknown, options: ErrorHandlerOptions = {}): void {
    const mergedOptions = {
      showAlert: true,
      logError: true,
      redirectOnUnauthorized: false,
      ...options,
    }

    const errorMessage = error instanceof Error ? error.message : '發生未知錯誤'

    if (mergedOptions.logError) {
      console.error('[Unknown Error]', error)
    }

    if (mergedOptions.showAlert !== false) {
      showErrorAlert(errorMessage, '錯誤')
    }
  }

  /**
   * 處理錯誤並返回訊息 (適用於需要顯示 toast 的場景)
   */
  static handleAndReturn(error: ApiError | NetworkError): string {
    const message = getErrorMessage(error)
    showErrorToast(message)
    return message
  }

  /**
   * 安全執行非同步函數
   *
   * @param fn 要執行的非同步函數
   * @param options 錯誤處理選項
   * @returns 包含結果或錯誤的元組
   */
  static async safeExecute<T>(
    fn: () => Promise<T>,
    options: ErrorHandlerOptions = {}
  ): Promise<[T, null] | [null, ApiError | NetworkError]> {
    try {
      const result = await fn()
      return [result, null]
    } catch (error) {
      // 判斷錯誤類型
      if (this.isApiError(error)) {
        this.handleApiError(error as ApiError, options)
        return [null, error as ApiError]
      } else if (this.isNetworkError(error)) {
        this.handleNetworkError(error as NetworkError, options)
        return [null, error as NetworkError]
      } else {
        this.handleUnknownError(error, options)
        return [null, { message: String(error) } as NetworkError]
      }
    }
  }

  /**
   * 檢查是否為 API 錯誤
   */
  static isApiError(error: unknown): error is ApiError {
    if (
      typeof error === 'object' &&
      error !== null &&
      'code' in error &&
      'message' in error &&
      (typeof (error as ApiError).code === 'string' || typeof (error as ApiError).code === 'number')
    ) {
      // 排除成功回應
      const code = (error as ApiError).code
      return !isSuccessCode(code)
    }
    return false
  }

  /**
   * 檢查是否為網路錯誤
   */
  static isNetworkError(error: unknown): error is NetworkError {
    return (
      typeof error === 'object' &&
      error !== null &&
      ('status' in error || 'message' in error)
    )
  }

  /**
   * 取得 HTTP 狀態碼對應的錯誤碼
   */
  static getErrorCodeFromStatus(status: number): string {
    const errorMap: Record<number, string> = {
      400: 'VALIDATION_ERROR',
      401: 'UNAUTHORIZED',
      403: 'FORBIDDEN',
      404: 'NOT_FOUND',
      409: 'CONFLICT',
      422: 'VALIDATION_ERROR',
      429: 'RATE_LIMIT_EXCEEDED',
      500: 'SYSTEM_ERROR',
      502: 'SYSTEM_ERROR',
      503: 'SYSTEM_ERROR',
    }
    return errorMap[status] || 'SYSTEM_ERROR'
  }
}

// ==================== 便利匯出函數 ====================

/**
 * 顯示錯誤 Toast (簡化 API)
 */
export function showErrorToast(message: string) {
  if (errorToastRef.value) {
    errorToastRef.value.show(message, 4000)
  } else {
    // Fallback 到 alert
    showErrorAlert(message, '提示')
  }
}

/**
 * 處理 API 錯誤 (簡化 API)
 */
export function handleApiError(error: ApiError, options?: ErrorHandlerOptions) {
  ErrorHandler.handleApiError(error, options)
}

/**
 * 處理網路錯誤 (簡化 API)
 */
export function handleNetworkError(error: NetworkError, options?: ErrorHandlerOptions) {
  ErrorHandler.handleNetworkError(error, options)
}

/**
 * 處理未知錯誤 (簡化 API)
 */
export function handleUnknownError(error: unknown, options?: ErrorHandlerOptions) {
  ErrorHandler.handleUnknownError(error, options)
}
