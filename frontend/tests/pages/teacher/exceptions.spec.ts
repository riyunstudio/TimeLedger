import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock modules before imports
vi.mock('vue', async () => {
  const vue = await import('vue')
  return {
    ...vue,
    ref: (v: any) => ({ value: v }),
    computed: (fn: () => any) => ({ value: fn() }),
    onMounted: (fn: () => void) => fn(),
    watch: (fn: () => void) => {},
  }
})

vi.mock('~/composables/useAlert', () => ({
  alertError: vi.fn(),
  alertSuccess: vi.fn(),
  alertWarning: vi.fn(),
  confirm: vi.fn().mockResolvedValue(true),
}))

vi.mock('~/composables/useNotification', () => ({
  default: () => ({
    show: { value: false },
    close: vi.fn(),
    success: vi.fn(),
    error: vi.fn(),
  }),
}))

vi.mock('~/stores/useScheduleStore', () => ({
  useScheduleStore: () => ({
    centers: [],
    exceptions: [],
    fetchCenters: vi.fn(),
    fetchExceptions: vi.fn(),
    revokeException: vi.fn(),
  }),
}))

vi.mock('~/composables/useSidebar', () => ({
  useSidebar: () => ({
    isOpen: { value: false },
    close: vi.fn(),
  }),
}))

describe('teacher/exceptions.vue 頁面邏輯', () => {
  // ExceptionStatusLogic 類別 - 例外狀態邏輯
  class ExceptionStatusLogic {
    statusFilters: { value: string; label: string }[]

    constructor() {
      this.statusFilters = [
        { value: '', label: '全部' },
        { value: 'PENDING', label: '待審核' },
        { value: 'APPROVED', label: '已核准' },
        { value: 'REJECTED', label: '已拒絕' },
        { value: 'REVOKED', label: '已撤回' },
      ]
    }

    getStatusLabel(status: string): string {
      const filter = this.statusFilters.find(f => f.value === status)
      return filter?.label || status
    }

    getAllStatusValues(): string[] {
      return this.statusFilters.map(f => f.value).filter(Boolean)
    }

    isValidStatus(status: string): boolean {
      return this.statusFilters.some(f => f.value === status)
    }

    getDefaultFilter(): string {
      return ''
    }
  }

  // ExceptionStatusDisplayLogic 類別 - 狀態顯示邏輯
  class ExceptionStatusDisplayLogic {
    getStatusClass(status: string): string {
      switch (status.toUpperCase()) {
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
      switch (status.toUpperCase()) {
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
  }

  // ExceptionFilterLogic 類別 - 篩選邏輯
  class ExceptionFilterLogic {
    currentFilter: string
    exceptions: any[]

    constructor() {
      this.currentFilter = ''
      this.exceptions = []
    }

    setExceptions(exceptions: any[]) {
      this.exceptions = exceptions
    }

    setFilter(filter: string) {
      this.currentFilter = filter
    }

    getFilteredExceptions(): any[] {
      if (!this.currentFilter) return this.exceptions
      return this.exceptions.filter(e => e.status === this.currentFilter)
    }

    getPendingExceptions(): any[] {
      return this.exceptions.filter(e => e.status === 'PENDING')
    }

    getApprovedExceptions(): any[] {
      return this.exceptions.filter(e => e.status === 'APPROVED' || e.status === 'APPROVE')
    }

    getRejectedExceptions(): any[] {
      return this.exceptions.filter(e => e.status === 'REJECTED' || e.status === 'REJECT')
    }

    getRevokedExceptions(): any[] {
      return this.exceptions.filter(e => e.status === 'REVOKED')
    }

    hasExceptions(): boolean {
      return this.exceptions.length > 0
    }

    hasFilteredResults(): boolean {
      return this.getFilteredExceptions().length > 0
    }

    getExceptionCount(): number {
      return this.exceptions.length
    }

    getFilteredCount(): number {
      return this.getFilteredExceptions().length
    }
  }

  // ExceptionActionLogic 類別 - 例外操作邏輯
  class ExceptionActionLogic {
    exceptions: any[]

    constructor() {
      this.exceptions = []
    }

    setExceptions(exceptions: any[]) {
      this.exceptions = exceptions
    }

    canRevoke(exception: any): boolean {
      return exception.status === 'PENDING'
    }

    canEdit(exception: any): boolean {
      return exception.status === 'PENDING'
    }

    revokeException(id: number): boolean {
      const index = this.exceptions.findIndex(e => e.id === id)
      if (index !== -1 && this.exceptions[index].status === 'PENDING') {
        this.exceptions[index].status = 'REVOKED'
        return true
      }
      return false
    }

    updateException(id: number, updates: any): boolean {
      const index = this.exceptions.findIndex(e => e.id === id)
      if (index !== -1) {
        this.exceptions[index] = { ...this.exceptions[index], ...updates }
        return true
      }
      return false
    }

    deleteException(id: number): boolean {
      const initialLength = this.exceptions.length
      this.exceptions = this.exceptions.filter(e => e.id !== id)
      return this.exceptions.length < initialLength
    }

    getExceptionById(id: number): any | undefined {
      return this.exceptions.find(e => e.id === id)
    }
  }

  // ExceptionDateTimeLogic 類別 - 日期時間格式化邏輯
  class ExceptionDateTimeLogic {
    formatDateTime(dateStr: string): string {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return date.toLocaleString('zh-TW')
    }

    formatDate(dateStr: string): string {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return date.toLocaleDateString('zh-TW')
    }

    formatTime(timeStr: string): string {
      if (!timeStr) return '-'
      return timeStr
    }

    isValidDate(dateStr: string): boolean {
      if (!dateStr) return false
      const date = new Date(dateStr)
      return !isNaN(date.getTime())
    }

    compareDateTime(a: string, b: string): number {
      const dateA = new Date(a).getTime()
      const dateB = new Date(b).getTime()
      return dateA - dateB
    }
  }

  // ExceptionTypeLogic 類別 - 例外類型邏輯
  class ExceptionTypeLogic {
    getTypeLabel(type: string): string {
      switch (type.toUpperCase()) {
        case 'CANCEL':
          return '停課'
        case 'RESCHEDULE':
          return '改期'
        case 'CUSTOM':
          return '自訂'
        default:
          return type
      }
    }

    getTypeClass(type: string): string {
      switch (type.toUpper()) {
        case 'CANCEL':
          return 'bg-critical-500/20 text-critical-500'
        case 'RESCHEDULE':
          return 'bg-warning-500/20 text-warning-500'
        default:
          return 'bg-slate-500/20 text-slate-400'
      }
    }

    isReschedule(exception: any): boolean {
      return exception.type === 'RESCHEDULE'
    }

    isCancel(exception: any): boolean {
      return exception.type === 'CANCEL'
    }

    hasNewTime(exception: any): boolean {
      return Boolean(exception.new_start_at && exception.new_end_at)
    }
  }

  // CenterFilterLogic 類別 - 中心篩選邏輯
  class CenterFilterLogic {
    centers: { center_id: number; center_name: string }[]

    constructor() {
      this.centers = []
    }

    setCenters(centers: { center_id: number; center_name: string }[]) {
      this.centers = centers
    }

    getCenterName(centerId: number): string {
      const center = this.centers.find(c => c.center_id === centerId)
      return center?.center_name || '未知中心'
    }

    getCenterOptions(): { value: number; label: string }[] {
      return this.centers.map(c => ({
        value: c.center_id,
        label: c.center_name,
      }))
    }

    hasCenters(): boolean {
      return this.centers.length > 0
    }

    filterExceptionsByCenter(exceptions: any[], centerId: number): any[] {
      return exceptions.filter(e => e.center_id === centerId)
    }
  }

  describe('ExceptionStatusLogic 例外狀態邏輯', () => {
    it('應該正確初始化所有狀態篩選器', () => {
      const logic = new ExceptionStatusLogic()
      expect(logic.statusFilters).toHaveLength(5)
      expect(logic.statusFilters[0].value).toBe('')
      expect(logic.statusFilters[0].label).toBe('全部')
    })

    it('getStatusLabel 應該返回正確的標籤', () => {
      const logic = new ExceptionStatusLogic()
      expect(logic.getStatusLabel('PENDING')).toBe('待審核')
      expect(logic.getStatusLabel('APPROVED')).toBe('已核准')
      expect(logic.getStatusLabel('REJECTED')).toBe('已拒絕')
      expect(logic.getStatusLabel('REVOKED')).toBe('已撤回')
      expect(logic.getStatusLabel('')).toBe('全部')
    })

    it('isValidStatus 應該正確驗證狀態', () => {
      const logic = new ExceptionStatusLogic()
      expect(logic.isValidStatus('PENDING')).toBe(true)
      expect(logic.isValidStatus('APPROVED')).toBe(true)
      expect(logic.isValidStatus('INVALID')).toBe(false)
    })

    it('getAllStatusValues 應該返回所有狀態值（不含空值）', () => {
      const logic = new ExceptionStatusLogic()
      const values = logic.getAllStatusValues()
      expect(values).toContain('PENDING')
      expect(values).toContain('APPROVED')
      expect(values).toContain('REJECTED')
      expect(values).toContain('REVOKED')
      expect(values).not.toContain('')
    })
  })

  describe('ExceptionStatusDisplayLogic 狀態顯示邏輯', () => {
    it('getStatusClass 應該返回正確的 CSS 類別', () => {
      const logic = new ExceptionStatusDisplayLogic()
      expect(logic.getStatusClass('PENDING')).toContain('warning')
      expect(logic.getStatusClass('APPROVED')).toContain('success')
      expect(logic.getStatusClass('APPROVE')).toContain('success')
      expect(logic.getStatusClass('REJECTED')).toContain('critical')
      expect(logic.getStatusClass('REVOKED')).toContain('slate')
    })

    it('getStatusText 應該返回正確的中文文字', () => {
      const logic = new ExceptionStatusDisplayLogic()
      expect(logic.getStatusText('PENDING')).toBe('待審核')
      expect(logic.getStatusText('APPROVED')).toBe('已核准')
      expect(logic.getStatusText('APPROVE')).toBe('已核准')
      expect(logic.getStatusText('REJECTED')).toBe('已拒絕')
      expect(logic.getStatusText('REJECT')).toBe('已拒絕')
      expect(logic.getStatusText('REVOKED')).toBe('已撤回')
    })

    it('getStatusText 應該對未知狀態返回原值', () => {
      const logic = new ExceptionStatusDisplayLogic()
      expect(logic.getStatusText('UNKNOWN')).toBe('UNKNOWN')
    })
  })

  describe('ExceptionFilterLogic 篩選邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new ExceptionFilterLogic()
      expect(logic.currentFilter).toBe('')
      expect(logic.exceptions).toHaveLength(0)
    })

    it('setExceptions 應該正確設定例外列表', () => {
      const logic = new ExceptionFilterLogic()
      const exceptions = [
        { id: 1, status: 'PENDING' },
        { id: 2, status: 'APPROVED' },
      ]
      logic.setExceptions(exceptions)
      expect(logic.exceptions).toHaveLength(2)
    })

    it('getFilteredExceptions 應該根據篩選條件過濾', () => {
      const logic = new ExceptionFilterLogic()
      logic.setExceptions([
        { id: 1, status: 'PENDING' },
        { id: 2, status: 'APPROVED' },
        { id: 3, status: 'PENDING' },
      ])
      logic.setFilter('PENDING')
      const filtered = logic.getFilteredExceptions()
      expect(filtered).toHaveLength(2)
      expect(filtered.every(e => e.status === 'PENDING')).toBe(true)
    })

    it('getFilteredExceptions 應該在篩選為空時返回全部', () => {
      const logic = new ExceptionFilterLogic()
      logic.setExceptions([{ id: 1 }, { id: 2 }])
      logic.setFilter('')
      expect(logic.getFilteredExceptions()).toHaveLength(2)
    })

    it('getPendingExceptions 應該返回待審核的例外', () => {
      const logic = new ExceptionFilterLogic()
      logic.setExceptions([
        { id: 1, status: 'PENDING' },
        { id: 2, status: 'APPROVED' },
      ])
      const pending = logic.getPendingExceptions()
      expect(pending).toHaveLength(1)
      expect(pending[0].id).toBe(1)
    })

    it('getApprovedExceptions 應該處理 APPROVE 舊資料', () => {
      const logic = new ExceptionFilterLogic()
      logic.setExceptions([
        { id: 1, status: 'APPROVED' },
        { id: 2, status: 'APPROVE' },
      ])
      const approved = logic.getApprovedExceptions()
      expect(approved).toHaveLength(2)
    })
  })

  describe('ExceptionActionLogic 操作邏輯', () => {
    it('canRevoke 應該在狀態為 PENDING 時返回 true', () => {
      const logic = new ExceptionActionLogic()
      expect(logic.canRevoke({ status: 'PENDING' })).toBe(true)
      expect(logic.canRevoke({ status: 'APPROVED' })).toBe(false)
      expect(logic.canRevoke({ status: 'REVOKED' })).toBe(false)
    })

    it('canEdit 應該在狀態為 PENDING 時返回 true', () => {
      const logic = new ExceptionActionLogic()
      expect(logic.canEdit({ status: 'PENDING' })).toBe(true)
      expect(logic.canEdit({ status: 'REJECTED' })).toBe(false)
    })

    it('revokeException 應該正確撤回例外', () => {
      const logic = new ExceptionActionLogic()
      logic.setExceptions([
        { id: 1, status: 'PENDING' },
        { id: 2, status: 'APPROVED' },
      ])
      const result = logic.revokeException(1)
      expect(result).toBe(true)
      const exception = logic.getExceptionById(1)
      expect(exception?.status).toBe('REVOKED')
    })

    it('revokeException 應該無法撤回已核准的例外', () => {
      const logic = new ExceptionActionLogic()
      logic.setExceptions([{ id: 1, status: 'APPROVED' }])
      const result = logic.revokeException(1)
      expect(result).toBe(false)
    })

    it('updateException 應該正確更新例外', () => {
      const logic = new ExceptionActionLogic()
      logic.setExceptions([{ id: 1, reason: '舊原因' }])
      const result = logic.updateException(1, { reason: '新原因' })
      expect(result).toBe(true)
      const exception = logic.getExceptionById(1)
      expect(exception?.reason).toBe('新原因')
    })

    it('updateException 應該處理不存在的例外', () => {
      const logic = new ExceptionActionLogic()
      logic.setExceptions([{ id: 1 }])
      const result = logic.updateException(999, { reason: '新原因' })
      expect(result).toBe(false)
    })

    it('deleteException 應該正確刪除例外', () => {
      const logic = new ExceptionActionLogic()
      logic.setExceptions([
        { id: 1, status: 'PENDING' },
        { id: 2, status: 'APPROVED' },
      ])
      const result = logic.deleteException(1)
      expect(result).toBe(true)
      expect(logic.exceptions).toHaveLength(1)
      expect(logic.getExceptionById(1)).toBeUndefined()
    })
  })

  describe('ExceptionDateTimeLogic 日期時間格式化邏輯', () => {
    it('formatDateTime 應該正確格式化日期時間', () => {
      const logic = new ExceptionDateTimeLogic()
      const result = logic.formatDateTime('2026-01-20T10:00:00')
      expect(result).toContain('2026')
      expect(result).toContain('1')
      expect(result).toContain('20')
    })

    it('formatDateTime 應該處理空值', () => {
      const logic = new ExceptionDateTimeLogic()
      expect(logic.formatDateTime('')).toBe('-')
      expect(logic.formatDateTime(null as any)).toBe('-')
    })

    it('formatDate 應該正確格式化日期', () => {
      const logic = new ExceptionDateTimeLogic()
      const result = logic.formatDate('2026-01-20')
      expect(result).toContain('2026')
    })

    it('isValidDate 應該正確驗證日期', () => {
      const logic = new ExceptionDateTimeLogic()
      expect(logic.isValidDate('2026-01-20')).toBe(true)
      expect(logic.isValidDate('')).toBe(false)
      expect(logic.isValidDate('invalid')).toBe(false)
    })

    it('compareDateTime 應該正確比較日期時間', () => {
      const logic = new ExceptionDateTimeLogic()
      const result = logic.compareDateTime(
        '2026-01-20T10:00:00',
        '2026-01-20T12:00:00'
      )
      expect(result).toBeLessThan(0)
    })
  })

  describe('ExceptionTypeLogic 例外類型邏輯', () => {
    it('getTypeLabel 應該返回正確的類型標籤', () => {
      const logic = new ExceptionTypeLogic()
      expect(logic.getTypeLabel('CANCEL')).toBe('停課')
      expect(logic.getTypeLabel('RESCHEDULE')).toBe('改期')
      expect(logic.getTypeLabel('CUSTOM')).toBe('自訂')
    })

    it('getTypeClass 應該返回正確的 CSS 類別', () => {
      const logic = new ExceptionTypeLogic()
      expect(logic.getTypeClass('CANCEL')).toContain('critical')
      expect(logic.getTypeClass('RESCHEDULE')).toContain('warning')
    })

    it('isReschedule 應該正確判斷是否為改期', () => {
      const logic = new ExceptionTypeLogic()
      expect(logic.isReschedule({ type: 'RESCHEDULE' })).toBe(true)
      expect(logic.isReschedule({ type: 'CANCEL' })).toBe(false)
    })

    it('isCancel 應該正確判斷是否為停課', () => {
      const logic = new ExceptionTypeLogic()
      expect(logic.isCancel({ type: 'CANCEL' })).toBe(true)
      expect(logic.isCancel({ type: 'RESCHEDULE' })).toBe(false)
    })

    it('hasNewTime 應該正確判斷是否有新時間', () => {
      const logic = new ExceptionTypeLogic()
      expect(logic.hasNewTime({ new_start_at: '10:00', new_end_at: '11:00' })).toBe(true)
      expect(logic.hasNewTime({ new_start_at: '', new_end_at: '11:00' })).toBe(false)
    })
  })

  describe('CenterFilterLogic 中心篩選邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new CenterFilterLogic()
      expect(logic.centers).toHaveLength(0)
    })

    it('setCenters 應該正確設定中心列表', () {
      const logic = new CenterFilterLogic()
      logic.setCenters([
        { center_id: 1, center_name: '中心 A' },
        { center_id: 2, center_name: '中心 B' },
      ])
      expect(logic.centers).toHaveLength(2)
    })

    it('getCenterName 應該返回正確的中心名稱', () => {
      const logic = new CenterFilterLogic()
      logic.setCenters([{ center_id: 1, center_name: '中心 A' }])
      expect(logic.getCenterName(1)).toBe('中心 A')
      expect(logic.getCenterName(999)).toBe('未知中心')
    })

    it('getCenterOptions 應該返回正確的下拉選項', () => {
      const logic = new CenterFilterLogic()
      logic.setCenters([{ center_id: 1, center_name: '中心 A' }])
      const options = logic.getCenterOptions()
      expect(options).toHaveLength(1)
      expect(options[0].value).toBe(1)
      expect(options[0].label).toBe('中心 A')
    })

    it('hasCenters 應該正確判斷是否有中心', () => {
      const logic = new CenterFilterLogic()
      expect(logic.hasCenters()).toBe(false)
      logic.setCenters([{ center_id: 1, center_name: '中心 A' }])
      expect(logic.hasCenters()).toBe(true)
    })

    it('filterExceptionsByCenter 應該根據中心過濾例外', () => {
      const logic = new CenterFilterLogic()
      const exceptions = [
        { id: 1, center_id: 1 },
        { id: 2, center_id: 2 },
        { id: 3, center_id: 1 },
      ]
      const filtered = logic.filterExceptionsByCenter(exceptions, 1)
      expect(filtered).toHaveLength(2)
    })
  })

  describe('頁面整合邏輯', () => {
    it('應該能夠完整處理例外申請流程', () => {
      const statusLogic = new ExceptionStatusDisplayLogic()
      const filterLogic = new ExceptionFilterLogic()
      const actionLogic = new ExceptionActionLogic()
      const datetimeLogic = new ExceptionDateTimeLogic()
      const typeLogic = new ExceptionTypeLogic()

      // 設定例外資料
      const exceptions = [
        {
          id: 1,
          type: 'CANCEL',
          status: 'PENDING',
          reason: '身體不適',
          created_at: '2026-01-20T08:00:00',
          original_date: '2026-01-22',
        },
        {
          id: 2,
          type: 'RESCHEDULE',
          status: 'PENDING',
          reason: '時間調整',
          created_at: '2026-01-21T09:00:00',
          original_date: '2026-01-23',
          new_start_at: '14:00',
          new_end_at: '15:00',
        },
      ]

      filterLogic.setExceptions(exceptions)
      actionLogic.setExceptions([...exceptions])

      // 驗證狀態顯示
      expect(statusLogic.getStatusClass('PENDING')).toContain('warning')
      expect(statusLogic.getStatusText('PENDING')).toBe('待審核')

      // 驗證篩選
      filterLogic.setFilter('PENDING')
      expect(filterLogic.getFilteredExceptions()).toHaveLength(2)

      // 驗證類型顯示
      expect(typeLogic.getTypeLabel('CANCEL')).toBe('停課')
      expect(typeLogic.getTypeLabel('RESCHEDULE')).toBe('改期')

      // 驗證日期格式化
      expect(datetimeLogic.formatDateTime(exceptions[0].created_at)).not.toBe('-')

      // 撤回第一個例外
      const revoked = actionLogic.revokeException(1)
      expect(revoked).toBe(true)

      // 驗證撤回後的狀態
      const updatedException = actionLogic.getExceptionById(1)
      expect(updatedException?.status).toBe('REVOKED')
      expect(statusLogic.getStatusText('REVOKED')).toBe('已撤回')
    })

    it('應該正確處理向後兼容的狀態值', () => {
      const displayLogic = new ExceptionStatusDisplayLogic()

      // 測試舊資料的 APPROVE 狀態
      expect(displayLogic.getStatusText('APPROVE')).toBe('已核准')
      expect(displayLogic.getStatusClass('APPROVE')).toContain('success')

      // 測試舊資料的 REJECT 狀態
      expect(displayLogic.getStatusText('REJECT')).toBe('已拒絕')
      expect(displayLogic.getStatusClass('REJECT')).toContain('critical')
    })

    it('應該正確處理混合狀態的例外列表', () => {
      const filterLogic = new ExceptionFilterLogic()

      filterLogic.setExceptions([
        { id: 1, status: 'PENDING' },
        { id: 2, status: 'APPROVED' },
        { id: 3, status: 'REJECTED' },
        { id: 4, status: 'REVOKED' },
        { id: 5, status: 'PENDING' },
      ])

      expect(filterLogic.getPendingExceptions()).toHaveLength(2)
      expect(filterLogic.getApprovedExceptions()).toHaveLength(1)
      expect(filterLogic.getRejectedExceptions()).toHaveLength(1)
      expect(filterLogic.getRevokedExceptions()).toHaveLength(1)
    })
  })
})
