<template>
  <Teleport to="body">
    <TransitionGroup
      name="error-toast"
      tag="div"
      class="fixed bottom-4 right-4 z-[9999] space-y-3"
    >
      <div
        v-for="item in errorToasts"
        :key="item.id"
        class="glass-card px-4 py-3 rounded-xl shadow-xl flex items-start gap-3 min-w-[320px] max-w-md border-l-4"
        :class="getToastClass(item.type)"
      >
        <!-- 錯誤圖示 -->
        <div class="shrink-0 mt-0.5">
          <svg
            v-if="item.type === 'error'"
            class="w-5 h-5 text-critical-500"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <svg
            v-else-if="item.type === 'warning'"
            class="w-5 h-5 text-warning-500"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
            />
          </svg>
          <svg
            v-else-if="item.type === 'info'"
            class="w-5 h-5 text-primary-500"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <svg
            v-else
            class="w-5 h-5 text-success-500"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        </div>

        <!-- 錯誤內容 -->
        <div class="flex-1 min-w-0">
          <p v-if="item.title" class="font-medium text-slate-100 text-sm mb-1">
            {{ item.title }}
          </p>
          <p class="text-slate-300 text-sm break-words">
            {{ item.message }}
          </p>

          <!-- 倒數計時指示器 -->
          <div
            v-if="item.autoClose && item.duration > 0"
            class="mt-2 h-1 rounded-full bg-white/10 overflow-hidden"
          >
            <div
              class="h-full transition-all ease-linear"
              :class="progressClass(item.type)"
              :style="{ width: `${(item.remainingTime / item.duration) * 100}%` }"
            />
          </div>

          <!-- 重試按鈕 -->
          <div
            v-if="item.onRetry"
            class="mt-3 flex items-center gap-2"
          >
            <button
              @click="handleRetry(item)"
              class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors"
              :class="retryButtonClass(item.type)"
            >
              <svg class="w-3.5 h-3.5 inline-block mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              重試
            </button>
            <button
              @click="dismiss(item.id)"
              class="px-3 py-1.5 text-xs font-medium text-slate-400 hover:text-slate-200 rounded-lg transition-colors"
            >
              關閉
            </button>
          </div>
        </div>

        <!-- 關閉按鈕 -->
        <button
          v-if="!item.autoClose || !item.onRetry"
          @click="dismiss(item.id)"
          class="shrink-0 p-1 rounded-lg hover:bg-white/10 transition-colors"
        >
          <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </TransitionGroup>
  </Teleport>
</template>

<script setup lang="ts">
/**
 * ErrorToast 組件
 *
 * 提供錯誤提示功能，支援：
 * - 錯誤、警告、資訊、成功類型
 * - 自動倒數關閉
 * - 重試按鈕
 * - 玻璃擬態 UI 設計
 */

interface ErrorToastItem {
  id: number
  type: 'error' | 'warning' | 'info' | 'success'
  message: string
  title?: string
  duration?: number
  autoClose?: boolean
  onRetry?: () => void
  remainingTime?: number
}

// 狀態
const errorToasts = ref<ErrorToastItem[]>([])
let toastId = 0

// 計時器映射
const timers = new Map<number, NodeJS.Timeout>()

/**
 * 取得 toast 樣式類別
 */
const getToastClass = (type: string): string => {
  switch (type) {
    case 'error':
      return 'border-l-critical-500 bg-critical-500/10'
    case 'warning':
      return 'border-l-warning-500 bg-warning-500/10'
    case 'info':
      return 'border-l-primary-500 bg-primary-500/10'
    case 'success':
      return 'border-l-success-500 bg-success-500/10'
    default:
      return 'border-l-primary-500 bg-primary-500/10'
  }
}

/**
 * 取得進度條顏色類別
 */
const progressClass = (type: string): string => {
  switch (type) {
    case 'error':
      return 'bg-critical-500'
    case 'warning':
      return 'bg-warning-500'
    case 'info':
      return 'bg-primary-500'
    case 'success':
      return 'bg-success-500'
    default:
      return 'bg-primary-500'
  }
}

/**
 * 取得重試按鈕類別
 */
const retryButtonClass = (type: string): string => {
  switch (type) {
    case 'error':
      return 'bg-critical-500/20 text-critical-400 hover:bg-critical-500/30'
    case 'warning':
      return 'bg-warning-500/20 text-warning-400 hover:bg-warning-500/30'
    case 'info':
      return 'bg-primary-500/20 text-primary-400 hover:bg-primary-500/30'
    case 'success':
      return 'bg-success-500/20 text-success-400 hover:bg-success-500/30'
    default:
      return 'bg-primary-500/20 text-primary-400 hover:bg-primary-500/30'
  }
}

/**
 * 關閉 toast
 */
const dismiss = (id: number) => {
  // 清除計時器
  const timer = timers.get(id)
  if (timer) {
    clearTimeout(timer)
    timers.delete(id)
  }

  // 移除 toast
  errorToasts.value = errorToasts.value.filter((t) => t.id !== id)
}

/**
 * 處理重試
 */
const handleRetry = (item: ErrorToastItem) => {
  if (item.onRetry) {
    item.onRetry()
  }
  dismiss(item.id)
}

/**
 * 新增錯誤提示
 */
const show = (
  message: string,
  options: {
    type?: 'error' | 'warning' | 'info' | 'success'
    title?: string
    duration?: number
    autoClose?: boolean
    onRetry?: () => void
  } = {}
) => {
  const {
    type = 'error',
    title,
    duration = 5000,
    autoClose = true,
    onRetry,
  } = options

  const id = ++toastId
  const remainingTime = duration

  const toast: ErrorToastItem = {
    id,
    type,
    message,
    title,
    duration,
    autoClose,
    onRetry,
    remainingTime,
  }

  errorToasts.value.push(toast)

  // 設定自動關閉
  if (autoClose) {
    const timer = setTimeout(() => {
      dismiss(id)
    }, duration)
    timers.set(id, timer)
  }

  // 註冊到全域
  registerErrorToast({
    show: (msg: string, dur?: number) => {
      show(msg, { type: 'error', duration: dur })
    },
  })

  return id
}

/**
 * 快速顯示錯誤
 */
const error = (message: string, title?: string) => {
  return show(message, { type: 'error', title })
}

/**
 * 快速顯示警告
 */
const warning = (message: string, title?: string) => {
  return show(message, { type: 'warning', title })
}

/**
 * 快速顯示資訊
 */
const info = (message: string, title?: string) => {
  return show(message, { type: 'info', title })
}

/**
 * 快速顯示成功
 */
const success = (message: string, title?: string) => {
  return show(message, { type: 'success', title })
}

/**
 * 顯示可重試的錯誤
 */
const retryable = (
  message: string,
  onRetry: () => void,
  title?: string
) => {
  return show(message, {
    type: 'error',
    title: title || '發生錯誤',
    autoClose: false,
    onRetry,
  })
}

// 暴露方法給外部使用
defineExpose({
  show,
  error,
  warning,
  info,
  success,
  retryable,
  dismiss,
})
</script>

<style scoped>
.error-toast-enter-active,
.error-toast-leave-active {
  transition: all 0.3s ease;
}

.error-toast-enter-from {
  opacity: 0;
  transform: translateX(100px);
}

.error-toast-leave-to {
  opacity: 0;
  transform: translateX(100px);
}

.error-toast-move {
  transition: transform 0.3s ease;
}
</style>
