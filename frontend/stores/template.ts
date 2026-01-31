import { defineStore } from 'pinia'

// ============================================
// 狀態類型定義
// ============================================

/**
 * 通用分頁參數
 */
export interface PaginationParams {
  page?: number
  limit?: number
  sort_by?: string
  sort_order?: 'ASC' | 'DESC'
}

/**
 * 通用分頁結果
 */
export interface PaginationResult {
  page: number
  limit: number
  total: number
  total_pages: number
  has_next: boolean
  has_prev: boolean
}

/**
 * API 響應基礎結構
 */
export interface ApiResponse<T> {
  code: number
  message: string
  datas?: T
  pagination?: PaginationResult
}

/**
 * 錯誤處理函式類型
 */
export type ErrorHandler = (error: unknown) => string

// ============================================
// Store 模板類型定義
// ============================================

/**
 * 標準 Store 狀態介面
 */
export interface BaseState<TItem> {
  // 資料
  items: TItem[]
  // 單一項目（用於編輯/詳情）
  item: TItem | null
  // 分頁資訊
  pagination: PaginationResult | null
  // 載入狀態
  loading: boolean
  isFetching: boolean
  // 操作狀態
  isCreating: boolean
  isUpdating: boolean
  isDeleting: boolean
  // 錯誤
  error: string | null
}

/**
 * 標準 Store 初始狀態
 */
export const createInitialState = <TItem>(): BaseState<TItem> => ({
  items: [],
  item: null,
  pagination: null,
  loading: false,
  isFetching: false,
  isCreating: false,
  isUpdating: false,
  isDeleting: false,
  error: null
})

/**
 * Store 選項介面
 */
export interface StoreOptions<TItem, TCreateDto, TUpdateDto> {
  /** API 端點基礎路徑 */
  endpoint: string
  /** 取得錯誤訊息的處理函式 */
  getErrorMessage?: ErrorHandler
}

// ============================================
// Store 模板工廠函式
// ============================================

/**
 * 建立標準化的 Pinia Store
 *
 * @example
 * ```typescript
 * interface User {
 *   id: number
 *   name: string
 *   email: string
 * }
 *
 * interface CreateUserDto {
 *   name: string
 *   email: string
 * }
 *
 * interface UpdateUserDto extends Partial<CreateUserDto> {}
 *
 * const useUserStore = createStoreTemplate<User, CreateUserDto, UpdateUserDto>({
 *   endpoint: '/users',
 *   getErrorMessage: (error) => getErrorMessage(error)
 * })
 * ```
 */
export function createStoreTemplate<
  TItem,
  TCreateDto,
  TUpdateDto
>(options: StoreOptions<TItem, TCreateDto, TUpdateDto>) {
  const { endpoint, getErrorMessage = (err: unknown) => String(err) } = options

  return defineStore(`store-${endpoint.replace(/\//g, '-')}`, () => {
    // 初始狀態
    const state = reactive(createInitialState<TItem>())

    // ============================================
    // Getters（計算屬性）
    // ============================================

    /**
     * 項目數量
     */
    const itemCount = computed(() => state.items.length)

    /**
     * 是否有資料
     */
    const hasItems = computed(() => state.items.length > 0)

    /**
     * 是否有載入中
     */
    const isLoading = computed(() => state.loading || state.isFetching)

    /**
     * 取得項目（安全取得，無資料時回傳 undefined）
     */
    const firstItem = computed(() => state.items[0] || null)

    /**
     * 是否有錯誤
     */
    const hasError = computed(() => state.error !== null)

    // ============================================
    // Actions（操作方法）
    // ============================================

    /**
     * 取得單一項目
     */
    async function fetchItem(id: number | string) {
      state.isFetching = true
      state.error = null

      try {
        const api = useApi()
        const response = await api.get<ApiResponse<TItem>>(`${endpoint}/${id}`)
        state.item = response.datas || null
        return state.item
      } catch (error) {
        state.error = getErrorMessage(error)
        throw error
      } finally {
        state.isFetching = false
      }
    }

    /**
     * 取得列表（可選分頁）
     */
    async function fetchItems(params?: PaginationParams) {
      state.loading = true
      state.error = null

      try {
        const api = useApi()
        const response = await api.get<ApiResponse<TItem[]>>(`${endpoint}`, { params })

        state.items = response.datas || []
        state.pagination = response.pagination || null

        return { items: state.items, pagination: state.pagination }
      } catch (error) {
        state.error = getErrorMessage(error)
        throw error
      } finally {
        state.loading = false
      }
    }

    /**
     * 建立新項目
     */
    async function createItem(data: TCreateDto): Promise<TItem> {
      state.isCreating = true
      state.error = null

      try {
        const api = useApi()
        const response = await api.post<ApiResponse<TItem>>(endpoint, data)
        const newItem = response.datas!

        state.items.push(newItem)
        return newItem
      } catch (error) {
        state.error = getErrorMessage(error)
        throw error
      } finally {
        state.isCreating = false
      }
    }

    /**
     * 更新項目
     */
    async function updateItem(id: number | string, data: TUpdateDto): Promise<TItem> {
      state.isUpdating = true
      state.error = null

      try {
        const api = useApi()
        const response = await api.put<ApiResponse<TItem>>(`${endpoint}/${id}`, data)
        const updatedItem = response.datas!

        // 更新列表中的項目
        const index = state.items.findIndex((item: any) => item.id === id)
        if (index !== -1) {
          state.items[index] = updatedItem
        }

        // 更新單一項目（如果相同）
        if (state.item && (state.item as any).id === id) {
          state.item = updatedItem
        }

        return updatedItem
      } catch (error) {
        state.error = getErrorMessage(error)
        throw error
      } finally {
        state.isUpdating = false
      }
    }

    /**
     * 刪除項目
     */
    async function deleteItem(id: number | string): Promise<void> {
      state.isDeleting = true
      state.error = null

      try {
        const api = useApi()
        await api.delete(`${endpoint}/${id}`)

        // 從列表中移除
        state.items = state.items.filter((item: any) => item.id !== id)

        // 如果刪除的是當前項目，清除它
        if (state.item && (state.item as any).id === id) {
          state.item = null
        }
      } catch (error) {
        state.error = getErrorMessage(error)
        throw error
      } finally {
        state.isDeleting = false
      }
    }

    /**
     * 清除錯誤
     */
    function clearError() {
      state.error = null
    }

    /**
     * 重置狀態
     */
    function reset() {
      Object.assign(state, createInitialState<TItem>())
    }

    /**
     * 設定單一項目（用於從外部設定）
     */
    function setItem(item: TItem | null) {
      state.item = item
    }

    /**
     * 批量更新列表
     */
    function setItems(items: TItem[]) {
      state.items = items
    }

    /**
     * 手動設定分頁資訊
     */
    function setPagination(pagination: PaginationResult | null) {
      state.pagination = pagination
    }

    // ============================================
    // 回傳介面
    // ============================================

    return {
      // 狀態（使用 toRefs 保持響應式）
      ...toRefs(state),

      // Getters
      itemCount,
      hasItems,
      isLoading,
      firstItem,
      hasError,

      // Actions
      fetchItem,
      fetchItems,
      createItem,
      updateItem,
      deleteItem,
      clearError,
      reset,
      setItem,
      setItems,
      setPagination
    }
  })
}

// ============================================
// 簡化版 Store（適用於簡單的唯讀列表）
// ============================================

/**
 * 簡化列表 Store 選項
 */
export interface SimpleListOptions<TItem> {
  /** API 端點基礎路徑 */
  endpoint: string
  /** 取得錯誤訊息的處理函式 */
  getErrorMessage?: ErrorHandler
  /** 額外的 fetch 參數 */
  defaultParams?: Record<string, any>
}

/**
 * 建立簡化的唯讀列表 Store
 *
 * @example
 * ```typescript
 * const useProductListStore = createSimpleListStore<Product>({
 *   endpoint: '/products',
 *   defaultParams: { status: 'active' }
 * })
 * ```
 */
export function createSimpleListStore<TItem>(options: SimpleListOptions<TItem>) {
  const { endpoint, getErrorMessage = (err: unknown) => String(err), defaultParams } = options

  return defineStore(`simple-list-${endpoint.replace(/\//g, '-')}`, () => {
    // 狀態
    const items = ref<TItem[]>([])
    const loading = ref(false)
    const error = ref<string | null>(null)
    const params = ref<Record<string, any>>(defaultParams || {})

    // Getters
    const hasItems = computed(() => items.value.length > 0)
    const count = computed(() => items.value.length)

    // Actions
    async function fetch() {
      loading.value = true
      error.value = null

      try {
        const api = useApi()
        const response = await api.get<ApiResponse<T[]>>(endpoint, { params: params.value })
        items.value = response.datas || []
      } catch (err) {
        error.value = getErrorMessage(err)
        throw err
      } finally {
        loading.value = false
      }
    }

    async function fetchById(id: number | string) {
      loading.value = true
      error.value = null

      try {
        const api = useApi()
        const response = await api.get<ApiResponse<TItem>>(`${endpoint}/${id}`)
        return response.datas || null
      } catch (err) {
        error.value = getErrorMessage(err)
        throw err
      } finally {
        loading.value = false
      }
    }

    function setParams(newParams: Record<string, any>) {
      params.value = { ...params.value, ...newParams }
    }

    function clearParams() {
      params.value = defaultParams || {}
    }

    function clearItems() {
      items.value = []
    }

    function clearError() {
      error.value = null
    }

    function reset() {
      items.value = []
      loading.value = false
      error.value = null
      params.value = defaultParams || {}
    }

    return {
      // 狀態
      items: readonly(items),
      loading: readonly(loading),
      error: readonly(error),
      params: readonly(params),

      // Getters
      hasItems,
      count,

      // Actions
      fetch,
      fetchById,
      setParams,
      clearParams,
      clearItems,
      clearError,
      reset
    }
  })
}
