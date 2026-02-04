/**
 * 全局錯誤處理 Composable
 *
 * 整合錯誤碼對應、使用者友善訊息和 UI 提示
 */

import { ref, readonly } from 'vue'
import { alertError, alertWarning, alertSuccess, alertInfo, alertConfirm, type AlertType } from './useAlert'
import { useToast, type ToastType } from './useToast'
import {
  ERROR_MESSAGES,
  type ErrorCode,
  isSuccessCode,
  isPermissionError,
  isValidationError,
  isUnauthorizedError,
} from '~/constants/errorCodes'
import type { ApiResponse } from '~/types/api'
import type { ErrorHandlerOptions, ErrorAction } from '~/types/errorHandler'

// ==================== 錯誤碼對應類型 ====================

// ==================== 全局錯誤狀態 ====================

/**
 * 當前是否正在顯示錯誤
 */
const isShowingError = ref(false)

/**
 * 當前錯誤佇列
 */
const errorQueue = ref<Array<{
  response: ApiResponse
  options: ErrorHandlerOptions
}>>([])

/**
 * 是否正在處理錯誤佇列
 */
const isProcessingQueue = ref(false)

/**
 * 取得當前錯誤狀態（唯讀）
 */
export function useErrorStatus() {
  return {
    hasError: readonly(isShowingError),
    queueLength: readonly(errorQueue),
    isProcessing: readonly(isProcessingQueue),
  }
}

// ==================== 錯誤處理核心函數 ====================

/**
 * 取得錯誤碼對應的使用者友善訊息
 */
export function getErrorMessage(code: string, fallback?: string): string {
  // 先檢查是否為已知的錯誤碼
  if (ERROR_MESSAGES[code]) {
    return ERROR_MESSAGES[code]
  }

  // 處理 HTTP 狀態碼格式
  const httpCode = parseInt(code)
  if (!isNaN(httpCode) && httpCode >= 100 && httpCode < 600) {
    return ERROR_MESSAGES.SYSTEM_ERROR
  }

  // 回退訊息
  return fallback || '發生錯誤，請稍後再試'
}

/**
 * 取得錯誤類型（用於 UI 顯示）
 */
export function getErrorAlertType(code: string): AlertType {
  if (isSuccessCode(code)) {
    return 'success'
  }

  if (isUnauthorizedError(code) || isPermissionError(code)) {
    return 'warning'
  }

  if (isValidationError(code)) {
    return 'info'
  }

  return 'error'
}

/**
 * 取得 Toast 類型
 */
function getErrorToastType(code: string): ToastType {
  if (isSuccessCode(code)) {
    return 'success'
  }

  if (isUnauthorizedError(code) || isPermissionError(code)) {
    return 'warning'
  }

  if (isValidationError(code)) {
    return 'info'
  }

  return 'error'
}

/**
 * 處理 API 回應
 */
export async function handleApiResponse<T>(
  response: ApiResponse<T>,
  options: ErrorHandlerOptions = {}
): Promise<{ success: boolean; data?: T; error?: string }> {
  const { code, message, data } = response

  // 成功處理
  if (isSuccessCode(code)) {
    options.onSuccess?.(data)
    return { success: true, data }
  }

  // 錯誤處理
  const userMessage = getErrorMessage(code, message)
  const alertType = getErrorAlertType(code)

  // 準備動作按鈕
  const buttons = options.actions?.map(action => ({
    text: action.label,
    style: action.style || 'primary',
    action: action.action,
  }))

  // 顯示錯誤
  isShowingError.value = true

  try {
    if (options.asToast) {
      // 使用 Toast 顯示（非阻斷式）
      const toast = useToast()
      const toastType = options.toastType || getErrorToastType(code)
      toast[toastType](userMessage, options.title)
    } else {
      // 使用 Alert 顯示（阻斷式）
      await alertError(userMessage, options.title || getDefaultTitle(code))
    }

    options.onError?.(response)
  } finally {
    isShowingError.value = false
  }

  return { success: false, error: userMessage }
}

/**
 * 取得錯誤預設標題
 */
function getDefaultTitle(code: string): string {
  if (isSuccessCode(code)) return '操作成功'
  if (isUnauthorizedError(code)) return '需要登入'
  if (isPermissionError(code)) return '權限不足'
  if (isValidationError(code)) return '資料驗證'
  return '發生錯誤'
}

/**
 * 處理錯誤物件（非 API 回應格式）
 */
export async function handleError(
  error: Error | string | any,
  options: ErrorHandlerOptions = {}
): Promise<void> {
  let message: string
  let code: string = 'UNKNOWN_ERROR'

  if (typeof error === 'string') {
    message = error
  } else if (error instanceof Error) {
    message = error.message || '發生未知錯誤'
  } else if (error?.message) {
    message = error.message
  } else if (error?.code) {
    code = error.code
    message = getErrorMessage(code)
  } else {
    message = '發生未知錯誤'
  }

  const alertType = getErrorAlertType(code)

  isShowingError.value = true

  try {
    if (options.asToast) {
      const toast = useToast()
      const toastType = options.toastType || getErrorToastType(code)
      toast[toastType](message, options.title)
    } else {
      await alertError(message, options.title)
    }

    options.onError?.({ code, message, data: null })
  } finally {
    isShowingError.value = false
  }
}

// ==================== 佇列處理 ====================

/**
 * 將錯誤加入佇列
 */
export function queueError(response: ApiResponse, options: ErrorHandlerOptions = {}): void {
  errorQueue.value.push({ response, options })
}

/**
 * 處理錯誤佇列
 */
export async function processErrorQueue(): Promise<void> {
  if (isProcessingQueue.value || errorQueue.value.length === 0) {
    return
  }

  isProcessingQueue.value = true

  try {
    while (errorQueue.value.length > 0) {
      const { response, options } = errorQueue.value.shift()!
      await handleApiResponse(response, options)
    }
  } finally {
    isProcessingQueue.value = false
  }
}

/**
 * 清空錯誤佇列
 */
export function clearErrorQueue(): void {
  errorQueue.value = []
}

// ==================== 常見錯誤情境處理 ====================

/**
 * 處理權限相關錯誤
 */
export async function handlePermissionError(
  code: string,
  message: string,
  options: ErrorHandlerOptions = {}
): Promise<void> {
  const actions: ErrorAction[] = []

  // 如果是需要登入的錯誤（401），清除登入資訊並跳轉到首頁
  if (isUnauthorizedError(code)) {
    // 清除登入資訊
    clearAuthOnUnauthorized()

    // 根據使用者類型決定跳轉頁面
    const redirectPath = determineLoginPath()

    actions.push({
      label: '前往登入',
      style: 'primary',
      action: () => {
        navigateTo(redirectPath)
      },
    })
  }

  await handleApiResponse(
    { code, message, data: null },
    {
      ...options,
      title: options.title || '權限不足',
      actions: [...(options.actions || []), ...actions],
    }
  )
}

/**
 * 清除未授權使用者的登入資訊
 */
function clearAuthOnUnauthorized(): void {
  try {
    // 清除所有可能的 token storage keys
    const tokenKeys = ['token', 'auth_token', 'admin_token', 'teacher_token']
    tokenKeys.forEach(key => {
      if (typeof localStorage !== 'undefined') {
        localStorage.removeItem(key)
      }
    })

    // 清除 sessionStorage
    if (typeof sessionStorage !== 'undefined') {
      sessionStorage.clear()
    }
  } catch (error) {
    console.error('Failed to clear auth storage:', error)
  }
}

/**
 * 根據使用者類型決定登入頁面路徑
 */
function determineLoginPath(): string {
  // 嘗試從 localStorage 判斷使用者類型
  if (typeof localStorage !== 'undefined') {
    // 檢查是否有 admin 相關的 token
    const adminToken = localStorage.getItem('admin_token')
    if (adminToken) {
      return '/admin/login'
    }

    // 檢查是否有 teacher 相關的 token
    const teacherToken = localStorage.getItem('teacher_token')
    if (teacherToken) {
      return '/teacher/login'
    }
  }

  // 預設跳轉到首頁
  return '/'
}

/**
 * 處理驗證錯誤
 */
export function handleValidationErrors(
  errors: Record<string, string[]>,
  options: ErrorHandlerOptions = {}
): void {
  // 組合多個驗證錯誤訊息
  const messages = Object.entries(errors)
    .map(([field, msgs]) => `${field}: ${msgs.join(', ')}`)
    .join('\n')

  const combinedMessage = messages || '輸入資料驗證失敗'

  alertWarning(combinedMessage, options.title || '資料驗證')
}

/**
 * 處理排課衝突錯誤
 */
export async function handleScheduleConflict(
  conflicts: Array<{
    type: 'TEACHER' | 'ROOM' | 'BUFFER'
    message: string
    details?: any
  }>,
  options: ErrorHandlerOptions = {}
): Promise<void> {
  const messages = conflicts
    .map(c => {
      const typeLabel = {
        TEACHER: '👤 老師',
        ROOM: '🏠 教室',
        BUFFER: '⏱️ 緩衝時間',
      }[c.type] || '⚠️'

      return `${typeLabel} ${c.message}`
    })
    .join('\n\n')

  const actions: ErrorAction[] = []

  // 提供查看課表動作
  actions.push({
    label: '查看課表',
    style: 'secondary',
    action: () => {
      navigateTo('/admin/schedule')
    },
  })

  await alertWarning(
    messages,
    options.title || '排課衝突',
    {
      ...options,
      actions: [...(options.actions || []), ...actions],
    }
  )
}

// ==================== 主要 Composable ====================

/**
 * 全局錯誤處理 Composable
 */
export function useErrorHandler() {
  const { hasError, queueLength, isProcessing } = useErrorStatus()

  /**
   * 處理 API 回應
   */
  const handle = async <T>(
    response: ApiResponse<T>,
    options?: ErrorHandlerOptions
  ): Promise<{ success: boolean; data?: T; error?: string }> => {
    return handleApiResponse(response, options)
  }

  /**
   * 處理錯誤
   */
  const catchError = async (
    error: Error | string | any,
    options?: ErrorHandlerOptions
  ): Promise<void> => {
    await handleError(error, options)
  }

  /**
   * 顯示成功訊息
   */
  const showSuccess = async (message: string, title?: string): Promise<void> => {
    await alertSuccess(message, title)
  }

  /**
   * 顯示警告訊息
   */
  const showWarning = async (message: string, title?: string): Promise<void> => {
    await alertWarning(message, title)
  }

  /**
   * 顯示資訊訊息
   */
  const showInfo = async (message: string, title?: string): Promise<void> => {
    await alertInfo(message, title)
  }

  /**
   * 顯示確認對話框
   */
  const confirm = async (
    message: string,
    title?: string
  ): Promise<boolean> => {
    return alertConfirm(message, title)
  }

  /**
   * 將錯誤加入佇列
   */
  const queue = (response: ApiResponse, options?: ErrorHandlerOptions): void => {
    queueError(response, options)
  }

  /**
   * 處理佇列中的錯誤
   */
  const processQueue = async (): Promise<void> => {
    await processErrorQueue()
  }

  /**
   * 清空錯誤佇列
   */
  const clearQueue = (): void => {
    clearErrorQueue()
  }

  /**
   * 處理權限錯誤
   */
  const handlePermission = async (
    code: string,
    message: string,
    options?: ErrorHandlerOptions
  ): Promise<void> => {
    await handlePermissionError(code, message, options)
  }

  /**
   * 處理驗證錯誤
   */
  const handleValidation = (
    errors: Record<string, string[]>,
    options?: ErrorHandlerOptions
  ): void => {
    handleValidationErrors(errors, options)
  }

  /**
   * 處理排課衝突
   */
  const handleConflicts = async (
    conflicts: Array<{
      type: 'TEACHER' | 'ROOM' | 'BUFFER'
      message: string
      details?: any
    }>,
    options?: ErrorHandlerOptions
  ): Promise<void> => {
    await handleScheduleConflict(conflicts, options)
  }

  return {
    // 狀態
    hasError,
    queueLength,
    isProcessing,

    // 核心方法
    handle,
    catchError,

    // 顯示方法
    showSuccess,
    showWarning,
    showInfo,
    confirm,

    // 佇列方法
    queue,
    processQueue,
    clearQueue,

    // 情境方法
    handlePermission,
    handleValidation,
    handleConflicts,

    // 工具函數
    getErrorMessage,
    getErrorAlertType,
    getErrorToastType,
    isSuccessCode,
    isPermissionError,
    isValidationError,
    isUnauthorizedError,
  }
}

// ==================== Axios 攔截器輔助 ====================

/**
 * 從 Axios 錯誤回應取得 API 回應格式
 */
export function extractApiResponse(error: any): ApiResponse {
  if (error.response?.data) {
    const { code, message, data } = error.response.data
    return {
      code: code || error.response.status?.toString() || 'SYSTEM_ERROR',
      message: message || getErrorMessage(code),
      data,
    }
  }

  if (error.request) {
    return {
      code: 'NETWORK_ERROR',
      message: '網路連線失敗，請檢查網路連線',
    }
  }

  return {
    code: 'UNKNOWN_ERROR',
    message: error.message || '發生未知錯誤',
  }
}

/**
 * 建立 Axios 錯誤處理器
 */
export function createAxiosErrorHandler(options: ErrorHandlerOptions = {}) {
  return async (error: any): Promise<void> => {
    const response = extractApiResponse(error)
    await handleApiResponse(response, options)
  }
}

// ==================== 類型重新匯出（向後相容） ====================

/**
 * API 回應結構（從 types/api.ts 重新匯出）
 *
 * @deprecated 請直接從 ~/types/api 匯入
 */
export type {
  ApiResponse,
} from '~/types/api'
