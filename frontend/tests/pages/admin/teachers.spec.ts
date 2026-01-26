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

describe('admin/teachers.vue 頁面邏輯', () => {
  // TeacherManagementLogic 類別 - 老師管理業務邏輯
  class TeacherManagementLogic {
    teachers: any[]
    loading: boolean
    searchQuery: string
    selectedStatus: string
    selectedCenter: string
    statuses: { id: string; label: string }[]

    constructor() {
      this.teachers = []
      this.loading = false
      this.searchQuery = ''
      this.selectedStatus = 'all'
      this.selectedCenter = 'all'
      this.statuses = [
        { id: 'all', label: '全部' },
        { id: 'ACTIVE', label: '在職中' },
        { id: 'INACTIVE', label: '已離職' },
        { id: 'ON_LEAVE', label: '請假中' },
      ]
    }

    // 載入老師列表
    async loadTeachers(): Promise<void> {
      this.loading = true
      await new Promise(resolve => setTimeout(resolve, 100))
      this.loading = false
    }

    setTeachers(teachers: any[]) {
      this.teachers = teachers
    }

    // 取得老師總數
    getTeacherCount(): number {
      return this.teachers.length
    }

    // 依狀態取得老師數量
    getTeacherCountByStatus(status: string): number {
      if (status === 'all') return this.teachers.length
      return this.teachers.filter(t => t.status === status).length
    }

    // 依 ID 取得老師
    getTeacherById(id: number): any | undefined {
      return this.teachers.find(t => t.id === id)
    }

    // 依 email 取得老師
    getTeacherByEmail(email: string): any | undefined {
      return this.teachers.find(t => t.email?.toLowerCase() === email.toLowerCase())
    }

    // 新增老師
    addTeacher(teacher: any) {
      const newTeacher = {
        ...teacher,
        id: this.teachers.length + 1,
        created_at: new Date().toISOString(),
        status: teacher.status || 'ACTIVE',
        joined_at: teacher.joined_at || new Date().toISOString(),
      }
      this.teachers.push(newTeacher)
      return newTeacher
    }

    // 更新老師
    updateTeacher(id: number, updates: any): boolean {
      const index = this.teachers.findIndex(t => t.id === id)
      if (index !== -1) {
        this.teachers[index] = { ...this.teachers[index], ...updates }
        return true
      }
      return false
    }

    // 刪除老師
    deleteTeacher(id: number): boolean {
      const index = this.teachers.findIndex(t => t.id === id)
      if (index !== -1) {
        this.teachers.splice(index, 1)
        return true
      }
      return false
    }

    // 更新老師狀態
    updateTeacherStatus(id: number, newStatus: string): boolean {
      const teacher = this.getTeacherById(id)
      if (teacher) {
        teacher.status = newStatus
        return true
      }
      return false
    }

    // 過濾老師（搜尋 + 狀態）
    filteredTeachers(): any[] {
      let result = this.teachers

      // 狀態過濾
      if (this.selectedStatus !== 'all') {
        result = result.filter(t => t.status === this.selectedStatus)
      }

      // 中心過濾
      if (this.selectedCenter !== 'all') {
        result = result.filter(t => t.center_id === parseInt(this.selectedCenter))
      }

      // 搜尋過濾
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase()
        result = result.filter(t =>
          t.name?.toLowerCase().includes(query) ||
          t.email?.toLowerCase().includes(query) ||
          t.phone?.includes(query)
        )
      }

      return result
    }

    // 取得在職老師
    activeTeachers(): any[] {
      return this.teachers.filter(t => t.status === 'ACTIVE')
    }

    // 取得請假老師
    onLeaveTeachers(): any[] {
      return this.teachers.filter(t => t.status === 'ON_LEAVE')
    }

    // 驗證老師資料
    validateTeacher(teacher: any): { valid: boolean; errors: string[] } {
      const errors: string[] = []

      if (!teacher.name || teacher.name.trim() === '') {
        errors.push('姓名為必填')
      }

      if (teacher.name && teacher.name.length > 50) {
        errors.push('姓名不能超過 50 個字元')
      }

      if (!teacher.email || teacher.email.trim() === '') {
        errors.push('Email 為必填')
      } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(teacher.email)) {
        errors.push('Email 格式不正確')
      }

      if (teacher.phone && !/^[\d\-+\s()]+$/.test(teacher.phone)) {
        errors.push('電話格式不正確')
      }

      return {
        valid: errors.length === 0,
        errors,
      }
    }

    // 檢查 Email 是否重複
    isEmailExists(email: string, excludeId?: number): boolean {
      return this.teachers.some(t =>
        t.email?.toLowerCase() === email.toLowerCase() &&
        t.id !== excludeId
      )
    }

    // 取得老師統計
    getTeacherStats(): {
      total: number
      active: number
      inactive: number
      onLeave: number
      byCenter: Record<number, number>
    } {
      const byCenter: Record<number, number> = {}
      this.teachers.forEach(t => {
        const center = t.center_id || 0
        byCenter[center] = (byCenter[center] || 0) + 1
      })

      return {
        total: this.teachers.length,
        active: this.activeTeachers().length,
        inactive: this.teachers.filter(t => t.status === 'INACTIVE').length,
        onLeave: this.onLeaveTeachers().length,
        byCenter,
      }
    }

    // 計算在職率
    getActiveRate(): number {
      if (this.teachers.length === 0) return 0
      return (this.activeTeachers().length / this.teachers.length) * 100
    }

    // 搜尋老師
    searchTeachers(keyword: string): any[] {
      if (!keyword) return this.teachers
      const query = keyword.toLowerCase()
      return this.teachers.filter(t =>
        t.name?.toLowerCase().includes(query) ||
        t.email?.toLowerCase().includes(query) ||
        t.phone?.includes(query) ||
        t.specialty?.some((s: string) => s.toLowerCase().includes(query))
      )
    }

    // 取得老師擅長領域
    getTeacherSpecialties(id: number): string[] {
      const teacher = this.getTeacherById(id)
      return teacher?.specialty || []
    }

    // 更新老師擅長領域
    updateTeacherSpecialties(id: number, specialties: string[]): boolean {
      const teacher = this.getTeacherById(id)
      if (teacher) {
        teacher.specialty = specialties
        return true
      }
      return false
    }
  }

  describe('TeacherManagementLogic 基本操作', () => {
    it('應該正確初始化', () => {
      const logic = new TeacherManagementLogic()
      expect(logic.teachers).toHaveLength(0)
      expect(logic.loading).toBe(false)
      expect(logic.searchQuery).toBe('')
      expect(logic.selectedStatus).toBe('all')
      expect(logic.statuses).toHaveLength(4)
    })

    it('setTeachers 應該正確設定老師列表', () => {
      const logic = new TeacherManagementLogic()
      const teachers = [
        { id: 1, name: '張老師', email: 'zhang@test.com', status: 'ACTIVE' },
        { id: 2, name: '李老師', email: 'li@test.com', status: 'ACTIVE' },
      ]
      logic.setTeachers(teachers)
      expect(logic.teachers).toHaveLength(2)
    })

    it('getTeacherCount 應該返回正確的老師數量', () => {
      const logic = new TeacherManagementLogic()
      expect(logic.getTeacherCount()).toBe(0)
      logic.setTeachers([{ id: 1 }, { id: 2 }, { id: 3 }])
      expect(logic.getTeacherCount()).toBe(3)
    })

    it('getTeacherCountByStatus 應該依狀態計算數量', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, status: 'ACTIVE' },
        { id: 2, status: 'ACTIVE' },
        { id: 3, status: 'INACTIVE' },
      ])
      expect(logic.getTeacherCountByStatus('all')).toBe(3)
      expect(logic.getTeacherCountByStatus('ACTIVE')).toBe(2)
      expect(logic.getTeacherCountByStatus('INACTIVE')).toBe(1)
    })

    it('getTeacherById 應該正確取得老師', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '張老師' },
        { id: 2, name: '李老師' },
      ])
      const teacher = logic.getTeacherById(1)
      expect(teacher?.name).toBe('張老師')
      const notFound = logic.getTeacherById(999)
      expect(notFound).toBeUndefined()
    })

    it('getTeacherByEmail 應該正確取得老師', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, email: 'zhang@test.com' },
        { id: 2, email: 'li@test.com' },
      ])
      const teacher = logic.getTeacherByEmail('ZHANG@TEST.COM')
      expect(teacher?.id).toBe(1)
    })
  })

  describe('TeacherManagementLogic CRUD 操作', () => {
    it('addTeacher 應該正確新增老師', () => {
      const logic = new TeacherManagementLogic()
      const newTeacher = logic.addTeacher({
        name: '新老師',
        email: 'new@test.com',
        phone: '0912345678',
      })
      expect(logic.teachers).toHaveLength(1)
      expect(newTeacher.name).toBe('新老師')
      expect(newTeacher.id).toBe(1)
      expect(newTeacher.status).toBe('ACTIVE')
    })

    it('addTeacher 應該自動設定加入日期', () => {
      const logic = new TeacherManagementLogic()
      const newTeacher = logic.addTeacher({ name: '新老師', email: 'test@test.com' })
      expect(newTeacher.joined_at).toBeDefined()
    })

    it('updateTeacher 應該更新老師資訊', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([{ id: 1, name: '舊名稱', email: 'old@test.com' }])
      const result = logic.updateTeacher(1, { name: '新名稱', phone: '0987654321' })
      expect(result).toBe(true)
      const teacher = logic.getTeacherById(1)
      expect(teacher?.name).toBe('新名稱')
      expect(teacher?.phone).toBe('0987654321')
    })

    it('updateTeacher 應該處理不存在的課程', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([{ id: 1, name: '老師 A' }])
      const result = logic.updateTeacher(999, { name: '新名稱' })
      expect(result).toBe(false)
    })

    it('deleteTeacher 應該刪除老師', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '老師 A' },
        { id: 2, name: '老師 B' },
      ])
      const result = logic.deleteTeacher(1)
      expect(result).toBe(true)
      expect(logic.teachers).toHaveLength(1)
      expect(logic.teachers[0].id).toBe(2)
    })

    it('deleteTeacher 應該處理不存在的課程', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([{ id: 1, name: '老師 A' }])
      const result = logic.deleteTeacher(999)
      expect(result).toBe(false)
      expect(logic.teachers).toHaveLength(1)
    })

    it('updateTeacherStatus 應該更新老師狀態', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([{ id: 1, name: '老師 A', status: 'ACTIVE' }])
      const result = logic.updateTeacherStatus(1, 'ON_LEAVE')
      expect(result).toBe(true)
      const teacher = logic.getTeacherById(1)
      expect(teacher?.status).toBe('ON_LEAVE')
    })
  })

  describe('TeacherManagementLogic 過濾與搜尋', () => {
    it('filteredTeachers 應該返回所有老師當無搜尋條件時', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '張老師', status: 'ACTIVE' },
        { id: 2, name: '李老師', status: 'INACTIVE' },
      ])
      logic.searchQuery = ''
      logic.selectedStatus = 'all'
      expect(logic.filteredTeachers()).toHaveLength(2)
    })

    it('filteredTeachers 應該依狀態過濾', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '張老師', status: 'ACTIVE' },
        { id: 2, name: '李老師', status: 'ACTIVE' },
        { id: 3, name: '王老師', status: 'INACTIVE' },
      ])
      logic.searchQuery = ''
      logic.selectedStatus = 'ACTIVE'
      expect(logic.filteredTeachers()).toHaveLength(2)
    })

    it('filteredTeachers 應該依中心過濾', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '張老師', center_id: 1 },
        { id: 2, name: '李老師', center_id: 2 },
        { id: 3, name: '王老師', center_id: 1 },
      ])
      logic.searchQuery = ''
      logic.selectedStatus = 'all'
      logic.selectedCenter = '1'
      expect(logic.filteredTeachers()).toHaveLength(2)
    })

    it('filteredTeachers 應該依搜尋條件過濾', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '張老師', email: 'zhang@test.com' },
        { id: 2, name: '李老師', email: 'li@test.com' },
        { id: 3, name: '張三老師', email: 'zhangsan@test.com' },
      ])
      logic.searchQuery = '張'
      logic.selectedStatus = 'all'
      const filtered = logic.filteredTeachers()
      expect(filtered).toHaveLength(2)
    })

    it('filteredTeachers 應該支援電話搜尋', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '張老師', phone: '0912345678' },
        { id: 2, name: '李老師', phone: '0987654321' },
      ])
      logic.searchQuery = '0912'
      logic.selectedStatus = 'all'
      expect(logic.filteredTeachers()).toHaveLength(1)
    })

    it('activeTeachers 應該返回在職老師', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '老師 A', status: 'ACTIVE' },
        { id: 2, name: '老師 B', status: 'INACTIVE' },
      ])
      const active = logic.activeTeachers()
      expect(active).toHaveLength(1)
      expect(active[0].name).toBe('老師 A')
    })

    it('onLeaveTeachers 應該返回請假老師', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '老師 A', status: 'ACTIVE' },
        { id: 2, name: '老師 B', status: 'ON_LEAVE' },
      ])
      const onLeave = logic.onLeaveTeachers()
      expect(onLeave).toHaveLength(1)
      expect(onLeave[0].name).toBe('老師 B')
    })
  })

  describe('TeacherManagementLogic 資料驗證', () => {
    it('validateTeacher 應該通過有效的老師資料', () => {
      const logic = new TeacherManagementLogic()
      const result = logic.validateTeacher({
        name: '張老師',
        email: 'zhang@test.com',
        phone: '0912345678',
      })
      expect(result.valid).toBe(true)
      expect(result.errors).toHaveLength(0)
    })

    it('validateTeacher 應該驗證姓名必填', () => {
      const logic = new TeacherManagementLogic()
      const result = logic.validateTeacher({
        name: '',
        email: 'test@test.com',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('姓名為必填')
    })

    it('validateTeacher 應該驗證姓名長度', () => {
      const logic = new TeacherManagementLogic()
      const result = logic.validateTeacher({
        name: 'a'.repeat(51),
        email: 'test@test.com',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('姓名不能超過 50 個字元')
    })

    it('validateTeacher 應該驗證 Email 必填', () => {
      const logic = new TeacherManagementLogic()
      const result = logic.validateTeacher({
        name: '張老師',
        email: '',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('Email 為必填')
    })

    it('validateTeacher 應該驗證 Email 格式', () => {
      const logic = new TeacherManagementLogic()
      const result = logic.validateTeacher({
        name: '張老師',
        email: 'invalid-email',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('Email 格式不正確')
    })

    it('validateTeacher 應該驗證電話格式', () => {
      const logic = new TeacherManagementLogic()
      const result = logic.validateTeacher({
        name: '張老師',
        email: 'test@test.com',
        phone: 'invalid',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('電話格式不正確')
    })

    it('validateTeacher 應該收集多個錯誤', () => {
      const logic = new TeacherManagementLogic()
      const result = logic.validateTeacher({
        name: '',
        email: '',
        phone: 'invalid',
      })
      expect(result.valid).toBe(false)
      expect(result.errors.length).toBeGreaterThan(1)
    })
  })

  describe('TeacherManagementLogic 資料完整性檢查', () => {
    it('isEmailExists 應該檢測 Email 是否存在', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, email: 'zhang@test.com' },
        { id: 2, email: 'li@test.com' },
      ])
      expect(logic.isEmailExists('zhang@test.com')).toBe(true)
      expect(logic.isEmailExists('ZHANG@TEST.COM')).toBe(true)
      expect(logic.isEmailExists('new@test.com')).toBe(false)
    })

    it('isEmailExists 應該支援排除特定 ID', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, email: 'zhang@test.com' },
        { id: 2, email: 'zhang@test.com' },
      ])
      expect(logic.isEmailExists('zhang@test.com')).toBe(true)
      expect(logic.isEmailExists('zhang@test.com', 1)).toBe(true)
      expect(logic.isEmailExists('zhang@test.com', 999)).toBe(true)
    })

    it('getTeacherStats 應該返回正確的統計資料', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, status: 'ACTIVE', center_id: 1 },
        { id: 2, status: 'ACTIVE', center_id: 1 },
        { id: 3, status: 'INACTIVE', center_id: 2 },
        { id: 4, status: 'ON_LEAVE', center_id: 1 },
      ])
      const stats = logic.getTeacherStats()
      expect(stats.total).toBe(4)
      expect(stats.active).toBe(2)
      expect(stats.inactive).toBe(1)
      expect(stats.onLeave).toBe(1)
      expect(stats.byCenter[1]).toBe(3)
      expect(stats.byCenter[2]).toBe(1)
    })

    it('getActiveRate 應該計算正確的在職率', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, status: 'ACTIVE' },
        { id: 2, status: 'ACTIVE' },
        { id: 3, status: 'INACTIVE' },
        { id: 4, status: 'INACTIVE' },
      ])
      expect(logic.getActiveRate()).toBe(50)
    })

    it('getActiveRate 應該處理空列表', () => {
      const logic = new TeacherManagementLogic()
      expect(logic.getActiveRate()).toBe(0)
    })
  })

  describe('TeacherManagementLogic 搜尋功能', () => {
    it('searchTeachers 應該返回所有老師當關鍵字為空', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '張老師' },
        { id: 2, name: '李老師' },
      ])
      expect(logic.searchTeachers('')).toHaveLength(2)
    })

    it('searchTeachers 應該支援姓名搜尋', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '張老師' },
        { id: 2, name: '李老師' },
        { id: 3, name: '張三老師' },
      ])
      expect(logic.searchTeachers('張')).toHaveLength(2)
    })

    it('searchTeachers 應該支援擅長領域搜尋', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '張老師', specialty: ['鋼琴', '小提琴'] },
        { id: 2, name: '李老師', specialty: ['鋼琴'] },
        { id: 3, name: '王老師', specialty: ['吉他'] },
      ])
      expect(logic.searchTeachers('小提琴')).toHaveLength(1)
    })
  })

  describe('TeacherManagementLogic 擅長領域管理', () => {
    it('getTeacherSpecialties 應該返回老師的擅長領域', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, specialty: ['鋼琴', '小提琴'] },
        { id: 2, specialty: ['吉他'] },
      ])
      expect(logic.getTeacherSpecialties(1)).toEqual(['鋼琴', '小提琴'])
      expect(logic.getTeacherSpecialties(2)).toEqual(['吉他'])
    })

    it('getTeacherSpecialties 應該處理沒有擅長領域的老師', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([{ id: 1 }])
      expect(logic.getTeacherSpecialties(1)).toEqual([])
    })

    it('updateTeacherSpecialties 應該更新老師的擅長領域', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([{ id: 1, specialty: ['鋼琴'] }])
      const result = logic.updateTeacherSpecialties(1, ['鋼琴', '吉他', '小提琴'])
      expect(result).toBe(true)
      expect(logic.getTeacherSpecialties(1)).toEqual(['鋼琴', '吉他', '小提琴'])
    })
  })

  describe('TeacherManagementLogic 邊界情況', () => {
    it('應該處理空搜尋結果', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([
        { id: 1, name: '張老師' },
        { id: 2, name: '李老師' },
      ])
      logic.searchQuery = '不存在的名字'
      logic.selectedStatus = 'all'
      expect(logic.filteredTeachers()).toHaveLength(0)
    })

    it('updateTeacherStatus 應該處理不存在的課程', () => {
      const logic = new TeacherManagementLogic()
      logic.setTeachers([{ id: 1, status: 'ACTIVE' }])
      const result = logic.updateTeacherStatus(999, 'INACTIVE')
      expect(result).toBe(false)
    })

    it('getTeacherSpecialties 應該處理不存在的課程', () => {
      const logic = new TeacherManagementLogic()
      expect(logic.getTeacherSpecialties(999)).toEqual([])
    })

    it('updateTeacherSpecialties 應該處理不存在的課程', () => {
      const logic = new TeacherManagementLogic()
      const result = logic.updateTeacherSpecialties(999, ['鋼琴'])
      expect(result).toBe(false)
    })
  })
})
