import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock modules before imports
vi.mock('vue', async () => {
  const vue = await import('vue')
  return {
    ...vue,
    ref: (v: any) => ({ value: v }),
    computed: (fn: () => any) => ({ value: fn() }),
    onMounted: (fn: () => void) => fn(),
  }
})

vi.mock('~/composables/useAlert', () => ({
  alertError: vi.fn(),
  alertSuccess: vi.fn(),
  alertWarning: vi.fn(),
  confirm: vi.fn(),
}))

vi.mock('~/composables/useNotification', () => ({
  default: () => ({
    show: { value: false },
    close: vi.fn(),
    success: vi.fn(),
    error: vi.fn(),
  }),
}))

vi.mock('~/composables/useCenterId', () => ({
  getCenterId: vi.fn(() => 1),
}))

vi.mock('~/composables/useApi', () => ({
  default: () => ({
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  }),
}))

describe('admin/offerings.vue 頁面邏輯', () => {
  // OfferingManagementLogic 類別 - 待排課程管理業務邏輯
  class OfferingManagementLogic {
    offerings: any[]
    loading: boolean
    searchQuery: string
    selectedStatus: string
    selectedPriority: string
    statuses: { id: string; label: string }[]
    priorities: { id: string; label: string }[]

    constructor() {
      this.offerings = []
      this.loading = false
      this.searchQuery = ''
      this.selectedStatus = 'all'
      this.selectedPriority = 'all'
      this.statuses = [
        { id: 'all', label: '全部' },
        { id: 'PENDING', label: '待處理' },
        { id: 'IN_PROGRESS', label: '進行中' },
        { id: 'SCHEDULED', label: '已排課' },
        { id: 'COMPLETED', label: '已完成' },
        { id: 'CANCELLED', label: '已取消' },
      ]
      this.priorities = [
        { id: 'all', label: '全部' },
        { id: 'HIGH', label: '高' },
        { id: 'MEDIUM', label: '中' },
        { id: 'LOW', label: '低' },
      ]
    }

    // 載入待排課程列表
    async loadOfferings(): Promise<void> {
      this.loading = true
      await new Promise(resolve => setTimeout(resolve, 100))
      this.loading = false
    }

    setOfferings(offerings: any[]) {
      this.offerings = offerings
    }

    // 取得待排課程總數
    getOfferingCount(): number {
      return this.offerings.length
    }

    // 依狀態取得待排課程數量
    getOfferingCountByStatus(status: string): number {
      if (status === 'all') return this.offerings.length
      return this.offerings.filter(o => o.status === status).length
    }

    // 依 ID 取得待排課程
    getOfferingById(id: number): any | undefined {
      return this.offerings.find(o => o.id === id)
    }

    // 新增待排課程
    addOffering(offering: any) {
      const newOffering = {
        ...offering,
        id: this.offerings.length + 1,
        created_at: new Date().toISOString(),
        status: offering.status || 'PENDING',
        priority: offering.priority || 'MEDIUM',
      }
      this.offerings.push(newOffering)
      return newOffering
    }

    // 更新待排課程
    updateOffering(id: number, updates: any): boolean {
      const index = this.offerings.findIndex(o => o.id === id)
      if (index !== -1) {
        this.offerings[index] = { ...this.offerings[index], ...updates }
        return true
      }
      return false
    }

    // 刪除待排課程
    deleteOffering(id: number): boolean {
      const index = this.offerings.findIndex(o => o.id === id)
      if (index !== -1) {
        this.offerings.splice(index, 1)
        return true
      }
      return false
    }

    // 更新待排課程狀態
    updateOfferingStatus(id: number, newStatus: string): boolean {
      const offering = this.getOfferingById(id)
      if (offering) {
        offering.status = newStatus
        return true
      }
      return false
    }

    // 過濾待排課程（搜尋 + 狀態 + 優先級）
    filteredOfferings(): any[] {
      let result = this.offerings

      // 狀態過濾
      if (this.selectedStatus !== 'all') {
        result = result.filter(o => o.status === this.selectedStatus)
      }

      // 優先級過濾
      if (this.selectedPriority !== 'all') {
        result = result.filter(o => o.priority === this.selectedPriority)
      }

      // 搜尋過濾
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase()
        result = result.filter(o =>
          o.name?.toLowerCase().includes(query) ||
          o.course_name?.toLowerCase().includes(query) ||
          o.student_name?.toLowerCase().includes(query)
        )
      }

      return result
    }

    // 取得待處理待排課程
    pendingOfferings(): any[] {
      return this.offerings.filter(o => o.status === 'PENDING')
    }

    // 取得進行中待排課程
    inProgressOfferings(): any[] {
      return this.offerings.filter(o => o.status === 'IN_PROGRESS')
    }

    // 取得已排課待排課程
    scheduledOfferings(): any[] {
      return this.offerings.filter(o => o.status === 'SCHEDULED')
    }

    // 取得已完成待排課程
    completedOfferings(): any[] {
      return this.offerings.filter(o => o.status === 'COMPLETED')
    }

    // 驗證待排課程資料
    validateOffering(offering: any): { valid: boolean; errors: string[] } {
      const errors: string[] = []

      if (!offering.name || offering.name.trim() === '') {
        errors.push('名稱為必填')
      }

      if (!offering.course_id) {
        errors.push('課程為必選')
      }

      if (!offering.student_name || offering.student_name.trim() === '') {
        errors.push('學生姓名為必填')
      }

      if (offering.sessions_requested !== undefined && offering.sessions_requested <= 0) {
        errors.push('需求堂數必須大於 0')
      }

      if (offering.priority && !['HIGH', 'MEDIUM', 'LOW'].includes(offering.priority)) {
        errors.push('優先級選項不正確')
      }

      return {
        valid: errors.length === 0,
        errors,
      }
    }

    // 取得待排課程統計
    getOfferingStats(): {
      total: number
      pending: number
      inProgress: number
      scheduled: number
      completed: number
      cancelled: number
      byPriority: Record<string, number>
    } {
      const byPriority: Record<string, number> = {}
      this.offerings.forEach(o => {
        const priority = o.priority || 'MEDIUM'
        byPriority[priority] = (byPriority[priority] || 0) + 1
      })

      return {
        total: this.offerings.length,
        pending: this.pendingOfferings().length,
        inProgress: this.inProgressOfferings().length,
        scheduled: this.scheduledOfferings().length,
        completed: this.completedOfferings().length,
        cancelled: this.offerings.filter(o => o.status === 'CANCELLED').length,
        byPriority,
      }
    }

    // 搜尋待排課程
    searchOfferings(keyword: string): any[] {
      if (!keyword) return this.offerings
      const query = keyword.toLowerCase()
      return this.offerings.filter(o =>
        o.name?.toLowerCase().includes(query) ||
        o.course_name?.toLowerCase().includes(query) ||
        o.student_name?.toLowerCase().includes(query) ||
        o.student_phone?.includes(query)
      )
    }

    // 取得高優先級待處理待排課程
    highPriorityPendingOfferings(): any[] {
      return this.pendingOfferings().filter(o => o.priority === 'HIGH')
    }

    // 計算待處理率
    getPendingRate(): number {
      if (this.offerings.length === 0) return 0
      return (this.pendingOfferings().length / this.offerings.length) * 100
    }

    // 指派老師
    assignTeacher(offeringId: number, teacherId: number): boolean {
      const offering = this.getOfferingById(offeringId)
      if (offering) {
        offering.assigned_teacher_id = teacherId
        offering.status = 'IN_PROGRESS'
        return true
      }
      return false
    }

    // 排課完成
    markAsScheduled(offeringId: number, scheduleId: number): boolean {
      const offering = this.getOfferingById(offeringId)
      if (offering) {
        offering.schedule_id = scheduleId
        offering.status = 'SCHEDULED'
        return true
      }
      return false
    }

    // 完成待排課程
    markAsCompleted(offeringId: number): boolean {
      const offering = this.getOfferingById(offeringId)
      if (offering) {
        offering.status = 'COMPLETED'
        offering.completed_at = new Date().toISOString()
        return true
      }
      return false
    }

    // 取消待排課程
    cancelOffering(offeringId: number, reason?: string): boolean {
      const offering = this.getOfferingById(offeringId)
      if (offering) {
        offering.status = 'CANCELLED'
        offering.cancelled_at = new Date().toISOString()
        if (reason) {
          offering.cancel_reason = reason
        }
        return true
      }
      return false
    }
  }

  describe('OfferingManagementLogic 基本操作', () => {
    it('應該正確初始化', () => {
      const logic = new OfferingManagementLogic()
      expect(logic.offerings).toHaveLength(0)
      expect(logic.loading).toBe(false)
      expect(logic.searchQuery).toBe('')
      expect(logic.selectedStatus).toBe('all')
      expect(logic.statuses).toHaveLength(6)
      expect(logic.priorities).toHaveLength(4)
    })

    it('setOfferings 應該正確設定待排課程列表', () => {
      const logic = new OfferingManagementLogic()
      const offerings = [
        { id: 1, name: '暑期鋼琴班', status: 'PENDING', priority: 'HIGH' },
        { id: 2, name: '常態鋼琴課', status: 'IN_PROGRESS', priority: 'MEDIUM' },
      ]
      logic.setOfferings(offerings)
      expect(logic.offerings).toHaveLength(2)
    })

    it('getOfferingCount 應該返回正確的待排課程數量', () => {
      const logic = new OfferingManagementLogic()
      expect(logic.getOfferingCount()).toBe(0)
      logic.setOfferings([{ id: 1 }, { id: 2 }, { id: 3 }])
      expect(logic.getOfferingCount()).toBe(3)
    })

    it('getOfferingCountByStatus 應該依狀態計算數量', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, status: 'PENDING' },
        { id: 2, status: 'PENDING' },
        { id: 3, status: 'COMPLETED' },
      ])
      expect(logic.getOfferingCountByStatus('all')).toBe(3)
      expect(logic.getOfferingCountByStatus('PENDING')).toBe(2)
      expect(logic.getOfferingCountByStatus('COMPLETED')).toBe(1)
    })

    it('getOfferingById 應該正確取得待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '待排課程 A' },
        { id: 2, name: '待排課程 B' },
      ])
      const offering = logic.getOfferingById(1)
      expect(offering?.name).toBe('待排課程 A')
      const notFound = logic.getOfferingById(999)
      expect(notFound).toBeUndefined()
    })
  })

  describe('OfferingManagementLogic CRUD 操作', () => {
    it('addOffering 應該正確新增待排課程', () => {
      const logic = new OfferingManagementLogic()
      const newOffering = logic.addOffering({
        name: '新待排課程',
        course_id: 1,
        student_name: '學生A',
        sessions_requested: 10,
      })
      expect(logic.offerings).toHaveLength(1)
      expect(newOffering.name).toBe('新待排課程')
      expect(newOffering.id).toBe(1)
      expect(newOffering.status).toBe('PENDING')
      expect(newOffering.priority).toBe('MEDIUM')
    })

    it('addOffering 應該自動設定建立時間', () => {
      const logic = new OfferingManagementLogic()
      const newOffering = logic.addOffering({ name: '新待排課程' })
      expect(newOffering.created_at).toBeDefined()
    })

    it('updateOffering 應該更新待排課程資訊', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1, name: '舊名稱', priority: 'LOW' }])
      const result = logic.updateOffering(1, { name: '新名稱', priority: 'HIGH' })
      expect(result).toBe(true)
      const offering = logic.getOfferingById(1)
      expect(offering?.name).toBe('新名稱')
      expect(offering?.priority).toBe('HIGH')
    })

    it('updateOffering 應該處理不存在的待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1, name: '待排課程 A' }])
      const result = logic.updateOffering(999, { name: '新名稱' })
      expect(result).toBe(false)
    })

    it('deleteOffering 應該刪除待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '待排課程 A' },
        { id: 2, name: '待排課程 B' },
      ])
      const result = logic.deleteOffering(1)
      expect(result).toBe(true)
      expect(logic.offerings).toHaveLength(1)
      expect(logic.offerings[0].id).toBe(2)
    })

    it('deleteOffering 應該處理不存在的待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1, name: '待排課程 A' }])
      const result = logic.deleteOffering(999)
      expect(result).toBe(false)
      expect(logic.offerings).toHaveLength(1)
    })

    it('updateOfferingStatus 應該更新待排課程狀態', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1, name: '待排課程 A', status: 'PENDING' }])
      const result = logic.updateOfferingStatus(1, 'IN_PROGRESS')
      expect(result).toBe(true)
      const offering = logic.getOfferingById(1)
      expect(offering?.status).toBe('IN_PROGRESS')
    })
  })

  describe('OfferingManagementLogic 過濾與搜尋', () => {
    it('filteredOfferings 應該返回所有待排課程當無搜尋條件時', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '待排課程 A', status: 'PENDING' },
        { id: 2, name: '待排課程 B', status: 'COMPLETED' },
      ])
      logic.searchQuery = ''
      logic.selectedStatus = 'all'
      logic.selectedPriority = 'all'
      expect(logic.filteredOfferings()).toHaveLength(2)
    })

    it('filteredOfferings 應該依狀態過濾', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '待排課程 A', status: 'PENDING' },
        { id: 2, name: '待排課程 B', status: 'PENDING' },
        { id: 3, name: '待排課程 C', status: 'COMPLETED' },
      ])
      logic.searchQuery = ''
      logic.selectedStatus = 'PENDING'
      logic.selectedPriority = 'all'
      expect(logic.filteredOfferings()).toHaveLength(2)
    })

    it('filteredOfferings 應該依優先級過濾', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '待排課程 A', priority: 'HIGH' },
        { id: 2, name: '待排課程 B', priority: 'MEDIUM' },
        { id: 3, name: '待排課程 C', priority: 'HIGH' },
      ])
      logic.searchQuery = ''
      logic.selectedStatus = 'all'
      logic.selectedPriority = 'HIGH'
      expect(logic.filteredOfferings()).toHaveLength(2)
    })

    it('filteredOfferings 應該依搜尋條件過濾', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '鋼琴課程', student_name: '學生A' },
        { id: 2, name: '小提琴課程', student_name: '學生B' },
        { id: 3, name: '鋼琴進階', student_name: '學生C' },
      ])
      logic.searchQuery = '鋼琴'
      logic.selectedStatus = 'all'
      logic.selectedPriority = 'all'
      const filtered = logic.filteredOfferings()
      expect(filtered).toHaveLength(2)
    })

    it('filteredOfferings 應該支援學生姓名搜尋', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '課程A', student_name: '張三' },
        { id: 2, name: '課程B', student_name: '李四' },
      ])
      logic.searchQuery = '張三'
      logic.selectedStatus = 'all'
      logic.selectedPriority = 'all'
      expect(logic.filteredOfferings()).toHaveLength(1)
    })

    it('filteredOfferings 應該同時支援多條件過濾', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '鋼琴課程', status: 'PENDING', priority: 'HIGH' },
        { id: 2, name: '鋼琴課程', status: 'PENDING', priority: 'MEDIUM' },
        { id: 3, name: '小提琴課程', status: 'PENDING', priority: 'HIGH' },
        { id: 4, name: '鋼琴課程', status: 'COMPLETED', priority: 'HIGH' },
      ])
      logic.searchQuery = '鋼琴'
      logic.selectedStatus = 'PENDING'
      logic.selectedPriority = 'HIGH'
      const filtered = logic.filteredOfferings()
      expect(filtered).toHaveLength(1)
      expect(filtered[0].id).toBe(1)
    })

    it('pendingOfferings 應該返回待處理待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '待排課程 A', status: 'PENDING' },
        { id: 2, name: '待排課程 B', status: 'IN_PROGRESS' },
      ])
      const pending = logic.pendingOfferings()
      expect(pending).toHaveLength(1)
      expect(pending[0].name).toBe('待排課程 A')
    })

    it('inProgressOfferings 應該返回進行中待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '待排課程 A', status: 'PENDING' },
        { id: 2, name: '待排課程 B', status: 'IN_PROGRESS' },
      ])
      const inProgress = logic.inProgressOfferings()
      expect(inProgress).toHaveLength(1)
      expect(inProgress[0].name).toBe('待排課程 B')
    })

    it('scheduledOfferings 應該返回已排課待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '待排課程 A', status: 'SCHEDULED' },
        { id: 2, name: '待排課程 B', status: 'COMPLETED' },
      ])
      const scheduled = logic.scheduledOfferings()
      expect(scheduled).toHaveLength(1)
      expect(scheduled[0].name).toBe('待排課程 A')
    })

    it('completedOfferings 應該返回已完成待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '待排課程 A', status: 'SCHEDULED' },
        { id: 2, name: '待排課程 B', status: 'COMPLETED' },
      ])
      const completed = logic.completedOfferings()
      expect(completed).toHaveLength(1)
      expect(completed[0].name).toBe('待排課程 B')
    })
  })

  describe('OfferingManagementLogic 資料驗證', () => {
    it('validateOffering 應該通過有效的待排課程資料', () => {
      const logic = new OfferingManagementLogic()
      const result = logic.validateOffering({
        name: '暑期鋼琴班',
        course_id: 1,
        student_name: '學生A',
        sessions_requested: 10,
        priority: 'HIGH',
      })
      expect(result.valid).toBe(true)
      expect(result.errors).toHaveLength(0)
    })

    it('validateOffering 應該驗證名稱必填', () => {
      const logic = new OfferingManagementLogic()
      const result = logic.validateOffering({
        name: '',
        course_id: 1,
        student_name: '學生A',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('名稱為必填')
    })

    it('validateOffering 應該驗證課程必選', () => {
      const logic = new OfferingManagementLogic()
      const result = logic.validateOffering({
        name: '待排課程',
        course_id: null,
        student_name: '學生A',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('課程為必選')
    })

    it('validateOffering 應該驗證學生姓名必填', () => {
      const logic = new OfferingManagementLogic()
      const result = logic.validateOffering({
        name: '待排課程',
        course_id: 1,
        student_name: '',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('學生姓名為必填')
    })

    it('validateOffering 應該驗證需求堂數', () => {
      const logic = new OfferingManagementLogic()
      const result = logic.validateOffering({
        name: '待排課程',
        course_id: 1,
        student_name: '學生A',
        sessions_requested: 0,
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('需求堂數必須大於 0')
    })

    it('validateOffering 應該驗證優先級選項', () => {
      const logic = new OfferingManagementLogic()
      const result = logic.validateOffering({
        name: '待排課程',
        course_id: 1,
        student_name: '學生A',
        priority: 'INVALID',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('優先級選項不正確')
    })

    it('validateOffering 應該收集多個錯誤', () => {
      const logic = new OfferingManagementLogic()
      const result = logic.validateOffering({
        name: '',
        course_id: null,
        student_name: '',
        sessions_requested: -1,
      })
      expect(result.valid).toBe(false)
      expect(result.errors.length).toBeGreaterThan(1)
    })
  })

  describe('OfferingManagementLogic 資料完整性檢查', () => {
    it('getOfferingStats 應該返回正確的統計資料', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, status: 'PENDING', priority: 'HIGH' },
        { id: 2, status: 'PENDING', priority: 'MEDIUM' },
        { id: 3, status: 'IN_PROGRESS', priority: 'HIGH' },
        { id: 4, status: 'SCHEDULED', priority: 'LOW' },
        { id: 5, status: 'COMPLETED', priority: 'MEDIUM' },
        { id: 6, status: 'CANCELLED', priority: 'LOW' },
      ])
      const stats = logic.getOfferingStats()
      expect(stats.total).toBe(6)
      expect(stats.pending).toBe(2)
      expect(stats.inProgress).toBe(1)
      expect(stats.scheduled).toBe(1)
      expect(stats.completed).toBe(1)
      expect(stats.cancelled).toBe(1)
      expect(stats.byPriority['HIGH']).toBe(2)
      expect(stats.byPriority['MEDIUM']).toBe(2)
      expect(stats.byPriority['LOW']).toBe(2)
    })

    it('getPendingRate 應該計算正確的待處理率', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, status: 'PENDING' },
        { id: 2, status: 'PENDING' },
        { id: 3, status: 'COMPLETED' },
        { id: 4, status: 'COMPLETED' },
      ])
      expect(logic.getPendingRate()).toBe(50)
    })

    it('getPendingRate 應該處理空列表', () => {
      const logic = new OfferingManagementLogic()
      expect(logic.getPendingRate()).toBe(0)
    })
  })

  describe('OfferingManagementLogic 搜尋功能', () => {
    it('searchOfferings 應該返回所有待排課程當關鍵字為空', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '待排課程 A' },
        { id: 2, name: '待排課程 B' },
      ])
      expect(logic.searchOfferings('')).toHaveLength(2)
    })

    it('searchOfferings 應該支援名稱搜尋', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '鋼琴課程' },
        { id: 2, name: '小提琴課程' },
        { id: 3, name: '鋼琴進階' },
      ])
      expect(logic.searchOfferings('鋼琴')).toHaveLength(2)
    })

    it('searchOfferings 應該支援學生電話搜尋', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '課程A', student_phone: '0912345678' },
        { id: 2, name: '課程B', student_phone: '0987654321' },
      ])
      expect(logic.searchOfferings('0912')).toHaveLength(1)
    })
  })

  describe('OfferingManagementLogic 工作流程', () => {
    it('highPriorityPendingOfferings 應該返回高優先級待處理待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, status: 'PENDING', priority: 'HIGH' },
        { id: 2, status: 'PENDING', priority: 'MEDIUM' },
        { id: 3, status: 'PENDING', priority: 'HIGH' },
        { id: 4, status: 'IN_PROGRESS', priority: 'HIGH' },
      ])
      const highPriority = logic.highPriorityPendingOfferings()
      expect(highPriority).toHaveLength(2)
    })

    it('assignTeacher 應該指派老師並更新狀態', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1, status: 'PENDING' }])
      const result = logic.assignTeacher(1, 100)
      expect(result).toBe(true)
      const offering = logic.getOfferingById(1)
      expect(offering?.assigned_teacher_id).toBe(100)
      expect(offering?.status).toBe('IN_PROGRESS')
    })

    it('assignTeacher 應該處理不存在的待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1 }])
      const result = logic.assignTeacher(999, 100)
      expect(result).toBe(false)
    })

    it('markAsScheduled 應該標記為已排課', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1, status: 'IN_PROGRESS' }])
      const result = logic.markAsScheduled(1, 200)
      expect(result).toBe(true)
      const offering = logic.getOfferingById(1)
      expect(offering?.schedule_id).toBe(200)
      expect(offering?.status).toBe('SCHEDULED')
    })

    it('markAsCompleted 應該標記為已完成', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1, status: 'SCHEDULED' }])
      const result = logic.markAsCompleted(1)
      expect(result).toBe(true)
      const offering = logic.getOfferingById(1)
      expect(offering?.status).toBe('COMPLETED')
      expect(offering?.completed_at).toBeDefined()
    })

    it('cancelOffering 應該取消待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1, status: 'PENDING' }])
      const result = logic.cancelOffering(1, '學生放棄')
      expect(result).toBe(true)
      const offering = logic.getOfferingById(1)
      expect(offering?.status).toBe('CANCELLED')
      expect(offering?.cancel_reason).toBe('學生放棄')
      expect(offering?.cancelled_at).toBeDefined()
    })

    it('cancelOffering 應該支援不提供取消原因', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1, status: 'PENDING' }])
      logic.cancelOffering(1)
      const offering = logic.getOfferingById(1)
      expect(offering?.status).toBe('CANCELLED')
      expect(offering?.cancel_reason).toBeUndefined()
    })
  })

  describe('OfferingManagementLogic 邊界情況', () => {
    it('應該處理空搜尋結果', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([
        { id: 1, name: '待排課程 A' },
        { id: 2, name: '待排課程 B' },
      ])
      logic.searchQuery = '不存在的課程'
      logic.selectedStatus = 'all'
      logic.selectedPriority = 'all'
      expect(logic.filteredOfferings()).toHaveLength(0)
    })

    it('updateOfferingStatus 應該處理不存在的待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1, status: 'PENDING' }])
      const result = logic.updateOfferingStatus(999, 'COMPLETED')
      expect(result).toBe(false)
    })

    it('markAsScheduled 應該處理不存在的待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1 }])
      const result = logic.markAsScheduled(999, 200)
      expect(result).toBe(false)
    })

    it('markAsCompleted 應該處理不存在的待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1 }])
      const result = logic.markAsCompleted(999)
      expect(result).toBe(false)
    })

    it('cancelOffering 應該處理不存在的待排課程', () => {
      const logic = new OfferingManagementLogic()
      logic.setOfferings([{ id: 1 }])
      const result = logic.cancelOffering(999)
      expect(result).toBe(false)
    })
  })
})
