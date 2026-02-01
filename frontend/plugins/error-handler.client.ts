/**
 * Nuxt 錯誤處理插件
 *
 * 自動註冊錯誤處理 composables 到全域
 */

import { useErrorHandler } from '~/composables/useErrorHandler'
import { useQuickError } from '~/composables/useQuickError'
import { useFormErrors } from '~/composables/useQuickError'
import { useAsyncOperation } from '~/composables/useQuickError'
import { useErrorBoundary } from '~/composables/useQuickError'

export default defineNuxtPlugin(() => {
  // 提供主要的錯誤處理器
  const errorHandler = useErrorHandler()

  // 提供簡化版錯誤處理器
  const quickError = useQuickError()

  // 提供表單錯誤處理
  const formErrors = useFormErrors()

  // 提供異步操作處理
  const asyncOperation = useAsyncOperation()

  // 提供錯誤邊界處理
  const errorBoundary = useErrorBoundary()

  return {
    provide: {
      // 主要錯誤處理器
      error: errorHandler,

      // 簡化版錯誤處理
      quickError,

      // 表單錯誤
      formErrors,

      // 異步操作
      asyncOperation,

      // 錯誤邊界
      errorBoundary,
    },
  }
})

// 類型擴展
declare module '#app' {
  interface NuxtApp {
    $error: ReturnType<typeof useErrorHandler>
    $quickError: ReturnType<typeof useQuickError>
    $formErrors: ReturnType<typeof useFormErrors>
    $asyncOperation: ReturnType<typeof useAsyncOperation>
    $errorBoundary: ReturnType<typeof useErrorBoundary>
  }
}
