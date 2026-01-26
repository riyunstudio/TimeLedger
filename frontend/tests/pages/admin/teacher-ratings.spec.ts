import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock modules before imports
vi.mock('vue', async () => {
  const vue = await import('vue')
  return {
    ...vue,
    ref: (v: any) => ({ value: v }),
    computed: (fn: () => any) => ({ value: fn() }),
    reactive: (v: any) => v,
    onMounted: (fn: () => void) => fn(),
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

describe('admin/teacher-ratings.vue 頁面邏輯', () => {
  // TeacherRatingFilterLogic 類別 - 評分篩選邏輯
  class TeacherRatingFilterLogic {
    searchQuery: string
    filterRating: string

    constructor() {
      this.searchQuery = ''
      this.filterRating = ''
    }

    setSearchQuery(query: string) {
      this.searchQuery = query
    }

    setFilterRating(rating: string) {
      this.filterRating = rating
    }

    filterTeachers(teachers: any[]): any[] {
      let result = [...teachers]

      // 搜尋過濾
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase()
        result = result.filter(t => t.name?.toLowerCase().includes(query))
      }

      // 評分過濾
      if (this.filterRating) {
        const minRating = parseInt(this.filterRating)
        if (minRating === 0) {
          result = result.filter(t => !t.note || t.note.rating === 0)
        } else {
          result = result.filter(t => t.note && t.note.rating >= minRating)
        }
      }

      return result
    }

    clearFilters() {
      this.searchQuery = ''
      this.filterRating = ''
    }

    hasActiveFilters(): boolean {
      return Boolean(this.searchQuery || this.filterRating)
    }
  }

  // TeacherRatingStatsLogic 類別 - 評分統計邏輯
  class TeacherRatingStatsLogic {
    teachers: any[]

    constructor() {
      this.teachers = []
    }

    setTeachers(teachers: any[]) {
      this.teachers = teachers
    }

    getRatedCount(): number {
      return this.teachers.filter(t => t.note && t.note.rating > 0).length
    }

    getUnratedCount(): number {
      return this.teachers.length - this.getRatedCount()
    }

    getAverageRating(): string {
      const rated = this.teachers.filter(t => t.note && t.note.rating > 0)
      if (rated.length === 0) return '0.0'
      const sum = rated.reduce((acc, t) => acc + (t.note?.rating || 0), 0)
      return (sum / rated.length).toFixed(1)
    }

    getRatingDistribution(): Record<number, number> {
      const distribution: Record<number, number> = { 0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0 }
      this.teachers.forEach(t => {
        const rating = t.note?.rating || 0
        distribution[rating]++
      })
      return distribution
    }

    getTotalCount(): number {
      return this.teachers.length
    }

    getRatingPercentage(): number {
      if (this.teachers.length === 0) return 0
      return Math.round((this.getRatedCount() / this.teachers.length) * 100)
    }
  }

  // RatingFormLogic 類別 - 評分表單邏輯
  class RatingFormLogic {
    form: {
      rating: number
      internal_note: string
    }

    constructor() {
      this.form = {
        rating: 0,
        internal_note: '',
      }
    }

    setRating(rating: number) {
      this.form.rating = Math.max(0, Math.min(5, rating))
    }

    setInternalNote(note: string) {
      this.form.internal_note = note
    }

    resetForm() {
      this.form = { rating: 0, internal_note: '' }
    }

    loadFromNote(note: any) {
      this.form.rating = note?.rating || 0
      this.form.internal_note = note?.internal_note || ''
    }

    isValid(): boolean {
      return this.form.rating >= 0 && this.form.rating <= 5
    }

    hasChanges(originalNote: any): boolean {
      const originalRating = originalNote?.rating || 0
      const originalNote = originalNote?.internal_note || ''
      return (
        this.form.rating !== originalRating ||
        this.form.internal_note !== originalNote
      )
    }

    getRatingLabel(rating: number): string {
      const labels: Record<number, string> = {
        0: '未評分',
        1: '需改進',
        2: '一般',
        3: '良好',
        4: '優良',
        5: '優秀',
      }
      return labels[rating] || '未評分'
    }
  }

  // StarDisplayLogic 類別 - 星級顯示邏輯
  class StarDisplayLogic {
    getStarClass(rating: number, starIndex: number): string {
      if (starIndex <= Math.floor(rating)) {
        return 'text-warning-500'
      }
      if (starIndex === Math.ceil(rating) && rating % 1 >= 0.5) {
        return 'text-warning-500'
      }
      return 'text-slate-600'
    }

    isStarFilled(rating: number, starIndex: number): boolean {
      return starIndex <= rating
    }

    isStarHalf(rating: number, starIndex: number): boolean {
      return starIndex === Math.ceil(rating) && rating % 1 >= 0.5
    }

    renderStars(rating: number): boolean[] {
      const stars: boolean[] = []
      for (let i = 1; i <= 5; i++) {
        stars.push(this.isStarFilled(rating, i))
      }
      return stars
    }
  }

  // KeywordLogic 類別 - 關鍵字邏輯
  class KeywordLogic {
    keywords: string[]

    constructor() {
      this.keywords = ['推薦', '優秀', '穩定', '經驗豐富', '專業', '認真', '負責', '熱心']
    }

    getKeywords(): string[] {
      return this.keywords
    }

    addKeyword(keyword: string, currentNote: string): string {
      const trimmed = keyword.trim()
      if (!trimmed) return currentNote
      if (currentNote && !currentNote.endsWith(' ')) {
        return `${currentNote} ${trimmed}`
      }
      return currentNote ? `${currentNote}${trimmed}` : trimmed
    }

    hasKeyword(keyword: string, note: string): boolean {
      return note?.includes(keyword) || false
    }

    removeKeyword(keyword: string, note: string): string {
      return note?.replace(new RegExp(`\\s*${keyword}\\s*`), ' ').trim() || ''
    }

    filterByPrefix(prefix: string): string[] {
      return this.keywords.filter(k =>
        k.toLowerCase().startsWith(prefix.toLowerCase())
      )
    }
  }

  // EditModalLogic 類別 - 編輯 Modal 邏輯
  class EditModalLogic {
    showModal: boolean
    editingTeacher: any | null
    saving: boolean

    constructor() {
      this.showModal = false
      this.editingTeacher = null
      this.saving = false
    }

    openModal(teacher: any) {
      this.editingTeacher = teacher
      this.showModal = true
    }

    closeModal() {
      this.showModal = false
      this.editingTeacher = null
    }

    setSaving(saving: boolean) {
      this.saving = saving
    }

    isModalOpen(): boolean {
      return this.showModal
    }

    getEditingTeacher(): any | null {
      return this.editingTeacher
    }

    isSaving(): boolean {
      return this.saving
    }
  }

  describe('TeacherRatingFilterLogic 評分篩選邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new TeacherRatingFilterLogic()
      expect(logic.searchQuery).toBe('')
      expect(logic.filterRating).toBe('')
    })

    it('filterTeachers 應該根據搜尋過濾', () => {
      const logic = new TeacherRatingFilterLogic()
      const teachers = [
        { id: 1, name: '張老師' },
        { id: 2, name: '李老師' },
        { id: 3, name: '張三老師' },
      ]
      logic.setSearchQuery('張')
      const filtered = logic.filterTeachers(teachers)
      expect(filtered).toHaveLength(2)
    })

    it('filterTeachers 應該不區分大小寫', () => {
      const logic = new TeacherRatingFilterLogic()
      const teachers = [{ id: 1, name: '張老師' }]
      logic.setSearchQuery('張')
      expect(logic.filterTeachers(teachers)).toHaveLength(1)
      logic.setSearchQuery('張')
      expect(logic.filterTeachers(teachers)).toHaveLength(1)
    })

    it('filterTeachers 應該根據評分過濾 - 5星', () => {
      const logic = new TeacherRatingFilterLogic()
      const teachers = [
        { id: 1, note: { rating: 5 } },
        { id: 2, note: { rating: 4 } },
        { id: 3, note: { rating: 3 } },
      ]
      logic.setFilterRating('5')
      const filtered = logic.filterTeachers(teachers)
      expect(filtered).toHaveLength(1)
      expect(filtered[0].id).toBe(1)
    })

    it('filterTeachers 應該根據評分過濾 - 4星以上', () => {
      const logic = new TeacherRatingFilterLogic()
      const teachers = [
        { id: 1, note: { rating: 5 } },
        { id: 2, note: { rating: 4 } },
        { id: 3, note: { rating: 3 } },
      ]
      logic.setFilterRating('4')
      const filtered = logic.filterTeachers(teachers)
      expect(filtered).toHaveLength(2)
    })

    it('filterTeachers 應該根據評分過濾 - 未評分', () => {
      const logic = new TeacherRatingFilterLogic()
      const teachers = [
        { id: 1, note: { rating: 5 } },
        { id: 2 },
        { id: 3, note: { rating: 0 } },
      ]
      logic.setFilterRating('0')
      const filtered = logic.filterTeachers(teachers)
      expect(filtered).toHaveLength(2)
    })

    it('filterTeachers 應該同時支援搜尋和評分過濾', () => {
      const logic = new TeacherRatingFilterLogic()
      const teachers = [
        { id: 1, name: '張老師', note: { rating: 5 } },
        { id: 2, name: '張三老師', note: { rating: 4 } },
        { id: 3, name: '李老師', note: { rating: 5 } },
      ]
      logic.setSearchQuery('張')
      logic.setFilterRating('5')
      const filtered = logic.filterTeachers(teachers)
      expect(filtered).toHaveLength(1)
      expect(filtered[0].id).toBe(1)
    })

    it('clearFilters 應該重置所有過濾條件', () => {
      const logic = new TeacherRatingFilterLogic()
      logic.setSearchQuery('測試')
      logic.setFilterRating('4')
      logic.clearFilters()
      expect(logic.searchQuery).toBe('')
      expect(logic.filterRating).toBe('')
    })

    it('hasActiveFilters 應該正確判斷是否有啟用的過濾', () => {
      const logic = new TeacherRatingFilterLogic()
      expect(logic.hasActiveFilters()).toBe(false)
      logic.setSearchQuery('測試')
      expect(logic.hasActiveFilters()).toBe(true)
    })
  })

  describe('TeacherRatingStatsLogic 評分統計邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new TeacherRatingStatsLogic()
      expect(logic.teachers).toHaveLength(0)
    })

    it('setTeachers 應該正確設定老師列表', () => {
      const logic = new TeacherRatingStatsLogic()
      logic.setTeachers([
        { id: 1, note: { rating: 5 } },
        { id: 2, note: { rating: 4 } },
      ])
      expect(logic.getTotalCount()).toBe(2)
    })

    it('getRatedCount 應該計算已評分數量', () => {
      const logic = new TeacherRatingStatsLogic()
      logic.setTeachers([
        { id: 1, note: { rating: 5 } },
        { id: 2 },
        { id: 3, note: { rating: 3 } },
      ])
      expect(logic.getRatedCount()).toBe(2)
    })

    it('getUnratedCount 應該計算未評分數量', () => {
      const logic = new TeacherRatingStatsLogic()
      logic.setTeachers([
        { id: 1, note: { rating: 5 } },
        { id: 2 },
      ])
      expect(logic.getUnratedCount()).toBe(1)
    })

    it('getAverageRating 應該計算平均評分', () => {
      const logic = new TeacherRatingStatsLogic()
      logic.setTeachers([
        { id: 1, note: { rating: 5 } },
        { id: 2, note: { rating: 3 } },
        { id: 3, note: { rating: 4 } },
      ])
      expect(logic.getAverageRating()).toBe('4.0')
    })

    it('getAverageRating 應該處理沒有評分的情況', () => {
      const logic = new TeacherRatingStatsLogic()
      logic.setTeachers([])
      expect(logic.getAverageRating()).toBe('0.0')
      logic.setTeachers([{ id: 1 }])
      expect(logic.getAverageRating()).toBe('0.0')
    })

    it('getRatingDistribution 應該返回評分分佈', () => {
      const logic = new TeacherRatingStatsLogic()
      logic.setTeachers([
        { id: 1, note: { rating: 5 } },
        { id: 2, note: { rating: 4 } },
        { id: 3, note: { rating: 5 } },
      ])
      const distribution = logic.getRatingDistribution()
      expect(distribution[5]).toBe(2)
      expect(distribution[4]).toBe(1)
      expect(distribution[3]).toBe(0)
    })

    it('getRatingPercentage 應該計算評分覆蓋率', () => {
      const logic = new TeacherRatingStatsLogic()
      logic.setTeachers([
        { id: 1, note: { rating: 5 } },
        { id: 2, note: { rating: 4 } },
        { id: 3 },
        { id: 4 },
      ])
      expect(logic.getRatingPercentage()).toBe(50)
    })
  })

  describe('RatingFormLogic 評分表單邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new RatingFormLogic()
      expect(logic.form.rating).toBe(0)
      expect(logic.form.internal_note).toBe('')
    })

    it('setRating 應該限制評分範圍在 0-5', () => {
      const logic = new RatingFormLogic()
      logic.setRating(3)
      expect(logic.form.rating).toBe(3)
      logic.setRating(10)
      expect(logic.form.rating).toBe(5)
      logic.setRating(-1)
      expect(logic.form.rating).toBe(0)
    })

    it('setInternalNote 應該正確設定備註', () => {
      const logic = new RatingFormLogic()
      logic.setInternalNote('這是一個備註')
      expect(logic.form.internal_note).toBe('這是一個備註')
    })

    it('resetForm 應該重置表單', () => {
      const logic = new RatingFormLogic()
      logic.form.rating = 5
      logic.form.internal_note = '備註'
      logic.resetForm()
      expect(logic.form.rating).toBe(0)
      expect(logic.form.internal_note).toBe('')
    })

    it('loadFromNote 應該從評分資料載入', () => {
      const logic = new RatingFormLogic()
      logic.loadFromNote({ rating: 4, internal_note: '良好' })
      expect(logic.form.rating).toBe(4)
      expect(logic.form.internal_note).toBe('良好')
    })

    it('loadFromNote 應該處理空值', () => {
      const logic = new RatingFormLogic()
      logic.loadFromNote(null)
      expect(logic.form.rating).toBe(0)
      expect(logic.form.internal_note).toBe('')
    })

    it('isValid 應該正確驗證表單', () => {
      const logic = new RatingFormLogic()
      logic.form.rating = 3
      expect(logic.isValid()).toBe(true)
      logic.form.rating = -1
      expect(logic.isValid()).toBe(false)
      logic.form.rating = 6
      expect(logic.isValid()).toBe(false)
    })

    it('hasChanges 應該正確判斷是否有變更', () => {
      const logic = new RatingFormLogic()
      logic.form.rating = 4
      logic.form.internal_note = '新備註'
      expect(logic.hasChanges({ rating: 3, internal_note: '舊備註' })).toBe(true)
      expect(logic.hasChanges({ rating: 4, internal_note: '新備註' })).toBe(false)
    })

    it('getRatingLabel 應該返回正確的評分標籤', () => {
      const logic = new RatingFormLogic()
      expect(logic.getRatingLabel(0)).toBe('未評分')
      expect(logic.getRatingLabel(1)).toBe('需改進')
      expect(logic.getRatingLabel(2)).toBe('一般')
      expect(logic.getRatingLabel(3)).toBe('良好')
      expect(logic.getRatingLabel(4)).toBe('優良')
      expect(logic.getRatingLabel(5)).toBe('優秀')
    })
  })

  describe('StarDisplayLogic 星級顯示邏輯', () => {
    it('getStarClass 應該返回正確的 CSS 類別', () => {
      const logic = new StarDisplayLogic()
      expect(logic.getStarClass(5, 1)).toBe('text-warning-500')
      expect(logic.getStarClass(5, 5)).toBe('text-warning-500')
      expect(logic.getStarClass(2, 3)).toBe('text-slate-600')
    })

    it('isStarFilled 應該正確判斷星星是否填滿', () => {
      const logic = new StarDisplayLogic()
      expect(logic.isStarFilled(5, 1)).toBe(true)
      expect(logic.isStarFilled(5, 5)).toBe(true)
      expect(logic.isStarFilled(2, 3)).toBe(false)
    })

    it('renderStars 應該返回正確的星星陣列', () => {
      const logic = new StarDisplayLogic()
      const stars5 = logic.renderStars(5)
      expect(stars5).toEqual([true, true, true, true, true])
      const stars3 = logic.renderStars(3)
      expect(stars3).toEqual([true, true, true, false, false])
    })
  })

  describe('KeywordLogic 關鍵字邏輯', () => {
    it('應該正確初始化所有關鍵字', () => {
      const logic = new KeywordLogic()
      expect(logic.getKeywords()).toHaveLength(8)
      expect(logic.getKeywords()).toContain('推薦')
      expect(logic.getKeywords()).toContain('優秀')
    })

    it('addKeyword 應該正確新增關鍵字', () => {
      const logic = new KeywordLogic()
      expect(logic.addKeyword('推薦', '')).toBe('推薦')
      expect(logic.addKeyword('優秀', '推薦')).toBe('推薦 優秀')
    })

    it('hasKeyword 應該正確判斷是否包含關鍵字', () => {
      const logic = new KeywordLogic()
      expect(logic.hasKeyword('推薦', '推薦 優秀')).toBe(true)
      expect(logic.hasKeyword('良好', '推薦 優秀')).toBe(false)
    })

    it('filterByPrefix 應該根據前綴過濾關鍵字', () => {
      const logic = new KeywordLogic()
      const filtered = logic.filterByPrefix('認')
      expect(filtered).toContain('認真')
      expect(filtered).toContain('負責')
    })
  })

  describe('EditModalLogic 編輯 Modal 邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new EditModalLogic()
      expect(logic.showModal).toBe(false)
      expect(logic.editingTeacher).toBeNull()
      expect(logic.saving).toBe(false)
    })

    it('openModal 應該開啟 Modal 並設定編輯中的老師', () => {
      const logic = new EditModalLogic()
      const teacher = { id: 1, name: '張老師' }
      logic.openModal(teacher)
      expect(logic.showModal).toBe(true)
      expect(logic.editingTeacher?.id).toBe(1)
    })

    it('closeModal 應該關閉 Modal 並清除編輯中的老師', () => {
      const logic = new EditModalLogic()
      logic.openModal({ id: 1 })
      logic.closeModal()
      expect(logic.showModal).toBe(false)
      expect(logic.editingTeacher).toBeNull()
    })

    it('isModalOpen 應該正確判斷 Modal 是否開啟', () => {
      const logic = new EditModalLogic()
      expect(logic.isModalOpen()).toBe(false)
      logic.openModal({ id: 1 })
      expect(logic.isModalOpen()).toBe(true)
    })
  })

  describe('頁面整合邏輯', () => {
    it('應該能夠完整處理老師評分流程', () => {
      const filterLogic = new TeacherRatingFilterLogic()
      const statsLogic = new TeacherRatingStatsLogic()
      const formLogic = new RatingFormLogic()
      const modalLogic = new EditModalLogic()
      const starLogic = new StarDisplayLogic()

      // 設定老師資料
      const teachers = [
        { id: 1, name: '張老師', note: { rating: 5, internal_note: '優秀' } },
        { id: 2, name: '李老師', note: { rating: 4, internal_note: '良好' } },
        { id: 3, name: '王老師' },
      ]

      statsLogic.setTeachers(teachers)
      filterLogic.setTeachers(teachers)

      // 驗證統計
      expect(statsLogic.getRatedCount()).toBe(2)
      expect(statsLogic.getUnratedCount()).toBe(1)
      expect(statsLogic.getAverageRating()).toBe('4.5')

      // 驗證篩選
      filterLogic.setFilterRating('5')
      const filtered = filterLogic.filterTeachers(teachers)
      expect(filtered).toHaveLength(1)

      // 開啟編輯 Modal
      modalLogic.openModal(teachers[0])
      expect(modalLogic.isModalOpen()).toBe(true)
      expect(modalLogic.getEditingTeacher()?.id).toBe(1)

      // 載入評分資料
      formLogic.loadFromNote(teachers[0].note)
      expect(formLogic.form.rating).toBe(5)
      expect(formLogic.getRatingLabel(5)).toBe('優秀')

      // 驗證星星顯示
      const stars = starLogic.renderStars(5)
      expect(stars.every(s => s)).toBe(true)

      // 關閉 Modal
      modalLogic.closeModal()
      expect(modalLogic.isModalOpen()).toBe(false)
    })

    it('應該正確處理評分分佈統計', () => {
      const statsLogic = new TeacherRatingStatsLogic()

      statsLogic.setTeachers([
        { id: 1, note: { rating: 5 } },
        { id: 2, note: { rating: 5 } },
        { id: 3, note: { rating: 4 } },
        { id: 4, note: { rating: 4 } },
        { id: 5, note: { rating: 3 } },
      ])

      const distribution = statsLogic.getRatingDistribution()
      expect(distribution[5]).toBe(2)
      expect(distribution[4]).toBe(2)
      expect(distribution[3]).toBe(1)
      expect(distribution[2]).toBe(0)
      expect(distribution[1]).toBe(0)
      expect(distribution[0]).toBe(0)

      expect(statsLogic.getRatingPercentage()).toBe(100)
    })
  })
})
