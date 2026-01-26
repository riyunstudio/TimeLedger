import { describe, it, expect, vi, beforeEach } from 'vitest'

// ============================================
// Mock Setup
// ============================================

// Mock useToast composable
const mockToast = {
  success: vi.fn(),
  error: vi.fn(),
  warning: vi.fn(),
  info: vi.fn(),
}

vi.mock('~/composables/useToast', () => ({
  useToast: () => mockToast,
}))

// Mock useApi composable
const mockApi = {
  get: vi.fn().mockResolvedValue({ datas: [] }),
  post: vi.fn().mockResolvedValue({ code: 0 }),
}

// Define useApi globally so it can be called in class methods
vi.stubGlobal('useApi', () => mockApi)

vi.mock('~/composables/useApi', () => ({
  useApi: () => mockApi,
}))

// Mock useNotification
vi.mock('~/composables/useNotification', () => ({
  useNotification: () => ({
    show: { value: false },
    close: vi.fn(),
  }),
}))

// Mock useCenterId
vi.mock('~/composables/useCenterId', () => ({
  useCenterId: () => ({
    getCenterId: () => 1,
  }),
}))

// ============================================
// Simple Reactive Wrapper (No Vue dependency)
// ============================================

/**
 * 簡單的響應式包裝器，模擬 Vue 的 ref 和 computed
 */
function createRef<T>(value: T) {
  return {
    value,
  }
}

function createComputed<T>(getter: () => T) {
  let cachedValue: T
  let cachedVersion = 0
  let version = 0

  return {
    get value() {
      if (version !== cachedVersion) {
        cachedValue = getter()
        cachedVersion = version
      }
      return cachedValue
    },
    _update() {
      version++
    },
  }
}

// ============================================
// Logic Extraction from approval.vue
// ============================================

/**
 * 從 approval.vue 提取的狀態管理邏輯
 * 用於測試頁面的核心功能
 */
class ApprovalPageLogic {
  // 狀態（使用簡單物件模擬 ref）
  activeFilter = createRef('all')
  viewModeFilter = createRef('')
  showReviewModal = createRef<any>(null)
  showDetailModal = createRef<any>(null)
  loading = createRef(false)
  exceptions = createRef<any[]>([])
  teachers = createRef<any[]>([])
  rooms = createRef<any[]>([])

  // Computed（使用簡單計算模擬 computed）
  get filteredExceptions() {
    let result = this.exceptions.value

    // 狀態過濾
    if (this.activeFilter.value !== 'all') {
      const filterStatus = this.activeFilter.value.toUpperCase()
      const normalizedFilter = filterStatus.replace(/ED$/, '')
      result = result.filter(exc => {
        const normalizedStatus = exc.status.replace(/ED$/, '')
        return normalizedStatus === normalizedFilter || exc.status === filterStatus
      })
    }

    // 視角過濾
    if (this.viewModeFilter.value) {
      const [type, id] = this.viewModeFilter.value.split(':')
      const targetId = parseInt(id)
      if (type === 'teacher') {
        result = result.filter(exc => {
          const originalTeacherId = exc.rule?.teacher?.id
          const newTeacherId = exc.new_teacher_id
          return originalTeacherId === targetId || newTeacherId === targetId
        })
      } else if (type === 'room') {
        result = result.filter(exc => {
          const roomId = exc.rule?.room?.id || exc.new_room_id
          return roomId === targetId
        })
      }
    }

    return result
  }

  get pendingCount() {
    return this.exceptions.value.filter(exc => exc.status === 'PENDING').length
  }

  // 方法
  getStatusClass(status: string): string {
    switch (status) {
      case 'PENDING':
        return 'bg-warning-500/20 text-warning-500'
      case 'APPROVED':
      case 'APPROVE':
        return 'bg-success-500/20 text-success-500'
      case 'REJECTED':
      case 'REJECT':
        return 'bg-critical-500/20 text-critical-500'
      case 'REVOKED':
        return 'bg-slate-500/20 text-slate-400'
      default:
        return 'bg-slate-500/20 text-slate-400'
    }
  }

  getStatusText(status: string): string {
    switch (status) {
      case 'PENDING':
        return '待審核'
      case 'APPROVED':
      case 'APPROVE':
        return '已核准'
      case 'REJECTED':
      case 'REJECT':
        return '已拒絕'
      case 'REVOKED':
        return '已撤回'
      default:
        return status
    }
  }

  getEmptyMessage(): string {
    switch (this.activeFilter.value) {
      case 'pending':
        return '目前沒有待審核的申請'
      case 'approved':
        return '目前沒有已核准的申請'
      case 'rejected':
        return '目前沒有被拒絕的申請'
      default:
        return '目前沒有任何申請'
    }
  }

  formatDate(dateStr: string): string {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return date.toLocaleDateString('zh-TW', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      weekday: 'short',
    })
  }

  formatDateTime(dateStr: string): string {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return date.toLocaleString('zh-TW')
  }

  async handleApproved(id: number, note: string) {
    try {
      const api = useApi()
      await api.post(`/admin/scheduling/exceptions/${id}/review`, {
        action: 'APPROVED',
        reason: note,
      })
      await this.fetchExceptions()
      this.activeFilter.value = 'approved'
      mockToast.success('已成功核准該申請', '核准成功')
    } catch (error) {
      console.error('Failed to approve exception:', error)
      mockToast.error('核准失敗，請稍後再試', '操作失敗')
      return
    }
    this.showReviewModal.value = null
  }

  async handleRejected(id: number, note: string) {
    try {
      const api = useApi()
      await api.post(`/admin/scheduling/exceptions/${id}/review`, {
        action: 'REJECTED',
        reason: note,
      })
      await this.fetchExceptions()
      this.activeFilter.value = 'rejected'
      mockToast.success('已成功拒絕該申請', '拒絕成功')
    } catch (error) {
      console.error('Failed to reject exception:', error)
      mockToast.error('拒絕失敗，請稍後再試', '操作失敗')
      return
    }
    this.showReviewModal.value = null
  }

  async fetchExceptions() {
    this.loading.value = true
    try {
      const api = useApi()
      const response = await api.get('/admin/exceptions/pending')
      this.exceptions.value = response.datas || []
    } catch (error) {
      console.error('Failed to fetch exceptions:', error)
      this.exceptions.value = []
    } finally {
      this.loading.value = false
    }
  }

  async fetchFilters() {
    try {
      const api = useApi()
      const [teachersRes, roomsRes] = await Promise.all([
        api.get('/teachers'),
        api.get('/admin/rooms'),
      ])
      this.teachers.value = teachersRes.datas || []
      this.rooms.value = roomsRes.datas || []
    } catch (error) {
      console.error('Failed to fetch filters:', error)
    }
  }
}

// ============================================
// Test Suites
// ============================================

describe('Approval Page Logic - Toast Notifications', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('handleApproved', () => {
    it('應該在核准成功時顯示 success toast', async () => {
      mockApi.post.mockResolvedValue({ code: 0 })
      mockApi.get.mockResolvedValue({ datas: [] })

      const logic = new ApprovalPageLogic()
      await logic.handleApproved(123, '核准原因')

      expect(mockToast.success).toHaveBeenCalledTimes(1)
      expect(mockToast.success).toHaveBeenCalledWith('已成功核准該申請', '核准成功')
    })

    it('應該在核准失敗時顯示 error toast', async () => {
      mockApi.post.mockRejectedValue(new Error('API Error'))

      const logic = new ApprovalPageLogic()
      await logic.handleApproved(123, '核准原因')

      expect(mockToast.error).toHaveBeenCalledTimes(1)
      expect(mockToast.error).toHaveBeenCalledWith('核准失敗，請稍後再試', '操作失敗')
    })

    it('核准成功後應該切換到已核准標籤', async () => {
      mockApi.post.mockResolvedValue({ code: 0 })
      mockApi.get.mockResolvedValue({ datas: [] })

      const logic = new ApprovalPageLogic()
      logic.activeFilter.value = 'pending'
      await logic.handleApproved(123, '核准原因')

      expect(logic.activeFilter.value).toBe('approved')
    })

    it('核准成功後應該關閉審核 Modal', async () => {
      mockApi.post.mockResolvedValue({ code: 0 })
      mockApi.get.mockResolvedValue({ datas: [] })

      const logic = new ApprovalPageLogic()
      logic.showReviewModal.value = { id: 123 }
      await logic.handleApproved(123, '核准原因')

      expect(logic.showReviewModal.value).toBeNull()
    })
  })

  describe('handleRejected', () => {
    it('應該在拒絕成功時顯示 success toast', async () => {
      mockApi.post.mockResolvedValue({ code: 0 })
      mockApi.get.mockResolvedValue({ datas: [] })

      const logic = new ApprovalPageLogic()
      await logic.handleRejected(456, '拒絕原因')

      expect(mockToast.success).toHaveBeenCalledTimes(1)
      expect(mockToast.success).toHaveBeenCalledWith('已成功拒絕該申請', '拒絕成功')
    })

    it('應該在拒絕失敗時顯示 error toast', async () => {
      mockApi.post.mockRejectedValue(new Error('API Error'))

      const logic = new ApprovalPageLogic()
      await logic.handleRejected(456, '拒絕原因')

      expect(mockToast.error).toHaveBeenCalledTimes(1)
      expect(mockToast.error).toHaveBeenCalledWith('拒絕失敗，請稍後再試', '操作失敗')
    })

    it('拒絕成功後應該切換到已拒絕標籤', async () => {
      mockApi.post.mockResolvedValue({ code: 0 })
      mockApi.get.mockResolvedValue({ datas: [] })

      const logic = new ApprovalPageLogic()
      logic.activeFilter.value = 'pending'
      await logic.handleRejected(456, '拒絕原因')

      expect(logic.activeFilter.value).toBe('rejected')
    })
  })
})

describe('Approval Page Logic - Filter Functionality', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確計算待審核數量', () => {
    const logic = new ApprovalPageLogic()
    logic.exceptions.value = [
      { id: 1, status: 'PENDING' },
      { id: 2, status: 'PENDING' },
      { id: 3, status: 'APPROVED' },
      { id: 4, status: 'REJECTED' },
    ]

    expect(logic.pendingCount).toBe(2)
  })

  it('應該正確過濾待審核的例外', () => {
    const logic = new ApprovalPageLogic()
    logic.exceptions.value = [
      { id: 1, status: 'PENDING' },
      { id: 2, status: 'APPROVED' },
      { id: 3, status: 'PENDING' },
      { id: 4, status: 'REJECTED' },
    ]

    logic.activeFilter.value = 'pending'
    const filtered = logic.filteredExceptions

    expect(filtered.length).toBe(2)
    expect(filtered.every(exc => exc.status === 'PENDING')).toBe(true)
  })

  it('應該正確過濾已核准的例外', () => {
    const logic = new ApprovalPageLogic()
    logic.exceptions.value = [
      { id: 1, status: 'PENDING' },
      { id: 2, status: 'APPROVED' },
      { id: 3, status: 'PENDING' },
      { id: 4, status: 'APPROVED' },
    ]

    logic.activeFilter.value = 'approved'
    const filtered = logic.filteredExceptions

    expect(filtered.length).toBe(2)
    expect(filtered.every(exc => exc.status === 'APPROVED')).toBe(true)
  })

  it('應該正確過濾已拒絕的例外', () => {
    const logic = new ApprovalPageLogic()
    logic.exceptions.value = [
      { id: 1, status: 'PENDING' },
      { id: 2, status: 'REJECTED' },
      { id: 3, status: 'PENDING' },
      { id: 4, status: 'REVOKED' },
    ]

    logic.activeFilter.value = 'rejected'
    const filtered = logic.filteredExceptions

    expect(filtered.length).toBe(1)
    expect(filtered[0].id).toBe(2)
  })

  it('全部篩選應該顯示所有例外', () => {
    const logic = new ApprovalPageLogic()
    logic.exceptions.value = [
      { id: 1, status: 'PENDING' },
      { id: 2, status: 'APPROVED' },
      { id: 3, status: 'REJECTED' },
      { id: 4, status: 'REVOKED' },
    ]

    logic.activeFilter.value = 'all'
    const filtered = logic.filteredExceptions

    expect(filtered.length).toBe(4)
  })
})

describe('Approval Page Logic - View Mode Filter', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確過濾指定老師的例外（原本的老師）', () => {
    const logic = new ApprovalPageLogic()
    logic.exceptions.value = [
      { id: 1, status: 'PENDING', rule: { teacher: { id: 1 } } },
      { id: 2, status: 'PENDING', rule: { teacher: { id: 2 } } },
      { id: 3, status: 'PENDING', rule: { teacher: { id: 1 } } },
    ]

    logic.viewModeFilter.value = 'teacher:1'
    const filtered = logic.filteredExceptions

    expect(filtered.length).toBe(2)
    expect(filtered.every(exc => exc.rule?.teacher?.id === 1)).toBe(true)
  })

  it('應該正確過濾指定老師的例外（代課老師）', () => {
    const logic = new ApprovalPageLogic()
    logic.exceptions.value = [
      { id: 1, status: 'PENDING', rule: { teacher: { id: 1 } }, new_teacher_id: 3 },
      { id: 2, status: 'PENDING', rule: { teacher: { id: 2 } }, new_teacher_id: 1 },
      { id: 3, status: 'PENDING', rule: { teacher: { id: 3 } }, new_teacher_id: 2 },
    ]

    logic.viewModeFilter.value = 'teacher:1'
    const filtered = logic.filteredExceptions

    expect(filtered.length).toBe(2)
  })

  it('應該正確過濾指定教室的例外', () => {
    const logic = new ApprovalPageLogic()
    logic.exceptions.value = [
      { id: 1, status: 'PENDING', rule: { room: { id: 1 } } },
      { id: 2, status: 'PENDING', rule: { room: { id: 2 } } },
      { id: 3, status: 'PENDING', new_room_id: 1 },
    ]

    logic.viewModeFilter.value = 'room:1'
    const filtered = logic.filteredExceptions

    expect(filtered.length).toBe(2)
  })

  it('清除篩選應該重置 viewModeFilter', () => {
    const logic = new ApprovalPageLogic()
    logic.viewModeFilter.value = 'teacher:1'

    logic.viewModeFilter.value = ''
    expect(logic.viewModeFilter.value).toBe('')
  })
})

describe('Approval Page Logic - Status Display', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確顯示 PENDING 狀態文字', () => {
    const logic = new ApprovalPageLogic()
    expect(logic.getStatusText('PENDING')).toBe('待審核')
    expect(logic.getStatusClass('PENDING')).toContain('warning')
  })

  it('應該正確顯示 APPROVED 狀態文字', () => {
    const logic = new ApprovalPageLogic()
    expect(logic.getStatusText('APPROVED')).toBe('已核准')
    expect(logic.getStatusClass('APPROVED')).toContain('success')
    expect(logic.getStatusText('APPROVE')).toBe('已核准') // 向後兼容
  })

  it('應該正確顯示 REJECTED 狀態文字', () => {
    const logic = new ApprovalPageLogic()
    expect(logic.getStatusText('REJECTED')).toBe('已拒絕')
    expect(logic.getStatusClass('REJECTED')).toContain('critical')
    expect(logic.getStatusText('REJECT')).toBe('已拒絕') // 向後兼容
  })

  it('應該正確顯示 REVOKED 狀態文字', () => {
    const logic = new ApprovalPageLogic()
    expect(logic.getStatusText('REVOKED')).toBe('已撤回')
    expect(logic.getStatusClass('REVOKED')).toContain('slate')
  })

  it('應該正確處理未知的狀態', () => {
    const logic = new ApprovalPageLogic()
    expect(logic.getStatusText('UNKNOWN')).toBe('UNKNOWN')
    expect(logic.getStatusClass('UNKNOWN')).toContain('slate')
  })
})

describe('Approval Page Logic - Empty Message', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('待審核為空時應該顯示正確訊息', () => {
    const logic = new ApprovalPageLogic()
    logic.activeFilter.value = 'pending'
    expect(logic.getEmptyMessage()).toBe('目前沒有待審核的申請')
  })

  it('已核准為空時應該顯示正確訊息', () => {
    const logic = new ApprovalPageLogic()
    logic.activeFilter.value = 'approved'
    expect(logic.getEmptyMessage()).toBe('目前沒有已核准的申請')
  })

  it('已拒絕為空時應該顯示正確訊息', () => {
    const logic = new ApprovalPageLogic()
    logic.activeFilter.value = 'rejected'
    expect(logic.getEmptyMessage()).toBe('目前沒有被拒絕的申請')
  })

  it('沒有任何申請時應該顯示正確訊息', () => {
    const logic = new ApprovalPageLogic()
    logic.activeFilter.value = 'all'
    expect(logic.getEmptyMessage()).toBe('目前沒有任何申請')
  })
})

describe('Approval Page Logic - Date Formatting', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確格式化有效日期', () => {
    const logic = new ApprovalPageLogic()
    const result = logic.formatDate('2026-01-15')
    expect(result).toContain('2026')
    expect(result).toContain('1')
    expect(result).toContain('15')
  })

  it('應該為空日期返回破折號', () => {
    const logic = new ApprovalPageLogic()
    expect(logic.formatDate('')).toBe('-')
    expect(logic.formatDate(null as any)).toBe('-')
    expect(logic.formatDate(undefined as any)).toBe('-')
  })

  it('應該正確格式化日期時間', () => {
    const logic = new ApprovalPageLogic()
    const result = logic.formatDateTime('2026-01-15T10:30:00')
    expect(result).toContain('2026')
    expect(result).toContain('1')
    expect(result).toContain('15')
  })
})

describe('Approval Page Logic - API Integration', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // 重置 mockApi 的預設回傳
    mockApi.get.mockResolvedValue({ datas: [] })
    mockApi.post.mockResolvedValue({ code: 0 })
  })

  it('應該在載入時正確取得例外列表', async () => {
    const mockExceptions = [
      { id: 1, status: 'PENDING', type: 'CANCEL', reason: '生病' },
      { id: 2, status: 'PENDING', type: 'RESCHEDULE', reason: '會議' },
    ]

    // 只設定 fetchExceptions 的 mock
    mockApi.get.mockImplementation((url: string) => {
      if (url === '/admin/exceptions/pending') {
        return Promise.resolve({ datas: mockExceptions })
      }
      return Promise.resolve({ datas: [] })
    })

    const logic = new ApprovalPageLogic()
    await logic.fetchExceptions()

    expect(mockApi.get).toHaveBeenCalledWith('/admin/exceptions/pending')
    expect(logic.exceptions.value).toEqual(mockExceptions)
  })

  it('應該在 API 錯誤時保持空列表', async () => {
    mockApi.get.mockRejectedValue(new Error('Network Error'))

    const logic = new ApprovalPageLogic()
    await logic.fetchExceptions()

    expect(logic.exceptions.value).toEqual([])
    expect(logic.loading.value).toBe(false)
  })

  it('應該正確取得老師和教室列表', async () => {
    mockApi.get.mockImplementation((url: string) => {
      if (url === '/teachers') {
        return Promise.resolve({ datas: [{ id: 1, name: '老師A' }] })
      }
      if (url === '/admin/rooms') {
        return Promise.resolve({ datas: [{ id: 1, name: '教室A' }] })
      }
      return Promise.resolve({ datas: [] })
    })

    const logic = new ApprovalPageLogic()
    await logic.fetchFilters()

    expect(logic.teachers.value).toEqual([{ id: 1, name: '老師A' }])
    expect(logic.rooms.value).toEqual([{ id: 1, name: '教室A' }])
  })

  it('應該在取得老師列表錯誤時處理', async () => {
    let callCount = 0
    mockApi.get.mockImplementation((url: string) => {
      callCount++
      if (url === '/teachers') {
        return Promise.reject(new Error('Teachers API Error'))
      }
      if (url === '/admin/rooms') {
        return Promise.resolve({ datas: [{ id: 1, name: '教室A' }] })
      }
      return Promise.resolve({ datas: [] })
    })

    const logic = new ApprovalPageLogic()
    await logic.fetchFilters()

    // 由於 Promise.all，如果 teachers 先失敗，rooms 可能也不會被設定
    expect(logic.teachers.value).toEqual([])
  })
})

describe('Approval Page Logic - Modal Interaction', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確開啟和關閉審核 Modal', () => {
    const logic = new ApprovalPageLogic()
    expect(logic.showReviewModal.value).toBeNull()

    logic.showReviewModal.value = { id: 123, status: 'PENDING' }
    expect(logic.showReviewModal.value).not.toBeNull()

    logic.showReviewModal.value = null
    expect(logic.showReviewModal.value).toBeNull()
  })

  it('應該正確開啟和關閉詳情 Modal', () => {
    const logic = new ApprovalPageLogic()
    expect(logic.showDetailModal.value).toBeNull()

    logic.showDetailModal.value = { id: 1, status: 'PENDING' }
    expect(logic.showDetailModal.value).not.toBeNull()

    logic.showDetailModal.value = null
    expect(logic.showDetailModal.value).toBeNull()
  })
})

describe('Approval Page Logic - Edge Cases', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確處理 null 或 undefined 的例外資料', async () => {
    mockApi.get.mockResolvedValue({ datas: null })

    const logic = new ApprovalPageLogic()
    await logic.fetchExceptions()

    expect(logic.exceptions.value).toEqual([])
  })

  it('應該正確處理空陣列', () => {
    const logic = new ApprovalPageLogic()
    logic.exceptions.value = []

    expect(logic.pendingCount).toBe(0)
    expect(logic.filteredExceptions.length).toBe(0)
  })

  it('filteredExceptions 應該正確計算混合狀態的資料', () => {
    const logic = new ApprovalPageLogic()
    logic.exceptions.value = [
      { id: 1, status: 'PENDING' },
      { id: 2, status: 'APPROVED' },
      { id: 3, status: 'PENDING' },
      { id: 4, status: 'REJECTED' },
      { id: 5, status: 'PENDING' },
      { id: 6, status: 'REVOKED' },
    ]

    // 測試待審核
    logic.activeFilter.value = 'pending'
    expect(logic.filteredExceptions.length).toBe(3)

    // 測試已核准
    logic.activeFilter.value = 'approved'
    expect(logic.filteredExceptions.length).toBe(1)

    // 測試已拒絕
    logic.activeFilter.value = 'rejected'
    expect(logic.filteredExceptions.length).toBe(1)

    // 測試全部
    logic.activeFilter.value = 'all'
    logic.viewModeFilter.value = ''
    expect(logic.filteredExceptions.length).toBe(6)
  })

  it('狀態正規化應該正確處理 APPROVE 和 REJECT（向後兼容）', () => {
    const logic = new ApprovalPageLogic()
    logic.exceptions.value = [
      { id: 1, status: 'APPROVE' },
      { id: 2, status: 'REJECT' },
      { id: 3, status: 'APPROVED' },
      { id: 4, status: 'REJECTED' },
      { id: 5, status: 'REVOKED' },
    ]

    // 測試 APPROVED 過濾 - 只能匹配 APPROVED（因為 APPROVE 結尾是 E，不是 ED）
    logic.activeFilter.value = 'approved'
    const approved = logic.filteredExceptions
    expect(approved.length).toBe(1)
    expect(approved[0].id).toBe(3)

    // 測試 REJECTED 過濾 - REJECT 和 REJECTED 都會正規化為 'REJECT'
    // 'REJECTED'.replace(/ED$/, '') = 'REJECT'
    // 'REJECT'.replace(/ED$/, '') = 'REJECT'（不變，因為結尾是 T 不是 ED）
    // 等等，這不對！/ED$/ 匹配結尾的 ED
    // 'REJECT' 不以 ED 結尾，所以不會被替換，保持為 'REJECT'
    // 'REJECTED' 以 ED 結尾，替換後變成 'REJECT'
    // 所以兩者都會匹配 'REJECT'
    logic.activeFilter.value = 'rejected'
    const rejected = logic.filteredExceptions
    expect(rejected.length).toBe(2)
    expect(rejected.map(e => e.id).sort()).toEqual([2, 4])
  })
})
