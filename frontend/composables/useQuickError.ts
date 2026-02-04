/**
 * 簡化版錯誤處理 Composable
 *
 * 提供快速、易用的錯誤處理接口
 * 整合 useErrorHandler、useToast 和 useAlert
 */

import { ref, computed, readonly } from 'vue'
import { useErrorHandler } from './useErrorHandler'
import type { ApiResponse } from '~/types/api'
import type { ErrorHandlerOptions, ErrorAction } from '~/types/errorHandler'

// ==================== 主要 Composable ====================

/**
 * 快速錯誤處理 Composable
 *
 * @example
 * ```typescript
 * const { handle, showSuccess, showError } = useQuickError()
 *
 * // 處理 API 回應
 * const result = await handle(response)
 * if (!result.success) return
 *
 * // 顯示成功訊息
 * await showSuccess('操作成功')
 *
 * // 顯示錯誤訊息
 * await showError('發生錯誤')
 * ```
 */
export function useQuickError() {
  const errorHandler = useErrorHandler()

  /**
   * 快速顯示錯誤訊息
   */
  const showError = async (
    message: string,
    title?: string
  ): Promise<void> => {
    await errorHandler.showWarning(message, title || '操作失敗')
  }

  /**
   * 快速顯示警告訊息
   */
  const showWarning = async (
    message: string,
    title?: string
  ): Promise<void> => {
    await errorHandler.showWarning(message, title)
  }

  /**
   * 快速顯示成功訊息
   */
  const showSuccess = async (
    message: string,
    title?: string
  ): Promise<void> => {
    await errorHandler.showSuccess(message, title)
  }

  /**
   * 快速顯示資訊訊息
   */
  const showInfo = async (
    message: string,
    title?: string
  ): Promise<void> => {
    await errorHandler.showInfo(message, title)
  }

  /**
   * 快速確認對話框
   */
  const confirm = async (
    message: string,
    title?: string
  ): Promise<boolean> => {
    return errorHandler.confirm(message, title)
  }

  /**
   * 處理 API 請求並自動顯示錯誤
   *
   * @example
   * ```typescript
   * const { request } = useQuickError()
   *
   * // 自動處理錯誤並顯示
   * const data = await request(() => api.getUser(id))
   * if (!data) return // 錯誤已處理
   *
   * // 或處理自訂錯誤
   * const data = await request(
   *   () => api.getUser(id),
   *   { onError: (err) => console.error(err) }
   * )
   * ```
   */
  const request = async <T>(
    fn: () => Promise<ApiResponse<T>>,
    options?: ErrorHandlerOptions
  ): Promise<T | null> => {
    try {
      const response = await fn()
      const result = await errorHandler.handle(response, options)

      if (result.success && result.data) {
        return result.data
      }

      return null
    } catch (error) {
      await errorHandler.catchError(error, options)
      return null
    }
  }

  /**
   * 處理 API 回應（不自動顯示錯誤）
   *
   * @example
   * ```typescript
   * const { handle } = useQuickError()
   *
   * const response = await api.getUser(id)
   * const result = handle(response)
   *
   * if (!result.success) {
   *   // 手動處理錯誤
   *   await showError(result.error)
   *   return
   * }
   *
   * // 成功處理
   * console.log(result.data)
   * ```
   */
  const handle = async <T>(
    response: ApiResponse<T>,
    options?: ErrorHandlerOptions
  ): Promise<{ success: boolean; data?: T; error?: string }> => {
    return errorHandler.handle(response, options)
  }

  /**
   * 處理錯誤（不自動顯示）
   */
  const catchError = async (
    error: Error | string | any,
    options?: ErrorHandlerOptions
  ): Promise<void> => {
    await errorHandler.catchError(error, options)
  }

  /**
   * 處理權限錯誤
   */
  const handlePermission = async (
    code: string,
    message: string,
    options?: ErrorHandlerOptions
  ): Promise<void> => {
    await errorHandler.handlePermission(code, message, options)
  }

  /**
   * 處理驗證錯誤
   */
  const handleValidation = (
    errors: Record<string, string[]>,
    options?: ErrorHandlerOptions
  ): void => {
    errorHandler.handleValidation(errors, options)
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
    await errorHandler.handleConflicts(conflicts, options)
  }

  return {
    // 顯示方法
    showError,
    showWarning,
    showSuccess,
    showInfo,
    confirm,

    // 處理方法
    request,
    handle,
    catchError,

    // 情境方法
    handlePermission,
    handleValidation,
    handleConflicts,

    // 狀態
    ...errorHandler,

    // 工具函數
    getErrorMessage: errorHandler.getErrorMessage,
    isSuccessCode: errorHandler.isSuccessCode,
    isPermissionError: errorHandler.isPermissionError,
    isValidationError: errorHandler.isValidationError,
    isUnauthorizedError: errorHandler.isUnauthorizedError,
  }
}

// ==================== API 請求包裝器 ====================

/**
 * API 請求配置
 */
export interface ApiRequestOptions<T> extends ErrorHandlerOptions {
  /** 成功時的回调 */
  onSuccess?: (data: T) => void
  /** 無論成功或失敗都執行 */
  finally?: () => void
}

/**
 * 創建帶錯誤處理的 API 請求
 *
 * @example
 * ```typescript
 * const fetchUser = createApiRequest(async () => {
 *   const response = await api.getUser(id)
 *   return response.data
 * })
 *
 * // 使用
 * const { loading, data, error } = await fetchUser({
 *   onSuccess: (user) => console.log('取得成功', user),
 * })
 * ```
 */
export function createApiRequest<T>(
  requestFn: () => Promise<ApiResponse<T>>
) {
  return async (options: ApiRequestOptions<T> = {}) => {
    const { onSuccess, finally: finallyCallback, ...errorOptions } = options

    const response = await requestFn()
    const result = await useErrorHandler().handle(response, errorOptions)

    if (result.success && result.data) {
      onSuccess?.(result.data)
    }

    finallyCallback?.()

    return {
      success: result.success,
      data: result.data,
      error: result.error,
    }
  }
}

/**
 * 創建異步操作包裝器（自動錯誤處理）
 *
 * @example
 * ```typescript
 * const { execute, loading, error } = useAsyncOperation()
 *
 * const handleSave = async () => {
 *   const result = await execute(async () => {
 *     await api.saveData(formData)
 *     return '儲存成功'
 *   })
 *
 *   if (result.success) {
 *     await showSuccess(result.data)
 *   }
 * }
 * ```
 */
export function useAsyncOperation() {
  const loading = ref(false)
  const error = ref<string | null>(null)

  const execute = async <T>(
    fn: () => Promise<T>,
    options: ErrorHandlerOptions = {}
  ): Promise<{ success: boolean; data?: T; error?: string }> => {
    loading.value = true
    error.value = null

    try {
      const data = await fn()
      return { success: true, data }
    } catch (err) {
      await useErrorHandler().catchError(err, options)
      return { success: false, error: String(err) }
    } finally {
      loading.value = false
    }
  }

  return {
    loading: readonly(loading),
    error: readonly(error),
    execute,
  }
}

// ==================== 表單錯誤處理 ====================

/**
 * 表單驗證錯誤處理
 *
 * @example
 * ```typescript
 * const { setFieldError, clearErrors, getErrorMessage } = useFormErrors()
 *
 * // 設置欄位錯誤
 * setFieldError('email', '電子郵件格式錯誤')
 *
 * // 清除所有錯誤
 * clearErrors()
 * ```
 */
export function useFormErrors() {
  const errors = ref<Record<string, string>>({})

  /**
   * 設置單一欄位錯誤
   */
  const setFieldError = (field: string, message: string): void => {
    errors.value[field] = message
  }

  /**
   * 設置多個欄位錯誤
   */
  const setErrors = (newErrors: Record<string, string>): void => {
    errors.value = { ...errors.value, ...newErrors }
  }

  /**
   * 清除單一欄位錯誤
   */
  const clearFieldError = (field: string): void => {
    delete errors.value[field]
  }

  /**
   * 清除所有錯誤
   */
  const clearErrors = (): void => {
    errors.value = {}
  }

  /**
   * 取得欄位錯誤訊息
   */
  const getFieldError = (field: string): string | null => {
    return errors.value[field] || null
  }

  /**
   * 是否有任何錯誤
   */
  const hasErrors = computed(() => Object.keys(errors.value).length > 0)

  /**
   * 取得所有錯誤訊息
   */
  const allErrors = computed(() => Object.values(errors.value))

  return {
    errors: readonly(errors),
    hasErrors: readonly(hasErrors),
    allErrors: readonly(allErrors),

    setFieldError,
    setErrors,
    clearFieldError,
    clearErrors,
    getFieldError,
  }
}

// ==================== 錯誤邊界處理 ====================

/**
 * 全局錯誤邊界 Composable
 *
 * @example
 * ```typescript
 * const { captureError, captureException } = useErrorBoundary()
 *
 * // 捕獲錯誤
 * captureError(new Error('Something went wrong'), {
 *   context: { userId: 123, action: 'save' }
 * })
 * ```
 */
export function useErrorBoundary() {
  const { showError } = useQuickError()

  /**
   * 捕獲錯誤並顯示
   */
  const captureError = async (
    error: Error | string,
    context?: Record<string, any>
  ): Promise<void> => {
    // 記錄到控制台（除錯用）
    console.error('Error captured:', error, context)

    // 顯示錯誤給使用者
    const message = error instanceof Error ? error.message : error
    await showError(message, '發生錯誤')
  }

  /**
   * 捕獲 Promise 錯誤
   */
  const capturePromise = <T>(
    promise: Promise<T>,
    errorOptions?: ErrorHandlerOptions
  ): Promise<T | null> => {
    return promise.catch(async (error) => {
      await useErrorHandler().catchError(error, errorOptions)
      return null
    }) as Promise<T | null>
  }

  /**
   * 安全的異步操作
   */
  const safeExecute = async <T>(
    fn: () => Promise<T>,
    fallback?: T,
    errorOptions?: ErrorHandlerOptions
  ): Promise<T> => {
    try {
      return await fn()
    } catch (error) {
      await useErrorHandler().catchError(error, errorOptions)
      return fallback as T
    }
  }

  return {
    captureError,
    capturePromise,
    safeExecute,
  }
}
