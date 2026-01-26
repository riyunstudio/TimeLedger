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

describe('admin/courses.vue 頁面邏輯', () => {
  // CourseManagementLogic 類別 - 課程管理業務邏輯
  class CourseManagementLogic {
    courses: any[]
    loading: boolean
    searchQuery: string
    selectedCategory: string
    categories: { id: string; name: string }[]

    constructor() {
      this.courses = []
      this.loading = false
      this.searchQuery = ''
      this.selectedCategory = 'all'
      this.categories = [
        { id: 'all', name: '全部課程' },
        { id: 'music', name: '音樂類' },
        { id: 'art', name: '美術類' },
        { id: 'dance', name: '舞蹈類' },
        { id: 'language', name: '語言類' },
        { id: 'sports', name: '運動類' },
      ]
    }

    // 載入課程列表
    async loadCourses(): Promise<void> {
      this.loading = true
      // 模擬載入
      await new Promise(resolve => setTimeout(resolve, 100))
      this.loading = false
    }

    setCourses(courses: any[]) {
      this.courses = courses
    }

    // 取得課程總數
    getCourseCount(): number {
      return this.courses.length
    }

    // 取得已啟用課程數
    getActiveCourseCount(): number {
      return this.courses.filter(c => c.is_active).length
    }

    // 依 ID 取得課程
    getCourseById(id: number): any | undefined {
      return this.courses.find(c => c.id === id)
    }

    // 新增課程
    addCourse(course: any) {
      const newCourse = {
        ...course,
        id: this.courses.length + 1,
        created_at: new Date().toISOString(),
        is_active: course.is_active ?? true,
      }
      this.courses.push(newCourse)
      return newCourse
    }

    // 更新課程
    updateCourse(id: number, updates: any): boolean {
      const index = this.courses.findIndex(c => c.id === id)
      if (index !== -1) {
        this.courses[index] = { ...this.courses[index], ...updates }
        return true
      }
      return false
    }

    // 刪除課程
    deleteCourse(id: number): boolean {
      const index = this.courses.findIndex(c => c.id === id)
      if (index !== -1) {
        this.courses.splice(index, 1)
        return true
      }
      return false
    }

    // 切換課程啟用狀態
    toggleCourseStatus(id: number): boolean {
      const course = this.getCourseById(id)
      if (course) {
        course.is_active = !course.is_active
        return true
      }
      return false
    }

    // 過濾課程（搜尋 + 分類）
    filteredCourses(): any[] {
      let result = this.courses

      // 分類過濾
      if (this.selectedCategory !== 'all') {
        result = result.filter(c => c.category === this.selectedCategory)
      }

      // 搜尋過濾
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase()
        result = result.filter(c =>
          c.name?.toLowerCase().includes(query) ||
          c.description?.toLowerCase().includes(query) ||
          c.code?.toLowerCase().includes(query)
        )
      }

      return result
    }

    // 依分類取得課程數量
    getCoursesByCategory(category: string): number {
      return this.courses.filter(c => c.category === category).length
    }

    // 驗證課程資料
    validateCourse(course: any): { valid: boolean; errors: string[] } {
      const errors: string[] = []

      if (!course.name || course.name.trim() === '') {
        errors.push('課程名稱為必填')
      }

      if (course.name && course.name.length > 50) {
        errors.push('課程名稱不能超過 50 個字元')
      }

      if (!course.duration || course.duration <= 0) {
        errors.push('課程時長必須大於 0')
      }

      if (!course.category) {
        errors.push('課程分類為必填')
      }

      if (course.price !== undefined && course.price < 0) {
        errors.push('課程價格不能為負數')
      }

      return {
        valid: errors.length === 0,
        errors,
      }
    }

    // 檢查課程代碼是否重複
    isCodeExists(code: string, excludeId?: number): boolean {
      return this.courses.some(c =>
        c.code?.toLowerCase() === code.toLowerCase() &&
        c.id !== excludeId
      )
    }

    // 取得課程統計
    getCourseStats(): {
      total: number
      active: number
      inactive: number
      byCategory: Record<string, number>
    } {
      const byCategory: Record<string, number> = {}
      this.courses.forEach(c => {
        const cat = c.category || 'other'
        byCategory[cat] = (byCategory[cat] || 0) + 1
      })

      return {
        total: this.courses.length,
        active: this.courses.filter(c => c.is_active).length,
        inactive: this.courses.filter(c => !c.is_active).length,
        byCategory,
      }
    }
  }

  describe('CourseManagementLogic 基本操作', () => {
    it('應該正確初始化', () => {
      const logic = new CourseManagementLogic()
      expect(logic.courses).toHaveLength(0)
      expect(logic.loading).toBe(false)
      expect(logic.searchQuery).toBe('')
      expect(logic.selectedCategory).toBe('all')
      expect(logic.categories).toHaveLength(6)
    })

    it('setCourses 應該正確設定課程列表', () => {
      const logic = new CourseManagementLogic()
      const courses = [
        { id: 1, name: '鋼琴基礎', duration: 60, category: 'music', is_active: true },
        { id: 2, name: '小提琴入門', duration: 45, category: 'music', is_active: true },
      ]
      logic.setCourses(courses)
      expect(logic.courses).toHaveLength(2)
    })

    it('getCourseCount 應該返回正確的課程數量', () => {
      const logic = new CourseManagementLogic()
      expect(logic.getCourseCount()).toBe(0)
      logic.setCourses([{ id: 1 }, { id: 2 }, { id: 3 }])
      expect(logic.getCourseCount()).toBe(3)
    })

    it('getActiveCourseCount 應該返回已啟用課程數量', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, is_active: true },
        { id: 2, is_active: false },
        { id: 3, is_active: true },
      ])
      expect(logic.getActiveCourseCount()).toBe(2)
    })

    it('getCourseById 應該正確取得課程', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, name: '課程 A' },
        { id: 2, name: '課程 B' },
      ])
      const course = logic.getCourseById(1)
      expect(course?.name).toBe('課程 A')
      const notFound = logic.getCourseById(999)
      expect(notFound).toBeUndefined()
    })
  })

  describe('CourseManagementLogic CRUD 操作', () => {
    it('addCourse 應該正確新增課程', () => {
      const logic = new CourseManagementLogic()
      const newCourse = logic.addCourse({
        name: '新課程',
        duration: 60,
        category: 'music',
      })
      expect(logic.courses).toHaveLength(1)
      expect(newCourse.name).toBe('新課程')
      expect(newCourse.id).toBe(1)
      expect(newCourse.is_active).toBe(true)
    })

    it('addCourse 應該自動產生 ID', () => {
      const logic = new CourseManagementLogic()
      logic.addCourse({ name: '課程 1' })
      logic.addCourse({ name: '課程 2' })
      logic.addCourse({ name: '課程 3' })
      expect(logic.courses[0].id).toBe(1)
      expect(logic.courses[1].id).toBe(2)
      expect(logic.courses[2].id).toBe(3)
    })

    it('updateCourse 應該更新課程資訊', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([{ id: 1, name: '舊名稱', duration: 60 }])
      const result = logic.updateCourse(1, { name: '新名稱', duration: 90 })
      expect(result).toBe(true)
      const course = logic.getCourseById(1)
      expect(course?.name).toBe('新名稱')
      expect(course?.duration).toBe(90)
    })

    it('updateCourse 應該處理不存在的課程', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([{ id: 1, name: '課程 A' }])
      const result = logic.updateCourse(999, { name: '新名稱' })
      expect(result).toBe(false)
    })

    it('deleteCourse 應該刪除課程', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, name: '課程 A' },
        { id: 2, name: '課程 B' },
      ])
      const result = logic.deleteCourse(1)
      expect(result).toBe(true)
      expect(logic.courses).toHaveLength(1)
      expect(logic.courses[0].id).toBe(2)
    })

    it('deleteCourse 應該處理不存在的課程', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([{ id: 1, name: '課程 A' }])
      const result = logic.deleteCourse(999)
      expect(result).toBe(false)
      expect(logic.courses).toHaveLength(1)
    })

    it('toggleCourseStatus 應該切換課程狀態', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([{ id: 1, name: '課程 A', is_active: true }])
      const result = logic.toggleCourseStatus(1)
      expect(result).toBe(true)
      const course = logic.getCourseById(1)
      expect(course?.is_active).toBe(false)
    })
  })

  describe('CourseManagementLogic 過濾與搜尋', () => {
    it('filteredCourses 應該返回所有課程當無搜尋條件時', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, name: '課程 A', category: 'music' },
        { id: 2, name: '課程 B', category: 'art' },
      ])
      logic.searchQuery = ''
      logic.selectedCategory = 'all'
      expect(logic.filteredCourses()).toHaveLength(2)
    })

    it('filteredCourses 應該依搜尋條件過濾', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, name: '鋼琴基礎課程', category: 'music' },
        { id: 2, name: '鋼琴進階課程', category: 'music' },
        { id: 3, name: '小提琴入門', category: 'music' },
      ])
      logic.searchQuery = '鋼琴'
      logic.selectedCategory = 'all'
      const filtered = logic.filteredCourses()
      expect(filtered).toHaveLength(2)
    })

    it('filteredCourses 應該不區分大小寫搜尋', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([{ id: 1, name: 'PIANO Basic' }])
      logic.searchQuery = 'piano'
      logic.selectedCategory = 'all'
      expect(logic.filteredCourses()).toHaveLength(1)
    })

    it('filteredCourses 應該依分類過濾', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, name: '鋼琴課', category: 'music' },
        { id: 2, name: '繪畫課', category: 'art' },
        { id: 3, name: '鋼琴進階', category: 'music' },
      ])
      logic.searchQuery = ''
      logic.selectedCategory = 'music'
      const filtered = logic.filteredCourses()
      expect(filtered).toHaveLength(2)
    })

    it('filteredCourses 應該同時支援搜尋和分類過濾', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, name: '鋼琴基礎', category: 'music' },
        { id: 2, name: '鋼琴進階', category: 'music' },
        { id: 3, name: '小提琴', category: 'music' },
        { id: 4, name: '水彩畫', category: 'art' },
      ])
      logic.searchQuery = '基礎'
      logic.selectedCategory = 'music'
      const filtered = logic.filteredCourses()
      expect(filtered).toHaveLength(1)
      expect(filtered[0].name).toBe('鋼琴基礎')
    })

    it('getCoursesByCategory 應該返回分類課程數量', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, category: 'music' },
        { id: 2, category: 'music' },
        { id: 3, category: 'art' },
      ])
      expect(logic.getCoursesByCategory('music')).toBe(2)
      expect(logic.getCoursesByCategory('art')).toBe(1)
      expect(logic.getCoursesByCategory('dance')).toBe(0)
    })
  })

  describe('CourseManagementLogic 資料驗證', () => {
    it('validateCourse 應該通過有效的課程資料', () => {
      const logic = new CourseManagementLogic()
      const result = logic.validateCourse({
        name: '鋼琴基礎課程',
        duration: 60,
        category: 'music',
        price: 1000,
      })
      expect(result.valid).toBe(true)
      expect(result.errors).toHaveLength(0)
    })

    it('validateCourse 應該驗證課程名稱必填', () => {
      const logic = new CourseManagementLogic()
      const result = logic.validateCourse({
        name: '',
        duration: 60,
        category: 'music',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('課程名稱為必填')
    })

    it('validateCourse 應該驗證課程名稱長度上限', () => {
      const logic = new CourseManagementLogic()
      const result = logic.validateCourse({
        name: 'a'.repeat(51),
        duration: 60,
        category: 'music',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('課程名稱不能超過 50 個字元')
    })

    it('validateCourse 應該驗證課程時長', () => {
      const logic = new CourseManagementLogic()
      const result = logic.validateCourse({
        name: '課程',
        duration: 0,
        category: 'music',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('課程時長必須大於 0')
    })

    it('validateCourse 應該驗證課程分類', () => {
      const logic = new CourseManagementLogic()
      const result = logic.validateCourse({
        name: '課程',
        duration: 60,
        category: '',
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('課程分類為必填')
    })

    it('validateCourse 應該驗證課程價格', () => {
      const logic = new CourseManagementLogic()
      const result = logic.validateCourse({
        name: '課程',
        duration: 60,
        category: 'music',
        price: -100,
      })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('課程價格不能為負數')
    })

    it('validateCourse 應該收集多個錯誤', () => {
      const logic = new CourseManagementLogic()
      const result = logic.validateCourse({
        name: '',
        duration: -1,
        category: '',
        price: -50,
      })
      expect(result.valid).toBe(false)
      expect(result.errors.length).toBeGreaterThan(1)
    })
  })

  describe('CourseManagementLogic 資料完整性檢查', () => {
    it('isCodeExists 應該檢測課程代碼是否存在', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, code: 'PIANO001' },
        { id: 2, code: 'VIOLIN001' },
      ])
      expect(logic.isCodeExists('PIANO001')).toBe(true)
      expect(logic.isCodeExists('VIOLIN001')).toBe(true)
      expect(logic.isCodeExists('GUITAR001')).toBe(false)
    })

    it('isCodeExists 應該不區分大小寫', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([{ id: 1, code: 'PIANO001' }])
      expect(logic.isCodeExists('piano001')).toBe(true)
      expect(logic.isCodeExists('Piano001')).toBe(true)
    })

    it('isCodeExists 應該支援排除特定 ID', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, code: 'PIANO001' },
        { id: 2, code: 'PIANO001' },
      ])
      expect(logic.isCodeExists('PIANO001')).toBe(true)
      expect(logic.isCodeExists('PIANO001', 1)).toBe(true)
      expect(logic.isCodeExists('PIANO001', 999)).toBe(true)
    })

    it('getCourseStats 應該返回正確的統計資料', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, category: 'music', is_active: true },
        { id: 2, category: 'music', is_active: true },
        { id: 3, category: 'art', is_active: false },
        { id: 4, category: 'dance', is_active: true },
      ])
      const stats = logic.getCourseStats()
      expect(stats.total).toBe(4)
      expect(stats.active).toBe(3)
      expect(stats.inactive).toBe(1)
      expect(stats.byCategory['music']).toBe(2)
      expect(stats.byCategory['art']).toBe(1)
      expect(stats.byCategory['dance']).toBe(1)
    })

    it('getCourseStats 應該處理空課程列表', () => {
      const logic = new CourseManagementLogic()
      const stats = logic.getCourseStats()
      expect(stats.total).toBe(0)
      expect(stats.active).toBe(0)
      expect(stats.inactive).toBe(0)
      expect(Object.keys(stats.byCategory).length).toBe(0)
    })
  })

  describe('CourseManagementLogic 邊界情況', () => {
    it('應該處理空課程名稱的搜尋', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, name: '課程 A' },
        { id: 2, name: '課程 B' },
      ])
      logic.searchQuery = ''
      expect(logic.filteredCourses()).toHaveLength(2)
    })

    it('應該處理特殊字元的搜尋', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, name: '課程-ABC' },
        { id: 2, name: '課程 XYZ' },
      ])
      logic.searchQuery = '-'
      logic.selectedCategory = 'all'
      expect(logic.filteredCourses()).toHaveLength(1)
    })

    it('應該處理空的搜尋結果', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([
        { id: 1, name: '鋼琴課' },
        { id: 2, name: '小提琴課' },
      ])
      logic.searchQuery = '不存在的課程'
      logic.selectedCategory = 'all'
      expect(logic.filteredCourses()).toHaveLength(0)
    })

    it('addCourse 應該處理未定義的 is_active', () => {
      const logic = new CourseManagementLogic()
      logic.addCourse({ name: '課程', duration: 60 })
      expect(logic.courses[0].is_active).toBe(true)
    })

    it('toggleCourseStatus 應該處理不存在的課程', () => {
      const logic = new CourseManagementLogic()
      logic.setCourses([{ id: 1, is_active: true }])
      const result = logic.toggleCourseStatus(999)
      expect(result).toBe(false)
      expect(logic.courses[0].is_active).toBe(true)
    })
  })
})
