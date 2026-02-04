import { ref } from 'vue'

// ============================================
// Loading 狀態管理
// ============================================

/**
 * 通用 loading 包裝器
 *
 * @deprecated 請使用 loadingHelper.ts 中的 withLoading
 */
export { withLoading } from './loadingHelper'

/**
 * 建立 Loading 狀態管理
 *
 * @example
 * ```typescript
 * const loadingState = createLoadingState({
 *   fetching: false,
 *   creating: false,
 *   updating: false,
 *   deleting: false
 * })
 *
 * // 使用
 * async function createItem(data: ItemData) {
 *   loadingState.creating = true
 *   try {
 *     // API 呼叫
 *   } finally {
 *     loadingState.creating = false
 *   }
 * }
 * ```
 */
export function createLoadingState<T extends Record<string, boolean>>(
  initialState: T
) {
  const state = ref<T>({ ...initialState })

  const isAnyLoading = computed(() =>
    Object.values(state.value).some(Boolean)
  )

  function setLoading(key: keyof T, value: boolean) {
    state.value[key] = value
  }

  function isLoading(key: keyof T) {
    return computed(() => state.value[key])
  }

  function reset() {
    Object.keys(state.value).forEach(key => {
      state.value[key as keyof T] = initialState[key as keyof T]
    })
  }

  return {
    state,
    isAnyLoading,
    setLoading,
    isLoading,
    reset
  }
}

// ============================================
// API 呼叫助手
// ============================================

/**
 * API 呼叫選項
 */
export interface ApiCallOptions<T> {
  /** 成功回調 */
  onSuccess?: (data: T) => void
  /** 錯誤回調 */
  onError?: (error: Error) => void
  /** 無論成功或失敗都會執行的回調 */
  onFinally?: () => void
  /** loading 狀態參考 */
  loadingRef?: { value: boolean }
  /** 錯誤時是否拋出異常 */
  throwOnError?: boolean
}

/**
 * 標準 API 呼叫範本
 *
 * @example
 * ```typescript
 * const loading = ref(false)
 *
 * // 簡單使用
 * await apiCall(async () => {
 *   const api = useApi()
 *   return await api.get('/users')
 * })
 *
 * // 完整選項
 * await apiCall(async () => {
 *   const api = useApi()
 *   return await api.post('/users', data)
 * }, {
 *   onSuccess: (data) => showToast('成功'),
 *   onError: (err) => showToast('失敗'),
 *   loadingRef: loading
 * })
 * ```
 */
export async function apiCall<T>(
  apiFn: () => Promise<T>,
  options?: ApiCallOptions<T>
): Promise<T | null> {
  const {
    onSuccess,
    onError,
    onFinally,
    loadingRef,
    throwOnError = false
  } = options || {}

  if (loadingRef) {
    loadingRef.value = true
  }

  try {
    const data = await apiFn()
    onSuccess?.(data)
    return data
  } catch (error) {
    const err = error instanceof Error ? error : new Error(String(error))
    onError?.(err)

    if (throwOnError) {
      throw err
    }

    return null
  } finally {
    if (loadingRef) {
      loadingRef.value = false
    }
    onFinally?.()
  }
}

/**
 * 建立帶 loading 狀態的 API 呼叫函式
 *
 * @example
 * ```typescript
 * const fetchUsers = createApiCaller({
 *   endpoint: '/users',
 *   loadingState: ref(false)
 * })
 *
 * // 使用
 * const result = await fetchUsers.get()
 * const result = await fetchUsers.post(data)
 * ```
 */
export function createApiCaller(options: {
  endpoint: string
  loadingState: { value: boolean }
}) {
  const { endpoint, loadingState } = options

  const api = useApi()

  async function get<T = any>(params?: Record<string, any>) {
    return apiCall(
      () => api.get<any>(endpoint, { params }),
      { loadingRef: loadingState }
    ) as Promise<T | null>
  }

  async function post<T = any>(data?: any) {
    return apiCall(
      () => api.post<any>(endpoint, data),
      { loadingRef: loadingState }
    ) as Promise<T | null>
  }

  async function put<T = any>(data?: any) {
    return apiCall(
      () => api.put<any>(endpoint, data),
      { loadingRef: loadingState }
    ) as Promise<T | null>
  }

  async function patch<T = any>(data?: any) {
    return apiCall(
      () => api.patch<any>(endpoint, data),
      { loadingRef: loadingState }
    ) as Promise<T | null>
  }

  async function deleteItem<T = any>() {
    return apiCall(
      () => api.delete<any>(endpoint),
      { loadingRef: loadingState }
    ) as Promise<T | null>
  }

  return { get, post, put, patch, delete: deleteItem }
}

// ============================================
// 分頁管理
// ============================================

/**
 * 分頁狀態管理器
 *
 * @example
 * ```typescript
 * const pagination = createPaginationState({
 *   page: 1,
 *   limit: 20
 * })
 *
 * // 下一頁
 * pagination.nextPage()
 *
 * // 上一頁
 * pagination.prevPage()
 *
 * // 跳轉到指定頁
 * pagination.goToPage(5)
 *
 * // 重置
 * pagination.reset()
 * ```
 */
export function createPaginationState(initial: { page?: number; limit?: number } = {}) {
  const page = ref(initial.page || 1)
  const limit = ref(initial.limit || 20)

  const offset = computed(() => (page.value - 1) * limit.value)

  const hasNext = computed(() => true) // 需配合實際資料判斷
  const hasPrev = computed(() => page.value > 1)

  function nextPage() {
    page.value++
  }

  function prevPage() {
    if (page.value > 1) {
      page.value--
    }
  }

  function goToPage(newPage: number) {
    if (newPage >= 1) {
      page.value = newPage
    }
  }

  function setLimit(newLimit: number) {
    limit.value = Math.max(1, Math.min(100, newLimit))
    page.value = 1 // 重置到第一頁
  }

  function reset() {
    page.value = initial.page || 1
    limit.value = initial.limit || 20
  }

  function getParams() {
    return {
      page: page.value,
      limit: limit.value
    }
  }

  return {
    page: readonly(page),
    limit: readonly(limit),
    offset,
    hasNext,
    hasPrev,
    nextPage,
    prevPage,
    goToPage,
    setLimit,
    reset,
    getParams
  }
}

// ============================================
// 狀態持久化
// ============================================

/**
 * 載入 localStorage 資料到 Store
 */
export function loadFromStorage<T>(key: string, defaultValue: T): T {
  try {
    const stored = localStorage.getItem(key)
    if (stored === null) {
      return defaultValue
    }
    return JSON.parse(stored) as T
  } catch {
    return defaultValue
  }
}

/**
 * 儲存 Store 資料到 localStorage
 */
export function saveToStorage<T>(key: string, value: T): void {
  try {
    localStorage.setItem(key, JSON.stringify(value))
  } catch (error) {
    console.error(`Failed to save to localStorage: ${key}`, error)
  }
}

/**
 * 從 localStorage 清除資料
 */
export function clearFromStorage(key: string): void {
  try {
    localStorage.removeItem(key)
  } catch (error) {
    console.error(`Failed to clear localStorage: ${key}`, error)
  }
}

// ============================================
// Debounce / Throttle
// ============================================

/**
 * 防抖動函式
 */
export function useDebounce<T extends (...args: any[]) => any>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: ReturnType<typeof setTimeout> | null = null

  return function (this: any, ...args: Parameters<T>) {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }

    timeoutId = setTimeout(() => {
      fn.apply(this, args)
    }, delay)
  }
}

/**
 * 節流函式
 */
export function useThrottle<T extends (...args: any[]) => any>(
  fn: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle = false

  return function (this: any, ...args: Parameters<T>) {
    if (!inThrottle) {
      fn.apply(this, args)
      inThrottle = true
      setTimeout(() => {
        inThrottle = false
      }, limit)
    }
  }
}
