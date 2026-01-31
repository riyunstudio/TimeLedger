import { ref, computed, type Ref } from 'vue'

/**
 * 自動管理 loading 狀態的工具函數
 * 確保 loading 狀態一定會被清除（無論成功或失敗）
 */
export async function withLoading<T>(
  loadingRef: Ref<boolean>,
  fn: () => Promise<T>
): Promise<T> {
  loadingRef.value = true
  try {
    return await fn()
  } finally {
    loadingRef.value = false
  }
}

/**
 * 平行請求的 loading 管理
 * 適用於多個並行發送的 API 請求
 */
export function useParallelLoading() {
  const pending = ref(0)
  const isLoading = computed(() => pending.value > 0)

  const track = <T>(promise: Promise<T>): Promise<T> => {
    pending.value++
    return promise.finally(() => pending.value--)
  }

  return { isLoading, track }
}

/**
 * 組合式 loading（多個狀態）
 * 當任一 loadingRef 為 true 時，isLoading 為 true
 */
export function useCompoundLoading(
  ...loadingRefs: Ref<boolean>[]
) {
  const isLoading = computed(() =>
    loadingRefs.some(ref => ref.value)
  )
  return { isLoading }
}

/**
 * 組合式 loading（傳入陣列）
 * 當傳入 Ref 陣列時使用
 */
export function useCompoundLoadingFromArray(
  loadingRefs: Ref<boolean>[]
) {
  const isLoading = computed(() =>
    loadingRefs.some(ref => ref.value)
  )
  return { isLoading }
}

/**
 * 創建帶 loading 狀態的異步操作鉤子
 * 適用於需要追蹤特定操作載入狀態的場景
 */
export function useLoadingState() {
  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  const execute = async <T>(
    fn: () => Promise<T>
  ): Promise<T | null> => {
    isLoading.value = true
    error.value = null
    try {
      return await fn()
    } catch (e) {
      error.value = e instanceof Error ? e : new Error(String(e))
      return null
    } finally {
      isLoading.value = false
    }
  }

  return { isLoading, error, execute }
}

/**
 * 延遲 loading 狀態
 * 避免閃爍的短暫 loading 狀態
 */
export function useDelayedLoading(delayMs: number = 300) {
  const isLoading = ref(false)
  const showLoading = ref(false)
  let timeoutId: ReturnType<typeof setTimeout> | null = null

  const start = () => {
    isLoading.value = true
    timeoutId = setTimeout(() => {
      showLoading.value = true
    }, delayMs)
  }

  const stop = () => {
    isLoading.value = false
    if (timeoutId) {
      clearTimeout(timeoutId)
      timeoutId = null
    }
    showLoading.value = false
  }

  const reset = () => {
    stop()
    start()
  }

  return {
    isLoading,
    showLoading,
    start,
    stop,
    reset
  }
}
